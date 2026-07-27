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

//go:generate mockgen -destination mocks/dynconfig_mock.go -source dynconfig.go -package mocks

package config

import (
	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"

	"d7y.io/dragonfly/v2/manager/types"
	managerclient "d7y.io/dragonfly/v2/pkg/rpc/manager/client"
)

// DynconfigInterface is the interface for dynconfig, which provides methods to get the dynamic configuration.
type DynconfigInterface interface {
	// GetApplications returns the applications config.
	GetApplications() ([]*managerv2.Application, error)

	// GetSeedPeerClusterConfig returns the seed peer cluster config.
	GetSeedPeerClusterConfig() (types.SeedPeerClusterConfig, error)

	// GetSchedulerClusterConfig returns the scheduler cluster config.
	GetSchedulerClusterConfig() (types.SchedulerClusterConfig, error)

	// GetSchedulerClusterClientConfig returns the client config.
	GetSchedulerClusterClientConfig() (types.SchedulerClusterClientConfig, error)
}

// NewDynconfig returns a new dynconfig instance. If the manager client is nil, it returns the local
// dynconfig, which loads the dynamic configuration from the local file specified by the dynconfig
// flag. Otherwise, it returns the remote dynconfig, which fetches the dynamic configuration from
// the manager.
func NewDynconfig(client managerclient.V2, cfg *Config) (DynconfigInterface, error) {
	if client == nil {
		return newLocalDynconfig(cfg)
	}

	return newRemoteDynconfig(client, cfg)
}
