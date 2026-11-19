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
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/events"
	"k8s.io/utils/clock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"d7y.io/dragonfly/v2/datacontroller/apis/data/v1alpha1"
	"d7y.io/dragonfly/v2/datacontroller/config"
	"d7y.io/dragonfly/v2/datacontroller/dragonfly"
	"d7y.io/dragonfly/v2/datacontroller/metrics"
	managertypes "d7y.io/dragonfly/v2/manager/types"
)

// Reasons reported in conditions, source status and events.
const (
	ReasonPending            = "Pending"
	ReasonInvalidSpec        = "InvalidSpec"
	ReasonInvalidLifecycle   = "InvalidLifecycle"
	ReasonManagerUnavailable = "ManagerUnavailable"
	ReasonSuspended          = "Suspended"
	ReasonDistributing       = "Distributing"
	ReasonDistributed        = "Distributed"
	ReasonDistributionFailed = "DistributionFailed"
	ReasonRetryScheduled     = "RetryScheduled"
	ReasonVerifying          = "Verifying"
	ReasonVerified           = "Verified"
	ReasonVerificationFailed = "VerificationFailed"
	ReasonDataMissing        = "DataMissing"
	ReasonExpired            = "Expired"
	ReasonPurging            = "Purging"
	ReasonPurged             = "Purged"
	ReasonPurgeFailed        = "PurgeFailed"
	ReasonNotVerified        = "NotVerified"
	ReasonNotExpired         = "NotExpired"
)

const (
	// maxJobsPerReconcile bounds the number of jobs created for one source in one
	// reconcile so a large image does not trip the manager's job rate limit.
	maxJobsPerReconcile = 10

	// minRequeue is the shortest requeue delay.
	minRequeue = time.Second
)

// ManagerClientFactory creates manager clients. It is a variable so tests can inject
// a fake manager.
type ManagerClientFactory func(opts dragonfly.Options) (dragonfly.Client, error)

// DatasetReconciler reconciles Dataset objects.
type DatasetReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder events.EventRecorder
	Config   *config.Config
	Clock    clock.PassiveClock

	// NewManagerClient creates manager clients. Defaults to dragonfly.New.
	NewManagerClient ManagerClientFactory

	clients sync.Map
}

// +kubebuilder:rbac:groups=data.d7y.io,resources=datasets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=data.d7y.io,resources=datasets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=data.d7y.io,resources=datasets/finalizers,verbs=update
// +kubebuilder:rbac:groups=data.d7y.io,resources=datalifecyclepolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// SetupWithManager registers the reconciler with the manager.
func (r *DatasetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Clock == nil {
		r.Clock = clock.RealClock{}
	}

	if r.NewManagerClient == nil {
		r.NewManagerClient = dragonfly.New
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("dataset").
		For(&v1alpha1.Dataset{}, builder.WithPredicates(builderPredicates())).
		Watches(&v1alpha1.DataLifecyclePolicy{}, handler.EnqueueRequestsFromMapFunc(r.datasetsForPolicy)).
		WithOptions(controller.Options{MaxConcurrentReconciles: r.Config.Controller.MaxConcurrentReconciles}).
		Complete(r)
}

// builderPredicates reconciles datasets on spec, label, finalizer and deletion changes
// but not on the status updates the reconciler itself makes.
func builderPredicates() predicate.Predicate {
	return predicate.Or(
		predicate.GenerationChangedPredicate{},
		predicate.LabelChangedPredicate{},
		predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				return !e.ObjectNew.GetDeletionTimestamp().IsZero() ||
					len(e.ObjectNew.GetFinalizers()) != len(e.ObjectOld.GetFinalizers())
			},
			CreateFunc:  func(event.CreateEvent) bool { return true },
			DeleteFunc:  func(event.DeleteEvent) bool { return true },
			GenericFunc: func(event.GenericEvent) bool { return true },
		},
	)
}

// datasetsForPolicy enqueues every dataset in the namespace of a changed policy. Matching
// is evaluated in the reconciler, so datasets that stop matching are updated too.
func (r *DatasetReconciler) datasetsForPolicy(ctx context.Context, obj client.Object) []reconcile.Request {
	var datasets v1alpha1.DatasetList
	if err := r.List(ctx, &datasets, client.InNamespace(obj.GetNamespace())); err != nil {
		ctrl.LoggerFrom(ctx).Error(err, "list datasets for policy", "policy", obj.GetName())
		return nil
	}

	requests := make([]reconcile.Request, 0, len(datasets.Items))
	for _, dataset := range datasets.Items {
		requests = append(requests, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(&dataset)})
	}

	return requests
}

// reconcileState carries the values shared by the steps of one reconcile.
type reconcileState struct {
	dataset   *v1alpha1.Dataset
	lifecycle *EffectiveLifecycle
	manager   dragonfly.Client
	secrets   *secretReader
	now       time.Time

	// dueTimes are the times the reconciler must run again at. Times in the past
	// are ignored: whatever was due has been acted upon or is blocked by an
	// operation in flight.
	dueTimes []time.Time

	// immediate is set when the reconciler must run again right away.
	immediate bool

	// poll is set when a job is in flight and must be polled.
	poll bool

	// errs are the transient errors hit; they are returned so the request is retried
	// with backoff.
	errs []error
}

func (s *reconcileState) due(t *time.Time) {
	if t != nil {
		s.dueTimes = append(s.dueTimes, *t)
	}
}

func (s *reconcileState) requeueNow() {
	s.immediate = true
}

func (s *reconcileState) fail(err error) {
	s.errs = append(s.errs, err)
}

// Reconcile drives a dataset towards its desired state.
func (r *DatasetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	var dataset v1alpha1.Dataset
	if err := r.Get(ctx, req.NamespacedName, &dataset); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if dataset.DeletionTimestamp.IsZero() && !controllerutil.ContainsFinalizer(&dataset, v1alpha1.DatasetFinalizer) {
		controllerutil.AddFinalizer(&dataset, v1alpha1.DatasetFinalizer)
		if err := r.Update(ctx, &dataset); err != nil {
			return ctrl.Result{}, fmt.Errorf("add finalizer: %w", err)
		}

		// The update triggers another reconcile.
		return ctrl.Result{}, nil
	}

	original := dataset.DeepCopy()
	state := &reconcileState{
		dataset: &dataset,
		now:     r.Clock.Now(),
		secrets: newSecretReader(r.Client, dataset.Namespace),
	}

	r.reconcile(ctx, state)

	dataset.Status.ObservedGeneration = dataset.Generation
	if !equality.Semantic.DeepEqual(original.Status, dataset.Status) {
		if err := r.Status().Patch(ctx, &dataset, client.MergeFromWithOptions(original, client.MergeFromWithOptimisticLock{})); err != nil {
			if client.IgnoreNotFound(err) == nil {
				return ctrl.Result{}, nil
			}

			return ctrl.Result{}, fmt.Errorf("patch status: %w", err)
		}
	}

	if !dataset.DeletionTimestamp.IsZero() && r.deletionComplete(state) {
		controllerutil.RemoveFinalizer(&dataset, v1alpha1.DatasetFinalizer)
		if err := r.Update(ctx, &dataset); err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(fmt.Errorf("remove finalizer: %w", err))
		}

		log.Info("dataset deleted")
		return ctrl.Result{}, nil
	}

	if len(state.errs) > 0 {
		return ctrl.Result{}, errors.Join(state.errs...)
	}

	return ctrl.Result{RequeueAfter: r.requeueAfter(state)}, nil
}

// reconcile runs the steps of a reconcile and records their outcome in the status.
func (r *DatasetReconciler) reconcile(ctx context.Context, state *reconcileState) {
	dataset := state.dataset

	lifecycle, err := r.resolveLifecycle(ctx, dataset)
	if err != nil {
		state.fail(err)
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, ReasonInvalidLifecycle, err.Error())
		dataset.Status.Phase = v1alpha1.DatasetPhaseFailed
		return
	}

	state.lifecycle = lifecycle
	dataset.Status.AppliedPolicy = lifecycle.PolicyName

	if err := validateSpec(dataset); err != nil {
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, ReasonInvalidSpec, err.Error())
		dataset.Status.Phase = v1alpha1.DatasetPhaseFailed
		if !dataset.DeletionTimestamp.IsZero() {
			// An invalid spec must not block deletion; purge what was distributed.
			r.reconcileDeletion(ctx, state)
		}

		return
	}

	manager, err := r.managerClient(ctx, state)
	if err != nil {
		state.fail(err)
		r.event(dataset, corev1.EventTypeWarning, ReasonManagerUnavailable, "%s", err.Error())
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, ReasonManagerUnavailable, err.Error())
		return
	}

	state.manager = manager

	if !dataset.DeletionTimestamp.IsZero() {
		r.reconcileDeletion(ctx, state)
		return
	}

	r.reconcileSources(ctx, state)
	r.summarize(state)
}

// validateSpec checks what the CRD schema cannot express.
func validateSpec(dataset *v1alpha1.Dataset) error {
	for i := range dataset.Spec.Sources {
		if err := validateSource(&dataset.Spec.Sources[i]); err != nil {
			return fmt.Errorf("source %q: %w", dataset.Spec.Sources[i].Name, err)
		}
	}

	return nil
}

// resolveLifecycle finds the policy governing the dataset and merges it with the
// inline lifecycle.
func (r *DatasetReconciler) resolveLifecycle(ctx context.Context, dataset *v1alpha1.Dataset) (*EffectiveLifecycle, error) {
	var policies v1alpha1.DataLifecyclePolicyList
	if err := r.List(ctx, &policies, client.InNamespace(dataset.Namespace)); err != nil {
		return nil, fmt.Errorf("list lifecycle policies: %w", err)
	}

	policy, err := SelectPolicy(policies.Items, dataset)
	if err != nil {
		return nil, err
	}

	return ResolveLifecycle(dataset.Spec.Lifecycle, policy)
}

// managerClient returns the client of the manager the dataset is served by.
func (r *DatasetReconciler) managerClient(ctx context.Context, state *reconcileState) (dragonfly.Client, error) {
	opts := dragonfly.Options{
		Endpoint:              r.Config.Manager.Endpoint,
		Token:                 r.Config.Manager.Token,
		InsecureSkipTLSVerify: r.Config.Manager.InsecureSkipTLSVerify,
		Timeout:               r.Config.Manager.RequestTimeout,
	}

	if ref := state.dataset.Spec.ManagerRef; ref != nil {
		if ref.Endpoint != "" {
			opts.Endpoint = ref.Endpoint
			opts.InsecureSkipTLSVerify = ref.InsecureSkipTLSVerify
		}

		if ref.TokenSecretRef != nil {
			key := ref.TokenSecretRef.Key
			if key == "" {
				key = SecretKeyToken
			}

			token, err := state.secrets.value(ctx, ref.TokenSecretRef.Name, key, true)
			if err != nil {
				return nil, fmt.Errorf("manager token: %w", err)
			}

			opts.Token = token
		}
	}

	if opts.Endpoint == "" {
		return nil, errors.New("no manager endpoint: set manager.endpoint on the controller or spec.managerRef.endpoint on the dataset")
	}

	if opts.Token == "" {
		return nil, errors.New("no manager token: set manager.token on the controller or spec.managerRef.tokenSecretRef on the dataset")
	}

	cacheKey := fmt.Sprintf("%s|%s|%t|%s", opts.Endpoint, opts.Token, opts.InsecureSkipTLSVerify, opts.Timeout)
	if cached, ok := r.clients.Load(cacheKey); ok {
		return cached.(dragonfly.Client), nil
	}

	c, err := r.NewManagerClient(opts)
	if err != nil {
		return nil, err
	}

	actual, _ := r.clients.LoadOrStore(cacheKey, c)
	return actual.(dragonfly.Client), nil
}

// reconcileSources plans and drives the operations of every source.
func (r *DatasetReconciler) reconcileSources(ctx context.Context, state *reconcileState) {
	dataset := state.dataset
	lifecycle := state.lifecycle

	// Decisions use the times derived from the current status; the summary recomputes
	// them afterwards.
	lastDistribution, lastVerification := datasetTimes(dataset)
	refreshDue := isDue(lifecycle.NextRefresh(lastDistribution), state.now)
	expireDue := isDue(lifecycle.Expiration(lastDistribution), state.now)
	verifyDue := isDue(lifecycle.NextVerification(lastDistribution, lastVerification), state.now)

	existing := make(map[string]*v1alpha1.DataSourceStatus, len(dataset.Status.Sources))
	for i := range dataset.Status.Sources {
		existing[dataset.Status.Sources[i].Name] = &dataset.Status.Sources[i]
	}

	statuses := make([]v1alpha1.DataSourceStatus, 0, len(dataset.Spec.Sources))
	for i := range dataset.Spec.Sources {
		source := &dataset.Spec.Sources[i]
		status, ok := existing[source.Name]
		if !ok {
			status = &v1alpha1.DataSourceStatus{Name: source.Name, State: v1alpha1.DataSourceStatePending, Reason: ReasonPending}
		}
		delete(existing, source.Name)

		r.reconcileSource(ctx, state, source, status, sourcePlan{
			refreshDue: refreshDue,
			expireDue:  expireDue,
			verifyDue:  verifyDue,
		})
		statuses = append(statuses, *status)
	}

	// Sources removed from the spec are purged according to the deletion policy before
	// their status is dropped.
	for i := range dataset.Status.Sources {
		status := &dataset.Status.Sources[i]
		if _, orphaned := existing[status.Name]; !orphaned {
			continue
		}

		if r.retireSource(ctx, state, status) {
			continue
		}

		statuses = append(statuses, *status)
	}

	dataset.Status.Sources = statuses
}

// sourcePlan is the set of dataset-wide events due for a source.
type sourcePlan struct {
	refreshDue bool
	expireDue  bool
	verifyDue  bool
}

// reconcileSource drives one source: it finishes the operation in flight, then starts
// the next one that is due.
func (r *DatasetReconciler) reconcileSource(ctx context.Context, state *reconcileState, source *v1alpha1.DataSource, status *v1alpha1.DataSourceStatus, plan sourcePlan) {
	dataset := state.dataset

	if status.Operation != nil {
		if !r.driveOperation(ctx, state, source, status) {
			return
		}

		// The operation finished; start the next one on the next reconcile so the
		// status reflects the outcome first.
		state.requeueNow()
		return
	}

	if status.NextRetryTime != nil {
		if state.now.Before(status.NextRetryTime.Time) {
			state.due(&status.NextRetryTime.Time)
			return
		}

		status.NextRetryTime = nil
		switch status.State {
		case v1alpha1.DataSourceStatePurging:
			r.startOperation(ctx, state, source, status, v1alpha1.OperationTypePurge)
		default:
			r.startOperation(ctx, state, source, status, v1alpha1.OperationTypeDistribute)
		}

		return
	}

	if dataset.Spec.Suspend {
		return
	}

	hash := SourceSpecHash(source, &dataset.Spec.Distribution)
	if status.SpecHash != hash {
		status.SpecHash = hash
		status.State = v1alpha1.DataSourceStatePending
		status.Reason = ReasonPending
		status.Message = "spec changed"
		status.Retries = 0
	}

	switch status.State {
	case v1alpha1.DataSourceStatePending:
		r.startOperation(ctx, state, source, status, v1alpha1.OperationTypeDistribute)
	case v1alpha1.DataSourceStateReady, v1alpha1.DataSourceStateDegraded:
		switch {
		case plan.refreshDue:
			r.startOperation(ctx, state, source, status, v1alpha1.OperationTypeDistribute)
		case plan.expireDue:
			r.expireSource(ctx, state, source, status)
		case plan.verifyDue:
			r.startOperation(ctx, state, source, status, v1alpha1.OperationTypeVerify)
		}
	case v1alpha1.DataSourceStateExpired:
		if plan.refreshDue {
			r.startOperation(ctx, state, source, status, v1alpha1.OperationTypeDistribute)
		}
	}
}

// expireSource applies the expire action to a source whose TTL elapsed.
func (r *DatasetReconciler) expireSource(ctx context.Context, state *reconcileState, source *v1alpha1.DataSource, status *v1alpha1.DataSourceStatus) {
	if state.lifecycle.ExpireAction == v1alpha1.ExpireActionPurge && len(status.Tasks) > 0 {
		r.startOperation(ctx, state, source, status, v1alpha1.OperationTypePurge)
		return
	}

	status.State = v1alpha1.DataSourceStateExpired
	status.Reason = ReasonExpired
	status.Message = "ttl elapsed, data retained on peers"
	r.event(state.dataset, corev1.EventTypeNormal, ReasonExpired, "source %s expired", source.Name)
}

// retireSource purges a source that was removed from the spec. It returns true once
// the status entry can be dropped.
func (r *DatasetReconciler) retireSource(ctx context.Context, state *reconcileState, status *v1alpha1.DataSourceStatus) bool {
	// The spec no longer describes the source; rebuild what is needed to address
	// its tasks from the recorded task IDs.
	source := &v1alpha1.DataSource{Name: status.Name}
	return r.purgeSource(ctx, state, source, status)
}

// reconcileDeletion purges every source according to the deletion policy.
func (r *DatasetReconciler) reconcileDeletion(ctx context.Context, state *reconcileState) {
	dataset := state.dataset
	dataset.Status.Phase = v1alpha1.DatasetPhasePurging
	r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, ReasonPurging, "dataset is being deleted")

	if state.lifecycle == nil || state.lifecycle.DeletionPolicy != v1alpha1.DeletionPolicyPurge || state.manager == nil {
		return
	}

	sources := make(map[string]*v1alpha1.DataSource, len(dataset.Spec.Sources))
	for i := range dataset.Spec.Sources {
		sources[dataset.Spec.Sources[i].Name] = &dataset.Spec.Sources[i]
	}

	for i := range dataset.Status.Sources {
		status := &dataset.Status.Sources[i]
		source, ok := sources[status.Name]
		if !ok {
			source = &v1alpha1.DataSource{Name: status.Name}
		}

		r.purgeSource(ctx, state, source, status)
	}
}

// purgeSource drives the purge of a source. It returns true when nothing is left to do,
// either because the data is gone or because purging failed for good.
func (r *DatasetReconciler) purgeSource(ctx context.Context, state *reconcileState, source *v1alpha1.DataSource, status *v1alpha1.DataSourceStatus) bool {
	if status.Operation != nil {
		if status.Operation.Type != v1alpha1.OperationTypePurge {
			// Let a distribution or verification finish so the tasks it produces are
			// known and can be purged.
			if !r.driveOperation(ctx, state, source, status) {
				return false
			}
		} else if !r.driveOperation(ctx, state, source, status) {
			return false
		}

		state.requeueNow()
		return false
	}

	if status.NextRetryTime != nil {
		if state.now.Before(status.NextRetryTime.Time) {
			state.due(&status.NextRetryTime.Time)
			return false
		}

		status.NextRetryTime = nil
	}

	if len(status.Tasks) == 0 || status.Reason == ReasonPurgeFailed {
		return true
	}

	if state.manager == nil {
		return false
	}

	r.startOperation(ctx, state, source, status, v1alpha1.OperationTypePurge)
	return false
}

// deletionComplete reports whether every source is purged or purging failed for good.
func (r *DatasetReconciler) deletionComplete(state *reconcileState) bool {
	if state.lifecycle == nil {
		// The lifecycle could not be resolved; do not hold the object hostage.
		return true
	}

	if state.lifecycle.DeletionPolicy != v1alpha1.DeletionPolicyPurge {
		return true
	}

	if state.manager == nil {
		return false
	}

	for _, status := range state.dataset.Status.Sources {
		if status.Operation != nil || status.NextRetryTime != nil {
			return false
		}

		if len(status.Tasks) > 0 && status.Reason != ReasonPurgeFailed {
			return false
		}
	}

	return true
}

// startOperation creates the jobs of an operation. Verify and purge create one job per
// task in batches, so the operation is recorded first and completed by driveOperation.
func (r *DatasetReconciler) startOperation(ctx context.Context, state *reconcileState, source *v1alpha1.DataSource, status *v1alpha1.DataSourceStatus, op v1alpha1.OperationType) {
	status.Operation = &v1alpha1.OperationStatus{Type: op, StartedAt: metav1.NewTime(state.now)}
	switch op {
	case v1alpha1.OperationTypeDistribute:
		status.State = v1alpha1.DataSourceStateDistributing
		status.Reason = ReasonDistributing
		status.Message = "preheating data onto peers"
	case v1alpha1.OperationTypeVerify:
		status.Reason = ReasonVerifying
		status.Message = "checking which peers hold the data"
	case v1alpha1.OperationTypePurge:
		status.State = v1alpha1.DataSourceStatePurging
		status.Reason = ReasonPurging
		status.Message = "removing data from peers"
	}

	r.driveOperation(ctx, state, source, status)
}

// driveOperation creates missing jobs, polls the ones in flight and, once every job is
// terminal, applies the outcome. It returns true when the operation finished.
func (r *DatasetReconciler) driveOperation(ctx context.Context, state *reconcileState, source *v1alpha1.DataSource, status *v1alpha1.DataSourceStatus) bool {
	op := status.Operation
	if state.manager == nil {
		state.poll = true
		return false
	}

	if err := r.createJobs(ctx, state, source, status); err != nil {
		var opErr *operationError
		if errors.As(err, &opErr) {
			r.finishOperation(state, source, status, false, opErr.Error())
			return true
		}

		state.fail(err)
		return false
	}

	pending := false
	for i := range op.Jobs {
		ref := &op.Jobs[i]
		if dragonfly.IsTerminalState(ref.State) {
			continue
		}

		job, err := state.manager.GetJob(ctx, ref.ID)
		if err != nil {
			if errors.Is(err, dragonfly.ErrNotFound) {
				ref.State = dragonfly.JobStateFailure
				continue
			}

			state.fail(fmt.Errorf("get job %d: %w", ref.ID, err))
			pending = true
			continue
		}

		ref.State = job.State
		if !dragonfly.IsTerminalState(job.State) {
			pending = true
		}
	}

	if pending || !operationJobsComplete(op, status) {
		state.poll = true
		return false
	}

	r.completeOperation(ctx, state, source, status)
	return true
}

// operationJobsComplete reports whether every job the operation needs was created.
func operationJobsComplete(op *v1alpha1.OperationStatus, status *v1alpha1.DataSourceStatus) bool {
	if op.Type == v1alpha1.OperationTypeDistribute {
		return len(op.Jobs) == 1
	}

	return len(op.Jobs) >= len(status.Tasks)
}

// createJobs creates the jobs an operation still needs, at most maxJobsPerReconcile.
func (r *DatasetReconciler) createJobs(ctx context.Context, state *reconcileState, source *v1alpha1.DataSource, status *v1alpha1.DataSourceStatus) error {
	op := status.Operation
	dataset := state.dataset

	if op.Type == v1alpha1.OperationTypeDistribute {
		if len(op.Jobs) > 0 {
			return nil
		}

		req, err := buildPreheatRequest(ctx, state.secrets, dataset, source)
		if err != nil {
			return &operationError{err: err}
		}

		job, err := state.manager.CreatePreheatJob(ctx, req)
		if err != nil {
			metrics.JobsCreated.WithLabelValues(dragonfly.JobTypePreheat, "error").Inc()
			if isRetryable(err) {
				return fmt.Errorf("create preheat job: %w", err)
			}

			return &operationError{err: fmt.Errorf("create preheat job: %w", err)}
		}

		metrics.JobsCreated.WithLabelValues(dragonfly.JobTypePreheat, "success").Inc()
		op.Jobs = append(op.Jobs, v1alpha1.JobReference{ID: job.ID, Type: job.Type, State: job.State})
		r.event(dataset, corev1.EventTypeNormal, ReasonDistributing, "created preheat job %d for source %s", job.ID, source.Name)
		return nil
	}

	created := make(map[string]struct{}, len(op.Jobs))
	for _, ref := range op.Jobs {
		created[ref.URL] = struct{}{}
	}

	var (
		count   int
		jobType = dragonfly.JobTypeGetTask
	)
	if op.Type == v1alpha1.OperationTypePurge {
		jobType = dragonfly.JobTypeDeleteTask
	}

	for i := range status.Tasks {
		task := &status.Tasks[i]
		if _, ok := created[task.URL]; ok {
			continue
		}

		if count >= maxJobsPerReconcile {
			state.requeueNow()
			return nil
		}

		var (
			job *dragonfly.Job
			err error
		)

		switch op.Type {
		case v1alpha1.OperationTypeVerify:
			var req managertypes.CreateGetTaskJobRequest
			req, err = buildGetTaskRequest(dataset, source, task)
			if err == nil {
				job, err = state.manager.CreateGetTaskJob(ctx, req)
			}
		case v1alpha1.OperationTypePurge:
			var req managertypes.CreateDeleteTaskJobRequest
			req, err = buildDeleteTaskRequest(dataset, source, task)
			if err == nil {
				job, err = state.manager.CreateDeleteTaskJob(ctx, req)
			}
		}

		if err != nil {
			metrics.JobsCreated.WithLabelValues(jobType, "error").Inc()
			if isRetryable(err) {
				return fmt.Errorf("create %s job for %s: %w", jobType, task.URL, err)
			}

			// Record the failure as a failed job so the operation completes and the
			// outcome is reported per task.
			op.Jobs = append(op.Jobs, v1alpha1.JobReference{Type: jobType, State: dragonfly.JobStateFailure, URL: task.URL})
			status.Message = truncate(fmt.Sprintf("create %s job for %s: %v", jobType, task.URL, err), 1024)
			count++
			continue
		}

		metrics.JobsCreated.WithLabelValues(jobType, "success").Inc()
		op.Jobs = append(op.Jobs, v1alpha1.JobReference{ID: job.ID, Type: job.Type, State: job.State, URL: task.URL})
		count++
	}

	return nil
}

// completeOperation applies the outcome of a finished operation to the source.
func (r *DatasetReconciler) completeOperation(ctx context.Context, state *reconcileState, source *v1alpha1.DataSource, status *v1alpha1.DataSourceStatus) {
	op := status.Operation
	switch op.Type {
	case v1alpha1.OperationTypeDistribute:
		job, err := state.manager.GetJob(ctx, op.Jobs[0].ID)
		if err != nil {
			if !errors.Is(err, dragonfly.ErrNotFound) {
				state.fail(fmt.Errorf("get job %d: %w", op.Jobs[0].ID, err))
				state.poll = true
				return
			}

			r.finishOperation(state, source, status, false, fmt.Sprintf("preheat job %d disappeared from the manager", op.Jobs[0].ID))
			return
		}

		outcome := evaluateDistribute(job, source)
		if !outcome.failed {
			status.Tasks = outcome.tasks
		}

		r.finishOperation(state, source, status, !outcome.failed, outcome.message)

	case v1alpha1.OperationTypeVerify:
		tasks := make(map[string]*v1alpha1.TaskStatus, len(status.Tasks))
		for i := range status.Tasks {
			tasks[status.Tasks[i].URL] = &status.Tasks[i]
		}

		var (
			failures []string
			missing  int
		)

		for _, ref := range op.Jobs {
			task, ok := tasks[ref.URL]
			if !ok {
				continue
			}

			if ref.ID == 0 {
				task.PeerCount = nil
				failures = append(failures, fmt.Sprintf("verification of %s was not started", task.URL))
				continue
			}

			job, err := state.manager.GetJob(ctx, ref.ID)
			if err != nil {
				task.PeerCount = nil
				failures = append(failures, fmt.Sprintf("get job %d: %v", ref.ID, err))
				continue
			}

			if ok, msg := evaluateVerify(job, task); !ok {
				failures = append(failures, msg)
				continue
			}

			if task.PeerCount != nil && *task.PeerCount == 0 {
				missing++
			}
		}

		r.finishOperation(state, source, status, len(failures) == 0, verifyMessage(len(status.Tasks), missing, failures))
		if len(failures) == 0 {
			if missing > 0 {
				status.State = v1alpha1.DataSourceStateDegraded
				status.Reason = ReasonDataMissing
			} else {
				status.State = v1alpha1.DataSourceStateReady
				status.Reason = ReasonVerified
			}
		}

	case v1alpha1.OperationTypePurge:
		tasks := make(map[string]*v1alpha1.TaskStatus, len(status.Tasks))
		for i := range status.Tasks {
			tasks[status.Tasks[i].URL] = &status.Tasks[i]
		}

		var failures []string
		for _, ref := range op.Jobs {
			task, ok := tasks[ref.URL]
			if !ok {
				continue
			}

			if ref.ID == 0 {
				failures = append(failures, fmt.Sprintf("purge of %s was not started", task.URL))
				continue
			}

			job, err := state.manager.GetJob(ctx, ref.ID)
			if err != nil {
				failures = append(failures, fmt.Sprintf("get job %d: %v", ref.ID, err))
				continue
			}

			if ok, msg := evaluatePurge(job, task); !ok {
				failures = append(failures, msg)
			}
		}

		r.finishOperation(state, source, status, len(failures) == 0, joinMessages(fmt.Sprintf("%d task(s) purged", len(status.Tasks)-len(failures)), failures))
	}
}

func verifyMessage(total, missing int, failures []string) string {
	if len(failures) > 0 {
		return joinMessages(fmt.Sprintf("verification of %d task(s) failed", len(failures)), failures)
	}

	if missing > 0 {
		return fmt.Sprintf("%d of %d task(s) missing from every peer", missing, total)
	}

	return fmt.Sprintf("%d task(s) present on peers", total)
}

// finishOperation records the outcome of an operation and schedules a retry when it
// failed and retries are left.
func (r *DatasetReconciler) finishOperation(state *reconcileState, source *v1alpha1.DataSource, status *v1alpha1.DataSourceStatus, succeeded bool, message string) {
	dataset := state.dataset
	op := status.Operation
	status.Operation = nil
	now := metav1.NewTime(state.now)
	metrics.OperationsCompleted.WithLabelValues(string(op.Type), outcomeLabel(succeeded)).Inc()

	if succeeded {
		status.Retries = 0
		status.NextRetryTime = nil
		status.Message = message
		switch op.Type {
		case v1alpha1.OperationTypeDistribute:
			status.State = v1alpha1.DataSourceStateReady
			status.Reason = ReasonDistributed
			status.LastDistributionTime = &now
			for i := range status.Tasks {
				status.Tasks[i].PeerCount = nil
			}
			r.event(dataset, corev1.EventTypeNormal, ReasonDistributed, "source %s: %s", source.Name, message)
		case v1alpha1.OperationTypeVerify:
			status.LastVerificationTime = &now
		case v1alpha1.OperationTypePurge:
			status.Tasks = nil
			status.Reason = ReasonPurged
			if !dataset.DeletionTimestamp.IsZero() || !sourceInSpec(dataset, source.Name) {
				status.State = v1alpha1.DataSourceStatePurged
			} else {
				status.State = v1alpha1.DataSourceStateExpired
				status.Message = "ttl elapsed, data purged from peers"
			}
			r.event(dataset, corev1.EventTypeNormal, ReasonPurged, "source %s: %s", source.Name, message)
		}

		return
	}

	switch op.Type {
	case v1alpha1.OperationTypeVerify:
		// A failed verification is reported but not retried on its own; the next
		// verification interval tries again.
		status.LastVerificationTime = &now
		status.Reason = ReasonVerificationFailed
		status.Message = message
		r.event(dataset, corev1.EventTypeWarning, ReasonVerificationFailed, "source %s: %s", source.Name, message)
		return

	case v1alpha1.OperationTypeDistribute, v1alpha1.OperationTypePurge:
		status.Retries++
		retryLimit := int32(0)
		if dataset.Spec.Distribution.RetryLimit != nil {
			retryLimit = *dataset.Spec.Distribution.RetryLimit
		}

		if status.Retries <= retryLimit {
			delay := retryBackoff(r.Config.Controller.RetryBackoff, r.Config.Controller.MaxRetryBackoff, status.Retries)
			next := metav1.NewTime(state.now.Add(delay))
			status.NextRetryTime = &next
			status.Reason = ReasonRetryScheduled
			status.Message = fmt.Sprintf("%s (retry %d/%d in %s)", message, status.Retries, retryLimit, delay.Round(time.Second))
			if op.Type == v1alpha1.OperationTypeDistribute {
				status.State = v1alpha1.DataSourceStatePending
			}

			state.due(&next.Time)
			r.event(dataset, corev1.EventTypeWarning, ReasonRetryScheduled, "source %s: %s", source.Name, status.Message)
			return
		}

		status.NextRetryTime = nil
		status.Message = message
		if op.Type == v1alpha1.OperationTypeDistribute {
			status.State = v1alpha1.DataSourceStateFailed
			status.Reason = ReasonDistributionFailed
			r.event(dataset, corev1.EventTypeWarning, ReasonDistributionFailed, "source %s: %s", source.Name, message)
			return
		}

		status.Reason = ReasonPurgeFailed
		if !dataset.DeletionTimestamp.IsZero() || !sourceInSpec(dataset, source.Name) {
			status.State = v1alpha1.DataSourceStateFailed
		} else {
			status.State = v1alpha1.DataSourceStateExpired
		}
		r.event(dataset, corev1.EventTypeWarning, ReasonPurgeFailed, "source %s: %s", source.Name, message)
	}
}

func outcomeLabel(succeeded bool) string {
	if succeeded {
		return "success"
	}

	return "failure"
}

func sourceInSpec(dataset *v1alpha1.Dataset, name string) bool {
	for _, source := range dataset.Spec.Sources {
		if source.Name == name {
			return true
		}
	}

	return false
}

// summarize derives the dataset-level phase, conditions and times from the sources.
func (r *DatasetReconciler) summarize(state *reconcileState) {
	dataset := state.dataset
	lifecycle := state.lifecycle

	lastDistribution, lastVerification := datasetTimes(dataset)
	dataset.Status.LastDistributionTime = toMetaTime(lastDistribution)
	dataset.Status.LastVerificationTime = toMetaTime(lastVerification)
	dataset.Status.NextRefreshTime = toMetaTime(lifecycle.NextRefresh(lastDistribution))
	dataset.Status.ExpirationTime = toMetaTime(lifecycle.Expiration(lastDistribution))
	dataset.Status.NextVerificationTime = toMetaTime(lifecycle.NextVerification(lastDistribution, lastVerification))

	if !dataset.Spec.Suspend {
		state.due(ptrTime(dataset.Status.NextRefreshTime))
		state.due(ptrTime(dataset.Status.ExpirationTime))
		state.due(ptrTime(dataset.Status.NextVerificationTime))
	}

	var counts struct {
		total, ready, degraded, failed, expired, pending, distributing, purging int
		verified, unverified, missing                                           int
	}

	for _, status := range dataset.Status.Sources {
		if !sourceInSpec(dataset, status.Name) {
			if status.State == v1alpha1.DataSourceStatePurging {
				counts.purging++
			}

			continue
		}

		counts.total++
		switch status.State {
		case v1alpha1.DataSourceStateReady:
			counts.ready++
		case v1alpha1.DataSourceStateDegraded:
			counts.degraded++
		case v1alpha1.DataSourceStateFailed:
			counts.failed++
		case v1alpha1.DataSourceStateExpired:
			counts.expired++
		case v1alpha1.DataSourceStateDistributing:
			counts.distributing++
		case v1alpha1.DataSourceStatePurging:
			counts.purging++
		default:
			counts.pending++
		}

		switch status.Reason {
		case ReasonVerified:
			counts.verified++
		case ReasonDataMissing:
			counts.missing++
		case ReasonVerificationFailed:
			counts.unverified++
		}
	}

	distributed := counts.total > 0 && counts.ready+counts.degraded+counts.expired == counts.total
	inFlight := counts.distributing > 0 || counts.purging > 0

	var phase v1alpha1.DatasetPhase
	switch {
	case dataset.Spec.Suspend && !inFlight:
		phase = v1alpha1.DatasetPhaseSuspended
	case counts.purging > 0:
		phase = v1alpha1.DatasetPhasePurging
	case counts.distributing > 0:
		phase = v1alpha1.DatasetPhaseDistributing
	case counts.total > 0 && counts.failed == counts.total:
		phase = v1alpha1.DatasetPhaseFailed
	case counts.total > 0 && counts.expired == counts.total:
		phase = v1alpha1.DatasetPhaseExpired
	case counts.failed > 0 || counts.degraded > 0 || (counts.expired > 0 && counts.ready > 0):
		phase = v1alpha1.DatasetPhaseDegraded
	case counts.total > 0 && counts.ready == counts.total:
		phase = v1alpha1.DatasetPhaseReady
	default:
		phase = v1alpha1.DatasetPhasePending
	}
	dataset.Status.Phase = phase

	switch {
	case dataset.Spec.Suspend && !inFlight:
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, ReasonSuspended, "dataset is suspended")
	case phase == v1alpha1.DatasetPhaseReady:
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionTrue, ReasonDistributed, fmt.Sprintf("%d source(s) distributed", counts.total))
	case phase == v1alpha1.DatasetPhaseFailed:
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, ReasonDistributionFailed, sourceSummary(dataset))
	case phase == v1alpha1.DatasetPhaseExpired:
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, ReasonExpired, "ttl elapsed")
	case phase == v1alpha1.DatasetPhaseDegraded:
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, degradedReason(counts.failed, counts.degraded), sourceSummary(dataset))
	case phase == v1alpha1.DatasetPhasePurging:
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, ReasonPurging, sourceSummary(dataset))
	default:
		r.setCondition(dataset, v1alpha1.ConditionReady, metav1.ConditionFalse, ReasonDistributing, sourceSummary(dataset))
	}

	if distributed {
		r.setCondition(dataset, v1alpha1.ConditionDistributed, metav1.ConditionTrue, ReasonDistributed, fmt.Sprintf("%d source(s) distributed", counts.total))
	} else if counts.failed > 0 && !inFlight {
		r.setCondition(dataset, v1alpha1.ConditionDistributed, metav1.ConditionFalse, ReasonDistributionFailed, sourceSummary(dataset))
	} else {
		r.setCondition(dataset, v1alpha1.ConditionDistributed, metav1.ConditionFalse, ReasonDistributing, sourceSummary(dataset))
	}

	switch {
	case lifecycle.VerifyInterval == nil:
		meta.RemoveStatusCondition(&dataset.Status.Conditions, v1alpha1.ConditionVerified)
	case counts.missing > 0:
		r.setCondition(dataset, v1alpha1.ConditionVerified, metav1.ConditionFalse, ReasonDataMissing, sourceSummary(dataset))
	case counts.unverified > 0:
		r.setCondition(dataset, v1alpha1.ConditionVerified, metav1.ConditionFalse, ReasonVerificationFailed, sourceSummary(dataset))
	case counts.total > 0 && counts.verified == counts.total:
		r.setCondition(dataset, v1alpha1.ConditionVerified, metav1.ConditionTrue, ReasonVerified, fmt.Sprintf("%d source(s) present on peers", counts.total))
	default:
		r.setCondition(dataset, v1alpha1.ConditionVerified, metav1.ConditionUnknown, ReasonNotVerified, "data has not been verified since the last distribution")
	}

	switch {
	case lifecycle.TTL == nil:
		meta.RemoveStatusCondition(&dataset.Status.Conditions, v1alpha1.ConditionExpired)
	case counts.expired > 0:
		r.setCondition(dataset, v1alpha1.ConditionExpired, metav1.ConditionTrue, ReasonExpired, fmt.Sprintf("%d of %d source(s) expired", counts.expired, counts.total))
	default:
		r.setCondition(dataset, v1alpha1.ConditionExpired, metav1.ConditionFalse, ReasonNotExpired, "ttl has not elapsed")
	}
}

func degradedReason(failed, degraded int) string {
	if failed > 0 {
		return ReasonDistributionFailed
	}

	if degraded > 0 {
		return ReasonDataMissing
	}

	return ReasonExpired
}

// sourceSummary lists the sources that are not ready with their reason.
func sourceSummary(dataset *v1alpha1.Dataset) string {
	var parts []string
	for _, status := range dataset.Status.Sources {
		if status.State == v1alpha1.DataSourceStateReady && status.Reason != ReasonDataMissing {
			continue
		}

		parts = append(parts, fmt.Sprintf("%s: %s (%s)", status.Name, status.State, status.Reason))
	}

	if len(parts) == 0 {
		return "all sources ready"
	}

	return truncate(strings.Join(parts, "; "), 1024)
}

// datasetTimes derives the dataset-level distribution and verification times. They are
// the latest source times, so that during a refresh a source that finished first does
// not become due again while the others are still running.
func datasetTimes(dataset *v1alpha1.Dataset) (lastDistribution, lastVerification *time.Time) {
	for _, status := range dataset.Status.Sources {
		if !sourceInSpec(dataset, status.Name) {
			continue
		}

		lastDistribution = laterOf(lastDistribution, status.LastDistributionTime)
		lastVerification = laterOf(lastVerification, status.LastVerificationTime)
	}

	return lastDistribution, lastVerification
}

func laterOf(current *time.Time, candidate *metav1.Time) *time.Time {
	if candidate == nil {
		return current
	}

	if current == nil || candidate.After(*current) {
		t := candidate.Time
		return &t
	}

	return current
}

func isDue(t *time.Time, now time.Time) bool {
	return t != nil && !now.Before(*t)
}

func toMetaTime(t *time.Time) *metav1.Time {
	if t == nil {
		return nil
	}

	mt := metav1.NewTime(*t)
	return &mt
}

func ptrTime(t *metav1.Time) *time.Time {
	if t == nil {
		return nil
	}

	return &t.Time
}

// event records an event on a dataset.
func (r *DatasetReconciler) event(dataset *v1alpha1.Dataset, eventType, reason, note string, args ...any) {
	if r.Recorder == nil {
		return
	}

	r.Recorder.Eventf(dataset, nil, eventType, reason, "Reconcile", note, args...)
}

func (r *DatasetReconciler) setCondition(dataset *v1alpha1.Dataset, conditionType string, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&dataset.Status.Conditions, metav1.Condition{
		Type:               conditionType,
		Status:             status,
		Reason:             reason,
		Message:            truncate(message, 32*1024),
		ObservedGeneration: dataset.Generation,
	})
}

// requeueAfter picks the next time the dataset must be looked at.
func (r *DatasetReconciler) requeueAfter(state *reconcileState) time.Duration {
	if state.immediate {
		return minRequeue
	}

	delay := r.Config.Controller.ResyncInterval
	if state.poll {
		delay = r.Config.Controller.PollInterval
	}

	for _, t := range state.dueTimes {
		if d := t.Sub(state.now); d > 0 && d < delay {
			delay = d
		}
	}

	if delay < minRequeue {
		delay = minRequeue
	}

	return delay
}
