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

package persistentcache

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
	"d7y.io/dragonfly/v2/scheduler/config"
)

func TestTaskManager_Load(t *testing.T) {
	type args struct {
		taskID string
	}

	tests := []struct {
		name           string
		args           args
		mockRedis      func(mock redismock.ClientMock)
		expectedTask   *Task
		expectedLoaded bool
	}{
		{
			name: "redis error",
			args: args{
				taskID: "foo",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "foo"),
				).SetErr(errors.New("redis error"))
			},
			expectedTask:   nil,
			expectedLoaded: false,
		},
		{
			name: "empty map from redis (not found)",
			args: args{
				taskID: "notfound",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "notfound"),
				).SetVal(map[string]string{})
			},
			expectedTask:   nil,
			expectedLoaded: false,
		},
		{
			name: "parsing error on persistent_replica_count",
			args: args{
				taskID: "badreplica",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "badreplica"),
				).SetVal(map[string]string{
					"id":                       "badreplica",
					"persistent_replica_count": "not_a_number",
				})
			},
			expectedTask:   nil,
			expectedLoaded: false,
		},
		{
			name: "parsing error on piece_length",
			args: args{
				taskID: "badpiece",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "badpiece"),
				).SetVal(map[string]string{
					"id":           "badpiece",
					"piece_length": "x",
				})
			},
			expectedTask:   nil,
			expectedLoaded: false,
		},
		{
			name: "parsing error on created_at",
			args: args{
				taskID: "badtime",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "badtime"),
				).SetVal(map[string]string{
					"id":         "badtime",
					"created_at": "invalid_time",
				})
			},
			expectedTask:   nil,
			expectedLoaded: false,
		},
		{
			name: "successful load",
			args: args{
				taskID: "goodtask",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":                       "goodtask",
					"tag":                      "tag_value",
					"application":              "app_value",
					"state":                    TaskStateSucceeded,
					"persistent_replica_count": "2",
					"piece_length":             "1024",
					"content_length":           "2048",
					"total_piece_count":        "2",
					"ttl":                      strconv.FormatInt((time.Second * 300).Nanoseconds(), 10),
					"created_at":               time.Now().Format(time.RFC3339),
					"updated_at":               time.Now().Format(time.RFC3339),
				}
				mock.ExpectHGetAll(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "goodtask"),
				).SetVal(mockData)
			},
			expectedTask:   NewTask("goodtask", "tag_value", "app_value", TaskStateSucceeded, 2, 1024, 2048, 2, 5*time.Minute, time.Now(), time.Now(), logger.WithTaskID("goodtask")),
			expectedLoaded: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			tm := &taskManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 42,
					},
				},
				rdb: rdb,
			}

			got, loaded := tm.Load(context.Background(), tt.args.taskID)
			assert.Equal(t, tt.expectedLoaded, loaded)

			if tt.expectedLoaded {
				assert.NotNil(t, got)
				assert.Equal(t, tt.expectedTask.ID, got.ID)
				assert.Equal(t, tt.expectedTask.Tag, got.Tag)
				assert.Equal(t, tt.expectedTask.Application, got.Application)
				assert.Equal(t, tt.expectedTask.PersistentReplicaCount, got.PersistentReplicaCount)
				assert.Equal(t, tt.expectedTask.PieceLength, got.PieceLength)
				assert.Equal(t, tt.expectedTask.ContentLength, got.ContentLength)
				assert.Equal(t, tt.expectedTask.TotalPieceCount, got.TotalPieceCount)
				assert.Equal(t, tt.expectedTask.FSM.Current(), got.FSM.Current())
			} else {
				assert.Nil(t, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestTaskManager_Store(t *testing.T) {
	type args struct {
		task *Task
	}

	tests := []struct {
		name        string
		args        args
		mockRedis   func(mock redismock.ClientMock, task *Task)
		expectedErr bool
	}{
		{
			name: "store success",
			args: args{
				task: NewTask(
					"store-success",
					"test-tag",
					"test-app",
					TaskStateSucceeded,
					1,
					1024,
					2048,
					2,
					5*time.Minute,
					time.Now().Add(-1*time.Minute),
					time.Now(),
					logger.WithTaskID("store-success"),
				),
			},
			mockRedis: func(mock redismock.ClientMock, task *Task) {
				mock.ExpectTxPipeline()
				mock.ExpectHSet(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, task.ID),
					"id", task.ID,
					"persistent_replica_count", task.PersistentReplicaCount,
					"tag", task.Tag,
					"application", task.Application,
					"piece_length", task.PieceLength,
					"content_length", task.ContentLength,
					"total_piece_count", task.TotalPieceCount,
					"state", task.FSM.Current(),
					"ttl", task.TTL.Nanoseconds(),
					"created_at", task.CreatedAt.Format(time.RFC3339),
					"updated_at", task.UpdatedAt.Format(time.RFC3339),
				).SetVal(int64(1))
				mock.ExpectExpire(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, task.ID),
					task.TTL-time.Since(task.CreatedAt),
				).SetVal(true)
				mock.ExpectTxPipelineExec()
			},
			expectedErr: false,
		},
		{
			name: "hset error",
			args: args{
				task: NewTask(
					"store-hset-error",
					"",
					"",
					TaskStatePending,
					0,
					0,
					0,
					0,
					time.Minute,
					time.Now(),
					time.Now(),
					logger.WithTaskID("store-hset-error"),
				),
			},
			mockRedis: func(mock redismock.ClientMock, task *Task) {
				mock.ExpectTxPipeline()
				mock.ExpectHSet(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, task.ID),
					"id", task.ID,
					"persistent_replica_count", task.PersistentReplicaCount,
					"tag", task.Tag,
					"application", task.Application,
					"piece_length", task.PieceLength,
					"content_length", task.ContentLength,
					"total_piece_count", task.TotalPieceCount,
					"state", task.FSM.Current(),
					"ttl", task.TTL.Nanoseconds(),
					"created_at", task.CreatedAt.Format(time.RFC3339),
					"updated_at", task.UpdatedAt.Format(time.RFC3339),
				).SetErr(errors.New("hset error"))
			},
			expectedErr: true,
		},
		{
			name: "expire error",
			args: args{
				task: NewTask(
					"store-expire-error",
					"",
					"",
					TaskStatePending,
					0,
					0,
					0,
					0,
					time.Minute,
					time.Now(),
					time.Now(),
					logger.WithTaskID("store-expire-error"),
				),
			},
			mockRedis: func(mock redismock.ClientMock, task *Task) {
				mock.ExpectTxPipeline()
				mock.ExpectHSet(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, task.ID),
					"id", task.ID,
					"persistent_replica_count", task.PersistentReplicaCount,
					"tag", task.Tag,
					"application", task.Application,
					"piece_length", task.PieceLength,
					"content_length", task.ContentLength,
					"total_piece_count", task.TotalPieceCount,
					"state", task.FSM.Current(),
					"ttl", task.TTL.Nanoseconds(),
					"created_at", task.CreatedAt.Format(time.RFC3339),
					"updated_at", task.UpdatedAt.Format(time.RFC3339),
				).SetVal(int64(1))
				mock.ExpectExpire(
					pkgredis.MakePersistentCacheTaskKeyInScheduler(42, task.ID),
					task.TTL-time.Since(task.CreatedAt),
				).SetErr(errors.New("expire error"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock, tt.args.task)

			tm := &taskManager{
				config: &config.Config{Manager: config.ManagerConfig{SchedulerClusterID: 42}},
				rdb:    rdb,
			}
			err := tm.Store(context.Background(), tt.args.task)
			assert.Equal(t, tt.expectedErr, err != nil, "error mismatch")
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskManager_LoadCorrentReplicaCount(t *testing.T) {
	type args struct {
		taskID string
	}

	tests := []struct {
		name          string
		args          args
		mockRedis     func(mock redismock.ClientMock)
		expectedCount uint64
		expectedErr   bool
	}{
		{
			name: "redis error",
			args: args{
				taskID: "foo",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSCard(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(42, "foo")).SetErr(errors.New("redis error"))
			},
			expectedCount: 0,
			expectedErr:   true,
		},
		{
			name: "successful count",
			args: args{
				taskID: "bar",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSCard(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(42, "bar")).SetVal(5)
			},
			expectedCount: 5,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			tm := &taskManager{
				config: &config.Config{Manager: config.ManagerConfig{SchedulerClusterID: 42}},
				rdb:    rdb,
			}

			cnt, err := tm.LoadCorrentReplicaCount(context.Background(), tt.args.taskID)
			assert.Equal(t, tt.expectedCount, cnt)
			assert.Equal(t, tt.expectedErr, err != nil, "error mismatch")
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskManager_LoadCurrentPersistentReplicaCount(t *testing.T) {
	type args struct {
		taskID string
	}

	tests := []struct {
		name          string
		args          args
		mockRedis     func(mock redismock.ClientMock)
		expectedCount uint64
		expectedErr   bool
	}{
		{
			name: "redis error",
			args: args{
				taskID: "foo",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSCard(pkgredis.MakePersistentPeersOfPersistentCacheTaskInScheduler(42, "foo")).SetErr(errors.New("redis error"))
			},
			expectedCount: 0,
			expectedErr:   true,
		},
		{
			name: "successful count",
			args: args{
				taskID: "bar",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSCard(pkgredis.MakePersistentPeersOfPersistentCacheTaskInScheduler(42, "bar")).SetVal(5)
			},
			expectedCount: 5,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			tm := &taskManager{
				config: &config.Config{Manager: config.ManagerConfig{SchedulerClusterID: 42}},
				rdb:    rdb,
			}

			cnt, err := tm.LoadCurrentPersistentReplicaCount(context.Background(), tt.args.taskID)
			assert.Equal(t, tt.expectedCount, cnt)
			assert.Equal(t, tt.expectedErr, err != nil, "error mismatch")
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskManager_Delete(t *testing.T) {
	type args struct {
		taskID string
	}

	tests := []struct {
		name        string
		args        args
		mockRedis   func(mock redismock.ClientMock)
		expectedErr bool
	}{
		{
			name: "delete success",
			args: args{
				taskID: "delete-success",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectDel(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "delete-success")).SetVal(1)
			},
			expectedErr: false,
		},
		{
			name: "delete error",
			args: args{
				taskID: "delete-error",
			},
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectDel(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "delete-error")).SetErr(errors.New("delete error"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			tm := &taskManager{
				config: &config.Config{Manager: config.ManagerConfig{SchedulerClusterID: 42}},
				rdb:    rdb,
			}

			err := tm.Delete(context.Background(), tt.args.taskID)
			assert.Equal(t, tt.expectedErr, err != nil, "error mismatch")
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskManager_LoadAll(t *testing.T) {
	tests := []struct {
		name        string
		mockRedis   func(mock redismock.ClientMock)
		expectedErr bool
		expectedLen int
	}{
		{
			name: "scan error",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, pkgredis.MakePersistentCacheTasksInScheduler(42), 10).SetErr(errors.New("scan error"))
			},
			expectedErr: true,
			expectedLen: 0,
		},
		{
			name: "load task error",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, pkgredis.MakePersistentCacheTasksInScheduler(42), 10).SetVal([]string{"task1"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "task1")).SetErr(errors.New("load error"))
			},
			expectedErr: false,
			expectedLen: 0,
		},
		{
			name: "successful load all",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, pkgredis.MakePersistentCacheTasksInScheduler(42), 10).SetVal([]string{"task1", "task2"}, 0)
				mockData := map[string]string{
					"id":                       "task1",
					"tag":                      "tag_value",
					"application":              "app_value",
					"state":                    TaskStateSucceeded,
					"persistent_replica_count": "2",
					"piece_length":             "1024",
					"content_length":           "2048",
					"total_piece_count":        "2",
					"ttl":                      strconv.FormatInt((time.Second * 300).Nanoseconds(), 10),
					"created_at":               time.Now().Format(time.RFC3339),
					"updated_at":               time.Now().Format(time.RFC3339),
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "task1")).SetVal(mockData)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "task2")).SetVal(mockData)
			},
			expectedErr: false,
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			tm := &taskManager{
				config: &config.Config{Manager: config.ManagerConfig{SchedulerClusterID: 42}},
				rdb:    rdb,
			}

			tasks, err := tm.LoadAll(context.Background())
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, tasks, tt.expectedLen)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
