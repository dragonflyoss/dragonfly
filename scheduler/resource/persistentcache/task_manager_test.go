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
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
	"d7y.io/dragonfly/v2/scheduler/config"
)

const mockInvalidFieldValue = "invalid"

var (
	mockTaskManagerConfig = &config.Config{Manager: config.ManagerConfig{SchedulerClusterID: 42}}

	mockTask = NewTask("task1", "tag", "application", TaskStateSucceeded, 2, 1024, 2048, 2, 5*time.Minute, time.Now().Add(-time.Minute), time.Now(), logger.WithTaskID("task1"))
)

func mockRawTaskFields(taskID string) map[string]string {
	return map[string]string{
		"id":                       taskID,
		"tag":                      mockTask.Tag,
		"application":              mockTask.Application,
		"state":                    mockTask.FSM.Current(),
		"persistent_replica_count": strconv.FormatUint(mockTask.PersistentReplicaCount, 10),
		"piece_length":             strconv.FormatUint(mockTask.PieceLength, 10),
		"content_length":           strconv.FormatUint(mockTask.ContentLength, 10),
		"total_piece_count":        strconv.FormatUint(uint64(mockTask.TotalPieceCount), 10),
		"ttl":                      strconv.FormatInt(mockTask.TTL.Nanoseconds(), 10),
		"created_at":               mockTask.CreatedAt.Format(time.RFC3339),
		"updated_at":               mockTask.UpdatedAt.Format(time.RFC3339),
	}
}

func matchScriptArgs(expected, actual []any) error {
	for i := range expected {
		if i == 1 || expected[i] == nil {
			continue
		}

		if !reflect.DeepEqual(expected[i], actual[i]) {
			return fmt.Errorf("script arg %d: expected %v, got %v", i, expected[i], actual[i])
		}
	}

	return nil
}

func mockStoreTaskScript(mock redismock.ClientMock) *redismock.ExpectedCmd {
	return mock.CustomMatch(matchScriptArgs).ExpectEvalSha("", []string{pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)},
		mockTask.ID,
		mockTask.PersistentReplicaCount,
		mockTask.Tag,
		mockTask.Application,
		mockTask.PieceLength,
		mockTask.ContentLength,
		mockTask.TotalPieceCount,
		mockTask.FSM.Current(),
		mockTask.CreatedAt.Format(time.RFC3339),
		mockTask.UpdatedAt.Format(time.RFC3339),
		mockTask.TTL.Nanoseconds(),
		nil,
	)
}

func TestTaskManager_Load(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, task *Task, loaded bool)
	}{
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(task)
			},
		},
		{
			name: "task not found",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(map[string]string{})
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(task)
			},
		},
		{
			name: "invalid persistent_replica_count value",
			mock: func(mock redismock.ClientMock) {
				fields := mockRawTaskFields(mockTask.ID)
				fields["persistent_replica_count"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(task)
			},
		},
		{
			name: "invalid piece_length value",
			mock: func(mock redismock.ClientMock) {
				fields := mockRawTaskFields(mockTask.ID)
				fields["piece_length"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(task)
			},
		},
		{
			name: "invalid content_length value",
			mock: func(mock redismock.ClientMock) {
				fields := mockRawTaskFields(mockTask.ID)
				fields["content_length"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(task)
			},
		},
		{
			name: "invalid total_piece_count value",
			mock: func(mock redismock.ClientMock) {
				fields := mockRawTaskFields(mockTask.ID)
				fields["total_piece_count"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(task)
			},
		},
		{
			name: "invalid ttl value",
			mock: func(mock redismock.ClientMock) {
				fields := mockRawTaskFields(mockTask.ID)
				fields["ttl"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(task)
			},
		},
		{
			name: "invalid created_at value",
			mock: func(mock redismock.ClientMock) {
				fields := mockRawTaskFields(mockTask.ID)
				fields["created_at"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(task)
			},
		},
		{
			name: "invalid updated_at value",
			mock: func(mock redismock.ClientMock) {
				fields := mockRawTaskFields(mockTask.ID)
				fields["updated_at"] = mockInvalidFieldValue
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(fields)
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.False(loaded)
				assert.Nil(task)
			},
		},
		{
			name: "successful load",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(mockRawTaskFields(mockTask.ID))
			},
			expect: func(t *testing.T, task *Task, loaded bool) {
				assert := assert.New(t)
				assert.True(loaded)
				assert.Equal(mockTask.ID, task.ID)
				assert.Equal(mockTask.Tag, task.Tag)
				assert.Equal(mockTask.Application, task.Application)
				assert.Equal(mockTask.PersistentReplicaCount, task.PersistentReplicaCount)
				assert.Equal(mockTask.PieceLength, task.PieceLength)
				assert.Equal(mockTask.ContentLength, task.ContentLength)
				assert.Equal(mockTask.TotalPieceCount, task.TotalPieceCount)
				assert.Equal(mockTask.TTL, task.TTL)
				assert.Equal(mockTask.FSM.Current(), task.FSM.Current())
				assert.Equal(mockTask.CreatedAt.Format(time.RFC3339), task.CreatedAt.Format(time.RFC3339))
				assert.Equal(mockTask.UpdatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339))
				assert.NotNil(task.Log)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(mock)

			tm := &taskManager{config: mockTaskManagerConfig, rdb: rdb}
			task, loaded := tm.Load(context.Background(), mockTask.ID)
			tc.expect(t, task, loaded)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestTaskManager_LoadCurrentReplicaCount(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, count uint64, err error)
	}{
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectSCard(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(42, mockTask.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, count uint64, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Zero(count)
			},
		},
		{
			name: "successful count",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectSCard(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(42, mockTask.ID)).SetVal(5)
			},
			expect: func(t *testing.T, count uint64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint64(5), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(mock)

			tm := &taskManager{config: mockTaskManagerConfig, rdb: rdb}
			count, err := tm.LoadCurrentReplicaCount(context.Background(), mockTask.ID)
			tc.expect(t, count, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestTaskManager_LoadCurrentPersistentReplicaCount(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, count uint64, err error)
	}{
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectSCard(pkgredis.MakePersistentPeersOfPersistentCacheTaskInScheduler(42, mockTask.ID)).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, count uint64, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Zero(count)
			},
		},
		{
			name: "successful count",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectSCard(pkgredis.MakePersistentPeersOfPersistentCacheTaskInScheduler(42, mockTask.ID)).SetVal(5)
			},
			expect: func(t *testing.T, count uint64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint64(5), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(mock)

			tm := &taskManager{config: mockTaskManagerConfig, rdb: rdb}
			count, err := tm.LoadCurrentPersistentReplicaCount(context.Background(), mockTask.ID)
			tc.expect(t, count, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestTaskManager_Store(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, err error)
	}{
		{
			name: "store succeeds",
			mock: func(mock redismock.ClientMock) {
				mockStoreTaskScript(mock).SetVal(true)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock) {
				mockStoreTaskScript(mock).SetErr(errors.New("redis error"))
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

			tm := &taskManager{config: mockTaskManagerConfig, rdb: rdb}
			tc.expect(t, tm.Store(context.Background(), mockTask))
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestTaskManager_Delete(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, err error)
	}{
		{
			name: "delete succeeds",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectDel(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetVal(1)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectDel(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, mockTask.ID)).SetErr(errors.New("delete error"))
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

			tm := &taskManager{config: mockTaskManagerConfig, rdb: rdb}
			tc.expect(t, tm.Delete(context.Background(), mockTask.ID))
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestTaskManager_LoadAll(t *testing.T) {
	prefix := fmt.Sprintf("%s:", pkgredis.MakePersistentCacheTasksInScheduler(42))
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, tasks []*Task, err error)
	}{
		{
			name: "scan error",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetErr(errors.New("scan error"))
			},
			expect: func(t *testing.T, tasks []*Task, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(tasks)
			},
		},
		{
			name: "invalid task key is skipped",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix}, 0)
			},
			expect: func(t *testing.T, tasks []*Task, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(tasks)
			},
		},
		{
			name: "task that fails to load is skipped",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix + "task1"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "task1")).SetErr(errors.New("load error"))
			},
			expect: func(t *testing.T, tasks []*Task, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(tasks)
			},
		},
		{
			name: "task keys with peer set suffixes resolve to one task",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix + "task1:persistent-cache-peers", prefix + "task1:persistent-peers", prefix + "task1"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "task1")).SetVal(mockRawTaskFields("task1"))
			},
			expect: func(t *testing.T, tasks []*Task, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(tasks, 1)
				assert.Equal("task1", tasks[0].ID)
			},
		},
		{
			name: "successful load all",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix + "task1", prefix + "task2"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "task1")).SetVal(mockRawTaskFields("task1"))
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "task2")).SetVal(mockRawTaskFields("task2"))
			},
			expect: func(t *testing.T, tasks []*Task, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(tasks, 2)
				assert.Equal("task1", tasks[0].ID)
				assert.Equal("task2", tasks[1].ID)
			},
		},
		{
			name: "keys spanning multiple scan cursors are all loaded",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix + "task1"}, 7)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "task1")).SetVal(mockRawTaskFields("task1"))
				mock.ExpectScan(7, prefix+"*", 10).SetVal([]string{prefix + "task2"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(42, "task2")).SetVal(mockRawTaskFields("task2"))
			},
			expect: func(t *testing.T, tasks []*Task, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(tasks, 2)
				assert.Equal("task1", tasks[0].ID)
				assert.Equal("task2", tasks[1].ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(mock)

			tm := &taskManager{config: mockTaskManagerConfig, rdb: rdb}
			tasks, err := tm.LoadAll(context.Background())
			tc.expect(t, tasks, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}
