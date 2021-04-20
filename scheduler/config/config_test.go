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
	"testing"

	testifyassert "github.com/stretchr/testify/assert"
)

func TestSchedulerConfig_Load(t *testing.T) {
	assert := testifyassert.New(t)

	config := &Config{
		Debug:   true,
		Console: true,
		Verbose: true,
		Scheduler: schedulerConfig{
			ABTest:     false,
			AScheduler: "a-scheduler",
			BScheduler: "b-scheduler",
		},
		Server: serverConfig{
			IP:   "127.0.0.1",
			Port: 8002,
		},
		Worker: schedulerWorkerConfig{
			WorkerNum:         8,
			WorkerJobPoolSize: 10000,
			SenderNum:         10,
			SenderJobPoolSize: 10000,
		},
		CDN: cdnConfig{
			Servers: []CdnServerConfig{
				{
					Name:         "cdn",
					IP:           "127.0.0.1",
					RpcPort:      8003,
					DownloadPort: 8001,
				},
			},
		},
		GC: gcConfig{
			TaskDelay:     3600 * 1000,
			PeerTaskDelay: 3600 * 1000,
		},
	}

	schedulerConfigYAML := &Config{}
	if err := schedulerConfigYAML.Load("./testdata/config/daemon.yaml"); err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(config, schedulerConfigYAML)

	schedulerConfigJSON := &Config{}
	if err := peerHostOptionJSON.Load("./testdata/config/daemon.json"); err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(config, schedulerConfigJSON)
}
