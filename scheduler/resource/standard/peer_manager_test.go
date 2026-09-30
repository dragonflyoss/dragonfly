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

package standard

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"

	"d7y.io/dragonfly/v2/pkg/gc"
	"d7y.io/dragonfly/v2/pkg/idgen"
	"d7y.io/dragonfly/v2/scheduler/config"
)

var (
	mockPeerGCConfig = &config.GCConfig{
		PeerGCInterval: 1 * time.Second,
		PeerTTL:        1 * time.Microsecond,
	}
)

func TestPeerManager_newPeerManager(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, peerManager PeerManager, err error)
	}{
		{
			name: "new peer manager",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("peerManager", reflect.TypeOf(peerManager).Elem().Name())
			},
		},
		{
			name: "new peer manager failed because of gc error",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			peerManager, err := newPeerManager(mockPeerGCConfig, gc)
			tc.expect(t, peerManager, err)
		})
	}
}

func TestPeerManager_Load(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, peerManager PeerManager, mockPeer *Peer)
	}{
		{
			name: "load peer",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(mockPeer.ID, peer.ID)
			},
		},
		{
			name: "peer does not exist",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer *Peer) {
				assert := assert.New(t)
				_, loaded := peerManager.Load(mockPeer.ID)
				assert.False(loaded)
			},
		},
		{
			name: "load key is empty",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer *Peer) {
				assert := assert.New(t)
				mockPeer.ID = ""
				peerManager.Store(mockPeer)
				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(mockPeer.ID, peer.ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			mockPeer := NewPeer(mockPeerID, mockTask, mockHost)
			peerManager, err := newPeerManager(mockPeerGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, peerManager, mockPeer)
		})
	}
}

func TestPeerManager_Store(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, peerManager PeerManager, mockPeer *Peer)
	}{
		{
			name: "store peer",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(mockPeer.ID, peer.ID)

				_, loaded = mockPeer.Task.LoadPeer(mockPeer.ID)
				assert.True(loaded)
				_, loaded = mockPeer.Host.LoadPeer(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(int32(1), mockPeer.Host.PeerCount.Load())
			},
		},
		{
			name: "store key is empty",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer *Peer) {
				assert := assert.New(t)
				mockPeer.ID = ""
				peerManager.Store(mockPeer)
				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(mockPeer.ID, peer.ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			mockPeer := NewPeer(mockPeerID, mockTask, mockHost)
			peerManager, err := newPeerManager(mockPeerGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, peerManager, mockPeer)
		})
	}
}

func TestPeerManager_LoadOrStore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, peerManager PeerManager, mockPeer *Peer)
	}{
		{
			name: "load peer exist",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				peer, loaded := peerManager.LoadOrStore(mockPeer)
				assert.True(loaded)
				assert.Equal(mockPeer.ID, peer.ID)
				assert.Equal(int32(1), mockPeer.Host.PeerCount.Load())
			},
		},
		{
			name: "load peer does not exist",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer *Peer) {
				assert := assert.New(t)
				peer, loaded := peerManager.LoadOrStore(mockPeer)
				assert.False(loaded)
				assert.Equal(mockPeer.ID, peer.ID)

				_, loaded = mockPeer.Task.LoadPeer(mockPeer.ID)
				assert.True(loaded)
				_, loaded = mockPeer.Host.LoadPeer(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(int32(1), mockPeer.Host.PeerCount.Load())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			mockPeer := NewPeer(mockPeerID, mockTask, mockHost)
			peerManager, err := newPeerManager(mockPeerGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, peerManager, mockPeer)
		})
	}
}

func TestPeerManager_Delete(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, peerManager PeerManager, mockPeer *Peer)
	}{
		{
			name: "delete peer",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				peerManager.Delete(mockPeer.ID)
				_, loaded := peerManager.Load(mockPeer.ID)
				assert.False(loaded)

				_, loaded = mockPeer.Task.LoadPeer(mockPeer.ID)
				assert.False(loaded)
				_, loaded = mockPeer.Host.LoadPeer(mockPeer.ID)
				assert.False(loaded)
				assert.Equal(int32(0), mockPeer.Host.PeerCount.Load())
			},
		},
		{
			name: "delete key does not exist",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer *Peer) {
				assert := assert.New(t)
				mockPeer.ID = ""
				peerManager.Store(mockPeer)
				peerManager.Delete(mockPeer.ID)
				_, loaded := peerManager.Load(mockPeer.ID)
				assert.False(loaded)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			mockPeer := NewPeer(mockPeerID, mockTask, mockHost)
			peerManager, err := newPeerManager(mockPeerGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, peerManager, mockPeer)
		})
	}
}

func TestPeerManager_DeleteAllByHostID(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, peerManager PeerManager, mockPeer, mockSeedPeer *Peer, hostID string)
	}{
		{
			name: "delete all peers by host id",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer, mockSeedPeer *Peer, hostID string) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				peerManager.Store(mockSeedPeer)
				peerManager.DeleteAllByHostID(hostID)
				_, loadedPeer := peerManager.Load(mockPeer.ID)
				_, loadedSeedPeer := peerManager.Load(mockSeedPeer.ID)
				assert.False(loadedPeer)
				assert.False(loadedSeedPeer)
			},
		},
		{
			name: "delete all peers with non-existent host id",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer, mockSeedPeer *Peer, hostID string) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				peerManager.Store(mockSeedPeer)
				peerManager.DeleteAllByHostID("non-existent-host-id")
				_, loadedPeer := peerManager.Load(mockPeer.ID)
				_, loadedSeedPeer := peerManager.Load(mockSeedPeer.ID)
				assert.True(loadedPeer)
				assert.True(loadedSeedPeer)
			},
		},
		{
			name: "delete all peers keeps peers of other hosts",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockPeer, mockSeedPeer *Peer, hostID string) {
				assert := assert.New(t)
				mockSeedHost := NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Name, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type)
				mockSeedPeer.Host = mockSeedHost
				peerManager.Store(mockPeer)
				peerManager.Store(mockSeedPeer)
				peerManager.DeleteAllByHostID(hostID)
				_, loadedPeer := peerManager.Load(mockPeer.ID)
				_, loadedSeedPeer := peerManager.Load(mockSeedPeer.ID)
				assert.False(loadedPeer)
				assert.True(loadedSeedPeer)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication,
				commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader,
				mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peerManager, err := newPeerManager(mockPeerGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			mockPeer := NewPeer(mockPeerID, mockTask, mockHost)
			mockSeedPeer := NewPeer(mockSeedPeerID, mockTask, mockHost)
			tc.expect(t, peerManager, mockPeer, mockSeedPeer, mockHost.ID)
		})
	}
}

func TestPeerManager_Range(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peerManager PeerManager, mockPeer, mockSeedPeer *Peer)
	}{
		{
			name: "range visits every peer",
			expect: func(t *testing.T, peerManager PeerManager, mockPeer, mockSeedPeer *Peer) {
				assert := assert.New(t)
				var ids []string
				peerManager.Range(func(_, value any) bool {
					ids = append(ids, value.(*Peer).ID)
					return true
				})

				assert.ElementsMatch([]string{mockPeer.ID, mockSeedPeer.ID}, ids)
			},
		},
		{
			name: "range stops when f returns false",
			expect: func(t *testing.T, peerManager PeerManager, mockPeer, mockSeedPeer *Peer) {
				assert := assert.New(t)
				var ids []string
				peerManager.Range(func(_, value any) bool {
					ids = append(ids, value.(*Peer).ID)
					return false
				})

				assert.Len(ids, 1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			gc.EXPECT().Add(gomock.Any()).Return(nil).Times(1)

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peerManager, err := newPeerManager(mockPeerGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			mockPeer := NewPeer(mockPeerID, mockTask, mockHost)
			mockSeedPeer := NewPeer(mockSeedPeerID, mockTask, mockHost)
			peerManager.Store(mockPeer)
			peerManager.Store(mockSeedPeer)
			tc.expect(t, peerManager, mockPeer, mockSeedPeer)
		})
	}
}

func TestPeerManager_RunGC(t *testing.T) {
	tests := []struct {
		name     string
		gcConfig *config.GCConfig
		mock     func(m *gc.MockGCMockRecorder)
		expect   func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer)
	}{
		{
			name: "peer leave",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 5 * time.Minute,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Microsecond,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateSucceeded)
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateLeave, peer.FSM.Current())
			},
		},
		{
			name: "peer reclaimed with disabled shared",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 5 * time.Minute,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Microsecond,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.Host.DisableShared = true
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateLeave, peer.FSM.Current())

				err = peerManager.RunGC(context.Background())
				assert.NoError(err)

				_, loaded = peerManager.Load(mockPeer.ID)
				assert.False(loaded)
			},
		},
		{
			name: "peer download piece timeout and peer state is PeerStateRunning",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 1 * time.Microsecond,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              5 * time.Minute,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateRunning)
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateLeave, peer.FSM.Current())

				err = peerManager.RunGC(context.Background())
				assert.NoError(err)

				_, loaded = peerManager.Load(mockPeer.ID)
				assert.False(loaded)
			},
		},
		{
			name: "peer download piece timeout and peer state is PeerStateBackToSource",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 1 * time.Microsecond,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              5 * time.Minute,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateBackToSource)
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateLeave, peer.FSM.Current())

				err = peerManager.RunGC(context.Background())
				assert.NoError(err)

				_, loaded = peerManager.Load(mockPeer.ID)
				assert.False(loaded)
			},
		},
		{
			name: "peer download piece timeout does not apply to peers that are not downloading",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 1 * time.Microsecond,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              5 * time.Minute,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateSucceeded)
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateSucceeded, peer.FSM.Current())
			},
		},
		{
			name: "peer reclaimed with peer ttl",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 5 * time.Minute,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Microsecond,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateSucceeded)
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateLeave, peer.FSM.Current())

				err = peerManager.RunGC(context.Background())
				assert.NoError(err)

				_, loaded = peerManager.Load(mockPeer.ID)
				assert.False(loaded)
			},
		},
		{
			name: "peer reclaimed with host ttl",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 5 * time.Minute,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              10 * time.Second,
				HostTTL:              1 * time.Microsecond,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateSucceeded)
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateLeave, peer.FSM.Current())

				err = peerManager.RunGC(context.Background())
				assert.NoError(err)

				_, loaded = peerManager.Load(mockPeer.ID)
				assert.False(loaded)
			},
		},
		{
			name: "peer within ttl is kept",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 1 * time.Hour,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Hour,
				HostTTL:              1 * time.Hour,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateRunning)
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateRunning, peer.FSM.Current())
			},
		},
		{
			name: "peer state is PeerStateFailed",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 5 * time.Minute,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Microsecond,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateFailed)
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateLeave, peer.FSM.Current())
			},
		},
		{
			name: "peer state is PeerStateFailed within ttl",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 1 * time.Hour,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Hour,
				HostTTL:              1 * time.Hour,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateFailed)
				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateLeave, peer.FSM.Current())
			},
		},
		{
			name: "peer gets degree failed",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 5 * time.Minute,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Hour,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateSucceeded)
				mockPeer.Task.DeletePeer(mockPeer.ID)

				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				_, loaded := peerManager.Load(mockPeer.ID)
				assert.False(loaded)
			},
		},
		{
			name: "peer reclaimed with PeerCountLimitForTask",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 5 * time.Minute,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Hour,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateSucceeded)
				for range PeerCountLimitForTask + 1 {
					peer := NewPeer(idgen.PeerID(), mockTask, mockHost)
					mockPeer.Task.StorePeer(peer)
				}

				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				_, loaded := peerManager.Load(mockPeer.ID)
				assert.False(loaded)
			},
		},
		{
			name: "peer with edges is kept when task exceeds PeerCountLimitForTask",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 5 * time.Minute,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Hour,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateSucceeded)
				mockChild := NewPeer(idgen.PeerID(), mockTask, mockHost)
				mockTask.StorePeer(mockChild)
				if err := mockTask.AddPeerEdge(mockPeer, mockChild); err != nil {
					t.Fatal(err)
				}

				for range PeerCountLimitForTask {
					mockTask.StorePeer(NewPeer(idgen.PeerID(), mockTask, mockHost))
				}

				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateSucceeded, peer.FSM.Current())
			},
		},
		{
			name: "peer that is not succeeded is kept when task exceeds PeerCountLimitForTask",
			gcConfig: &config.GCConfig{
				PieceDownloadTimeout: 5 * time.Minute,
				PeerGCInterval:       1 * time.Second,
				PeerTTL:              1 * time.Hour,
				HostTTL:              10 * time.Second,
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peerManager PeerManager, mockHost *Host, mockTask *Task, mockPeer *Peer) {
				assert := assert.New(t)
				peerManager.Store(mockPeer)
				mockPeer.FSM.SetState(PeerStateRunning)
				for range PeerCountLimitForTask + 1 {
					mockTask.StorePeer(NewPeer(idgen.PeerID(), mockTask, mockHost))
				}

				err := peerManager.RunGC(context.Background())
				assert.NoError(err)

				peer, loaded := peerManager.Load(mockPeer.ID)
				assert.True(loaded)
				assert.Equal(PeerStateRunning, peer.FSM.Current())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			mockPeer := NewPeer(mockPeerID, mockTask, mockHost)
			peerManager, err := newPeerManager(tc.gcConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, peerManager, mockHost, mockTask, mockPeer)
		})
	}
}
