/*
 *     Copyright 2025 The Dragonfly Authors
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

package persistent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/bits-and-blooms/bitset"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
	"d7y.io/dragonfly/v2/scheduler/config"
)

var (
	mockPeerManagerConfig = &config.Config{Manager: config.ManagerConfig{SchedulerClusterID: 42}}

	mockPeer = NewPeer("peer1", PeerStateSucceeded, true, bitset.New(2).Set(1), []string{"parent1", "parent2"}, mockTask, &mockRawHost, time.Second, time.Now(), time.Now(), nil, WithConcurrentPieceCount(1))
)

func mockRawPeerFields(peerID string) map[string]string {
	finishedPieces, _ := mockPeer.FinishedPieces.MarshalBinary()
	blockParents, _ := json.Marshal(mockPeer.BlockParents)
	return map[string]string{
		"id":                     peerID,
		"state":                  mockPeer.FSM.Current(),
		"persistent":             strconv.FormatBool(mockPeer.Persistent),
		"concurrent_piece_count": strconv.FormatUint(uint64(mockPeer.ConcurrentPieceCount), 10),
		"finished_pieces":        string(finishedPieces),
		"block_parents":          string(blockParents),
		"task_id":                mockPeer.Task.ID,
		"host_id":                mockPeer.Host.ID,
		"cost":                   strconv.FormatInt(mockPeer.Cost.Nanoseconds(), 10),
		"created_at":             mockPeer.CreatedAt.Format(time.RFC3339),
		"updated_at":             mockPeer.UpdatedAt.Format(time.RFC3339),
	}
}

func mockLoadPeer(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder, peerID string) {
	mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, peerID)).SetVal(mockRawPeerFields(peerID))
	mockHostManager.Load(gomock.Any(), mockPeer.Host.ID).Return(&mockRawHost, true).Times(1)
	mockTaskManager.Load(gomock.Any(), mockPeer.Task.ID).Return(mockTask, true).Times(1)
}

func mockStorePeerScript(mock redismock.ClientMock) *redismock.ExpectedCmd {
	fields := mockRawPeerFields(mockPeer.ID)
	return mock.CustomMatch(matchScriptArgs).ExpectEvalSha("", []string{
		pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID),
		pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockPeer.Task.ID),
		pkgredis.MakePersistentPeersOfPersistentTaskInScheduler(42, mockPeer.Task.ID),
		pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockPeer.Host.ID),
	},
		mockPeer.ID,
		mockPeer.Persistent,
		fields["finished_pieces"],
		mockPeer.FSM.Current(),
		fields["block_parents"],
		mockPeer.Task.ID,
		mockPeer.Host.ID,
		mockPeer.Cost.Nanoseconds(),
		fields["created_at"],
		fields["updated_at"],
		mockPeer.ConcurrentPieceCount,
		nil,
	)
}

func mockDeletePeerScript(mock redismock.ClientMock, peerID string) *redismock.ExpectedCmd {
	return mock.CustomMatch(matchScriptArgs).ExpectEvalSha("", []string{
		pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, peerID),
		pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockPeer.Task.ID),
		pkgredis.MakePersistentPeersOfPersistentTaskInScheduler(42, mockPeer.Task.ID),
		pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockPeer.Host.ID),
	}, peerID, mockPeer.Persistent, mockPeer.Task.ID, mockPeer.Host.ID)
}

func TestPeerManager_Load(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder)
		expect func(t *testing.T, peer *Peer, loaded bool)
	}{
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "peer not found",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(map[string]string{})
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "invalid persistent value",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				fields := mockRawPeerFields(mockPeer.ID)
				fields["persistent"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "invalid finished_pieces value",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				fields := mockRawPeerFields(mockPeer.ID)
				fields["finished_pieces"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "invalid block_parents value",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				fields := mockRawPeerFields(mockPeer.ID)
				fields["block_parents"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "invalid cost value",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				fields := mockRawPeerFields(mockPeer.ID)
				fields["cost"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "invalid created_at value",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				fields := mockRawPeerFields(mockPeer.ID)
				fields["created_at"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "invalid updated_at value",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				fields := mockRawPeerFields(mockPeer.ID)
				fields["updated_at"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "host not found",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(mockRawPeerFields(mockPeer.ID))
				mockHostManager.Load(gomock.Any(), mockPeer.Host.ID).Return(nil, false).Times(1)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "task not found",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(mockRawPeerFields(mockPeer.ID))
				mockHostManager.Load(gomock.Any(), mockPeer.Host.ID).Return(&mockRawHost, true).Times(1)
				mockTaskManager.Load(gomock.Any(), mockPeer.Task.ID).Return(nil, false).Times(1)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "invalid concurrent_piece_count value",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				fields := mockRawPeerFields(mockPeer.ID)
				fields["concurrent_piece_count"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(fields)
				mockHostManager.Load(gomock.Any(), mockPeer.Host.ID).Return(&mockRawHost, true).Times(1)
				mockTaskManager.Load(gomock.Any(), mockPeer.Task.ID).Return(mockTask, true).Times(1)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(peer)
			},
		},
		{
			name: "successful load",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mockLoadPeer(mock, mockHostManager, mockTaskManager, mockPeer.ID)
			},
			expect: func(t *testing.T, peer *Peer, loaded bool) {
				assert := assert.New(t)
				assert.True(loaded)
				assert.Equal(mockPeer.ID, peer.ID)
				assert.Equal(mockPeer.FSM.Current(), peer.FSM.Current())
				assert.Equal(mockPeer.Persistent, peer.Persistent)
				assert.Equal(mockPeer.ConcurrentPieceCount, peer.ConcurrentPieceCount)
				assert.Equal(mockPeer.FinishedPieces, peer.FinishedPieces)
				assert.Equal(mockPeer.BlockParents, peer.BlockParents)
				assert.Equal(mockPeer.Task.ID, peer.Task.ID)
				assert.Equal(mockPeer.Host.ID, peer.Host.ID)
				assert.Equal(mockPeer.Cost, peer.Cost)
				assert.Equal(mockPeer.CreatedAt.Format(time.RFC3339), peer.CreatedAt.Format(time.RFC3339))
				assert.Equal(mockPeer.UpdatedAt.Format(time.RFC3339), peer.UpdatedAt.Format(time.RFC3339))
				assert.NotNil(peer.Log)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)
			tc.mock(mock, hostManager.EXPECT(), taskManager.EXPECT())

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb, hostManager: hostManager, taskManager: taskManager}
			peer, loaded := pm.Load(context.Background(), mockPeer.ID)
			tc.expect(t, peer, loaded)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_Store(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, err error)
	}{
		{
			name: "store succeeds",
			mock: func(mock redismock.ClientMock) {
				mockStorePeerScript(mock).SetVal(true)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock) {
				mockStorePeerScript(mock).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(mock)

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb}
			tc.expect(t, pm.Store(context.Background(), mockPeer))
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_Delete(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder)
		expect func(t *testing.T, err error)
	}{
		{
			name: "peer not found",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, mockPeer.ID)).SetVal(map[string]string{})
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mockLoadPeer(mock, mockHostManager, mockTaskManager, mockPeer.ID)
				mockDeletePeerScript(mock, mockPeer.ID).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "delete succeeds",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mockLoadPeer(mock, mockHostManager, mockTaskManager, mockPeer.ID)
				mockDeletePeerScript(mock, mockPeer.ID).SetVal(true)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)
			tc.mock(mock, hostManager.EXPECT(), taskManager.EXPECT())

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb, hostManager: hostManager, taskManager: taskManager}
			tc.expect(t, pm.Delete(context.Background(), mockPeer.ID))
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_LoadAll(t *testing.T) {
	prefix := fmt.Sprintf("%s:", pkgredis.MakePersistentCachePeersForPersistentTaskInScheduler(42))
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder)
		expect func(t *testing.T, peers []*Peer, err error)
	}{
		{
			name: "redis scan error",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectScan(0, prefix+"*", 10).SetErr(errors.New("redis scan error"))
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(peers)
			},
		},
		{
			name: "invalid peer key is skipped",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix}, 0)
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(peers)
			},
		},
		{
			name: "peer that fails to load is skipped",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix + "peer1"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1")).SetErr(errors.New("redis hgetall error"))
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(peers)
			},
		},
		{
			name: "successful load",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix + "peer1"}, 0)
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer1")
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peers, 1)
				assert.Equal("peer1", peers[0].ID)
			},
		},
		{
			name: "keys spanning multiple scan cursors are all loaded",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix + "peer1"}, 7)
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer1")
				mock.ExpectScan(7, prefix+"*", 10).SetVal([]string{prefix + "peer2"}, 0)
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer2")
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peers, 2)
				assert.Equal("peer1", peers[0].ID)
				assert.Equal("peer2", peers[1].ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)
			tc.mock(mock, hostManager.EXPECT(), taskManager.EXPECT())

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb, hostManager: hostManager, taskManager: taskManager}
			peers, err := pm.LoadAll(context.Background())
			tc.expect(t, peers, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_LoadAllByTaskID(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder)
		expect func(t *testing.T, peers []*Peer, err error)
	}{
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(peers)
			},
		},
		{
			name: "peer that fails to load is skipped",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetVal([]string{"peer1"})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1")).SetErr(errors.New("redis hgetall error"))
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(peers)
			},
		},
		{
			name: "successful load",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetVal([]string{"peer1"})
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer1")
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peers, 1)
				assert.Equal("peer1", peers[0].ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)
			tc.mock(mock, hostManager.EXPECT(), taskManager.EXPECT())

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb, hostManager: hostManager, taskManager: taskManager}
			peers, err := pm.LoadAllByTaskID(context.Background(), mockTask.ID)
			tc.expect(t, peers, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_LoadAllIDsByTaskID(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, peerIDs []string, err error)
	}{
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, peerIDs []string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(peerIDs)
			},
		},
		{
			name: "successful load",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetVal([]string{"peer1", "peer2"})
			},
			expect: func(t *testing.T, peerIDs []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"peer1", "peer2"}, peerIDs)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(mock)

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb}
			peerIDs, err := pm.LoadAllIDsByTaskID(context.Background(), mockTask.ID)
			tc.expect(t, peerIDs, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_LoadPersistentAllByTaskID(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder)
		expect func(t *testing.T, peers []*Peer, err error)
	}{
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(peers)
			},
		},
		{
			name: "peer that fails to load is skipped",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetVal([]string{"peer1"})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1")).SetErr(errors.New("redis hgetall error"))
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(peers)
			},
		},
		{
			name: "successful load",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetVal([]string{"peer1"})
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer1")
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peers, 1)
				assert.Equal("peer1", peers[0].ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)
			tc.mock(mock, hostManager.EXPECT(), taskManager.EXPECT())

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb, hostManager: hostManager, taskManager: taskManager}
			peers, err := pm.LoadPersistentAllByTaskID(context.Background(), mockTask.ID)
			tc.expect(t, peers, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_DeleteAllByTaskID(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder)
		expect func(t *testing.T, err error)
	}{
		{
			name: "load peer ids error",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "deletes every peer of the task",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetVal([]string{"peer1", "peer2"})
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer1")
				mockDeletePeerScript(mock, "peer1").SetVal(true)
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer2")
				mockDeletePeerScript(mock, "peer2").SetVal(true)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "peer that fails to delete is skipped",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, mockTask.ID)).SetVal([]string{"peer1", "peer2"})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1")).SetErr(errors.New("redis hgetall error"))
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer2")
				mockDeletePeerScript(mock, "peer2").SetVal(true)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)
			tc.mock(mock, hostManager.EXPECT(), taskManager.EXPECT())

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb, hostManager: hostManager, taskManager: taskManager}
			tc.expect(t, pm.DeleteAllByTaskID(context.Background(), mockTask.ID))
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_LoadAllByHostID(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder)
		expect func(t *testing.T, peers []*Peer, err error)
	}{
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockRawHost.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(peers)
			},
		},
		{
			name: "peer that fails to load is skipped",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockRawHost.ID)).SetVal([]string{"peer1"})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1")).SetErr(errors.New("redis hgetall error"))
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(peers)
			},
		},
		{
			name: "successful load",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockRawHost.ID)).SetVal([]string{"peer1"})
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer1")
			},
			expect: func(t *testing.T, peers []*Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peers, 1)
				assert.Equal("peer1", peers[0].ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)
			tc.mock(mock, hostManager.EXPECT(), taskManager.EXPECT())

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb, hostManager: hostManager, taskManager: taskManager}
			peers, err := pm.LoadAllByHostID(context.Background(), mockRawHost.ID)
			tc.expect(t, peers, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_LoadAllIDsByHostID(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, peerIDs []string, err error)
	}{
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockRawHost.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, peerIDs []string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(peerIDs)
			},
		},
		{
			name: "successful load",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockRawHost.ID)).SetVal([]string{"peer1", "peer2"})
			},
			expect: func(t *testing.T, peerIDs []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"peer1", "peer2"}, peerIDs)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(mock)

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb}
			peerIDs, err := pm.LoadAllIDsByHostID(context.Background(), mockRawHost.ID)
			tc.expect(t, peerIDs, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestPeerManager_DeleteAllByHostID(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder)
		expect func(t *testing.T, err error)
	}{
		{
			name: "load peer ids error",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockRawHost.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "deletes every peer of the host",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockRawHost.ID)).SetVal([]string{"peer1", "peer2"})
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer1")
				mockDeletePeerScript(mock, "peer1").SetVal(true)
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer2")
				mockDeletePeerScript(mock, "peer2").SetVal(true)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "peer that fails to delete is skipped",
			mock: func(mock redismock.ClientMock, mockHostManager *MockHostManagerMockRecorder, mockTaskManager *MockTaskManagerMockRecorder) {
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, mockRawHost.ID)).SetVal([]string{"peer1", "peer2"})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1")).SetErr(errors.New("redis hgetall error"))
				mockLoadPeer(mock, mockHostManager, mockTaskManager, "peer2")
				mockDeletePeerScript(mock, "peer2").SetVal(true)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)
			tc.mock(mock, hostManager.EXPECT(), taskManager.EXPECT())

			pm := &peerManager{config: mockPeerManagerConfig, rdb: rdb, hostManager: hostManager, taskManager: taskManager}
			tc.expect(t, pm.DeleteAllByHostID(context.Background(), mockRawHost.ID))
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}
