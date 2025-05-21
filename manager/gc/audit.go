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
	// DefaultAuditGCBatchSize is the default batch size for deleting jobs.
	DefaultAuditGCBatchSize = 5000

	// DefaultAuditGCInterval is the default interval for running audit GC.
	DefaultAuditGCInterval = time.Hour * 6

	// DefaultAuditGCTimeout is the default timeout for running audit GC.
	DefaultAuditGCTimeout = time.Hour * 2

	// AuditGCTaskID is the ID of the audit GC task.
	AuditGCTaskID = "audit"
)

// NewAuditGCTask returns a new audit GC task.
func NewAuditGCTask(db *gorm.DB) pkggc.Task {
	return pkggc.Task{
		ID:       AuditGCTaskID,
		Interval: DefaultAuditGCInterval,
		Timeout:  DefaultAuditGCTimeout,
		Runner:   &audit{db: db, recorder: newJobRecorder(db)},
	}
}

// audit is the struct for cleaning up audits which implements the gc Runner interface.
type audit struct {
	db       *gorm.DB
	recorder *jobRecorder
}

// RunGC implements the gc Runner interface.
func (a *audit) RunGC(ctx context.Context) error {
	ttl, err := a.getTTL()
	if err != nil {
		return err
	}

	args := models.JSONMap{
		"ttl":        ttl,
		"batch_size": DefaultAuditGCBatchSize,
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

	if err := a.recorder.Init(userID, taskID, args); err != nil {
		return err
	}

	var gcResult Result
	defer func() {
		if err := a.recorder.Record(gcResult); err != nil {
			logger.Errorf("failed to record audit GC result: %v", err)
		}
	}()

	for {
		result := a.db.Where("created_at < ?", time.Now().Add(-ttl)).Limit(DefaultAuditGCBatchSize).Unscoped().Delete(&models.Audit{})
		if result.Error != nil {
			gcResult.Error = result.Error
			return result.Error
		}

		if result.RowsAffected == 0 {
			break
		}

		gcResult.Purged += result.RowsAffected
		logger.Infof("gc audit deleted %d audits", result.RowsAffected)
	}

	return nil
}

func (a *audit) getTTL() (time.Duration, error) {
	var config models.Config
	if err := a.db.Model(models.Config{}).First(&config, &models.Config{Name: models.ConfigGC}).Error; err != nil {
		return 0, err
	}

	var gcConfig models.GCConfig
	if err := json.Unmarshal([]byte(config.Value), &gcConfig); err != nil {
		return 0, err
	}

	return gcConfig.Audit.TTL, nil
}
