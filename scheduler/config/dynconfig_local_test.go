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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"

	"d7y.io/dragonfly/v2/manager/types"
)

func TestLocalDynconfig_New(t *testing.T) {
	tests := []struct {
		name       string
		configPath func(t *testing.T) string
		expect     func(t *testing.T, configPath string, d DynconfigInterface, err error)
	}{
		{
			name: "new local dynconfig success",
			configPath: func(t *testing.T) string {
				return filepath.Join("testdata", "dynconfig.yaml")
			},
			expect: func(t *testing.T, configPath string, d DynconfigInterface, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.NotNil(d)
			},
		},
		{
			name: "generate local dynconfig with default values when file does not exist",
			configPath: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "dynconfig.yaml")
			},
			expect: func(t *testing.T, configPath string, d DynconfigInterface, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.FileExists(configPath)

				_, err = d.GetApplications()
				assert.EqualError(err, "application not found")

				seedPeerConfig, err := d.GetSeedPeerClusterConfig()
				assert.NoError(err)
				assert.Equal(uint32(DefaultSeedPeerConcurrentUploadLimit), seedPeerConfig.LoadLimit)

				config, err := d.GetSchedulerClusterConfig()
				assert.NoError(err)
				assert.Equal(uint32(DefaultSchedulerCandidateParentLimit), config.CandidateParentLimit)
				assert.Equal(uint32(DefaultSchedulerFilterParentLimit), config.FilterParentLimit)

				clientConfig, err := d.GetSchedulerClusterClientConfig()
				assert.NoError(err)
				assert.Equal(uint32(DefaultPeerConcurrentUploadLimit), clientConfig.LoadLimit)

				ld, ok := d.(*localDynconfig)
				if !ok {
					t.Fatal("invalid local dynconfig type")
				}
				assert.Equal(DefaultLocalDynconfigRefreshInterval, ld.refreshInterval())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configPath := tc.configPath(t)
			d, err := newLocalDynconfig(configPath)
			tc.expect(t, configPath, d, err)
		})
	}
}

func TestLocalDynconfig_GetApplications(t *testing.T) {
	d, err := newLocalDynconfig(filepath.Join("testdata", "dynconfig.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	applications, err := d.GetApplications()
	assert.NoError(err)
	assert.EqualValues(applications, []*managerv2.Application{
		{
			Id:   1,
			Name: "foo",
			Url:  "example.com",
			Bio:  "bar",
			Priority: &managerv2.ApplicationPriority{
				Value: commonv2.Priority_LEVEL1,
				Urls: []*managerv2.URLPriority{
					{
						Regex: "blobs*",
						Value: commonv2.Priority_LEVEL2,
					},
				},
			},
		},
	})
}

func TestLocalDynconfig_GetSeedPeerClusterConfig(t *testing.T) {
	d, err := newLocalDynconfig(filepath.Join("testdata", "dynconfig.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	config, err := d.GetSeedPeerClusterConfig()
	assert.NoError(err)
	assert.Equal(types.SeedPeerClusterConfig{
		LoadLimit: 2000,
	}, config)
}

func TestLocalDynconfig_GetSchedulerClusterConfig(t *testing.T) {
	d, err := newLocalDynconfig(filepath.Join("testdata", "dynconfig.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	config, err := d.GetSchedulerClusterConfig()
	assert.NoError(err)
	assert.Equal(types.SchedulerClusterConfig{
		CandidateParentLimit: 5,
		FilterParentLimit:    10,
		JobRateLimit:         10,
	}, config)
}

func TestLocalDynconfig_GetSchedulerClusterClientConfig(t *testing.T) {
	d, err := newLocalDynconfig(filepath.Join("testdata", "dynconfig.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	config, err := d.GetSchedulerClusterClientConfig()
	assert.NoError(err)
	assert.Equal(types.SchedulerClusterClientConfig{
		LoadLimit: 100,
	}, config)
}

func TestLocalDynconfig_RefreshInterval(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "dynconfig.yaml")
	if err := os.WriteFile(configPath, []byte("schedulerClusterConfig:\n  candidateParentLimit: 5\n"), 0600); err != nil {
		t.Fatal(err)
	}

	d, err := newLocalDynconfig(configPath)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	ld, ok := d.(*localDynconfig)
	if !ok {
		t.Fatal("invalid local dynconfig type")
	}
	assert.Equal(DefaultLocalDynconfigRefreshInterval, ld.refreshInterval())

	if err := os.WriteFile(configPath, []byte("refreshInterval: 10s\nschedulerClusterConfig:\n  candidateParentLimit: 5\n"), 0600); err != nil {
		t.Fatal(err)
	}
	assert.NoError(ld.load())
	assert.Equal(10*time.Second, ld.refreshInterval())
}

func TestLocalDynconfig_Serve(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "dynconfig.yaml")
	if err := os.WriteFile(configPath, []byte("refreshInterval: 100ms\nschedulerClusterConfig:\n  candidateParentLimit: 5\n"), 0600); err != nil {
		t.Fatal(err)
	}

	d, err := newLocalDynconfig(configPath)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.NoError(d.Serve())

	if err := os.WriteFile(configPath, []byte("refreshInterval: 100ms\nschedulerClusterConfig:\n  candidateParentLimit: 10\n"), 0600); err != nil {
		t.Fatal(err)
	}

	assert.Eventually(func() bool {
		config, err := d.GetSchedulerClusterConfig()
		return err == nil && config.CandidateParentLimit == 10
	}, 3*time.Second, 100*time.Millisecond)

	assert.NoError(d.Stop())
}
