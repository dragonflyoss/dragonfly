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
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/yaml"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var logger = logf.Log.WithName("injector")

const (
	// InjectConfigMapPath is the path where the injector configuration is mounted
	InjectConfigMapPath string = "/etc/dragonfly-injector"

	// Namespace labels for injection control
	NamespaceInjectLabelName  string = "dragonfly.io/inject"
	NamespaceInjectLabelValue string = "true"

	// Pod annotation for injection control
	PodInjectAnnotationName  string = "dragonfly.io/inject"
	PodInjectAnnotationValue string = "true"

	// Environment variable control
	NodeNameEnvName   string = "NODE_NAME"
	ProxyPortEnvName  string = "DRAGONFLY_PROXY_PORT"
	ProxyPortEnvValue int    = 4001 // Default port of dragonfly dfdaemon proxy
	ProxyEnvName      string = "DRAGONFLY_INJECT_PROXY"

	// Dfdaemon unix socket volume control
	DfdaemonUnixSockVolumeName string = "dfdaemon-unix-sock"
	DfdaemonUnixSockPath       string = "/var/run/dragonfly/dfdaemon.sock"

	// CLI tools initContainer control
	CliToolsImageAnnotation   string = "dragonfly.io/cli-tools-image"
	CliToolsImage             string = "dragonflyoss/toolkits:latest"
	CliToolsInitContainerName string = "d7y-cli-tools"
	CliToolsVolumeName        string = CliToolsInitContainerName + "-volume"
	CliToolsDirPath           string = "/dragonfly-tools"
	CliToolsPathEnvName       string = "DRAGONFLY_TOOLS_PATH"

	// ConfigReloadWaitTime is the interval for checking configuration updates
	ConfigReloadWaitTime time.Duration = 15 * time.Second
)

// InjectConf represents the configuration for the dragonfly injector webhook
type InjectConf struct {
	Enable          bool   `yaml:"enable" json:"enable"`
	ProxyPort       int    `yaml:"proxy_port" json:"proxy_port"`
	CliToolsImage   string `yaml:"cli_tools_image" json:"cli_tools_image"`
	CliToolsDirPath string `yaml:"cli_tools_dir_path" json:"cli_tools_dir_path"`
}

// NewDefaultInjectConf returns a new InjectConf with default values
func NewDefaultInjectConf() *InjectConf {
	return &InjectConf{
		Enable:          true,
		ProxyPort:       ProxyPortEnvValue,
		CliToolsImage:   CliToolsImage,
		CliToolsDirPath: CliToolsDirPath,
	}
}

// ConfigManager manages the injector configuration with hot-reload capability
type ConfigManager struct {
	mu         sync.RWMutex
	config     *InjectConf
	configPath string
}

// NewConfigManager creates a new ConfigManager instance
func NewConfigManager(injectConfigMapPath string) *ConfigManager {
	configPath := filepath.Join(injectConfigMapPath, "config.yaml")
	return &ConfigManager{
		mu:         sync.RWMutex{},
		config:     LoadInjectConf(configPath),
		configPath: configPath,
	}
}

// GetConfig returns a copy of the current configuration
func (cm *ConfigManager) GetConfig() *InjectConf {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	copiedConf := *cm.config
	logger.V(1).Info("Retrieved configuration", "config", copiedConf)
	return &copiedConf
}

// Start begins the configuration file watcher
func (cm *ConfigManager) Start(ctx context.Context) error {
	logger.Info("Starting configuration file watcher")

	ticker := time.NewTicker(ConfigReloadWaitTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping configuration file watcher")
			return nil
		case <-ticker.C:
			logger.V(1).Info("Performing periodic configuration reload check")
			cm.reload()
		}
	}
}

// reload reloads the configuration from the file
func (cm *ConfigManager) reload() {
	config := LoadInjectConf(cm.configPath)
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config = config
	logger.Info("Configuration reloaded successfully")
}

// LoadInjectConf loads the injector configuration from the specified path
func LoadInjectConf(injectConfigMapPath string) *InjectConf {
	ic, err := LoadInjectConfFromFile(injectConfigMapPath)
	if err != nil {
		logger.Error(err, "Failed to load configuration from file, using default configuration")
		ic = NewDefaultInjectConf()
	}
	return ic
}

// LoadInjectConfFromFile loads the injector configuration from a YAML file
func LoadInjectConfFromFile(injectConfigMapPath string) (*InjectConf, error) {
	data, err := os.ReadFile(injectConfigMapPath)
	if err != nil {
		return nil, err
	}

	injectConf := &InjectConf{}
	if err := yaml.Unmarshal(data, injectConf); err != nil {
		return nil, err
	}

	return injectConf, nil
}
