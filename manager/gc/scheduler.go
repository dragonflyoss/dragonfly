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

	// DefaultSchedulerKeepAliveTimeout is the fallback keepalive timeout used when
	// a scheduler does not report manager_keep_alive_interval in its config
	// (e.g., legacy v1 registrations).
	DefaultSchedulerKeepAliveTimeout = time.Minute * 30

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
	if err := s.mark(ctx); err != nil {
		// Only log error if failed in mark phase,
		// because this can wait to be executed in the next gc window.
		logger.Errorf("failed to mark inactive schedulers: %v", err)
	}

	return s.sweep(ctx)
}

// mark running the mark operation for marking inactive schedulers.
func (s *scheduler) mark(ctx context.Context) error {
	logger.Info("running scheduler GC mark")

	var lastID uint
	for {
		var schedulers []*models.Scheduler
		if err := s.db.WithContext(ctx).Model(&models.Scheduler{}).
			Where("state = ? AND id > ?", models.SchedulerStateActive, lastID).
			Order("id ASC").
			Limit(DefaultSchedulerGCBatchSize).
			Find(&schedulers).Error; err != nil {
			return err
		}
		if len(schedulers) == 0 {
			break
		}
		lastID = schedulers[len(schedulers)-1].ID

		now := time.Now()
		schedulerIDs := make([]uint, 0, len(schedulers))
		for _, scheduler := range schedulers {
			timeout := DefaultSchedulerKeepAliveTimeout
			if scheduler.Config != nil {
				// Retrieve the keep alive interval from the scheduler's configuration.
				if keepAliveInterval, ok := scheduler.Config["manager_keep_alive_interval"].(float64); ok && keepAliveInterval > 0 {
					timeout = time.Duration(keepAliveInterval) * 3
				}
			}

			// Fall back to UpdatedAt for legacy rows where LastKeepAliveAt was never populated.
			lastKeepAliveAt := scheduler.LastKeepAliveAt
			if lastKeepAliveAt.IsZero() {
				lastKeepAliveAt = scheduler.UpdatedAt
			}

			// Check whether the last keep alive time exceeds the timeout,
			// indicating that the scheduler is inactive.
			if now.Sub(lastKeepAliveAt) > timeout {
				schedulerIDs = append(schedulerIDs, scheduler.ID)
			}
		}

		if len(schedulerIDs) > 0 {
			if err := s.db.WithContext(ctx).Model(&models.Scheduler{}).
				Where("id IN (?)", schedulerIDs).
				Update("state", models.SchedulerStateInactive).
				Error; err != nil {
				return err
			}
		}

		logger.Infof("scheduler GC marks %d schedulers to inactive", len(schedulerIDs))

		// If this batch is not full, break the loop as it indicates that this is the last page.
		if len(schedulers) < DefaultSchedulerGCBatchSize {
			break
		}
	}

	return nil
}

// sweep running the sweep operation for cleaning up inactive schedulers.
func (s *scheduler) sweep(ctx context.Context) error {
	logger.Info("running scheduler GC sweep")

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
		result := s.db.Where("updated_at < ?", time.Now().
			Add(-DefaultSchedulerGCTTL)).
			Where("state = ?", models.SchedulerStateInactive).
			Limit(DefaultSchedulerGCBatchSize).
			Unscoped().
			Delete(&models.Scheduler{})
		if result.Error != nil {
			gcResult.Error = result.Error
			return result.Error
		}

		if result.RowsAffected == 0 {
			break
		}

		gcResult.Purged += result.RowsAffected
		logger.Infof("scheduler GC deleted %d inactive schedulers", result.RowsAffected)
	}

	return nil
}
