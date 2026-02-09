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
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/bits-and-blooms/bitset"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
	"d7y.io/dragonfly/v2/scheduler/config"
)

func TestPeerManager_Load(t *testing.T) {
	type args struct {
		peerID string
	}

	tests := []struct {
		name           string
		args           args
		mock           func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder)
		mockRedis      func(mock redismock.ClientMock)
		expectedPeer   *Peer
		expectedLoaded bool
	}{
		{
			name: "redis error",
			args: args{
				peerID: "foo",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "foo"),
				).SetErr(errors.New("redis error"))
			},
			expectedPeer:   nil,
			expectedLoaded: false,
		},
		{
			name: "host not found",
			args: args{
				peerID: "nohost",
			},
			mock: func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder) {
				mockHostManager.Load(gomock.Any(), gomock.Any()).Return(nil, false).Times(1)
			},
			mockRedis: func(mock redismock.ClientMock) {
				finishedPieces, err := bitset.New(2).Set(1).MarshalBinary()
				if err != nil {
					t.Fatalf("failed to marshal bitset: %v", err)
				}

				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "nohost"),
				).SetVal(map[string]string{
					"id":                     "nohost",
					"state":                  PeerStateSucceeded,
					"persistent":             "true",
					"concurrent_piece_count": "1",
					"finished_pieces":        string(finishedPieces),
					"block_parents":          `["parent1", "parent2"]`,
					"task_id":                "task1",
					"host_id":                "host1",
					"cost":                   strconv.FormatUint(uint64(time.Second.Nanoseconds()), 10),
					"created_at":             time.Now().Format(time.RFC3339),
					"updated_at":             time.Now().Format(time.RFC3339),
				})
			},
			expectedPeer:   nil,
			expectedLoaded: false,
		},
		{
			name: "task not found",
			args: args{
				peerID: "notask",
			},
			mock: func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder) {
				mockHostManager.Load(gomock.Any(), gomock.Any()).Return(&mockRawHost, true).Times(1)
				mockTaskManager.Load(gomock.Any(), gomock.Any()).Return(nil, false).Times(1)
			},
			mockRedis: func(mock redismock.ClientMock) {
				finishedPieces, err := bitset.New(2).Set(1).MarshalBinary()
				if err != nil {
					t.Fatalf("failed to marshal bitset: %v", err)
				}

				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "notask"),
				).SetVal(map[string]string{
					"id":                     "notask",
					"state":                  PeerStateSucceeded,
					"persistent":             "true",
					"concurrent_piece_count": "1",
					"finished_pieces":        string(finishedPieces),
					"block_parents":          `["parent1", "parent2"]`,
					"task_id":                "task1",
					"host_id":                "host1",
					"cost":                   strconv.FormatUint(uint64(time.Second.Nanoseconds()), 10),
					"created_at":             time.Now().Format(time.RFC3339),
					"updated_at":             time.Now().Format(time.RFC3339),
				})
			},
			expectedPeer:   nil,
			expectedLoaded: false,
		},
		{
			name: "successful load",
			args: args{
				peerID: "goodpeer",
			},
			mockRedis: func(mock redismock.ClientMock) {
				finishedPieces, err := bitset.New(2).Set(1).MarshalBinary()
				if err != nil {
					t.Fatalf("failed to marshal bitset: %v", err)
				}

				mockData := map[string]string{
					"id":                     "goodpeer",
					"state":                  PeerStateSucceeded,
					"persistent":             "true",
					"concurrent_piece_count": "1",
					"finished_pieces":        string(finishedPieces),
					"block_parents":          `["parent1", "parent2"]`,
					"task_id":                "task1",
					"host_id":                "127.0.0.1-foo",
					"cost":                   strconv.FormatUint(uint64(time.Second.Nanoseconds()), 10),
					"created_at":             time.Now().Format(time.RFC3339),
					"updated_at":             time.Now().Format(time.RFC3339),
				}
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "goodpeer"),
				).SetVal(mockData)
			},
			mock: func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder) {
				mockHostManager.Load(gomock.Any(), gomock.Any()).Return(&mockRawHost, true).Times(1)
				mockTaskManager.Load(gomock.Any(), gomock.Any()).Return(NewTask(
					"task1",
					"url",
					"region",
					"endpoint",
					TaskStateSucceeded,
					1,
					2048,
					2,
					5*time.Minute,
					time.Now().Add(-1*time.Minute),
					time.Now(),
					logger.WithTaskID("store-success"),
				), true).Times(1)
			},
			expectedPeer: NewPeer(
				"goodpeer",
				PeerStateSucceeded,
				true,
				bitset.New(2).Set(1),
				[]string{"parent1", "parent2"},
				&Task{ID: "task1"},
				&Host{ID: "127.0.0.1-foo"},
				time.Second,
				time.Now(),
				time.Now(),
				logger.WithPeer("host1", "task1", "goodpeer"),
				WithConcurrentPieceCount(1),
			),
			expectedLoaded: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)

			if tt.mock != nil {
				tt.mock(hostManager, hostManager.EXPECT(), taskManager, taskManager.EXPECT())
			}

			pm := &peerManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb:         rdb,
				hostManager: hostManager,
				taskManager: taskManager,
			}

			got, loaded := pm.Load(context.Background(), tt.args.peerID)
			assert.Equal(t, tt.expectedLoaded, loaded)

			if tt.expectedLoaded {
				assert.NotNil(t, got)
				assert.Equal(t, tt.expectedPeer.ID, got.ID)
				assert.Equal(t, tt.expectedPeer.FSM.Current(), got.FSM.Current())
				assert.Equal(t, tt.expectedPeer.Persistent, got.Persistent)
				assert.Equal(t, tt.expectedPeer.ConcurrentPieceCount, got.ConcurrentPieceCount)
				assert.Equal(t, tt.expectedPeer.FinishedPieces, got.FinishedPieces)
				assert.Equal(t, tt.expectedPeer.BlockParents, got.BlockParents)
				assert.Equal(t, tt.expectedPeer.Task.ID, got.Task.ID)
				assert.Equal(t, tt.expectedPeer.Host.ID, got.Host.ID)
				assert.Equal(t, tt.expectedPeer.Cost, got.Cost)
				assert.Equal(t, tt.expectedPeer.CreatedAt.Format(time.RFC3339), got.CreatedAt.Format(time.RFC3339))
				assert.Equal(t, tt.expectedPeer.UpdatedAt.Format(time.RFC3339), got.UpdatedAt.Format(time.RFC3339))
			} else {
				assert.Nil(t, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestPeerManager_LoadAll(t *testing.T) {
	tests := []struct {
		name          string
		mockRedis     func(mock redismock.ClientMock)
		mock          func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder)
		expectedPeers []*Peer
		expectedErr   bool
	}{
		{
			name: "redis scan error",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, fmt.Sprintf("%s:*", pkgredis.MakePersistentCachePeersForPersistentTaskInScheduler(42)), 10).SetErr(errors.New("redis scan error"))
			},
			expectedPeers: nil,
			expectedErr:   true,
		},
		{
			name: "invalid peer key",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, fmt.Sprintf("%s:*", pkgredis.MakePersistentCachePeersForPersistentTaskInScheduler(42)), 10).SetVal([]string{fmt.Sprintf("%s:", pkgredis.MakePersistentCachePeersForPersistentTaskInScheduler(42))}, 0)
			},
			expectedPeers: nil,
			expectedErr:   false,
		},
		{
			name: "load peer error",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, fmt.Sprintf("%s:*", pkgredis.MakePersistentCachePeersForPersistentTaskInScheduler(42)), 10).SetVal([]string{fmt.Sprintf("%s:peer1", pkgredis.MakePersistentCachePeersForPersistentTaskInScheduler(42))}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1")).SetErr(errors.New("redis hgetall error"))
			},
			expectedPeers: nil,
			expectedErr:   false,
		},
		{
			name: "successful load",
			mockRedis: func(mock redismock.ClientMock) {
				finishedPieces, err := bitset.New(2).Set(1).MarshalBinary()
				if err != nil {
					t.Fatalf("failed to marshal bitset: %v", err)
				}

				mock.ExpectScan(0, fmt.Sprintf("%s:*", pkgredis.MakePersistentCachePeersForPersistentTaskInScheduler(42)), 10).SetVal([]string{fmt.Sprintf("%s:peer1", pkgredis.MakePersistentCachePeersForPersistentTaskInScheduler(42))}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1")).SetVal(map[string]string{
					"id":                     "peer1",
					"state":                  PeerStateSucceeded,
					"persistent":             "true",
					"concurrent_piece_count": "1",
					"finished_pieces":        string(finishedPieces),
					"block_parents":          `["parent1", "parent2"]`,
					"task_id":                "task1",
					"host_id":                "host1",
					"cost":                   strconv.FormatUint(uint64(time.Second.Nanoseconds()), 10),
					"created_at":             time.Now().Format(time.RFC3339),
					"updated_at":             time.Now().Format(time.RFC3339),
				})
			},
			mock: func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder) {
				mockHostManager.Load(gomock.Any(), gomock.Any()).Return(&mockRawHost, true).Times(1)
				mockTaskManager.Load(gomock.Any(), gomock.Any()).Return(NewTask(
					"task1",
					"url",
					"region",
					"endpoint",
					TaskStateSucceeded,
					1,
					2048,
					2,
					5*time.Minute,
					time.Now().Add(-1*time.Minute),
					time.Now(),
					logger.WithTaskID("store-success"),
				), true).Times(1)
			},
			expectedPeers: []*Peer{
				NewPeer(
					"peer1",
					PeerStateSucceeded,
					true,
					bitset.New(2).Set(1),
					[]string{"parent1", "parent2"},
					&Task{ID: "task1"},
					&Host{ID: "host1"},
					time.Second,
					time.Now(),
					time.Now(),
					logger.WithPeer("host1", "task1", "peer1"),
					WithConcurrentPieceCount(1),
				),
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)

			if tt.mock != nil {
				tt.mock(hostManager, hostManager.EXPECT(), taskManager, taskManager.EXPECT())
			}

			pm := &peerManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb:         rdb,
				hostManager: hostManager,
				taskManager: taskManager,
			}

			got, err := pm.LoadAll(context.Background())
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedPeers), len(got))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestPeerManager_LoadAllByTaskID(t *testing.T) {
	type args struct {
		taskID string
	}

	tests := []struct {
		name          string
		args          args
		mock          func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder)
		mockRedis     func(mock redismock.ClientMock)
		expectedPeers []*Peer
		expectedErr   bool
	}{
		{
			name: "redis error",
			args: args{
				taskID: "task1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, "task1"),
				).SetErr(errors.New("redis error"))
			},
			expectedPeers: nil,
			expectedErr:   true,
		},
		{
			name: "load peer error",
			args: args{
				taskID: "task1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, "task1"),
				).SetVal([]string{"peer1"})
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1"),
				).SetErr(errors.New("redis hgetall error"))
			},
			expectedPeers: nil,
			expectedErr:   false,
		},
		{
			name: "successful load",
			args: args{
				taskID: "task1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				finishedPieces, err := bitset.New(2).Set(1).MarshalBinary()
				if err != nil {
					t.Fatalf("failed to marshal bitset: %v", err)
				}

				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, "task1"),
				).SetVal([]string{"peer1"})
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1"),
				).SetVal(map[string]string{
					"id":                     "peer1",
					"state":                  PeerStateSucceeded,
					"persistent":             "true",
					"concurrent_piece_count": "1",
					"finished_pieces":        string(finishedPieces),
					"block_parents":          `["parent1", "parent2"]`,
					"task_id":                "task1",
					"host_id":                "host1",
					"cost":                   strconv.FormatUint(uint64(time.Second.Nanoseconds()), 10),
					"created_at":             time.Now().Format(time.RFC3339),
					"updated_at":             time.Now().Format(time.RFC3339),
				})
			},
			mock: func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder) {
				mockHostManager.Load(gomock.Any(), gomock.Any()).Return(&mockRawHost, true).Times(1)
				mockTaskManager.Load(gomock.Any(), gomock.Any()).Return(NewTask(
					"task1",
					"url",
					"region",
					"endpoint",
					TaskStateSucceeded,
					1,
					2048,
					2,
					5*time.Minute,
					time.Now().Add(-1*time.Minute),
					time.Now(),
					logger.WithTaskID("store-success"),
				), true).Times(1)
			},
			expectedPeers: []*Peer{
				NewPeer(
					"peer1",
					PeerStateSucceeded,
					true,
					bitset.New(2).Set(1),
					[]string{"parent1", "parent2"},
					&Task{ID: "task1"},
					&Host{ID: "host1"},
					time.Second,
					time.Now(),
					time.Now(),
					logger.WithPeer("host1", "task1", "peer1"),
					WithConcurrentPieceCount(1),
				),
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)

			if tt.mock != nil {
				tt.mock(hostManager, hostManager.EXPECT(), taskManager, taskManager.EXPECT())
			}

			pm := &peerManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb:         rdb,
				hostManager: hostManager,
				taskManager: taskManager,
			}

			got, err := pm.LoadAllByTaskID(context.Background(), tt.args.taskID)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedPeers), len(got))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestPeerManager_LoadAllIDsByTaskID(t *testing.T) {
	type args struct {
		taskID string
	}

	tests := []struct {
		name        string
		args        args
		mockRedis   func(mock redismock.ClientMock)
		expectedIDs []string
		expectedErr bool
	}{
		{
			name: "redis error",
			args: args{
				taskID: "task1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, "task1"),
				).SetErr(errors.New("redis error"))
			},
			expectedIDs: nil,
			expectedErr: true,
		},
		{
			name: "successful load",
			args: args{
				taskID: "task1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, "task1"),
				).SetVal([]string{"peer1", "peer2"})
			},
			expectedIDs: []string{"peer1", "peer2"},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)

			pm := &peerManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb:         rdb,
				hostManager: hostManager,
				taskManager: taskManager,
			}

			got, err := pm.LoadAllIDsByTaskID(context.Background(), tt.args.taskID)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIDs, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestPeerManager_LoadPersistentAllByTaskID(t *testing.T) {
	type args struct {
		taskID string
	}

	tests := []struct {
		name          string
		args          args
		mock          func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder)
		mockRedis     func(mock redismock.ClientMock)
		expectedPeers []*Peer
		expectedErr   bool
	}{
		{
			name: "redis error",
			args: args{
				taskID: "task1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, "task1"),
				).SetErr(errors.New("redis error"))
			},
			expectedPeers: nil,
			expectedErr:   true,
		},
		{
			name: "load peer error",
			args: args{
				taskID: "task1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, "task1"),
				).SetVal([]string{"peer1"})
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1"),
				).SetErr(errors.New("redis hgetall error"))
			},
			expectedPeers: nil,
			expectedErr:   false,
		},
		{
			name: "successful load",
			args: args{
				taskID: "task1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				finishedPieces, err := bitset.New(2).Set(1).MarshalBinary()
				if err != nil {
					t.Fatalf("failed to marshal bitset: %v", err)
				}

				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, "task1"),
				).SetVal([]string{"peer1"})
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1"),
				).SetVal(map[string]string{
					"id":                     "peer1",
					"state":                  PeerStateSucceeded,
					"persistent":             "true",
					"concurrent_piece_count": "1",
					"finished_pieces":        string(finishedPieces),
					"block_parents":          `["parent1", "parent2"]`,
					"task_id":                "task1",
					"host_id":                "host1",
					"cost":                   strconv.FormatUint(uint64(time.Second.Nanoseconds()), 10),
					"created_at":             time.Now().Format(time.RFC3339),
					"updated_at":             time.Now().Format(time.RFC3339),
				})
			},
			mock: func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder) {
				mockHostManager.Load(gomock.Any(), gomock.Any()).Return(&mockRawHost, true).Times(1)
				mockTaskManager.Load(gomock.Any(), gomock.Any()).Return(NewTask(
					"task1",
					"url",
					"region",
					"endpoint",
					TaskStateSucceeded,
					1,
					2048,
					2,
					5*time.Minute,
					time.Now().Add(-1*time.Minute),
					time.Now(),
					logger.WithTaskID("store-success"),
				), true).Times(1)
			},
			expectedPeers: []*Peer{
				NewPeer(
					"peer1",
					PeerStateSucceeded,
					true,
					bitset.New(2).Set(1),
					[]string{"parent1", "parent2"},
					&Task{ID: "task1"},
					&Host{ID: "host1"},
					time.Second,
					time.Now(),
					time.Now(),
					logger.WithPeer("host1", "task1", "peer1"),
					WithConcurrentPieceCount(1),
				),
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)

			if tt.mock != nil {
				tt.mock(hostManager, hostManager.EXPECT(), taskManager, taskManager.EXPECT())
			}

			pm := &peerManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb:         rdb,
				hostManager: hostManager,
				taskManager: taskManager,
			}

			got, err := pm.LoadPersistentAllByTaskID(context.Background(), tt.args.taskID)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedPeers), len(got))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestPeerManager_DeleteAllByTaskID(t *testing.T) {
	type args struct {
		taskID string
	}

	tests := []struct {
		name        string
		args        args
		mock        func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder)
		mockRedis   func(mock redismock.ClientMock)
		expectedErr bool
	}{
		{
			name: "load peers error",
			args: args{
				taskID: "task1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentTaskInScheduler(42, "task1"),
				).SetErr(errors.New("redis error"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)

			if tt.mock != nil {
				tt.mock(hostManager, hostManager.EXPECT(), taskManager, taskManager.EXPECT())
			}

			pm := &peerManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb:         rdb,
				hostManager: hostManager,
				taskManager: taskManager,
			}

			err := pm.DeleteAllByTaskID(context.Background(), tt.args.taskID)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestPeerManager_LoadAllByHostID(t *testing.T) {
	type args struct {
		hostID string
	}

	tests := []struct {
		name          string
		args          args
		mock          func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder)
		mockRedis     func(mock redismock.ClientMock)
		expectedPeers []*Peer
		expectedErr   bool
	}{
		{
			name: "redis error",
			args: args{
				hostID: "host1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, "host1"),
				).SetErr(errors.New("redis error"))
			},
			expectedPeers: nil,
			expectedErr:   true,
		},
		{
			name: "load peer error",
			args: args{
				hostID: "host1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, "host1"),
				).SetVal([]string{"peer1"})
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1"),
				).SetErr(errors.New("redis hgetall error"))
			},
			expectedPeers: nil,
			expectedErr:   false,
		},
		{
			name: "successful load",
			args: args{
				hostID: "host1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				finishedPieces, err := bitset.New(2).Set(1).MarshalBinary()
				if err != nil {
					t.Fatalf("failed to marshal bitset: %v", err)
				}

				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, "host1"),
				).SetVal([]string{"peer1"})
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCachePeerKeyForPersistentTaskInScheduler(42, "peer1"),
				).SetVal(map[string]string{
					"id":                     "peer1",
					"state":                  PeerStateSucceeded,
					"persistent":             "true",
					"concurrent_piece_count": "1",
					"finished_pieces":        string(finishedPieces),
					"block_parents":          `["parent1", "parent2"]`,
					"task_id":                "task1",
					"host_id":                "host1",
					"cost":                   strconv.FormatUint(uint64(time.Second.Nanoseconds()), 10),
					"created_at":             time.Now().Format(time.RFC3339),
					"updated_at":             time.Now().Format(time.RFC3339),
				})
			},
			mock: func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder) {
				mockHostManager.Load(gomock.Any(), gomock.Any()).Return(&mockRawHost, true).Times(1)
				mockTaskManager.Load(gomock.Any(), gomock.Any()).Return(NewTask(
					"task1",
					"url",
					"region",
					"endpoint",
					TaskStateSucceeded,
					1,
					2048,
					2,
					5*time.Minute,
					time.Now().Add(-1*time.Minute),
					time.Now(),
					logger.WithTaskID("store-success"),
				), true).Times(1)
			},
			expectedPeers: []*Peer{
				NewPeer(
					"peer1",
					PeerStateSucceeded,
					true,
					bitset.New(2).Set(1),
					[]string{"parent1", "parent2"},
					&Task{ID: "task1"},
					&Host{ID: "host1"},
					time.Second,
					time.Now(),
					time.Now(),
					logger.WithPeer("host1", "task1", "peer1"),
					WithConcurrentPieceCount(1),
				),
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)

			if tt.mock != nil {
				tt.mock(hostManager, hostManager.EXPECT(), taskManager, taskManager.EXPECT())
			}

			pm := &peerManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb:         rdb,
				hostManager: hostManager,
				taskManager: taskManager,
			}

			got, err := pm.LoadAllByHostID(context.Background(), tt.args.hostID)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedPeers), len(got))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestPeerManager_LoadAllIDsByHostID(t *testing.T) {
	type args struct {
		hostID string
	}

	tests := []struct {
		name        string
		args        args
		mockRedis   func(mock redismock.ClientMock)
		expectedIDs []string
		expectedErr bool
	}{
		{
			name: "redis error",
			args: args{
				hostID: "host1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, "host1"),
				).SetErr(errors.New("redis error"))
			},
			expectedIDs: nil,
			expectedErr: true,
		},
		{
			name: "successful load",
			args: args{
				hostID: "host1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, "host1"),
				).SetVal([]string{"peer1", "peer2"})
			},
			expectedIDs: []string{"peer1", "peer2"},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)

			pm := &peerManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb:         rdb,
				hostManager: hostManager,
				taskManager: taskManager,
			}

			got, err := pm.LoadAllIDsByHostID(context.Background(), tt.args.hostID)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIDs, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestPeerManager_DeleteAllByHostID(t *testing.T) {
	type args struct {
		hostID string
	}

	tests := []struct {
		name        string
		args        args
		mock        func(hostManager *MockHostManager, mockHostManager *MockHostManagerMockRecorder, taskManager *MockTaskManager, mockTaskManager *MockTaskManagerMockRecorder)
		mockRedis   func(mock redismock.ClientMock)
		expectedErr bool
	}{
		{
			name: "load peers error",
			args: args{
				hostID: "host1",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(
					pkgredis.MakePersistentCachePeersOfPersistentHostInScheduler(42, "host1"),
				).SetErr(errors.New("redis error"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			hostManager := NewMockHostManager(ctrl)
			taskManager := NewMockTaskManager(ctrl)

			if tt.mock != nil {
				tt.mock(hostManager, hostManager.EXPECT(), taskManager, taskManager.EXPECT())
			}

			pm := &peerManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb:         rdb,
				hostManager: hostManager,
				taskManager: taskManager,
			}

			err := pm.DeleteAllByHostID(context.Background(), tt.args.hostID)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}
