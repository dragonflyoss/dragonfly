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
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"

	"d7y.io/dragonfly/v2/pkg/rpc"
	"d7y.io/dragonfly/v2/pkg/rpc/manager/client/mocks"
)

func TestNewDynconfig(t *testing.T) {
	assert := assert.New(t)
	viper.Set("dynconfig", filepath.Join(t.TempDir(), "dynconfig.yaml"))
	defer viper.Set("dynconfig", "")

	d, err := NewDynconfig(nil, t.TempDir(), &Config{}, rpc.NewInsecureCredentials())
	assert.NoError(err)
	_, ok := d.(*localDynconfig)
	assert.True(ok)

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockManagerClient := mocks.NewMockV2(ctl)
	mockManagerClient.EXPECT().GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{}, nil).Times(1)
	mockManagerClient.EXPECT().ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)

	d, err = NewDynconfig(mockManagerClient, t.TempDir(), &Config{
		Server: ServerConfig{
			Host: "localhost",
		},
		DynConfig: DynConfig{
			RefreshInterval: 10 * time.Second,
		},
		Manager: ManagerConfig{
			Addr:               "localhost",
			SchedulerClusterID: 1,
		},
	}, rpc.NewInsecureCredentials())
	assert.NoError(err)

	_, ok = d.(*remoteDynconfig)
	assert.True(ok)
}
