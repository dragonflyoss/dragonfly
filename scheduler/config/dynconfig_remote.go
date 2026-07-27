/*
 *     Copyright 2020 The Dragonfly Authors
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
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"

	dc "d7y.io/dragonfly/v2/internal/dynconfig"
	"d7y.io/dragonfly/v2/manager/types"
	managerclient "d7y.io/dragonfly/v2/pkg/rpc/manager/client"
	"d7y.io/dragonfly/v2/pkg/slices"
)

var (
	// cacheFileName is the cache file name of the remote dynamic configuration.
	cacheFileName = "scheduler"
)

// DynconfigData is the dynamic configuration fetched from the manager.
type DynconfigData struct {
	// Scheduler is the scheduler config from manager.
	Scheduler *managerv2.Scheduler

	// Applications is the applications config from manager.
	Applications []*managerv2.Application
}

// remoteDynconfig is the remote dynconfig, which fetches the dynamic
// configuration from the manager.
type remoteDynconfig struct {
	dc.Dynconfig[DynconfigData]
	done                 chan struct{}
	cachePath            string
	transportCredentials credentials.TransportCredentials
	mu                   *sync.Mutex
}

// newRemoteDynconfig returns a new remote dynconfig instance.
func newRemoteDynconfig(rawManagerClient managerclient.V2, cacheDir string, cfg *Config, transportCredentials credentials.TransportCredentials) (DynconfigInterface, error) {
	cachePath := filepath.Join(cacheDir, cacheFileName)
	d := &remoteDynconfig{
		done:                 make(chan struct{}),
		cachePath:            cachePath,
		transportCredentials: transportCredentials,
		mu:                   &sync.Mutex{},
	}

	if rawManagerClient != nil {
		client, err := dc.New[DynconfigData](
			newManagerClient(rawManagerClient, cfg),
			cachePath,
			cfg.DynConfig.RefreshInterval,
		)
		if err != nil {
			return nil, err
		}

		d.Dynconfig = client
	}

	return d, nil
}

// GetApplications returns the applications config from manager.
func (d *remoteDynconfig) GetApplications() ([]*managerv2.Application, error) {
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
func (d *remoteDynconfig) GetSeedPeerClusterConfig() (types.SeedPeerClusterConfig, error) {
	seedPeers, err := d.getSeedPeers()
	if err != nil {
		return types.SeedPeerClusterConfig{}, err
	}

	if len(seedPeers) == 0 {
		return types.SeedPeerClusterConfig{}, errors.New("seed peer not found ")
	}

	var config types.SeedPeerClusterConfig
	if err := json.Unmarshal(seedPeers[0].SeedPeerCluster.Config, &config); err != nil {
		return types.SeedPeerClusterConfig{}, err
	}

	return config, nil
}

// GetSchedulerClusterConfig returns the scheduler cluster config.
func (d *remoteDynconfig) GetSchedulerClusterConfig() (types.SchedulerClusterConfig, error) {
	schedulerCluster, err := d.getSchedulerCluster()
	if err != nil {
		return types.SchedulerClusterConfig{}, err
	}

	var config types.SchedulerClusterConfig
	if err := json.Unmarshal(schedulerCluster.Config, &config); err != nil {
		return types.SchedulerClusterConfig{}, err
	}

	return config, nil
}

// GetSchedulerClusterClientConfig returns the client config.
func (d *remoteDynconfig) GetSchedulerClusterClientConfig() (types.SchedulerClusterClientConfig, error) {
	schedulerCluster, err := d.getSchedulerCluster()
	if err != nil {
		return types.SchedulerClusterClientConfig{}, err
	}

	var config types.SchedulerClusterClientConfig
	if err := json.Unmarshal(schedulerCluster.ClientConfig, &config); err != nil {
		return types.SchedulerClusterClientConfig{}, err
	}

	return config, nil
}

// Serve the dynconfig listening service.
func (d *remoteDynconfig) Serve() error {
	return nil
}

// Stop the dynconfig listening service.
func (d *remoteDynconfig) Stop() error {
	close(d.done)
	if err := os.Remove(d.cachePath); err != nil {
		return err
	}

	return nil
}

// getScheduler returns the scheduler config from manager.
func (d *remoteDynconfig) getScheduler() (*managerv2.Scheduler, error) {
	data, err := d.Get()
	if err != nil {
		return nil, err
	}

	if data.Scheduler == nil {
		return nil, errors.New("invalid scheduler")
	}

	return data.Scheduler, nil
}

// getSeedPeers returns the seed peers config from manager.
func (d *remoteDynconfig) getSeedPeers() ([]*managerv2.SeedPeer, error) {
	scheduler, err := d.getScheduler()
	if err != nil {
		return nil, err
	}

	if len(scheduler.SeedPeers) == 0 {
		return nil, errors.New("seed peer not found ")
	}

	return scheduler.SeedPeers, nil
}

// getSchedulerCluster returns the scheduler cluster config from manager.
func (d *remoteDynconfig) getSchedulerCluster() (*managerv2.SchedulerCluster, error) {
	scheduler, err := d.getScheduler()
	if err != nil {
		return nil, err
	}

	if scheduler.SchedulerCluster == nil {
		return nil, errors.New("invalid scheduler cluster")
	}

	return scheduler.SchedulerCluster, nil
}

// Manager client for dynconfig.
type managerClient struct {
	managerClient managerclient.V2
	config        *Config
}

// New the manager client used by dynconfig.
func newManagerClient(client managerclient.V2, cfg *Config) dc.ManagerClient {
	return &managerClient{
		managerClient: client,
		config:        cfg,
	}
}

func (mc *managerClient) Get() (any, error) {
	getSchedulerResp, err := mc.managerClient.GetScheduler(context.Background(), &managerv2.GetSchedulerRequest{
		SourceType:         managerv2.SourceType_SCHEDULER_SOURCE,
		Hostname:           mc.config.Server.Host,
		Ip:                 mc.config.Server.AdvertiseIP.String(),
		SchedulerClusterId: uint64(mc.config.Manager.SchedulerClusterID),
	})
	if err != nil {
		return nil, err
	}

	listApplicationsResp, err := mc.managerClient.ListApplications(context.Background(), &managerv2.ListApplicationsRequest{
		SourceType: managerv2.SourceType_SCHEDULER_SOURCE,
		Hostname:   mc.config.Server.Host,
		Ip:         mc.config.Server.AdvertiseIP.String(),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// TODO Compatible with old version manager.
			if slices.Contains([]codes.Code{codes.Unimplemented, codes.NotFound}, st.Code()) {
				return DynconfigData{
					Scheduler:    getSchedulerResp,
					Applications: nil,
				}, nil
			}
		}

		return nil, err
	}

	return DynconfigData{
		Scheduler:    getSchedulerResp,
		Applications: listApplicationsResp.Applications,
	}, nil
}
