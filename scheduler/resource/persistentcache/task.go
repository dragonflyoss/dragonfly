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

package persistentcache

import (
	"context"
	"time"

	"github.com/looplab/fsm"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

const (
	// Tiny file size is 128 bytes.
	TinyFileSize = 128

	// Empty file size is 0 bytes.
	EmptyFileSize = 0
)

const (
	// Task has been created but did not start uploading.
	TaskStatePending = "Pending"

	// Task is uploading resources for p2p cluster.
	TaskStateUploading = "Uploading"

	// Task has been uploaded successfully.
	TaskStateSucceeded = "Succeeded"

	// Task has been uploaded failed.
	TaskStateFailed = "Failed"
)

const (
	// Task is uploading.
	TaskEventUpload = "Upload"

	// Task uploaded successfully.
	TaskEventSucceeded = "Succeeded"

	// Task uploaded failed.
	TaskEventFailed = "Failed"
)

// Task contains content for persistent cache task.
type Task struct {
	// ID is task id.
	ID string

	// Replica count of the persistent cache task. The persistent cache task will
	// not be deleted when dfdamon runs garbage collection. It only be deleted
	// when the task is deleted by the user.
	PersistentReplicaCount uint64

	// Tag is used to distinguish different persistent cache tasks.
	Tag string

	// Application of persistent cache task.
	Application string

	// Persistet cache task piece length.
	PieceLength int32

	// ContentLength is persistent cache task total content length.
	ContentLength int64

	// TotalPieceCount is total piece count.
	TotalPieceCount int32

	// Persistent cache task state machine.
	FSM *fsm.FSM

	// TTL is persistent cache task time to live.
	TTL time.Duration

	// CreatedAt is persistent cache task create time.
	CreatedAt time.Time

	// UpdatedAt is persistent cache task update time.
	UpdatedAt time.Time

	// Persistent cache task log.
	Log *logger.SugaredLoggerOnWith
}

// New persistent cache task instance.
func NewTask(id, tag, application, state string, persistentReplicaCount uint64, pieceLength int32,
	contentLength int64, totalPieceCount int32, ttl time.Duration, createdAt, updatedAt time.Time,
	log *logger.SugaredLoggerOnWith) *Task {
	t := &Task{
		ID:                     id,
		PersistentReplicaCount: persistentReplicaCount,
		Tag:                    tag,
		Application:            application,
		ContentLength:          contentLength,
		TotalPieceCount:        totalPieceCount,
		TTL:                    ttl,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
		Log:                    logger.WithTaskID(id),
	}

	// Initialize state machine.
	t.FSM = fsm.NewFSM(
		TaskStatePending,
		fsm.Events{
			fsm.EventDesc{Name: TaskEventUpload, Src: []string{TaskStatePending, TaskStateFailed}, Dst: TaskStateUploading},
			fsm.EventDesc{Name: TaskEventSucceeded, Src: []string{TaskStateUploading}, Dst: TaskStateSucceeded},
			fsm.EventDesc{Name: TaskEventFailed, Src: []string{TaskStateUploading}, Dst: TaskStateFailed},
		},
		fsm.Callbacks{
			TaskEventUpload: func(ctx context.Context, e *fsm.Event) {
				t.Log.Infof("task state is %s", e.FSM.Current())
			},
			TaskEventSucceeded: func(ctx context.Context, e *fsm.Event) {
				t.Log.Infof("task state is %s", e.FSM.Current())
			},
			TaskEventFailed: func(ctx context.Context, e *fsm.Event) {
				t.Log.Infof("task state is %s", e.FSM.Current())
			},
		},
	)
	t.FSM.SetState(state)

	return t
}

// SizeScope return task size scope type.
func (t *Task) SizeScope() commonv2.SizeScope {
	if t.ContentLength < 0 {
		return commonv2.SizeScope_UNKNOW
	}

	if t.TotalPieceCount < 0 {
		return commonv2.SizeScope_UNKNOW
	}

	if t.ContentLength == EmptyFileSize {
		return commonv2.SizeScope_EMPTY
	}

	if t.ContentLength <= TinyFileSize {
		return commonv2.SizeScope_TINY
	}

	if t.TotalPieceCount == 1 {
		return commonv2.SizeScope_SMALL
	}

	return commonv2.SizeScope_NORMAL
}
