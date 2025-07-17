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
	// DefaultSeedPeerGCBatchSize is the default batch size for deleting seed peers.
	DefaultSeedPeerGCBatchSize = 5000

	// DefaultSeedPeerTTL is the default TTL for seed peer.
	DefaultSeedPeerGCTTL = time.Minute * 30

	// DefaultSeedPeerGCInterval is the default interval for running seed peer GC.
	DefaultSeedPeerGCInterval = time.Hour * 1

	// DefaultSeedPeerGCTimeout is the default timeout for running seed peer GC.
	DefaultSeedPeerGCTimeout = time.Hour * 1

	// SeedPeerGCTaskID is the ID of the seed peer GC task.
	SeedPeerGCTaskID = "seed_peer"
)

// NewSeedPeerGCTask returns a new seed peer GC task.
func NewSeedPeerGCTask(db *gorm.DB) pkggc.Task {
	return pkggc.Task{
		ID:       SeedPeerGCTaskID,
		Interval: DefaultSeedPeerGCInterval,
		Timeout:  DefaultSeedPeerGCTimeout,
		Runner:   &seedPeer{db: db, recorder: newJobRecorder(db)},
	}
}

// seedPeer is the struct for cleaning up inactive seed peers which implements the gc Runner interface.
type seedPeer struct {
	db       *gorm.DB
	recorder *jobRecorder
}

// RunGC implements the gc Runner interface.
func (s *seedPeer) RunGC(ctx context.Context) error {
	args := models.JSONMap{
		"type":       SeedPeerGCTaskID,
		"ttl":        DefaultSeedPeerGCTTL,
		"batch_size": DefaultSeedPeerGCBatchSize,
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
		taskID = SeedPeerGCTaskID
	}

	if err := s.recorder.Init(userID, taskID, args); err != nil {
		return err
	}

	var gcResult Result
	defer func() {
		if err := s.recorder.Record(gcResult); err != nil {
			logger.Errorf("failed to record seed peer GC result: %v", err)
		}
	}()

	for {
		result := s.db.Where("updated_at < ?", time.Now().Add(-DefaultSeedPeerGCTTL)).Where("state = ?", models.SeedPeerStateInactive).Limit(DefaultSeedPeerGCBatchSize).Unscoped().Delete(&models.SeedPeer{})
		if result.Error != nil {
			gcResult.Error = result.Error
			return result.Error
		}

		if result.RowsAffected == 0 {
			break
		}

		gcResult.Purged += result.RowsAffected
		logger.Infof("seed peer GC deleted %d inactive seed peers", result.RowsAffected)
	}

	return nil
}
