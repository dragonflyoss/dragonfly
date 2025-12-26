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

import corev1 "k8s.io/api/core/v1"

// UnixSocketInjector injects dfdaemon unix socket volume mounts into pod containers
type UnixSocketInjector struct{}

// NewUnixSocketInjector creates a new UnixSocketInjector instance
func NewUnixSocketInjector() *UnixSocketInjector {
	return &UnixSocketInjector{}
}

// Inject adds dfdaemon unix socket volume and mounts to the pod
func (usi *UnixSocketInjector) Inject(pod *corev1.Pod, config *InjectConf) {
	logger.V(1).Info("Injecting dfdaemon unix socket volume", "pod", pod.Name)

	// Check if volume already exists
	volumeExists := false
	for _, volume := range pod.Spec.Volumes {
		if volume.Name == DfdaemonUnixSockVolumeName {
			volumeExists = true
			logger.V(1).Info("Dfdaemon socket volume already exists", "volume", volume.Name)
			break
		}
	}

	// Add volume if it doesn't exist
	if !volumeExists {
		hostPathType := corev1.HostPathSocket
		dfdaemonSocketVolume := corev1.Volume{
			Name: DfdaemonUnixSockVolumeName,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: DfdaemonUnixSockPath,
					Type: &hostPathType,
				},
			},
		}
		pod.Spec.Volumes = append(pod.Spec.Volumes, dfdaemonSocketVolume)
	}

	// Inject volume mount to all containers
	for i := range pod.Spec.Containers {
		usi.injectContainerVolumeMount(&pod.Spec.Containers[i])
	}
}

// injectContainerVolumeMount adds dfdaemon socket volume mount to a container if it doesn't exist
func (usi *UnixSocketInjector) injectContainerVolumeMount(container *corev1.Container) {
	// Check if volume mount already exists
	exists := false
	for _, volumeMount := range container.VolumeMounts {
		if volumeMount.Name == DfdaemonUnixSockVolumeName {
			exists = true
			break
		}
	}

	// Add volume mount if it doesn't exist
	if !exists {
		dfdaemonSocketVolumeMount := corev1.VolumeMount{
			Name:      DfdaemonUnixSockVolumeName,
			MountPath: DfdaemonUnixSockPath,
		}
		container.VolumeMounts = append(container.VolumeMounts, dfdaemonSocketVolumeMount)
	}
}
