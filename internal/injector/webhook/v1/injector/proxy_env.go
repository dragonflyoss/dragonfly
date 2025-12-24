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

package injector

import (
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

// ProxyEnvInjector injects proxy environment variables into pod containers
type ProxyEnvInjector struct{}

// NewProxyEnvInjector creates a new ProxyEnvInjector instance
func NewProxyEnvInjector() *ProxyEnvInjector {
	return &ProxyEnvInjector{}
}

// Inject adds proxy environment variables to all containers in the pod
func (pei *ProxyEnvInjector) Inject(pod *corev1.Pod, config *InjectConf) {
	logger.V(1).Info("Injecting proxy environment variables", "pod", pod.Name)

	envs := envsFromConfig(config)
	
	// Inject environment variables to all containers
	for i := range pod.Spec.Containers {
		injectContainerEnv(&pod.Spec.Containers[i], envs)
	}
}

// envsFromConfig creates environment variables from the configuration
func envsFromConfig(config *InjectConf) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: NodeNameEnvName,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "spec.nodeName",
				},
			},
		},
		{
			Name:  ProxyPortEnvName,
			Value: strconv.Itoa(config.ProxyPort),
		},
		{
			Name:  ProxyEnvName,
			Value: "http://$(" + NodeNameEnvName + "):$(" + ProxyPortEnvName + ")",
		},
	}
}

// injectContainerEnv injects environment variables into a container if they don't already exist
func injectContainerEnv(container *corev1.Container, envs []corev1.EnvVar) {
	for _, env := range envs {
		// Check if environment variable already exists
		exists := false
		for _, existingEnv := range container.Env {
			if env.Name == existingEnv.Name {
				exists = true
				break
			}
		}
		
		// Only inject if it doesn't exist
		if !exists {
			container.Env = append(container.Env, env)
		}
	}
}
