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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name                   string
		id                     string
		tag                    string
		application            string
		state                  string
		persistentReplicaCount uint64
		pieceLength            uint64
		contentLength          uint64
		totalPieceCount        uint32
		ttl                    time.Duration
		createdAt              time.Time
		updatedAt              time.Time
		log                    *logger.SugaredLoggerOnWith
		expectedState          string
	}{
		{
			name:                   "new task with pending state",
			id:                     "task-1",
			tag:                    "tag-1",
			application:            "app-1",
			state:                  TaskStatePending,
			persistentReplicaCount: 3,
			pieceLength:            1024 * 1024,
			contentLength:          1024 * 1024 * 10,
			totalPieceCount:        10,
			ttl:                    time.Hour,
			createdAt:              time.Now(),
			updatedAt:              time.Now(),
			expectedState:          TaskStatePending,
		},
		{
			name:                   "new task with uploading state",
			id:                     "task-2",
			tag:                    "tag-2",
			application:            "app-2",
			state:                  TaskStateUploading,
			persistentReplicaCount: 5,
			pieceLength:            1024 * 1024,
			contentLength:          1024 * 1024 * 20,
			totalPieceCount:        20,
			ttl:                    2 * time.Hour,
			createdAt:              time.Now(),
			updatedAt:              time.Now(),
			expectedState:          TaskStateUploading,
		},
		{
			name:                   "new task with tiny file",
			id:                     "task-3",
			tag:                    "tag-3",
			application:            "app-3",
			state:                  TaskStateSucceeded,
			persistentReplicaCount: 2,
			pieceLength:            128,
			contentLength:          TinyFileSize,
			totalPieceCount:        1,
			ttl:                    30 * time.Minute,
			createdAt:              time.Now(),
			updatedAt:              time.Now(),
			expectedState:          TaskStateSucceeded,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := NewTask(
				tc.id,
				tc.tag,
				tc.application,
				tc.state,
				tc.persistentReplicaCount,
				tc.pieceLength,
				tc.contentLength,
				tc.totalPieceCount,
				tc.ttl,
				tc.createdAt,
				tc.updatedAt,
				tc.log,
			)

			assert.Equal(t, tc.id, task.ID)
			assert.Equal(t, tc.tag, task.Tag)
			assert.Equal(t, tc.application, task.Application)
			assert.Equal(t, tc.persistentReplicaCount, task.PersistentReplicaCount)
			assert.Equal(t, tc.pieceLength, task.PieceLength)
			assert.Equal(t, tc.contentLength, task.ContentLength)
			assert.Equal(t, tc.totalPieceCount, task.TotalPieceCount)
			assert.Equal(t, tc.ttl, task.TTL)
			assert.Equal(t, tc.createdAt, task.CreatedAt)
			assert.Equal(t, tc.updatedAt, task.UpdatedAt)
			assert.Equal(t, tc.expectedState, task.FSM.Current())
			assert.NotNil(t, task.Log)
		})
	}
}

func TestTask_SizeScope(t *testing.T) {
	tests := []struct {
		name              string
		contentLength     uint64
		totalPieceCount   uint32
		expectedSizeScope commonv2.SizeScope
	}{
		{
			name:              "empty file",
			contentLength:     EmptyFileSize,
			totalPieceCount:   0,
			expectedSizeScope: commonv2.SizeScope_EMPTY,
		},
		{
			name:              "tiny file",
			contentLength:     TinyFileSize,
			totalPieceCount:   1,
			expectedSizeScope: commonv2.SizeScope_TINY,
		},
		{
			name:              "small file",
			contentLength:     TinyFileSize + 1,
			totalPieceCount:   1,
			expectedSizeScope: commonv2.SizeScope_SMALL,
		},
		{
			name:              "normal file",
			contentLength:     1024 * 1024,
			totalPieceCount:   10,
			expectedSizeScope: commonv2.SizeScope_NORMAL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := &Task{
				ContentLength:   tc.contentLength,
				TotalPieceCount: tc.totalPieceCount,
			}
			got := task.SizeScope()
			assert.Equal(t, tc.expectedSizeScope, got)
		})
	}
}
