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
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	dfdaemonv2 "d7y.io/api/v2/pkg/apis/dfdaemon/v2"
	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"
	logger "d7y.io/dragonfly/v2/internal/dflog"
	dfdaemonclientmocks "d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/client/mocks"
)

func mockSeedHost(port int32) *Host {
	return NewHost(
		mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Name, mockRawSeedHost.Hostname,
		port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type)
}

func mockSeedHostServer(t *testing.T) *Host {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	svr := grpc.NewServer()
	healthpb.RegisterHealthServer(svr, health.NewServer())
	go func() {
		if err := svr.Serve(lis); err != nil {
			logger.Errorf("failed to serve the health service: %v", err)
		}
	}()
	t.Cleanup(svr.Stop)

	return mockSeedHost(int32(lis.Addr().(*net.TCPAddr).Port))
}

func mockUnreachableSeedHost(t *testing.T) *Host {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	port := int32(lis.Addr().(*net.TCPAddr).Port)
	if err := lis.Close(); err != nil {
		t.Fatal(err)
	}

	return mockSeedHost(port)
}

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
		hosts  func(t *testing.T) []*Host
		mock   func(m *MockHostManagerMockRecorder, hosts []*Host)
		expect func(t *testing.T, seedPeer *seedPeer, hosts []*Host)
	}{
		{
			name: "refresh healthy seed peer",
			hosts: func(t *testing.T) []*Host {
				return []*Host{mockSeedHostServer(t)}
			},
			mock: func(m *MockHostManagerMockRecorder, hosts []*Host) {
				m.LoadAllSeeds().Return(hosts)
			},
			expect: func(t *testing.T, seedPeer *seedPeer, hosts []*Host) {
				assert := assert.New(t)
				seedPeer.refresh(context.Background())
				assert.True(seedPeer.HasAvailable())

				host, err := seedPeer.Select(context.Background(), mockTaskID)
				assert.NoError(err)
				assert.Equal(hosts[0].ID, host.ID)
			},
		},
		{
			name: "filter unhealthy seed peer",
			hosts: func(t *testing.T) []*Host {
				return []*Host{mockUnreachableSeedHost(t)}
			},
			mock: func(m *MockHostManagerMockRecorder, hosts []*Host) {
				m.LoadAllSeeds().Return(hosts)
			},
			expect: func(t *testing.T, seedPeer *seedPeer, hosts []*Host) {
				assert := assert.New(t)
				seedPeer.refresh(context.Background())
				assert.False(seedPeer.HasAvailable())

				_, err := seedPeer.Select(context.Background(), mockTaskID)
				assert.EqualError(err, "no available seed peer")
			},
		},
		{
			name: "clear stale seed peers when host manager is empty",
			hosts: func(t *testing.T) []*Host {
				return []*Host{}
			},
			mock: func(m *MockHostManagerMockRecorder, hosts []*Host) {
				m.LoadAllSeeds().Return(hosts)
			},
			expect: func(t *testing.T, seedPeer *seedPeer, hosts []*Host) {
				assert := assert.New(t)
				mockAddr := "127.0.0.1:4000"
				seedPeer.hosts.Store(mockAddr, &Host{})
				seedPeer.hashring.Add(mockAddr)

				seedPeer.refresh(context.Background())
				_, loaded := seedPeer.hosts.Load(mockAddr)
				assert.False(loaded)
				assert.False(seedPeer.HasAvailable())

				_, err := seedPeer.Select(context.Background(), mockTaskID)
				assert.EqualError(err, "no available seed peer")
			},
		},
		{
			name: "refresh and select concurrently",
			hosts: func(t *testing.T) []*Host {
				return []*Host{}
			},
			mock: func(m *MockHostManagerMockRecorder, hosts []*Host) {
				m.LoadAllSeeds().Return(hosts).AnyTimes()
			},
			expect: func(t *testing.T, seedPeer *seedPeer, hosts []*Host) {
				var wg sync.WaitGroup
				wg.Add(2)
				go func() {
					defer wg.Done()
					for range 100 {
						seedPeer.refresh(context.Background())
					}
				}()

				go func() {
					defer wg.Done()
					for range 100 {
						seedPeer.HasAvailable()
						if _, err := seedPeer.Select(context.Background(), mockTaskID); err == nil {
							assert.Fail(t, "select should fail when no seed peer is available")
						}
					}
				}()

				wg.Wait()
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

			hosts := tc.hosts(t)
			tc.mock(hostManager.EXPECT(), hosts)

			seedPeer := newSeedPeer(peerManager, hostManager, clientPool,
				grpc.WithTransportCredentials(insecure.NewCredentials())).(*seedPeer)
			tc.expect(t, seedPeer, hosts)
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
				assert.EqualError(err, "no available seed peer")
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
				assert.EqualError(err, "no available seed peer")
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
