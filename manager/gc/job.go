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
	"encoding/json"
	"time"

	"gorm.io/gorm"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/models"
	libgc "d7y.io/dragonfly/v2/pkg/gc"
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

func NewJobGCTask(db *gorm.DB) libgc.Task {
	return libgc.Task{
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
func (j *job) RunGC() error {
	ttl, err := j.getTTL()
	if err != nil {
		return err
	}

	if err = j.recorder.Init(JobGCTaskID, models.JSONMap{
		"ttl":        ttl,
		"batch_size": DefaultJobGCBatchSize,
	}); err != nil {
		return err
	}

	var gcResult Result
	defer func() {
		if err := j.recorder.Record(gcResult); err != nil {
			logger.Errorf("failed to record job GC result: %v", err)
		}
	}()

	for {
		result := j.db.Where("created_at < ?", time.Now().Add(-ttl)).Limit(DefaultJobGCBatchSize).Unscoped().Delete(&models.Job{})
		if result.Error != nil {
			gcResult.Error = result.Error
			return result.Error
		}

		if result.RowsAffected == 0 {
			break
		}

		gcResult.Purged += result.RowsAffected
		logger.Infof("gc job deleted %d jobs", result.RowsAffected)
	}

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
