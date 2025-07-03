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

package gc

import (
	"context"
	"time"

	"gorm.io/gorm"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/models"
	pkggc "d7y.io/dragonfly/v2/pkg/gc"
)

const (
	// DefaultSchedulerGCBatchSize is the default batch size for deleting schedulers.
	DefaultSchedulerGCBatchSize = 5000

	// DefaultSchedulerTTL is the default TTL for scheduler.
	DefaultSchedulerGCTTL = time.Minute * 30

	// DefaultSchedulerGCInterval is the default interval for running scheduler GC.
	DefaultSchedulerGCInterval = time.Hour * 1

	// DefaultSchedulerGCTimeout is the default timeout for running scheduler GC.
	DefaultSchedulerGCTimeout = time.Hour * 1

	// SchedulerGCTaskID is the ID of the scheduler GC task.
	SchedulerGCTaskID = "scheduler"
)

// NewSchedulerGCTask returns a new scheduler GC task.
func NewSchedulerGCTask(db *gorm.DB) pkggc.Task {
	return pkggc.Task{
		ID:       SchedulerGCTaskID,
		Interval: DefaultSchedulerGCInterval,
		Timeout:  DefaultSchedulerGCTimeout,
		Runner:   &scheduler{db: db, recorder: newJobRecorder(db)},
	}
}

// scheduler is the struct for cleaning up inactive schedulers which implements the gc Runner interface.
type scheduler struct {
	db       *gorm.DB
	recorder *jobRecorder
}

// RunGC implements the gc Runner interface.
func (s *scheduler) RunGC(ctx context.Context) error {
	args := models.JSONMap{
		"type":       SchedulerGCTaskID,
		"ttl":        DefaultSchedulerGCTTL,
		"batch_size": DefaultSchedulerGCBatchSize,
	}

	var userID uint
	if id, ok := ctx.Value(pkggc.ContextKeyUserID).(uint); ok {
		userID = id
	}

	var taskID string
	if id, ok := ctx.Value(pkggc.ContextKeyTaskID).(string); ok {
		taskID = id
	} else {
		// Use the default task ID if taskID is not provided. (applied to background periodic execution scenarios)
		taskID = SchedulerGCTaskID
	}

	if err := s.recorder.Init(userID, taskID, args); err != nil {
		return err
	}

	var gcResult Result
	defer func() {
		if err := s.recorder.Record(gcResult); err != nil {
			logger.Errorf("failed to record scheduler GC result: %v", err)
		}
	}()

	for {
		result := s.db.Where("updated_at < ?", time.Now().Add(-DefaultSchedulerGCTTL)).Where("state = ?", models.SchedulerStateInactive).Limit(DefaultSchedulerGCBatchSize).Unscoped().Delete(&models.Scheduler{})
		if result.Error != nil {
			gcResult.Error = result.Error
			return result.Error
		}

		if result.RowsAffected == 0 {
			break
		}

		gcResult.Purged += result.RowsAffected
		logger.Infof("gc scheduler deleted %d inactive schedulers", result.RowsAffected)
	}

	return nil
}
