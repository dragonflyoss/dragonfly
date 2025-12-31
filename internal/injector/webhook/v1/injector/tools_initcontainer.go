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
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
)

// ToolsInitContainerInjector injects CLI tools via init container into pods
type ToolsInitContainerInjector struct{}

// NewToolsInitcontainerInjector creates a new ToolsInitContainerInjector instance
func NewToolsInitcontainerInjector() *ToolsInitContainerInjector {
	return &ToolsInitContainerInjector{}
}

// Inject adds CLI tools init container, volume, and environment variables to the pod
func (tii *ToolsInitContainerInjector) Inject(pod *corev1.Pod, config *InjectConf) {
	logger.V(1).Info("Injecting CLI tools init container", "pod", pod.Name)

	cliToolsVolumeMountPath := filepath.Clean(config.CliToolsDirPath) + "-mount"
	initContainerCmd := []string{
		"cp",
		"-rf",
		config.CliToolsDirPath + "/.",
		cliToolsVolumeMountPath + "/",
	}

	// Get init container image from annotation or use default
	initContainerImage := config.CliToolsImage
	if pod.Annotations != nil {
		if image, ok := pod.Annotations[CliToolsImageAnnotation]; ok {
			initContainerImage = image
		}
	}

	// Add init container if it doesn't exist
	if !tii.initContainerExists(pod) {
		toolContainer := corev1.Container{
			Name:            CliToolsInitContainerName,
			Image:           initContainerImage,
			ImagePullPolicy: corev1.PullIfNotPresent,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      CliToolsVolumeName,
					MountPath: cliToolsVolumeMountPath,
				},
			},
			Command: initContainerCmd,
		}
		pod.Spec.InitContainers = append(pod.Spec.InitContainers, toolContainer)
	}

	// Add volume if it doesn't exist
	if !tii.volumeExists(pod) {
		toolsVolume := corev1.Volume{
			Name: CliToolsVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
		pod.Spec.Volumes = append(pod.Spec.Volumes, toolsVolume)
	}

	// Add volume mount and environment variable to all containers
	for i := range pod.Spec.Containers {
		container := &pod.Spec.Containers[i]

		if !tii.volumeMountExists(container) {
			container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
				Name:      CliToolsVolumeName,
				MountPath: cliToolsVolumeMountPath,
			})
		}

		if !tii.envExists(container) {
			container.Env = append(container.Env, corev1.EnvVar{
				Name:  CliToolsPathEnvName,
				Value: cliToolsVolumeMountPath,
			})
		}
	}
}

// initContainerExists checks if the CLI tools init container already exists
func (tii *ToolsInitContainerInjector) initContainerExists(pod *corev1.Pod) bool {
	if pod == nil {
		return false
	}

	for _, initContainer := range pod.Spec.InitContainers {
		if initContainer.Name == CliToolsInitContainerName {
			return true
		}
	}
	return false
}

// volumeExists checks if the CLI tools volume already exists
func (tii *ToolsInitContainerInjector) volumeExists(pod *corev1.Pod) bool {
	if pod == nil {
		return false
	}

	for _, volume := range pod.Spec.Volumes {
		if volume.Name == CliToolsVolumeName {
			return true
		}
	}
	return false
}

// volumeMountExists checks if the CLI tools volume mount already exists in a container
func (tii *ToolsInitContainerInjector) volumeMountExists(container *corev1.Container) bool {
	if container == nil {
		return false
	}

	for _, volumeMount := range container.VolumeMounts {
		if volumeMount.Name == CliToolsVolumeName {
			return true
		}
	}
	return false
}

// envExists checks if the CLI tools path environment variable already exists in a container
func (tii *ToolsInitContainerInjector) envExists(container *corev1.Container) bool {
	if container == nil {
		return false
	}

	for _, env := range container.Env {
		if env.Name == CliToolsPathEnvName {
			return true
		}
	}
	return false
}
