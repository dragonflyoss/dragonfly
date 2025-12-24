/*
Copyright 2025 The Dragonfly Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// +kubebuilder:rbac:groups="",resources=namespaces;pods,verbs=get;list;watch

package v1

import (
	"context"
	"fmt"

	"d7y.io/dragonfly/v2/internal/injector/webhook/v1/injector"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var webhookLogger = logf.Log.WithName("pod-webhook")

// Injector defines the interface for pod injection components
type Injector interface {
	Inject(pod *corev1.Pod, config *injector.InjectConf)
}

// SetupPodWebhookWithManager registers the webhook for Pod in the manager
func SetupPodWebhookWithManager(mgr ctrl.Manager) error {
	configManager := injector.NewConfigManager(injector.InjectConfigMapPath)
	if err := mgr.Add(configManager); err != nil {
		return fmt.Errorf("failed to add config manager to manager: %w", err)
	}

	defaulter := NewPodCustomDefaulter(mgr.GetClient(), configManager)

	return ctrl.NewWebhookManagedBy(mgr).For(&corev1.Pod{}).
		WithDefaulter(defaulter).
		Complete()
}

// +kubebuilder:webhook:path=/mutate--v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod-v1.d7y.io,admissionReviewVersions=v1

// PodCustomDefaulter is responsible for mutating pods to inject dragonfly capabilities
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type PodCustomDefaulter struct {
	configManager *injector.ConfigManager
	kubeClient    client.Client
	injectors     []Injector
}

var _ webhook.CustomDefaulter = &PodCustomDefaulter{}

// NewPodCustomDefaulter creates a new PodCustomDefaulter instance
func NewPodCustomDefaulter(c client.Client, configManager *injector.ConfigManager) *PodCustomDefaulter {
	return &PodCustomDefaulter{
		kubeClient:    c,
		configManager: configManager,
		injectors: []Injector{
			injector.NewProxyEnvInjector(),
			injector.NewUnixSocketInjector(),
			injector.NewToolsInitcontainerInjector(),
		},
	}
}

// Default implements webhook.CustomDefaulter to mutate pods
func (d *PodCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod object but got %T", obj)
	}

	webhookLogger.V(1).Info("Processing pod for injection", "pod", pod.Name, "namespace", pod.Namespace)
	d.applyDefaults(ctx, pod)
	return nil
}

// applyDefaults applies dragonfly injection to the pod if required
func (d *PodCustomDefaulter) applyDefaults(ctx context.Context, pod *corev1.Pod) {
	config := d.configManager.GetConfig()
	if config == nil || !config.Enable {
		webhookLogger.V(1).Info("Injection disabled in configuration, skipping", "pod", pod.Name)
		return
	}

	if !d.injectRequired(ctx, pod) {
		webhookLogger.V(1).Info("Injection not required for pod", "pod", pod.Name)
		return
	}

	webhookLogger.Info("Injecting dragonfly capabilities into pod", "pod", pod.Name, "namespace", pod.Namespace)
	for _, injector := range d.injectors {
		injector.Inject(pod, config)
	}
}

// injectRequired determines if injection is required for the pod
func (d *PodCustomDefaulter) injectRequired(ctx context.Context, pod *corev1.Pod) bool {
	// Pod-level annotation takes precedence over namespace-level label
	if d.isPodInjectionEnabled(ctx, pod) {
		return true
	}
	return d.isNamespaceInjectionEnabled(ctx, pod)
}

// isNamespaceInjectionEnabled checks if injection is enabled at the namespace level
func (d *PodCustomDefaulter) isNamespaceInjectionEnabled(ctx context.Context, pod *corev1.Pod) bool {
	nsName := pod.GetNamespace()
	ns := &corev1.Namespace{}
	if err := d.kubeClient.Get(ctx, client.ObjectKey{Name: nsName}, ns); err != nil {
		webhookLogger.Error(err, "Failed to get namespace", "namespace", nsName)
		return false
	}

	labels := ns.GetLabels()
	if len(labels) == 0 {
		webhookLogger.V(1).Info("Namespace has no labels", "namespace", nsName)
		return false
	}

	value, exists := labels[injector.NamespaceInjectLabelName]
	if !exists || value != injector.NamespaceInjectLabelValue {
		webhookLogger.V(1).Info("Namespace injection not enabled",
			"namespace", nsName,
			"label", injector.NamespaceInjectLabelName,
			"value", value)
		return false
	}

	webhookLogger.V(1).Info("Namespace injection enabled", "namespace", nsName)
	return true
}

// isPodInjectionEnabled checks if injection is enabled at the pod level
func (d *PodCustomDefaulter) isPodInjectionEnabled(_ context.Context, pod *corev1.Pod) bool {
	annotations := pod.GetAnnotations()
	if len(annotations) == 0 {
		webhookLogger.V(1).Info("Pod has no annotations", "pod", pod.Name)
		return false
	}

	value, exists := annotations[injector.PodInjectAnnotationName]
	if !exists || value != injector.PodInjectAnnotationValue {
		webhookLogger.V(1).Info("Pod injection not enabled",
			"pod", pod.Name,
			"annotation", injector.PodInjectAnnotationName,
			"value", value)
		return false
	}

	webhookLogger.V(1).Info("Pod injection enabled", "pod", pod.Name)
	return true
}
