/*
 *     Copyright 2024 The Dragonfly Authors
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

package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	resource "d7y.io/dragonfly/v2/scheduler/resource/standard"
)

func TestJob_selectPeers(t *testing.T) {
	tests := []struct {
		name       string
		hosts      []*resource.Host
		percentage uint32
		expect     func(t *testing.T, hosts []*resource.Host, err error)
	}{
		{
			name:       "percentage over 100 is clamped to available peers",
			hosts:      []*resource.Host{{}, {}, {}},
			percentage: 200,
			expect: func(t *testing.T, hosts []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(hosts, 3)
			},
		},
		{
			name:       "valid percentage selects a proportional number of peers",
			hosts:      []*resource.Host{{}, {}, {}, {}},
			percentage: 50,
			expect: func(t *testing.T, hosts []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(hosts, 2)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			res.EXPECT().HostManager().Return(hostManager).Times(1)
			hostManager.EXPECT().LoadAllNormals().Return(tc.hosts).Times(1)

			j := &job{resource: res}
			percentage := tc.percentage
			hosts, err := j.selectPeers(nil, nil, &percentage, logger.WithPeerID("test"))
			tc.expect(t, hosts, err)
		})
	}
}

func TestJob_selectSeedPeers(t *testing.T) {
	tests := []struct {
		name       string
		seedPeers  []*resource.Host
		percentage uint32
		expect     func(t *testing.T, seedPeers []*resource.Host, err error)
	}{
		{
			name:       "percentage over 100 is clamped to available seed peers",
			seedPeers:  []*resource.Host{{}, {}},
			percentage: 200,
			expect: func(t *testing.T, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(seedPeers, 2)
			},
		},
		{
			name:       "valid percentage selects a proportional number of seed peers",
			seedPeers:  []*resource.Host{{}, {}, {}, {}},
			percentage: 50,
			expect: func(t *testing.T, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(seedPeers, 2)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			seedPeer := resource.NewMockSeedPeer(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			res.EXPECT().SeedPeer().Return(seedPeer).Times(1)
			seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
			res.EXPECT().HostManager().Return(hostManager).Times(1)
			hostManager.EXPECT().LoadAllSeeds().Return(tc.seedPeers).Times(1)

			j := &job{resource: res}
			percentage := tc.percentage
			seedPeers, err := j.selectSeedPeers(nil, nil, &percentage, logger.WithPeerID("test"))
			tc.expect(t, seedPeers, err)
		})
	}
}
