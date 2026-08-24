/*
 *     Copyright 2022 The Dragonfly Authors
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

package standard

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	dfdaemonv2 "d7y.io/api/v2/pkg/apis/dfdaemon/v2"
	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"
	dfdaemonclientmocks "d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/client/mocks"
)

func TestSeedPeer_newSeedPeer(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, s SeedPeer)
	}{
		{
			name: "new seed peer",
			expect: func(t *testing.T, s SeedPeer) {
				assert := assert.New(t)
				assert.Equal(reflect.TypeOf(s).Elem().Name(), "seedPeer")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			hostManager := NewMockHostManager(ctl)
			peerManager := NewMockPeerManager(ctl)
			clientPool := dfdaemonclientmocks.NewMockPool(ctl)

			tc.expect(t, newSeedPeer(peerManager, hostManager, clientPool))
		})
	}
}

func TestSeedPeer_refresh(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T)
	}{
		{
			name: "clears cached seed peers when host manager is empty",
			expect: func(t *testing.T) {
				// Given a seed peer cache containing a host that is no longer in HostManager.
				ctl := gomock.NewController(t)
				hostManager := NewMockHostManager(ctl)
				peerManager := NewMockPeerManager(ctl)
				clientPool := dfdaemonclientmocks.NewMockPool(ctl)
				hostManager.EXPECT().LoadAllSeeds().Return([]*Host{})

				seedPeer := newSeedPeer(peerManager, hostManager, clientPool).(*seedPeer)
				const staleAddress = "127.0.0.1:4000"
				seedPeer.hosts.Store(staleAddress, &Host{})
				seedPeer.hashring.Add(staleAddress)

				// When the seed peer cache is refreshed from the empty HostManager.
				seedPeer.refresh(context.Background())

				// Then the stale seed peer is no longer available or selectable.
				seedPeer.snapshotMutex.RLock()
				hosts := seedPeer.hosts
				seedPeer.snapshotMutex.RUnlock()
				_, found := hosts.Load(staleAddress)
				assert.False(t, found)
				assert.False(t, seedPeer.HasAvailable())
				_, err := seedPeer.Select(context.Background(), mockTaskID)
				assert.EqualError(t, err, "no seed peer available")
			},
		},
		{
			name: "is safe during concurrent selection",
			expect: func(t *testing.T) {
				// Given a seed peer cache that repeatedly refreshes to an empty snapshot.
				ctl := gomock.NewController(t)
				hostManager := NewMockHostManager(ctl)
				peerManager := NewMockPeerManager(ctl)
				clientPool := dfdaemonclientmocks.NewMockPool(ctl)
				hostManager.EXPECT().LoadAllSeeds().Return([]*Host{}).AnyTimes()
				seedPeer := newSeedPeer(peerManager, hostManager, clientPool).(*seedPeer)

				// When refresh and selection run concurrently.
				var waitGroup sync.WaitGroup
				waitGroup.Add(2)
				go func() {
					defer waitGroup.Done()
					for range 100 {
						seedPeer.refresh(context.Background())
					}
				}()
				go func() {
					defer waitGroup.Done()
					for range 100 {
						seedPeer.HasAvailable()
						_, _ = seedPeer.Select(context.Background(), mockTaskID)
					}
				}()

				// Then all concurrent operations complete without a data race.
				waitGroup.Wait()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t)
		})
	}
}

func TestSeedPeer_TriggerDownloadTask(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, err error)
	}{
		{
			name: "trigger download task failed",
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "no seed peer available")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			hostManager := NewMockHostManager(ctl)
			peerManager := NewMockPeerManager(ctl)
			clientPool := dfdaemonclientmocks.NewMockPool(ctl)

			seedPeer := newSeedPeer(peerManager, hostManager, clientPool)
			tc.expect(t, seedPeer.TriggerDownloadTask(context.Background(), mockTaskID, &dfdaemonv2.DownloadTaskRequest{}))
		})
	}
}

func TestSeedPeer_TriggerTask(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, result *schedulerv1.PeerResult, err error)
	}{
		{
			name: "start obtain seed stream failed",
			expect: func(t *testing.T, peer *Peer, result *schedulerv1.PeerResult, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "no seed peer available")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			hostManager := NewMockHostManager(ctl)
			peerManager := NewMockPeerManager(ctl)
			clientPool := dfdaemonclientmocks.NewMockPool(ctl)

			seedPeer := newSeedPeer(peerManager, hostManager, clientPool)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer, result, err := seedPeer.TriggerTask(context.Background(), nil, mockTask)
			tc.expect(t, peer, result, err)
		})
	}
}
