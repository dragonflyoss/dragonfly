/*
 *     Copyright 2026 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"d7y.io/dragonfly/v2/datacontroller/apis/data/v1alpha1"
)

// Condition types and reasons reported on a DataLifecyclePolicy.
const (
	ConditionPolicyValid = "Valid"

	ReasonPolicyValid   = "Valid"
	ReasonPolicyInvalid = "Invalid"
)

// DataLifecyclePolicyReconciler validates DataLifecyclePolicy objects and reports how
// many datasets they govern. The policies are applied by the Dataset reconciler.
type DataLifecyclePolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=data.d7y.io,resources=datalifecyclepolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=data.d7y.io,resources=datalifecyclepolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=data.d7y.io,resources=datasets,verbs=get;list;watch

// SetupWithManager registers the reconciler with the manager.
func (r *DataLifecyclePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("datalifecyclepolicy").
		For(&v1alpha1.DataLifecyclePolicy{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(&v1alpha1.Dataset{}, handler.EnqueueRequestsFromMapFunc(r.policiesForDataset),
			builder.WithPredicates(predicate.Or(predicate.LabelChangedPredicate{}, predicate.Funcs{
				UpdateFunc: func(event.UpdateEvent) bool { return false },
			}))).
		Complete(r)
}

// policiesForDataset enqueues every policy in the namespace of a dataset so their
// matched counts stay current.
func (r *DataLifecyclePolicyReconciler) policiesForDataset(ctx context.Context, obj client.Object) []reconcile.Request {
	var policies v1alpha1.DataLifecyclePolicyList
	if err := r.List(ctx, &policies, client.InNamespace(obj.GetNamespace())); err != nil {
		ctrl.LoggerFrom(ctx).Error(err, "list policies for dataset", "dataset", obj.GetName())
		return nil
	}

	requests := make([]reconcile.Request, 0, len(policies.Items))
	for _, policy := range policies.Items {
		requests = append(requests, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(&policy)})
	}

	return requests
}

// Reconcile validates the policy and counts the datasets it governs.
func (r *DataLifecyclePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var policy v1alpha1.DataLifecyclePolicy
	if err := r.Get(ctx, req.NamespacedName, &policy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	original := policy.DeepCopy()
	policy.Status.ObservedGeneration = policy.Generation

	if _, err := ResolveLifecycle(nil, &policy); err != nil {
		policy.Status.MatchedDatasets = 0
		meta.SetStatusCondition(&policy.Status.Conditions, metav1.Condition{
			Type:               ConditionPolicyValid,
			Status:             metav1.ConditionFalse,
			Reason:             ReasonPolicyInvalid,
			Message:            err.Error(),
			ObservedGeneration: policy.Generation,
		})

		return ctrl.Result{}, r.patchStatus(ctx, original, &policy)
	}

	var policies v1alpha1.DataLifecyclePolicyList
	if err := r.List(ctx, &policies, client.InNamespace(policy.Namespace)); err != nil {
		return ctrl.Result{}, fmt.Errorf("list policies: %w", err)
	}

	var datasets v1alpha1.DatasetList
	if err := r.List(ctx, &datasets, client.InNamespace(policy.Namespace)); err != nil {
		return ctrl.Result{}, fmt.Errorf("list datasets: %w", err)
	}

	var matched int32
	for i := range datasets.Items {
		selected, err := SelectPolicy(policies.Items, &datasets.Items[i])
		if err != nil {
			// Another policy in the namespace is invalid; it is reported on that policy.
			continue
		}

		if selected != nil && selected.Name == policy.Name {
			matched++
		}
	}

	policy.Status.MatchedDatasets = matched
	meta.SetStatusCondition(&policy.Status.Conditions, metav1.Condition{
		Type:               ConditionPolicyValid,
		Status:             metav1.ConditionTrue,
		Reason:             ReasonPolicyValid,
		Message:            fmt.Sprintf("policy governs %d dataset(s)", matched),
		ObservedGeneration: policy.Generation,
	})

	return ctrl.Result{}, r.patchStatus(ctx, original, &policy)
}

func (r *DataLifecyclePolicyReconciler) patchStatus(ctx context.Context, original, policy *v1alpha1.DataLifecyclePolicy) error {
	if equality.Semantic.DeepEqual(original.Status, policy.Status) {
		return nil
	}

	if err := r.Status().Patch(ctx, policy, client.MergeFrom(original)); err != nil {
		return client.IgnoreNotFound(fmt.Errorf("patch status: %w", err))
	}

	return nil
}
