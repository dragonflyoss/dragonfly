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

package v1

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"d7y.io/dragonfly/v2/internal/injector"
)

// mockInjector is a mock implementation of the Injector interface for testing purposes.
// It records whether its Inject method has been called.
type mockInjector struct {
	called bool
	config *injector.Config
}

func (m *mockInjector) Inject(pod *corev1.Pod, config *injector.Config) {
	m.called = true
	m.config = config
}

func (m *mockInjector) Reset() {
	m.called = false
	m.config = nil
}

var _ = Describe("Pod Webhook", func() {
	var (
		defaulter   *PodCustomDefaulter
		mockInj     *mockInjector
		ctx         context.Context
		testPod     *corev1.Pod
		configMgr   *injector.ConfigManager
		tempDir     string
		fakeClient  client.Client
		scheme      *runtime.Scheme
		testNsName  string
		testPodName string
	)

	BeforeEach(func() {
		ctx = context.Background()
		testNsName = "test-namespace"
		testPodName = "test-pod"

		// Create a mock injector to verify the injection logic
		mockInj = &mockInjector{}

		tempDir = GinkgoT().TempDir()

		// Write a predictable config file for the test
		testConfig := &injector.Config{
			InitContainerImage: injector.InitContainerImage{
				Registry:   "docker.io",
				Repository: "dragonflyoss/client",
				Tag:        "latest",
			},
		}
		yamlData, err := yaml.Marshal(testConfig)
		Expect(err).NotTo(HaveOccurred())

		configPath := filepath.Join(tempDir, "config.yaml")
		err = os.WriteFile(configPath, yamlData, 0644)
		Expect(err).NotTo(HaveOccurred())

		// Initialize ConfigManager with the temp config path
		configMgr = injector.NewConfigManager(tempDir)

		// Initialize the scheme and fake client
		scheme = runtime.NewScheme()
		err = corev1.AddToScheme(scheme)
		Expect(err).NotTo(HaveOccurred())

		// Base test pod object
		testPod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        testPodName,
				Namespace:   testNsName,
				Annotations: make(map[string]string),
			},
		}
	})

	// This function centralizes the setup of the defaulter for each context
	setupDefaulter := func(initObjs ...client.Object) {
		// Create a new fake client for each test scenario to ensure isolation
		fakeClient = fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjs...).Build()
		defaulter = NewPodCustomDefaulter(fakeClient, configMgr)
		// CRITICAL: Replace the real injectors with our mock for testing purposes
		defaulter.injectors = []Injector{mockInj}
	}

	Context("When evaluating if injection is required", func() {

		Context("and injection is enabled by Namespace label", func() {
			It("should inject the pod", func() {
				By("creating a namespace with the injection label")
				labeledNs := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: testNsName,
						Labels: map[string]string{
							injector.NamespaceInjectLabelName: injector.NamespaceInjectLabelValue,
						},
					},
				}
				setupDefaulter(labeledNs)

				By("calling the Default method")
				err := defaulter.Default(ctx, testPod)
				Expect(err).NotTo(HaveOccurred())

				By("verifying the injector was called")
				Expect(mockInj.called).To(BeTrue())
				Expect(mockInj.config).NotTo(BeNil())
			})
		})

		Context("and injection is enabled by Pod annotation", func() {
			It("should inject the pod", func() {
				By("creating a namespace without the injection label")
				unlabeledNs := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{Name: testNsName},
				}
				setupDefaulter(unlabeledNs)

				By("annotating the pod to enable injection")
				testPod.Annotations[injector.PodInjectAnnotationName] = injector.PodInjectAnnotationValue

				By("calling the Default method")
				err := defaulter.Default(ctx, testPod)
				Expect(err).NotTo(HaveOccurred())

				By("verifying the injector was called")
				Expect(mockInj.called).To(BeTrue())
			})
		})

		Context("and injection is enabled by both Namespace and Pod", func() {
			It("should inject the pod once", func() {
				By("creating a namespace with the injection label")
				labeledNs := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: testNsName,
						Labels: map[string]string{
							injector.NamespaceInjectLabelName: injector.NamespaceInjectLabelValue,
						},
					},
				}
				setupDefaulter(labeledNs)

				By("annotating the pod to enable injection")
				testPod.Annotations[injector.PodInjectAnnotationName] = injector.PodInjectAnnotationValue

				By("calling the Default method")
				err := defaulter.Default(ctx, testPod)
				Expect(err).NotTo(HaveOccurred())

				By("verifying the injector was called")
				Expect(mockInj.called).To(BeTrue())
			})
		})

		Context("and injection is not required", func() {
			It("should not inject the pod unless label or annotation is present", func() {
				By("creating a namespace without the injection label")
				unlabeledNs := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testNsName}}
				setupDefaulter(unlabeledNs)

				By("calling the Default method on a pod with no relevant annotation")
				err := defaulter.Default(ctx, testPod)
				Expect(err).NotTo(HaveOccurred())

				By("verifying the injector was NOT called")
				Expect(mockInj.called).To(BeFalse())
			})

			It("should not inject if the namespace label has the wrong value", func() {
				By("creating a namespace with an incorrect injection label value")
				incorrectlyLabeledNs := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name:   testNsName,
						Labels: map[string]string{injector.NamespaceInjectLabelName: "disabled"},
					},
				}
				setupDefaulter(incorrectlyLabeledNs)

				By("calling the Default method")
				err := defaulter.Default(ctx, testPod)
				Expect(err).NotTo(HaveOccurred())

				By("verifying the injector was NOT called")
				Expect(mockInj.called).To(BeFalse())
			})

			It("should not inject if the pod annotation has the wrong value", func() {
				By("creating a namespace without the injection label")
				unlabeledNs := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testNsName}}
				setupDefaulter(unlabeledNs)

				By("annotating the pod with an incorrect value")
				testPod.Annotations[injector.PodInjectAnnotationName] = "false"

				By("calling the Default method")
				err := defaulter.Default(ctx, testPod)
				Expect(err).NotTo(HaveOccurred())

				By("verifying the injector was NOT called")
				Expect(mockInj.called).To(BeFalse())
			})

			It("should not inject if the pod's namespace cannot be found", func() {
				By("creating a fake client with NO namespace object")
				setupDefaulter() // No objects added to client

				By("calling the Default method")
				err := defaulter.Default(ctx, testPod)
				Expect(err).NotTo(HaveOccurred()) // The method logs error but returns nil

				By("verifying the injector was NOT called")
				Expect(mockInj.called).To(BeFalse())
			})
		})

		Context("when the object is not a Pod", func() {
			It("should return an error", func() {
				By("creating a non-pod object")
				notAPod := &corev1.ConfigMap{}
				setupDefaulter()

				By("calling the Default method")
				err := defaulter.Default(ctx, notAPod)

				By("verifying that an error is returned")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("expected an Pod object but got"))
				Expect(mockInj.called).To(BeFalse())
			})
		})
	})
})
