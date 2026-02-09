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

package ratelimiter

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/config"
	"d7y.io/dragonfly/v2/manager/database"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

const (
	// jobRateLimiterSuffix is the suffix of the job rate limiter key.
	jobRateLimiterSuffix = "job"

	// defaultRefreshInterval is the default interval to refresh the rate limiters.
	defaultRefreshInterval = 3 * time.Minute
)

// JobRateLimiter is an interface for a job rate limiter.
type JobRateLimiter interface {
	// AllowByClusterID checks if a request is allowed based on the rate limit for a specific cluster ID.
	AllowByClusterID(ctx context.Context, clusterID uint) bool

	// AllowByClusterIDs checks if a request is allowed based on the rate limit for multiple cluster IDs.
	// If any cluster ID is not allowed, it returns false.
	AllowByClusterIDs(ctx context.Context, clusterIDs []uint) bool

	// Serve started job rate limiter server.
	Serve()

	// Stop job rate limiter server.
	Stop()
}

// jobRateLimiter is an implementation of JobRateLimiter.
type jobRateLimiter struct {
	// database used to store the rate limit.
	database *database.Database

	// clusters is a map of rate limiters for each cluster.
	clusters *sync.Map

	// refreshInterval is the interval to refresh the rate limiters.
	refreshInterval time.Duration

	// done is the channel to stop the rate limiter server.
	done chan struct{}
}

// NewJobRateLimiter creates a new instance of JobRateLimiter.
func NewJobRateLimiter(database *database.Database) (JobRateLimiter, error) {
	j := &jobRateLimiter{
		database:        database,
		clusters:        &sync.Map{},
		refreshInterval: defaultRefreshInterval,
		done:            make(chan struct{}),
	}

	if err := j.refresh(context.Background()); err != nil {
		return nil, err
	}

	return j, nil
}

// AllowByClusterID checks if a request is allowed based on the rate limit for a specific cluster ID.
func (j *jobRateLimiter) AllowByClusterID(ctx context.Context, clusterID uint) bool {
	rawLimiter, loaded := j.clusters.Load(clusterID)
	if !loaded {
		logger.Errorf("[job-rate-limiter]: cluster %d not found", clusterID)
		return false
	}

	limiter, ok := rawLimiter.(DistributedRateLimiter)
	if !ok {
		logger.Errorf("[job-rate-limiter]: cluster %d is not a distributed rate limiter", clusterID)
		return false
	}

	result, err := limiter.Allow(ctx)
	if err != nil {
		logger.Errorf("[job-rate-limiter]: cluster %d allow failed: %v", clusterID, err)
		return false
	}

	if result.Allowed == 0 {
		logger.Errorf("[job-rate-limiter]: cluster %d rate limit exceeded", clusterID)
		return false
	}

	return true
}

// AllowByClusterIDs checks if a request is allowed based on the rate limit for multiple cluster IDs.
// If any cluster ID is not allowed, it returns false.
func (j *jobRateLimiter) AllowByClusterIDs(ctx context.Context, clusterIDs []uint) bool {
	for _, clusterID := range clusterIDs {
		if allowed := j.AllowByClusterID(ctx, clusterID); !allowed {
			return false
		}
	}

	return true
}

// Serve started rate limiter server.
func (j *jobRateLimiter) Serve() {
	tick := time.NewTicker(j.refreshInterval)
	for {
		select {
		case <-tick.C:
			logger.Infof("[job-rate-limiter]: refresh job rate limiter started")
			if err := j.refresh(context.Background()); err != nil {
				logger.Errorf("[job-rate-limiter]: refresh job rate limiter failed: %v", err)
			}
		case <-j.done:
			return
		}
	}
}

// Stop rate limiter server.
func (j *jobRateLimiter) Stop() {
	close(j.done)
}

// refresh refreshes the rate limiters for all scheduler clusters.
func (j *jobRateLimiter) refresh(ctx context.Context) error {
	var schedulerClusters []models.SchedulerCluster
	if err := j.database.DB.WithContext(ctx).Find(&schedulerClusters).Error; err != nil {
		return err
	}

	j.clusters.Clear()
	for _, schedulerCluster := range schedulerClusters {
		b, err := schedulerCluster.Config.MarshalJSON()
		if err != nil {
			logger.Errorf("[job-rate-limiter]: marshal scheduler cluster %d config failed: %v", schedulerCluster.ID, err)
			return err
		}

		var schedulerClusterConfig types.SchedulerClusterConfig
		if err := json.Unmarshal(b, &schedulerClusterConfig); err != nil {
			logger.Errorf("[job-rate-limiter]: unmarshal scheduler cluster %d config failed: %v", schedulerCluster.ID, err)
			return err
		}

		// Use the default rate limit if the rate limit is not set.
		jobRateLimit := config.DefaultClusterJobRateLimit
		if schedulerClusterConfig.JobRateLimit != 0 {
			jobRateLimit = schedulerClusterConfig.JobRateLimit
		}

		logger.Debugf("[job-rate-limiter]: create job rate limiter for scheduler cluster %d with rate limit %d", schedulerCluster.ID, jobRateLimit)
		j.clusters.Store(schedulerCluster.ID,
			NewDistributedRateLimiter(j.database.RDB, j.key(schedulerCluster.ID), jobRateLimit))
	}

	return nil
}

// key is the rate limiter key for storing value in the database.
func (j *jobRateLimiter) key(clusterID uint) string {
	return fmt.Sprintf("%d-%s", clusterID, jobRateLimiterSuffix)
}
