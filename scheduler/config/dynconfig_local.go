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
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/atomic"

	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/pkg/dfpath"
)

var (
	// localDynconfigFileName is the file name of the local dynamic configuration.
	localDynconfigFileName = "dynconfig.yaml"
)

// LocalDynconfigData is the dynamic configuration loaded from the local file.
type LocalDynconfigData struct {
	// RefreshInterval is the interval for refreshing the local dynamic configuration.
	RefreshInterval time.Duration `yaml:"refreshInterval" mapstructure:"refreshInterval"`

	// SchedulerClusterConfig is the scheduler cluster configuration.
	SchedulerClusterConfig types.SchedulerClusterConfig `yaml:"schedulerClusterConfig" mapstructure:"schedulerClusterConfig"`
}

// localDynconfig is the local dynconfig, which loads the dynamic configuration
// from the local file instead of the manager. It is used when the manager
// is not configured.
type localDynconfig struct {
	configPath string
	data       *atomic.Pointer[LocalDynconfigData]
	done       chan struct{}
}

// newLocalDynconfig returns a new local dynconfig instance.
func newLocalDynconfig(configPath string) (DynconfigInterface, error) {
	d := &localDynconfig{
		configPath: configPath,
		data:       atomic.NewPointer[LocalDynconfigData](nil),
		done:       make(chan struct{}),
	}

	if err := d.load(); err != nil {
		return nil, err
	}

	return d, nil
}

// localDynconfigPath returns the path of the local dynamic configuration file,
// which is located in the same directory as the config file.
func localDynconfigPath() string {
	configDir := dfpath.DefaultConfigDir
	if configFile := viper.ConfigFileUsed(); configFile != "" {
		configDir = filepath.Dir(configFile)
	}

	return filepath.Join(configDir, localDynconfigFileName)
}

// GetScheduler returns the scheduler config from manager. It is not supported
// when the manager is not configured.
func (d *localDynconfig) GetScheduler() (*managerv2.Scheduler, error) {
	return nil, errors.New("manager is not configured")
}

// GetApplications returns the applications config from manager. It is not
// supported when the manager is not configured.
func (d *localDynconfig) GetApplications() ([]*managerv2.Application, error) {
	return nil, errors.New("manager is not configured")
}

// GetSeedPeers returns the seed peers config from manager. It is not supported
// when the manager is not configured.
func (d *localDynconfig) GetSeedPeers() ([]*managerv2.SeedPeer, error) {
	return nil, errors.New("manager is not configured")
}

// GetSeedPeerClusterConfig returns the seed peer cluster config. It is not
// supported when the manager is not configured.
func (d *localDynconfig) GetSeedPeerClusterConfig() (types.SeedPeerClusterConfig, error) {
	return types.SeedPeerClusterConfig{}, errors.New("manager is not configured")
}

// GetSchedulerCluster returns the scheduler cluster config from manager. It is
// not supported when the manager is not configured.
func (d *localDynconfig) GetSchedulerCluster() (*managerv2.SchedulerCluster, error) {
	return nil, errors.New("manager is not configured")
}

// GetSchedulerClusterConfig returns the scheduler cluster config.
func (d *localDynconfig) GetSchedulerClusterConfig() (types.SchedulerClusterConfig, error) {
	data := d.data.Load()
	if data == nil {
		return types.SchedulerClusterConfig{}, errors.New("invalid data")
	}

	return data.SchedulerClusterConfig, nil
}

// GetSchedulerClusterClientConfig returns the client config. It is not
// supported when the manager is not configured.
func (d *localDynconfig) GetSchedulerClusterClientConfig() (types.SchedulerClusterClientConfig, error) {
	return types.SchedulerClusterClientConfig{}, errors.New("manager is not configured")
}

// Get returns the dynamic config from manager. It is not supported when the
// manager is not configured.
func (d *localDynconfig) Get() (*DynconfigData, error) {
	return nil, errors.New("manager is not configured")
}

// Serve the dynconfig listening service.
func (d *localDynconfig) Serve() error {
	go d.refresh()
	return nil
}

// Stop the dynconfig listening service.
func (d *localDynconfig) Stop() error {
	close(d.done)
	return nil
}

// refresh reloads the local dynamic configuration periodically.
func (d *localDynconfig) refresh() {
	for {
		select {
		case <-time.After(d.refreshInterval()):
			if err := d.load(); err != nil {
				logger.Warnf("refresh local dynconfig %s failed: %s", d.configPath, err.Error())
			}
		case <-d.done:
			return
		}
	}
}

// refreshInterval returns the interval for refreshing the local dynamic
// configuration.
func (d *localDynconfig) refreshInterval() time.Duration {
	if data := d.data.Load(); data != nil && data.RefreshInterval > 0 {
		return data.RefreshInterval
	}

	return DefaultLocalDynconfigRefreshInterval
}

// load loads the dynamic configuration from the local file. The local file is
// reread on each load, so updating the mounted file (e.g. a Kubernetes
// ConfigMap) propagates the new configuration within one refresh interval.
func (d *localDynconfig) load() error {
	v := viper.New()
	v.SetConfigFile(d.configPath)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var data LocalDynconfigData
	if err := v.Unmarshal(&data); err != nil {
		return err
	}
	d.data.Store(&data)

	return nil
}
