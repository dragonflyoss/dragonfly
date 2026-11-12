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
	"testing"
	"time"

	"github.com/looplab/fsm"
	"github.com/stretchr/testify/assert"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		expect func(t *testing.T, task *Task)
	}{
		{
			name:  "new task with pending state",
			state: TaskStatePending,
			expect: func(t *testing.T, task *Task) {
				assert := assert.New(t)
				assert.Equal("task-1", task.ID)
				assert.Equal("url", task.URL)
				assert.Equal("region", task.ObjectStorageRegion)
				assert.Equal("endpoint", task.ObjectStorageEndpoint)
				assert.Equal(uint64(3), task.PersistentReplicaCount)
				assert.Equal(uint64(1024*1024*10), task.ContentLength)
				assert.Equal(uint32(10), task.TotalPieceCount)
				assert.Equal(time.Hour, task.TTL)
				assert.Equal(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), task.CreatedAt)
				assert.Equal(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), task.UpdatedAt)
				assert.Equal(TaskStatePending, task.FSM.Current())
				assert.NotNil(task.Log)
			},
		},
		{
			name:  "new task with uploading state",
			state: TaskStateUploading,
			expect: func(t *testing.T, task *Task) {
				assert := assert.New(t)
				assert.Equal("task-1", task.ID)
				assert.Equal("url", task.URL)
				assert.Equal("region", task.ObjectStorageRegion)
				assert.Equal("endpoint", task.ObjectStorageEndpoint)
				assert.Equal(uint64(3), task.PersistentReplicaCount)
				assert.Equal(uint64(1024*1024*10), task.ContentLength)
				assert.Equal(uint32(10), task.TotalPieceCount)
				assert.Equal(time.Hour, task.TTL)
				assert.Equal(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), task.CreatedAt)
				assert.Equal(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), task.UpdatedAt)
				assert.Equal(TaskStateUploading, task.FSM.Current())
				assert.NotNil(task.Log)
			},
		},
		{
			name:  "new task with succeeded state",
			state: TaskStateSucceeded,
			expect: func(t *testing.T, task *Task) {
				assert := assert.New(t)
				assert.Equal("task-1", task.ID)
				assert.Equal("url", task.URL)
				assert.Equal("region", task.ObjectStorageRegion)
				assert.Equal("endpoint", task.ObjectStorageEndpoint)
				assert.Equal(uint64(3), task.PersistentReplicaCount)
				assert.Equal(uint64(1024*1024*10), task.ContentLength)
				assert.Equal(uint32(10), task.TotalPieceCount)
				assert.Equal(time.Hour, task.TTL)
				assert.Equal(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), task.CreatedAt)
				assert.Equal(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), task.UpdatedAt)
				assert.Equal(TaskStateSucceeded, task.FSM.Current())
				assert.NotNil(task.Log)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, NewTask(
				"task-1",
				"url",
				"region",
				"endpoint",
				tc.state,
				3,
				1024*1024*10,
				10,
				time.Hour,
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				nil,
			))
		})
	}
}

func TestTask_FSM(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		event  string
		expect func(t *testing.T, task *Task, err error)
	}{
		{
			name:  "upload from pending",
			state: TaskStatePending,
			event: TaskEventUpload,
			expect: func(t *testing.T, task *Task, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(TaskStateUploading, task.FSM.Current())
			},
		},
		{
			name:  "upload from failed",
			state: TaskStateFailed,
			event: TaskEventUpload,
			expect: func(t *testing.T, task *Task, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(TaskStateUploading, task.FSM.Current())
			},
		},
		{
			name:  "upload from uploading is rejected",
			state: TaskStateUploading,
			event: TaskEventUpload,
			expect: func(t *testing.T, task *Task, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(TaskStateUploading, task.FSM.Current())
			},
		},
		{
			name:  "upload from succeeded is rejected",
			state: TaskStateSucceeded,
			event: TaskEventUpload,
			expect: func(t *testing.T, task *Task, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(TaskStateSucceeded, task.FSM.Current())
			},
		},
		{
			name:  "succeeded from uploading",
			state: TaskStateUploading,
			event: TaskEventSucceeded,
			expect: func(t *testing.T, task *Task, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(TaskStateSucceeded, task.FSM.Current())
			},
		},
		{
			name:  "succeeded from pending is rejected",
			state: TaskStatePending,
			event: TaskEventSucceeded,
			expect: func(t *testing.T, task *Task, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(TaskStatePending, task.FSM.Current())
			},
		},
		{
			name:  "failed from uploading",
			state: TaskStateUploading,
			event: TaskEventFailed,
			expect: func(t *testing.T, task *Task, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(TaskStateFailed, task.FSM.Current())
			},
		},
		{
			name:  "failed from pending is rejected",
			state: TaskStatePending,
			event: TaskEventFailed,
			expect: func(t *testing.T, task *Task, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(TaskStatePending, task.FSM.Current())
			},
		},
		{
			name:  "failed from succeeded is rejected",
			state: TaskStateSucceeded,
			event: TaskEventFailed,
			expect: func(t *testing.T, task *Task, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(TaskStateSucceeded, task.FSM.Current())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := NewTask("task-1", "url", "region", "endpoint", tc.state, 1, 1024, 1, time.Hour, time.Now(), time.Now(), nil)
			err := task.FSM.Event(context.Background(), tc.event)
			tc.expect(t, task, err)
		})
	}
}

func TestTask_SizeScope(t *testing.T) {
	tests := []struct {
		name            string
		contentLength   uint64
		totalPieceCount uint32
		expect          func(t *testing.T, sizeScope commonv2.SizeScope)
	}{
		{
			name:            "empty file",
			contentLength:   EmptyFileSize,
			totalPieceCount: 0,
			expect: func(t *testing.T, sizeScope commonv2.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv2.SizeScope_EMPTY, sizeScope)
			},
		},
		{
			name:            "empty file ignores piece count",
			contentLength:   EmptyFileSize,
			totalPieceCount: 5,
			expect: func(t *testing.T, sizeScope commonv2.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv2.SizeScope_EMPTY, sizeScope)
			},
		},
		{
			name:            "tiny file",
			contentLength:   TinyFileSize,
			totalPieceCount: 1,
			expect: func(t *testing.T, sizeScope commonv2.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv2.SizeScope_TINY, sizeScope)
			},
		},
		{
			name:            "tiny file ignores piece count",
			contentLength:   TinyFileSize,
			totalPieceCount: 3,
			expect: func(t *testing.T, sizeScope commonv2.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv2.SizeScope_TINY, sizeScope)
			},
		},
		{
			name:            "small file",
			contentLength:   TinyFileSize + 1,
			totalPieceCount: 1,
			expect: func(t *testing.T, sizeScope commonv2.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv2.SizeScope_SMALL, sizeScope)
			},
		},
		{
			name:            "normal file",
			contentLength:   1024 * 1024,
			totalPieceCount: 10,
			expect: func(t *testing.T, sizeScope commonv2.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv2.SizeScope_NORMAL, sizeScope)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := &Task{ContentLength: tc.contentLength, TotalPieceCount: tc.totalPieceCount}
			tc.expect(t, task.SizeScope())
		})
	}
}
