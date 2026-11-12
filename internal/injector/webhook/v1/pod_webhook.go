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
// +kubebuilder:rbac:groups="",resources=namespaces;pods,verbs=get;list;watch
package v1

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"d7y.io/dragonfly/v2/internal/injector"
)

// nolint:unused
// log is for logging in this package.
var logger = logf.Log.WithName("pod-resource")

// SetupPodWebhookWithManager registers the webhook for Pod in the manager.
func SetupPodWebhookWithManager(mgr ctrl.Manager) error {
	configManager := injector.NewConfigManager(injector.ConfigMapPath)
	if err := mgr.Add(configManager); err != nil {
		return fmt.Errorf("failed to add config manager to manager: %w", err)
	}

	defaulter := NewPodCustomDefaulter(mgr.GetClient(), configManager)

	return ctrl.NewWebhookManagedBy(mgr).For(&corev1.Pod{}).
		WithDefaulter(defaulter).
		Complete()
}

type Injector interface {
	Inject(pod *corev1.Pod, config *injector.Config)
}

// +kubebuilder:webhook:path=/mutate--v1-pod,mutating=true,failurePolicy=Ignore,sideEffects=None,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod-v1.d7y.io,admissionReviewVersions=v1

// PodCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Pod when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type PodCustomDefaulter struct {
	configManager *injector.ConfigManager
	kubeClient    client.Client
	injectors     []Injector
}

var _ webhook.CustomDefaulter = &PodCustomDefaulter{}

func NewPodCustomDefaulter(c client.Client, configManager *injector.ConfigManager) *PodCustomDefaulter {
	return &PodCustomDefaulter{
		kubeClient:    c,
		configManager: configManager,
		injectors: []Injector{
			injector.NewUnixSocket(),
			injector.NewBinaries(),
		},
	}
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Pod.
func (d *PodCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	pod, ok := obj.(*corev1.Pod)

	if !ok {
		return fmt.Errorf("expected an Pod object but got %T", obj)
	}
	logger.Info("Defaulting for Pod", "pod namespace", pod.GetNamespace(), "pod name", pod.GetName())

	d.applyDefaults(ctx, pod)
	return nil
}

func (d *PodCustomDefaulter) applyDefaults(ctx context.Context, pod *corev1.Pod) {
	config := d.configManager.GetConfig()
	if config == nil {
		logger.Info("Config disabled, skip inject", "pod namespace", pod.GetNamespace(), "pod name", pod.GetName())
		return
	}

	// Check if need inject.
	if !d.needInject(ctx, pod) {
		logger.Info("Pod not inject", "pod namespace", pod.GetNamespace(), "pod name", pod.GetName())
		return
	}

	for _, ij := range d.injectors {
		ij.Inject(pod, config)
	}
}

func (d *PodCustomDefaulter) needInject(ctx context.Context, pod *corev1.Pod) bool {
	annotations := pod.GetAnnotations()
	if _, ok := annotations[injector.PodInjectAnnotationName]; ok {
		return d.isPodInjectionEnabled(ctx, pod)
	}

	return d.isNamespaceInjectionEnabled(ctx, pod)
}

func (d *PodCustomDefaulter) isNamespaceInjectionEnabled(ctx context.Context, pod *corev1.Pod) bool {
	namespace := pod.GetNamespace()
	ns := &corev1.Namespace{}
	if err := d.kubeClient.Get(ctx, client.ObjectKey{Name: namespace}, ns); err != nil {
		logger.Error(err, "failed to get namespace", "namespace", namespace)
		return false
	}

	labels := ns.GetLabels()
	logger.Info(
		"pod namespace labels",
		"namespace", namespace,
		"pod", pod.GetName(),
		"labels", labels,
	)
	if len(labels) == 0 {
		logger.Info(
			"namespace missing required injection label",
			"namespace", namespace,
			"requiredLabel", injector.NamespaceInjectLabelName,
			"pod", pod.GetName(),
		)
		return false
	}

	if v, ok := labels[injector.NamespaceInjectLabelName]; !ok ||
		v != injector.NamespaceInjectLabelValue {
		logger.Info(
			"namespace skipped injection: label not enabled",
			"namespace", namespace,
			"label", fmt.Sprintf("%s: %s", injector.NamespaceInjectLabelName, v),
			"pod", pod.GetName(),
		)
		return false
	}

	return true
}

func (d *PodCustomDefaulter) isPodInjectionEnabled(_ context.Context, pod *corev1.Pod) bool {
	annotations := pod.GetAnnotations()
	if len(annotations) == 0 {
		logger.Info(
			"pod missing required injection annotation, skip inject",
			"namespace", pod.GetNamespace(),
			"pod", pod.GetName(),
			"annotation", injector.PodInjectAnnotationName,
		)
		return false
	}

	logger.Info(
		"pod annotations",
		"namespace", pod.GetNamespace(),
		"pod", pod.GetName(),
		"annotations", annotations,
	)
	if v, ok := annotations[injector.PodInjectAnnotationName]; !ok ||
		v != injector.PodInjectAnnotationValue {
		logger.Info(
			"pod skipped injection: annotation not true, skip inject",
			"namespace", pod.GetNamespace(),
			"pod", pod.GetName(),
			"annotation", injector.PodInjectAnnotationName,
		)
		return false
	}

	return true
}
