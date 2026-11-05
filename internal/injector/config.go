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
	"os"
	"path/filepath"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var logger = log.Log.WithName("dragonfly-inject")

// Config is the configuration for dragonfly injection.
type Config struct {
	InitContainerImage InitContainerImage `yaml:"initContainerImage" json:"initContainerImage"`
}

// InitContainerImage is the image configuration for the init container.
type InitContainerImage struct {
	Registry    string                        `yaml:"registry" json:"registry"`
	Repository  string                        `yaml:"repository" json:"repository"`
	Tag         string                        `yaml:"tag" json:"tag"`
	Digest      string                        `yaml:"digest" json:"digest"`
	PullPolicy  corev1.PullPolicy             `yaml:"pullPolicy" json:"pullPolicy"`
	PullSecrets []corev1.LocalObjectReference `yaml:"pullSecrets" json:"pullSecrets"`
}

// GetInitContainerImageReference returns the init container image reference.
func (c *Config) GetInitContainerImageReference() string {
	reference := c.InitContainerImage.Registry + "/" + c.InitContainerImage.Repository
	if c.InitContainerImage.Tag != "" {
		reference += ":" + c.InitContainerImage.Tag
	}
	if c.InitContainerImage.Digest != "" {
		reference += "@" + c.InitContainerImage.Digest
	}

	return reference
}

func DefaultConfig() *Config {
	return &Config{
		InitContainerImage: InitContainerImage{
			Registry:    "docker.io",
			Repository:  "dragonflyoss/client",
			Tag:         "v1.3.0",
			PullPolicy:  "IfNotPresent",
			PullSecrets: []corev1.LocalObjectReference{},
		},
	}
}

type ConfigManager struct {
	mu         sync.RWMutex
	config     *Config
	configPath string
}

func NewConfigManager(injectConfigMapPath string) *ConfigManager {
	configPath := filepath.Join(injectConfigMapPath, "injector.yaml")
	return &ConfigManager{
		mu:         sync.RWMutex{},
		config:     LoadConfig(configPath),
		configPath: configPath,
	}
}

func (cm *ConfigManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	copiedConf := *cm.config
	logger.Info("Get config", "config", copiedConf)
	return &copiedConf
}

func (cm *ConfigManager) Start(ctx context.Context) error {
	logger.Info("Starting config file watcher")

	ticker := time.NewTicker(ConfigReloadInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping config file watcher")
			return nil
		case <-ticker.C:
			logger.Info("Periodic reload check")
			cm.reload()
		}
	}
}

func (cm *ConfigManager) reload() {
	config := LoadConfig(cm.configPath)
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config = config
	logger.Info("Configuration reloaded successfully")
}

func LoadConfig(injectConfigMapPath string) *Config {
	c, err := LoadConfigFromFile(injectConfigMapPath)
	if err != nil {
		logger.Error(err, "load config from file failed")
		logger.Info("use default config")
		c = DefaultConfig()
	}
	return c
}

// LoadConfigFromFile loads config from file.
func LoadConfigFromFile(injectConfigMapPath string) (*Config, error) {
	cf, err := os.ReadFile(injectConfigMapPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(cf, config); err != nil {
		return nil, err
	}

	return config, nil
}
