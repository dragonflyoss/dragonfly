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
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	format.MaxLength = 0 // 0 means no limit
}

var _ = Describe("BinariesInjector", func() {
	var (
		binariesInjector *Binaries
		defaultConfig    *Config
		annotationImage  string
	)

	BeforeEach(func() {
		binariesInjector = NewBinaries()
		defaultConfig = &Config{
			InitContainerImage: InitContainerImage{
				Registry:    "docker.io",
				Repository:  "dragonflyoss/client",
				Tag:         "v1.3.0",
				PullPolicy:  corev1.PullIfNotPresent,
				PullSecrets: []corev1.LocalObjectReference{},
			},
		}
		annotationImage = "custom-registry.io/custom-client:v2.0.0"
	})

	// Helper function to create a clean pod object for each test.
	makePod := func(name string, containers int, annotations map[string]string) *corev1.Pod {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Annotations: annotations,
			},
			Spec: corev1.PodSpec{},
		}
		for i := 0; i < containers; i++ {
			pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{
				Name: fmt.Sprintf("container-%d", i+1),
			})
		}
		return pod
	}

	// Helper function to create the expected volume.
	makeExpectedVolume := func() corev1.Volume {
		return corev1.Volume{
			Name:         BinaryVolumeName,
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		}
	}

	// Helper function to create the expected volume mounts for a container.
	makeExpectedVolumeMounts := func() []corev1.VolumeMount {
		return []corev1.VolumeMount{
			{
				Name:      BinaryVolumeName,
				MountPath: BinaryVolumeMountPath,
			},
			{
				Name:      BinaryVolumeName,
				MountPath: filepath.Join(BinaryMountDirPath, DfgetBinaryName),
				SubPath:   DfgetBinaryName,
				ReadOnly:  true,
			},
			{
				Name:      BinaryVolumeName,
				MountPath: filepath.Join(BinaryMountDirPath, DfcacheBinaryName),
				SubPath:   DfcacheBinaryName,
				ReadOnly:  true,
			},
			{
				Name:      BinaryVolumeName,
				MountPath: filepath.Join(BinaryMountDirPath, DfstoreBinaryName),
				SubPath:   DfstoreBinaryName,
				ReadOnly:  true,
			},
			{
				Name:      BinaryVolumeName,
				MountPath: filepath.Join(BinaryMountDirPath, DfctlBinaryName),
				SubPath:   DfctlBinaryName,
				ReadOnly:  true,
			},
			{
				Name:      BinaryVolumeName,
				MountPath: filepath.Join(BinaryMountDirPath, DfdaemonBinaryName),
				SubPath:   DfdaemonBinaryName,
				ReadOnly:  true,
			},
		}
	}

	// Helper function to create the expected init container.
	makeExpectedInitContainer := func(image string, pullPolicy corev1.PullPolicy) corev1.Container {
		return corev1.Container{
			Name:            InitContainerName,
			Image:           image,
			ImagePullPolicy: pullPolicy,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      BinaryVolumeName,
					MountPath: BinaryVolumeMountPath,
				},
			},
			Command: []string{
				"install",
				"-D",
				filepath.Join(BinaryDirPath, DfgetBinaryName),
				filepath.Join(BinaryDirPath, DfcacheBinaryName),
				filepath.Join(BinaryDirPath, DfstoreBinaryName),
				filepath.Join(BinaryDirPath, DfctlBinaryName),
				filepath.Join(BinaryDirPath, DfdaemonBinaryName),
				"-t",
				BinaryVolumeMountPath + "/",
			},
		}
	}

	Describe("Inject", func() {
		Context("when injecting into a pod with a single container", func() {
			It("should inject init container, volume, and volume mounts", func() {
				By("creating a simple pod with one container")
				pod := makePod("test-pod", 1, nil)

				By("performing injection")
				binariesInjector.Inject(pod, defaultConfig)

				By("verifying the init container is added")
				expectedImage := defaultConfig.GetInitContainerImageReference()
				expectedInitContainer := makeExpectedInitContainer(expectedImage, corev1.PullIfNotPresent)
				Expect(pod.Spec.InitContainers).To(HaveLen(1))
				Expect(pod.Spec.InitContainers[0]).To(Equal(expectedInitContainer))

				By("verifying the volume is added")
				Expect(pod.Spec.Volumes).To(HaveLen(1))
				Expect(pod.Spec.Volumes[0]).To(Equal(makeExpectedVolume()))

				By("verifying volume mounts are added to the container")
				Expect(pod.Spec.Containers[0].VolumeMounts).To(Equal(makeExpectedVolumeMounts()))
			})
		})

		Context("when injecting into a pod with multiple containers", func() {
			It("should inject volume mounts into all containers", func() {
				By("creating a pod with three containers")
				pod := makePod("multi-container-pod", 3, nil)

				By("performing injection")
				binariesInjector.Inject(pod, defaultConfig)

				By("verifying each container has the correct volume mounts")
				expectedVolumeMounts := makeExpectedVolumeMounts()
				for i := range pod.Spec.Containers {
					Expect(pod.Spec.Containers[i].VolumeMounts).To(Equal(expectedVolumeMounts),
						"container-%d should have the expected volume mounts", i+1)
				}

				By("verifying only one init container and volume are added")
				Expect(pod.Spec.InitContainers).To(HaveLen(1))
				Expect(pod.Spec.Volumes).To(HaveLen(1))
			})
		})

		Context("when the pod has an init container image annotation", func() {
			It("should use the image from the annotation instead of the config", func() {
				By("creating a pod with the init container image annotation")
				pod := makePod("annotated-pod", 1, map[string]string{
					InitContainerImageAnnotation: annotationImage,
				})

				By("performing injection")
				binariesInjector.Inject(pod, defaultConfig)

				By("verifying the init container uses the annotated image")
				Expect(pod.Spec.InitContainers).To(HaveLen(1))
				Expect(pod.Spec.InitContainers[0].Image).To(Equal(annotationImage))
			})
		})

		Context("when the pod has no annotation", func() {
			It("should use the image reference from config", func() {
				By("creating a pod without annotations")
				pod := makePod("no-annotation-pod", 1, nil)

				By("performing injection")
				binariesInjector.Inject(pod, defaultConfig)

				By("verifying the init container uses the config image reference")
				expectedImage := defaultConfig.GetInitContainerImageReference()
				Expect(pod.Spec.InitContainers[0].Image).To(Equal(expectedImage))
			})
		})

		Context("when the config has image pull secrets", func() {
			It("should append pull secrets to the pod spec", func() {
				By("creating a config with pull secrets")
				configWithSecrets := &Config{
					InitContainerImage: InitContainerImage{
						Registry:   "private-registry.io",
						Repository: "dragonflyoss/client",
						Tag:        "v1.3.0",
						PullPolicy: corev1.PullAlways,
						PullSecrets: []corev1.LocalObjectReference{
							{Name: "registry-secret-1"},
							{Name: "registry-secret-2"},
						},
					},
				}

				By("creating a pod")
				pod := makePod("secret-pod", 1, nil)

				By("performing injection")
				binariesInjector.Inject(pod, configWithSecrets)

				By("verifying image pull secrets are appended")
				Expect(pod.Spec.ImagePullSecrets).To(HaveLen(2))
				Expect(pod.Spec.ImagePullSecrets).To(ContainElements(
					corev1.LocalObjectReference{Name: "registry-secret-1"},
					corev1.LocalObjectReference{Name: "registry-secret-2"},
				))
			})
		})

		Context("when the config has a digest in the image reference", func() {
			It("should use the image reference with digest", func() {
				By("creating a config with digest")
				configWithDigest := &Config{
					InitContainerImage: InitContainerImage{
						Registry:    "docker.io",
						Repository:  "dragonflyoss/client",
						Tag:         "v1.3.0",
						Digest:      "sha256:abcdef1234567890",
						PullPolicy:  corev1.PullIfNotPresent,
						PullSecrets: []corev1.LocalObjectReference{},
					},
				}

				By("creating a pod")
				pod := makePod("digest-pod", 1, nil)

				By("performing injection")
				binariesInjector.Inject(pod, configWithDigest)

				By("verifying the init container uses the image reference with digest")
				expectedImage := "docker.io/dragonflyoss/client:v1.3.0@sha256:abcdef1234567890"
				Expect(pod.Spec.InitContainers[0].Image).To(Equal(expectedImage))
			})
		})

		Context("when the pod already has everything injected (idempotency)", func() {
			It("should not duplicate init container, volume, or volume mounts", func() {
				By("creating a pod that is already fully injected")
				expectedImage := defaultConfig.GetInitContainerImageReference()
				pod := makePod("idempotent-pod", 1, nil)
				pod.Spec.InitContainers = []corev1.Container{
					makeExpectedInitContainer(expectedImage, corev1.PullIfNotPresent),
				}
				pod.Spec.Volumes = []corev1.Volume{makeExpectedVolume()}
				pod.Spec.Containers[0].VolumeMounts = makeExpectedVolumeMounts()

				By("creating expected pod (should be unchanged)")
				expectedPod := makePod("idempotent-pod", 1, nil)
				expectedPod.Spec.InitContainers = []corev1.Container{
					makeExpectedInitContainer(expectedImage, corev1.PullIfNotPresent),
				}
				expectedPod.Spec.Volumes = []corev1.Volume{makeExpectedVolume()}
				expectedPod.Spec.Containers[0].VolumeMounts = makeExpectedVolumeMounts()

				By("performing injection")
				binariesInjector.Inject(pod, defaultConfig)

				By("verifying the pod is unchanged")
				Expect(pod).To(Equal(expectedPod))
			})
		})

		Context("when some containers already have volume mounts", func() {
			It("should only inject into containers without existing volume mounts", func() {
				By("creating a pod with two containers where one already has the volume mount")
				pod := makePod("partial-pod", 2, nil)
				pod.Spec.Containers[0].VolumeMounts = makeExpectedVolumeMounts()

				By("performing injection")
				binariesInjector.Inject(pod, defaultConfig)

				By("verifying container-1 is unchanged")
				Expect(pod.Spec.Containers[0].VolumeMounts).To(Equal(makeExpectedVolumeMounts()))

				By("verifying container-2 receives the volume mounts")
				Expect(pod.Spec.Containers[1].VolumeMounts).To(Equal(makeExpectedVolumeMounts()))
			})
		})

		Context("when the pod has no containers", func() {
			It("should inject init container and volume but no volume mounts", func() {
				By("creating a pod with no containers")
				pod := makePod("empty-pod", 0, nil)

				By("performing injection")
				binariesInjector.Inject(pod, defaultConfig)

				By("verifying init container is added")
				Expect(pod.Spec.InitContainers).To(HaveLen(1))

				By("verifying volume is added")
				Expect(pod.Spec.Volumes).To(HaveLen(1))
				Expect(pod.Spec.Volumes[0]).To(Equal(makeExpectedVolume()))

				By("verifying no containers are affected")
				Expect(pod.Spec.Containers).To(BeEmpty())
			})
		})

		Context("when the pod has existing unrelated init containers and volumes", func() {
			It("should append without disturbing existing resources", func() {
				By("creating a pod with existing init containers and volumes")
				pod := makePod("existing-resources-pod", 1, nil)
				pod.Spec.InitContainers = []corev1.Container{
					{Name: "existing-init", Image: "busybox:latest"},
				}
				pod.Spec.Volumes = []corev1.Volume{
					{Name: "existing-volume", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
				}
				pod.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
					{Name: "existing-mount", MountPath: "/data"},
				}

				By("performing injection")
				binariesInjector.Inject(pod, defaultConfig)

				By("verifying existing init container is preserved")
				Expect(pod.Spec.InitContainers).To(HaveLen(2))
				Expect(pod.Spec.InitContainers[0].Name).To(Equal("existing-init"))
				Expect(pod.Spec.InitContainers[1].Name).To(Equal(InitContainerName))

				By("verifying existing volume is preserved")
				Expect(pod.Spec.Volumes).To(HaveLen(2))
				Expect(pod.Spec.Volumes[0].Name).To(Equal("existing-volume"))
				Expect(pod.Spec.Volumes[1].Name).To(Equal(BinaryVolumeName))

				By("verifying existing volume mount is preserved and new ones are appended")
				Expect(pod.Spec.Containers[0].VolumeMounts[0].Name).To(Equal("existing-mount"))
				Expect(pod.Spec.Containers[0].VolumeMounts).To(HaveLen(7)) // 1 existing + 5 new
			})
		})

		Context("when the pod has existing image pull secrets", func() {
			It("should append new pull secrets to existing ones", func() {
				By("creating a pod with existing image pull secrets")
				pod := makePod("secrets-pod", 1, nil)
				pod.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
					{Name: "existing-secret"},
				}

				By("creating config with pull secrets")
				configWithSecrets := &Config{
					InitContainerImage: InitContainerImage{
						Registry:   "docker.io",
						Repository: "dragonflyoss/client",
						Tag:        "v1.3.0",
						PullPolicy: corev1.PullIfNotPresent,
						PullSecrets: []corev1.LocalObjectReference{
							{Name: "new-secret"},
						},
					},
				}

				By("performing injection")
				binariesInjector.Inject(pod, configWithSecrets)

				By("verifying both existing and new pull secrets are present")
				Expect(pod.Spec.ImagePullSecrets).To(HaveLen(2))
				Expect(pod.Spec.ImagePullSecrets[0].Name).To(Equal("existing-secret"))
				Expect(pod.Spec.ImagePullSecrets[1].Name).To(Equal("new-secret"))
			})
		})
	})

	Describe("hasInitContainer", func() {
		Context("when the init container exists", func() {
			It("should return true", func() {
				pod := &corev1.Pod{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{
							{Name: InitContainerName},
						},
					},
				}
				Expect(binariesInjector.hasInitContainer(pod)).To(BeTrue())
			})
		})

		Context("when the init container does not exist", func() {
			It("should return false", func() {
				pod := &corev1.Pod{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{
							{Name: "other-init-container"},
						},
					},
				}
				Expect(binariesInjector.hasInitContainer(pod)).To(BeFalse())
			})
		})

		Context("when there are no init containers", func() {
			It("should return false", func() {
				pod := &corev1.Pod{
					Spec: corev1.PodSpec{},
				}
				Expect(binariesInjector.hasInitContainer(pod)).To(BeFalse())
			})
		})
	})

	Describe("hasVolume", func() {
		Context("when the volume exists", func() {
			It("should return true", func() {
				pod := &corev1.Pod{
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{
							{Name: BinaryVolumeName},
						},
					},
				}
				Expect(binariesInjector.hasVolume(pod)).To(BeTrue())
			})
		})

		Context("when the volume does not exist", func() {
			It("should return false", func() {
				pod := &corev1.Pod{
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{
							{Name: "other-volume"},
						},
					},
				}
				Expect(binariesInjector.hasVolume(pod)).To(BeFalse())
			})
		})

		Context("when there are no volumes", func() {
			It("should return false", func() {
				pod := &corev1.Pod{
					Spec: corev1.PodSpec{},
				}
				Expect(binariesInjector.hasVolume(pod)).To(BeFalse())
			})
		})
	})

	Describe("hasVolumeMount", func() {
		Context("when the volume mount exists", func() {
			It("should return true", func() {
				container := &corev1.Container{
					VolumeMounts: []corev1.VolumeMount{
						{Name: BinaryVolumeName, MountPath: BinaryVolumeMountPath},
					},
				}
				Expect(binariesInjector.hasVolumeMount(container)).To(BeTrue())
			})
		})

		Context("when the volume mount does not exist", func() {
			It("should return false", func() {
				container := &corev1.Container{
					VolumeMounts: []corev1.VolumeMount{
						{Name: "other-volume", MountPath: "/other"},
					},
				}
				Expect(binariesInjector.hasVolumeMount(container)).To(BeFalse())
			})
		})

		Context("when there are no volume mounts", func() {
			It("should return false", func() {
				container := &corev1.Container{}
				Expect(binariesInjector.hasVolumeMount(container)).To(BeFalse())
			})
		})
	})

	Describe("NewBinaries", func() {
		It("should return a non-nil Binaries instance", func() {
			b := NewBinaries()
			Expect(b).NotTo(BeNil())
		})
	})
})
