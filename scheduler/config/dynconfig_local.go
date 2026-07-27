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

package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	dc "d7y.io/dragonfly/v2/internal/dynconfig"
	"d7y.io/dragonfly/v2/manager/types"
)

// LocalDynconfigData is the dynamic configuration loaded from the local file.
type LocalDynconfigData struct {
	// Applications is the applications configuration.
	Applications []*managerv2.Application `yaml:"applications" mapstructure:"applications"`

	// SeedPeerClusterConfig is the seed peer cluster configuration.
	SeedPeerClusterConfig types.SeedPeerClusterConfig `yaml:"seedPeerClusterConfig" mapstructure:"seedPeerClusterConfig"`

	// SchedulerClusterConfig is the scheduler cluster configuration.
	SchedulerClusterConfig types.SchedulerClusterConfig `yaml:"schedulerClusterConfig" mapstructure:"schedulerClusterConfig"`

	// SchedulerClusterClientConfig is the client configuration.
	SchedulerClusterClientConfig types.SchedulerClusterClientConfig `yaml:"schedulerClusterClientConfig" mapstructure:"schedulerClusterClientConfig"`
}

// NewLocalDynconfigData returns the default local dynamic configuration data.
func NewLocalDynconfigData() *LocalDynconfigData {
	return &LocalDynconfigData{
		SeedPeerClusterConfig: types.SeedPeerClusterConfig{
			LoadLimit: DefaultSeedPeerConcurrentUploadLimit,
		},
		SchedulerClusterConfig: types.SchedulerClusterConfig{
			CandidateParentLimit: DefaultSchedulerCandidateParentLimit,
			FilterParentLimit:    DefaultSchedulerFilterParentLimit,
		},
		SchedulerClusterClientConfig: types.SchedulerClusterClientConfig{
			LoadLimit: DefaultPeerConcurrentUploadLimit,
		},
	}
}

// localDynconfig is the local dynconfig, which loads the dynamic configuration
// from the local file instead of the manager. It is used when the manager addr
// is not configured.
type localDynconfig struct {
	dc.Dynconfig[LocalDynconfigData]
}

// newLocalDynconfig returns a new local dynconfig instance, which loads the
// dynamic configuration from the local file specified by configPath.
func newLocalDynconfig(configPath string, cfg *Config) (DynconfigInterface, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		content, err := yaml.Marshal(NewLocalDynconfigData())
		if err != nil {
			return nil, err
		}

		if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
			return nil, err
		}

		if err := os.WriteFile(configPath, content, 0600); err != nil {
			return nil, err
		}

		logger.Infof("generate local dynconfig %s with default values", configPath)
	}

	client, err := dc.New[LocalDynconfigData](
		&localFileClient{configPath: configPath},
		cfg.DynConfig.RefreshInterval,
	)
	if err != nil {
		return nil, err
	}

	return &localDynconfig{Dynconfig: client}, nil
}

// GetApplications returns the applications config.
func (d *localDynconfig) GetApplications() ([]*managerv2.Application, error) {
	data, err := d.Get()
	if err != nil {
		return nil, err
	}

	if len(data.Applications) == 0 {
		return nil, errors.New("application not found")
	}

	return data.Applications, nil
}

// GetSeedPeerClusterConfig returns the seed peer cluster config.
func (d *localDynconfig) GetSeedPeerClusterConfig() (types.SeedPeerClusterConfig, error) {
	data, err := d.Get()
	if err != nil {
		return types.SeedPeerClusterConfig{}, err
	}

	return data.SeedPeerClusterConfig, nil
}

// GetSchedulerClusterConfig returns the scheduler cluster config.
func (d *localDynconfig) GetSchedulerClusterConfig() (types.SchedulerClusterConfig, error) {
	data, err := d.Get()
	if err != nil {
		return types.SchedulerClusterConfig{}, err
	}

	return data.SchedulerClusterConfig, nil
}

// GetSchedulerClusterClientConfig returns the client config.
func (d *localDynconfig) GetSchedulerClusterClientConfig() (types.SchedulerClusterClientConfig, error) {
	data, err := d.Get()
	if err != nil {
		return types.SchedulerClusterClientConfig{}, err
	}

	return data.SchedulerClusterClientConfig, nil
}

// localFileClient is the client that loads the dynamic configuration from the
// local file. The file is reread on each load, so updating the mounted file
// (e.g. a Kubernetes ConfigMap) propagates the new configuration within one
// refresh interval.
type localFileClient struct {
	configPath string
}

// Get returns the dynamic configuration loaded from the local file.
func (c *localFileClient) Get() (any, error) {
	v := viper.New()
	v.SetConfigFile(c.configPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var data LocalDynconfigData
	if err := v.Unmarshal(&data); err != nil {
		return nil, err
	}

	return data, nil
}
