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
	"encoding/json"
	"time"

	"gorm.io/gorm"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/models"
	pkggc "d7y.io/dragonfly/v2/pkg/gc"
)

const (
	// DefaultJobGCBatchSize is the default batch size for deleting jobs.
	DefaultJobGCBatchSize = 5000

	// DefaultJobGCInterval is the default interval for running job GC.
	DefaultJobGCInterval = time.Hour * 3

	// DefaultJobGCTimeout is the default timeout for running job GC.
	DefaultJobGCTimeout = time.Hour * 1

	// JohGCTaskID is the ID of the job GC task.
	JobGCTaskID = "job"
)

func NewJobGCTask(db *gorm.DB) pkggc.Task {
	return pkggc.Task{
		ID:       JobGCTaskID,
		Interval: DefaultJobGCInterval,
		Timeout:  DefaultJobGCTimeout,
		Runner:   &job{db: db, recorder: newJobRecorder(db)},
	}
}

// job is the struct for cleaning up jobs which implements the gc Runner interface.
type job struct {
	db       *gorm.DB
	recorder *jobRecorder
}

// RunGC implements the gc Runner interface.
func (j *job) RunGC(ctx context.Context) error {
	ttl, err := j.getTTL()
	if err != nil {
		return err
	}

	args := models.JSONMap{
		"type":       JobGCTaskID,
		"ttl":        ttl,
		"batch_size": DefaultJobGCBatchSize,
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
		taskID = AuditGCTaskID
	}

	if err = j.recorder.Init(userID, taskID, args); err != nil {
		return err
	}

	var gcResult Result
	defer func() {
		if err := j.recorder.Record(gcResult); err != nil {
			logger.Errorf("failed to record job GC result: %v", err)
		}
	}()

	for {
		var currentBatchAffectedRows int64
		if err := j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var jobIDs []uint
			if err := tx.Model(&models.Job{}).
				Where("created_at < ?", time.Now().Add(-ttl)).
				Limit(DefaultJobGCBatchSize).
				Pluck("id", &jobIDs).Error; err != nil {
				return err
			}

			if len(jobIDs) == 0 {
				currentBatchAffectedRows = 0
				return nil
			}

			// 1. Explicitly delete from the first join table (e.g., job_seed_peer_cluster).
			//    Using tx.Exec for raw SQL deletion from the join table.
			if err := tx.Exec("DELETE FROM job_seed_peer_cluster WHERE job_id IN (?)", jobIDs).Error; err != nil {
				return err
			}

			// 2. Explicitly delete from the second join table (e.g., job_scheduler_cluster).
			if err := tx.Exec("DELETE FROM job_scheduler_cluster WHERE job_id IN (?)", jobIDs).Error; err != nil {
				return err
			}

			// 3. Delete the jobs themselves using the collected IDs.
			//    Use Unscoped() for a hard delete if your models.Job uses GORM's soft delete (gorm.DeletedAt).
			result := tx.Where("id IN (?)", jobIDs).Unscoped().Delete(&models.Job{})
			if result.Error != nil {
				return result.Error
			}
			currentBatchAffectedRows = result.RowsAffected

			return nil
		}); err != nil {
			gcResult.Error = err
			logger.Errorf("gc job batch processing failed: %v", err)
			return err
		}

		if currentBatchAffectedRows == 0 {
			logger.Info("gc job finished: no more jobs to delete in this iteration.")
			break
		}

		gcResult.Purged += currentBatchAffectedRows
		logger.Infof("gc job deleted %d jobs in this batch", currentBatchAffectedRows)
	}

	logger.Infof("gc job completed: total %d jobs purged", gcResult.Purged)
	return nil
}

func (j *job) getTTL() (time.Duration, error) {
	var config models.Config
	if err := j.db.Model(models.Config{}).First(&config, &models.Config{Name: models.ConfigGC}).Error; err != nil {
		return 0, err
	}

	var gcConfig models.GCConfig
	if err := json.Unmarshal([]byte(config.Value), &gcConfig); err != nil {
		return 0, err
	}

	return gcConfig.Job.TTL, nil
}
