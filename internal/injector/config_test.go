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
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Config", func() {
	var (
		tempDir string
	)

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
	})

	// Helper to write a Config as a JSON file (k8s yaml.Unmarshal supports JSON).
	writeConfigFile := func(path string, config *Config) {
		data, err := json.Marshal(config)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		err = os.WriteFile(path, data, 0644)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}

	Describe("GetInitContainerImageReference", func() {
		Context("when tag and digest are both set", func() {
			It("should return registry/repository:tag@digest", func() {
				config := &Config{
					InitContainerImage: InitContainerImage{
						Registry:   "docker.io",
						Repository: "dragonflyoss/client",
						Tag:        "v1.3.0",
						Digest:     "sha256:abc123",
					},
				}

				ref := config.GetInitContainerImageReference()
				Expect(ref).To(Equal("docker.io/dragonflyoss/client:v1.3.0@sha256:abc123"))
			})
		})

		Context("when only tag is set", func() {
			It("should return registry/repository:tag", func() {
				config := &Config{
					InitContainerImage: InitContainerImage{
						Registry:   "docker.io",
						Repository: "dragonflyoss/client",
						Tag:        "v1.3.0",
					},
				}

				ref := config.GetInitContainerImageReference()
				Expect(ref).To(Equal("docker.io/dragonflyoss/client:v1.3.0"))
			})
		})

		Context("when only digest is set", func() {
			It("should return registry/repository@digest", func() {
				config := &Config{
					InitContainerImage: InitContainerImage{
						Registry:   "docker.io",
						Repository: "dragonflyoss/client",
						Digest:     "sha256:abc123",
					},
				}

				ref := config.GetInitContainerImageReference()
				Expect(ref).To(Equal("docker.io/dragonflyoss/client@sha256:abc123"))
			})
		})

		Context("when neither tag nor digest is set", func() {
			It("should return registry/repository", func() {
				config := &Config{
					InitContainerImage: InitContainerImage{
						Registry:   "docker.io",
						Repository: "dragonflyoss/client",
					},
				}

				ref := config.GetInitContainerImageReference()
				Expect(ref).To(Equal("docker.io/dragonflyoss/client"))
			})
		})
	})

	Describe("DefaultConfig", func() {
		It("should return the correct default configuration", func() {
			By("creating a new default config")
			defaultConfig := DefaultConfig()

			By("verifying the default values")
			Expect(defaultConfig.InitContainerImage.Registry).To(Equal("docker.io"))
			Expect(defaultConfig.InitContainerImage.Repository).To(Equal("dragonflyoss/client"))
			Expect(defaultConfig.InitContainerImage.Tag).To(Equal("v1.3.0"))
			Expect(defaultConfig.InitContainerImage.Digest).To(BeEmpty())
			Expect(defaultConfig.InitContainerImage.PullPolicy).To(Equal(corev1.PullPolicy("IfNotPresent")))
			Expect(defaultConfig.InitContainerImage.PullSecrets).To(BeEmpty())
		})
	})

	Describe("LoadConfigFromFile", func() {
		Context("when loading configuration from file", func() {
			It("should load valid config file successfully", func() {
				By("creating a valid config file")
				configPath := filepath.Join(tempDir, "valid-config.yaml")
				configData := &Config{
					InitContainerImage: InitContainerImage{
						Registry:   "custom-registry.io",
						Repository: "custom/client",
						Tag:        "v2.0.0",
						PullPolicy: corev1.PullAlways,
						PullSecrets: []corev1.LocalObjectReference{
							{Name: "my-secret"},
						},
					},
				}
				writeConfigFile(configPath, configData)

				By("loading the config from file")
				loadedConfig, err := LoadConfigFromFile(configPath)
				Expect(err).NotTo(HaveOccurred())

				By("verifying the loaded configuration")
				Expect(loadedConfig.InitContainerImage.Registry).To(Equal("custom-registry.io"))
				Expect(loadedConfig.InitContainerImage.Repository).To(Equal("custom/client"))
				Expect(loadedConfig.InitContainerImage.Tag).To(Equal("v2.0.0"))
				Expect(loadedConfig.InitContainerImage.PullPolicy).To(Equal(corev1.PullAlways))
				Expect(loadedConfig.InitContainerImage.PullSecrets).To(HaveLen(1))
				Expect(loadedConfig.InitContainerImage.PullSecrets[0].Name).To(Equal("my-secret"))
			})

			It("should return error for non-existent file", func() {
				By("attempting to load a non-existent file")
				_, err := LoadConfigFromFile(filepath.Join(tempDir, "non-existent-file.yaml"))
				Expect(err).To(HaveOccurred())
				Expect(os.IsNotExist(err)).To(BeTrue())
			})

			It("should return error for invalid YAML content", func() {
				By("creating a file with invalid YAML content")
				configPath := filepath.Join(tempDir, "invalid.yaml")
				invalidYAML := "invalid: yaml: content: ["
				err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
				Expect(err).NotTo(HaveOccurred())

				By("attempting to load the invalid YAML file")
				_, err = LoadConfigFromFile(configPath)
				Expect(err).To(HaveOccurred())
			})

			It("should handle partial config with zero values", func() {
				By("creating a partial config file with only registry set")
				configPath := filepath.Join(tempDir, "partial-config.yaml")
				partialConfig := &Config{
					InitContainerImage: InitContainerImage{
						Registry: "partial-registry.io",
					},
				}
				writeConfigFile(configPath, partialConfig)

				By("loading the partial config")
				loadedConfig, err := LoadConfigFromFile(configPath)
				Expect(err).NotTo(HaveOccurred())

				By("verifying the loaded configuration with zero values")
				Expect(loadedConfig.InitContainerImage.Registry).To(Equal("partial-registry.io"))
				Expect(loadedConfig.InitContainerImage.Repository).To(BeEmpty())
				Expect(loadedConfig.InitContainerImage.Tag).To(BeEmpty())
				Expect(loadedConfig.InitContainerImage.Digest).To(BeEmpty())
				Expect(loadedConfig.InitContainerImage.PullPolicy).To(BeEmpty())
				Expect(loadedConfig.InitContainerImage.PullSecrets).To(BeNil())
			})

			It("should handle empty JSON file", func() {
				By("creating an empty config file")
				configPath := filepath.Join(tempDir, "empty-config.yaml")
				err := os.WriteFile(configPath, []byte("{}"), 0644)
				Expect(err).NotTo(HaveOccurred())

				By("loading the empty config")
				loadedConfig, err := LoadConfigFromFile(configPath)
				Expect(err).NotTo(HaveOccurred())

				By("verifying the config has zero values")
				Expect(loadedConfig.InitContainerImage.Registry).To(BeEmpty())
				Expect(loadedConfig.InitContainerImage.Repository).To(BeEmpty())
			})
		})
	})

	Describe("LoadConfig", func() {
		Context("when loading configuration with fallback behavior", func() {
			It("should load config from file when file exists", func() {
				By("creating an existing config file")
				configPath := filepath.Join(tempDir, "existing-config.yaml")
				configData := &Config{
					InitContainerImage: InitContainerImage{
						Registry:   "my-registry.io",
						Repository: "my-org/my-image",
						Tag:        "v3.0.0",
						PullPolicy: corev1.PullNever,
					},
				}
				writeConfigFile(configPath, configData)

				By("loading the config")
				loadedConfig := LoadConfig(configPath)

				By("verifying the loaded configuration")
				Expect(loadedConfig.InitContainerImage.Registry).To(Equal("my-registry.io"))
				Expect(loadedConfig.InitContainerImage.Repository).To(Equal("my-org/my-image"))
				Expect(loadedConfig.InitContainerImage.Tag).To(Equal("v3.0.0"))
				Expect(loadedConfig.InitContainerImage.PullPolicy).To(Equal(corev1.PullNever))
			})

			It("should return default config when file does not exist", func() {
				By("loading a non-existent file")
				loadedConfig := LoadConfig(filepath.Join(tempDir, "non-existent-file.yaml"))
				expected := DefaultConfig()

				By("verifying the default configuration is returned")
				Expect(loadedConfig.InitContainerImage.Registry).To(Equal(expected.InitContainerImage.Registry))
				Expect(loadedConfig.InitContainerImage.Repository).To(Equal(expected.InitContainerImage.Repository))
				Expect(loadedConfig.InitContainerImage.Tag).To(Equal(expected.InitContainerImage.Tag))
				Expect(loadedConfig.InitContainerImage.PullPolicy).To(Equal(expected.InitContainerImage.PullPolicy))
			})

			It("should return default config when file is invalid", func() {
				By("creating an invalid config file")
				configPath := filepath.Join(tempDir, "invalid.yaml")
				invalidYAML := "invalid: yaml: content: ["
				err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
				Expect(err).NotTo(HaveOccurred())

				By("loading the invalid file")
				loadedConfig := LoadConfig(configPath)
				expected := DefaultConfig()

				By("verifying the default configuration is returned")
				Expect(loadedConfig.InitContainerImage.Registry).To(Equal(expected.InitContainerImage.Registry))
				Expect(loadedConfig.InitContainerImage.Repository).To(Equal(expected.InitContainerImage.Repository))
				Expect(loadedConfig.InitContainerImage.Tag).To(Equal(expected.InitContainerImage.Tag))
				Expect(loadedConfig.InitContainerImage.PullPolicy).To(Equal(expected.InitContainerImage.PullPolicy))
			})
		})
	})

	Describe("ConfigManager", func() {
		var (
			configManager *ConfigManager
		)

		Context("with basic functionality", func() {
			BeforeEach(func() {
				By("creating initial configuration")
				configPath := filepath.Join(tempDir, "injector.yaml")
				initialConfig := &Config{
					InitContainerImage: InitContainerImage{
						Registry:    "docker.io",
						Repository:  "dragonflyoss/client",
						Tag:         "v1.0.0",
						PullPolicy:  corev1.PullIfNotPresent,
						PullSecrets: []corev1.LocalObjectReference{},
					},
				}
				writeConfigFile(configPath, initialConfig)

				By("creating the ConfigManager")
				configManager = NewConfigManager(tempDir)
				Expect(configManager).NotTo(BeNil())
			})

			It("should get the correct configuration", func() {
				By("retrieving the configuration")
				config := configManager.GetConfig()

				By("verifying the configuration values")
				Expect(config.InitContainerImage.Registry).To(Equal("docker.io"))
				Expect(config.InitContainerImage.Repository).To(Equal("dragonflyoss/client"))
				Expect(config.InitContainerImage.Tag).To(Equal("v1.0.0"))
				Expect(config.InitContainerImage.PullPolicy).To(Equal(corev1.PullIfNotPresent))
			})

			It("should return a copy of the config, not the original", func() {
				By("retrieving the configuration twice")
				config1 := configManager.GetConfig()
				config2 := configManager.GetConfig()

				By("modifying the first copy")
				config1.InitContainerImage.Tag = "modified"

				By("verifying the second copy is unaffected")
				Expect(config2.InitContainerImage.Tag).To(Equal("v1.0.0"))
			})

			It("should reload configuration correctly", func() {
				By("updating the configuration file")
				updatedConfig := &Config{
					InitContainerImage: InitContainerImage{
						Registry:   "updated-registry.io",
						Repository: "updated/client",
						Tag:        "v2.0.0",
						PullPolicy: corev1.PullAlways,
					},
				}
				configPath := filepath.Join(tempDir, "injector.yaml")
				writeConfigFile(configPath, updatedConfig)

				By("triggering configuration reload")
				configManager.reload()

				By("verifying the updated configuration")
				config := configManager.GetConfig()
				Expect(config.InitContainerImage.Registry).To(Equal("updated-registry.io"))
				Expect(config.InitContainerImage.Repository).To(Equal("updated/client"))
				Expect(config.InitContainerImage.Tag).To(Equal("v2.0.0"))
				Expect(config.InitContainerImage.PullPolicy).To(Equal(corev1.PullAlways))
			})
		})

		Context("when configuration file does not exist", func() {
			It("should use default configuration", func() {
				By("creating ConfigManager without config file")
				configManager = NewConfigManager(tempDir)

				By("verifying default configuration is used")
				config := configManager.GetConfig()
				expected := DefaultConfig()
				Expect(config.InitContainerImage.Registry).To(Equal(expected.InitContainerImage.Registry))
				Expect(config.InitContainerImage.Repository).To(Equal(expected.InitContainerImage.Repository))
				Expect(config.InitContainerImage.Tag).To(Equal(expected.InitContainerImage.Tag))
				Expect(config.InitContainerImage.PullPolicy).To(Equal(expected.InitContainerImage.PullPolicy))
			})
		})

		Context("Start and Stop functionality", func() {
			It("should start and stop gracefully", func() {
				By("creating ConfigManager")
				configManager = NewConfigManager(tempDir)

				By("creating a cancellable context")
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				By("starting the ConfigManager")
				done := make(chan error)
				go func() {
					done <- configManager.Start(ctx)
				}()

				By("waiting for startup")
				time.Sleep(100 * time.Millisecond)

				By("cancelling the context")
				cancel()

				By("waiting for graceful shutdown")
				Eventually(done, 5*time.Second).Should(Receive(BeNil()))
			})
		})

		Context("concurrent access", func() {
			BeforeEach(func() {
				By("creating initial configuration for concurrent testing")
				configPath := filepath.Join(tempDir, "injector.yaml")
				configData := &Config{
					InitContainerImage: InitContainerImage{
						Registry:    "docker.io",
						Repository:  "dragonflyoss/client",
						Tag:         "v1.0.0",
						PullPolicy:  corev1.PullIfNotPresent,
						PullSecrets: []corev1.LocalObjectReference{},
					},
				}
				writeConfigFile(configPath, configData)

				configManager = NewConfigManager(tempDir)
			})

			It("should handle concurrent reads and reloads safely", func() {
				By("starting concurrent readers and reloaders")
				done := make(chan bool, 2)

				go func() {
					defer GinkgoRecover()
					for range 100 {
						config := configManager.GetConfig()
						Expect(config).NotTo(BeNil())
						Expect(config.InitContainerImage.Registry).NotTo(BeEmpty())
					}
					done <- true
				}()

				go func() {
					defer GinkgoRecover()
					for range 100 {
						configManager.reload()
					}
					done <- true
				}()

				By("waiting for all goroutines to complete")
				Eventually(done, 5*time.Second).Should(Receive(BeTrue()))
				Eventually(done, 5*time.Second).Should(Receive(BeTrue()))
			})
		})
	})
})
