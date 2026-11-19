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

package injector

import (
	corev1 "k8s.io/api/core/v1"
)

type UnixSocket struct{}

func NewUnixSocket() *UnixSocket {
	return &UnixSocket{}
}

func (us *UnixSocket) Inject(pod *corev1.Pod, config *Config) {
	logger.Info("UnixSocket inject", "pod namespace", pod.GetNamespace(), "pod name", pod.GetName())

	if pod.Annotations != nil && pod.Annotations[SkipUnixSockInjectAnnotationName] == SkipUnixSockInjectAnnotationValue {
		logger.Info("UnixSocket inject skipped", "pod namespace", pod.GetNamespace(), "pod name", pod.GetName())
		return
	}

	if !us.hasVolume(pod) {
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

	for i, c := range pod.Spec.Containers {
		if !us.hasVolumeMount(&c) {
			dfdaemonSocketVolumeMount := corev1.VolumeMount{
				Name:      DfdaemonUnixSockVolumeName,
				MountPath: DfdaemonUnixSockPath,
			}
			pod.Spec.Containers[i].VolumeMounts = append(pod.Spec.Containers[i].VolumeMounts, dfdaemonSocketVolumeMount)
		}
	}
}

// hasVolume checks if the pod has the volume.
func (us *UnixSocket) hasVolume(pod *corev1.Pod) bool {
	for _, v := range pod.Spec.Volumes {
		if v.Name == DfdaemonUnixSockVolumeName {
			return true
		}
	}

	return false
}

// hasVolumeMount checks if the container has the volume mount.
func (us *UnixSocket) hasVolumeMount(c *corev1.Container) bool {
	for _, vm := range c.VolumeMounts {
		if vm.Name == DfdaemonUnixSockVolumeName {
			return true
		}
	}

	return false
}
