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

	"d7y.io/dragonfly/v2/manager/types"
)

func TestLocalDynconfig_New(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		expect     func(t *testing.T, d DynconfigInterface, err error)
	}{
		{
			name:       "new local dynconfig success",
			configPath: filepath.Join("testdata", "dynconfig.yaml"),
			expect: func(t *testing.T, d DynconfigInterface, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.NotNil(d)
			},
		},
		{
			name:       "local dynconfig file not found",
			configPath: filepath.Join("testdata", "foo.yaml"),
			expect: func(t *testing.T, d DynconfigInterface, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(d)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, err := newLocalDynconfig(tc.configPath)
			tc.expect(t, d, err)
		})
	}
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

func TestLocalDynconfig_GetWithoutManager(t *testing.T) {
	d, err := newLocalDynconfig(filepath.Join("testdata", "dynconfig.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	_, err = d.GetScheduler()
	assert.EqualError(err, "manager is not configured")

	_, err = d.GetApplications()
	assert.EqualError(err, "manager is not configured")

	_, err = d.GetSeedPeers()
	assert.EqualError(err, "manager is not configured")

	_, err = d.GetSeedPeerClusterConfig()
	assert.EqualError(err, "manager is not configured")

	_, err = d.GetSchedulerCluster()
	assert.EqualError(err, "manager is not configured")

	_, err = d.GetSchedulerClusterClientConfig()
	assert.EqualError(err, "manager is not configured")

	_, err = d.Get()
	assert.EqualError(err, "manager is not configured")
}

func TestLocalDynconfig_RefreshInterval(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), localDynconfigFileName)
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
	configPath := filepath.Join(t.TempDir(), localDynconfigFileName)
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
