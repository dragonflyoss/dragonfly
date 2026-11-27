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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("UnixSocketInjector", func() {
	var (
		injector *UnixSocket
	)

	BeforeEach(func() {
		injector = NewUnixSocket()
	})

	// Helper function to create the expected Volume
	makeExpectedVolume := func() corev1.Volume {
		hostPathType := corev1.HostPathSocket
		return corev1.Volume{
			Name: DfdaemonUnixSockVolumeName,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: DfdaemonUnixSockPath,
					Type: &hostPathType,
				},
			},
		}
	}

	// Helper function to create the expected VolumeMount
	makeExpectedVolumeMount := func() corev1.VolumeMount {
		return corev1.VolumeMount{
			Name:      DfdaemonUnixSockVolumeName,
			MountPath: DfdaemonUnixSockPath,
		}
	}

	Context("when injecting unix socket volume and mounts", func() {
		It("should inject into a pod with no existing volume or volume mounts", func() {
			By("creating a simple pod")
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-1"},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "container-1"}},
				},
			}

			By("creating expected pod")
			expectedVolume := makeExpectedVolume()
			expectedVolumeMount := makeExpectedVolumeMount()
			expectedPod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-1"},
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{expectedVolume},
					Containers: []corev1.Container{{Name: "container-1", VolumeMounts: []corev1.VolumeMount{expectedVolumeMount}}},
				},
			}

			By("performing injection")
			injector.Inject(pod, &Config{})

			By("verifying the result")
			Expect(pod).To(Equal(expectedPod))
		})

		It("should not inject when volume already exists", func() {
			By("creating a pod with existing volume")
			expectedVolume := makeExpectedVolume()
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-2"},
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{expectedVolume},
					Containers: []corev1.Container{{Name: "container-1"}},
				},
			}

			By("creating expected pod (volume remains unchanged)")
			expectedVolumeMount := makeExpectedVolumeMount()
			expectedPod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-2"},
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{expectedVolume},
					Containers: []corev1.Container{{Name: "container-1", VolumeMounts: []corev1.VolumeMount{expectedVolumeMount}}},
				},
			}

			By("performing injection")
			injector.Inject(pod, &Config{})

			By("verifying the result")
			Expect(pod).To(Equal(expectedPod))
		})

		It("should inject into containers that don't have volume mount", func() {
			By("creating a pod where one container already has the volume mount")
			expectedVolume := makeExpectedVolume()
			expectedVolumeMount := makeExpectedVolumeMount()
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-3"},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "container-a"}, // Needs injection
						{
							Name:         "container-b", // Already exists, should not be injected again
							VolumeMounts: []corev1.VolumeMount{expectedVolumeMount},
						},
					},
				},
			}

			By("creating expected pod")
			expectedPod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-3"},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{expectedVolume},
					Containers: []corev1.Container{
						{
							Name:         "container-a",
							VolumeMounts: []corev1.VolumeMount{expectedVolumeMount},
						},
						{
							Name:         "container-b", // Remains unchanged
							VolumeMounts: []corev1.VolumeMount{expectedVolumeMount},
						},
					},
				},
			}

			By("performing injection")
			injector.Inject(pod, &Config{})

			By("verifying the result")
			Expect(pod).To(Equal(expectedPod))
		})

		It("should be idempotent when both volume and volume mounts already exist", func() {
			By("creating a pod with everything already injected")
			expectedVolume := makeExpectedVolume()
			expectedVolumeMount := makeExpectedVolumeMount()
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-4"},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{expectedVolume},
					Containers: []corev1.Container{
						{
							Name:         "container-1",
							VolumeMounts: []corev1.VolumeMount{expectedVolumeMount},
						},
					},
				},
			}

			By("creating expected pod (should be unchanged)")
			expectedPod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-4"},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{expectedVolume},
					Containers: []corev1.Container{
						{
							Name:         "container-1",
							VolumeMounts: []corev1.VolumeMount{expectedVolumeMount},
						},
					},
				},
			}

			By("performing injection")
			injector.Inject(pod, &Config{})

			By("verifying the result")
			Expect(pod).To(Equal(expectedPod))
		})

		It("should handle pods with no containers gracefully", func() {
			By("creating a pod with no containers")
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-5"},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{}, // No containers
				},
			}

			By("creating expected pod")
			expectedVolume := makeExpectedVolume()
			expectedPod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-5"},
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{expectedVolume}, // Volume is still injected
					Containers: []corev1.Container{},            // Container list remains empty
				},
			}

			By("performing injection")
			injector.Inject(pod, &Config{})

			By("verifying the result")
			Expect(pod).To(Equal(expectedPod))
		})

		It("should inject into a pod with existing unrelated volumes and mounts", func() {
			By("creating a pod with existing unrelated volumes and mounts")
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-6"},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{{Name: "other-volume", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}},
					Containers: []corev1.Container{
						{
							Name:         "container-1",
							VolumeMounts: []corev1.VolumeMount{{Name: "other-mount", MountPath: "/data"}},
						},
					},
				},
			}

			By("creating expected pod")
			expectedVolume := makeExpectedVolume()
			expectedVolumeMount := makeExpectedVolumeMount()
			expectedPod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod-6"},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{Name: "other-volume", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
						expectedVolume, // New volume is appended
					},
					Containers: []corev1.Container{
						{
							Name: "container-1",
							VolumeMounts: []corev1.VolumeMount{
								{Name: "other-mount", MountPath: "/data"},
								expectedVolumeMount, // New volume mount is appended
							},
						},
					},
				},
			}

			By("performing injection")
			injector.Inject(pod, &Config{})

			By("verifying the result")
			Expect(pod).To(Equal(expectedPod))
		})
	})
})
