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

	redis_rate "github.com/go-redis/redis_rate/v10"
	redis "github.com/redis/go-redis/v9"
)

// DistributedRateLimiter is an interface for a distributed rate limiter.
type DistributedRateLimiter interface {
	// Allow checks if a request is allowed based on the rate limit.
	Allow(ctx context.Context) (*redis_rate.Result, error)
}

// distributedRateLimiter is an implementation of DistributedRateLimiter.
type distributedRateLimiter struct {
	// key used to store the rate limit in the database.
	key string

	// limiter is the rate limiter instance.
	limiter *redis_rate.Limiter

	// limit is the rate limit in requests per second.
	limit uint32
}

// NewDistributedRateLimiter creates a new instance of DistributedRateLimiter. Parameters:
//   - rdb: Redis client for distributed rate limiting.
//   - key: Unique key for the rate limit.
//   - limit: Rate limit in requests per second.
func NewDistributedRateLimiter(rdb redis.UniversalClient, key string, limit uint32) DistributedRateLimiter {
	limiter := redis_rate.NewLimiter(rdb)
	return &distributedRateLimiter{key, limiter, limit}
}

// Allow checks if a request is allowed based on the rate limit.
func (d *distributedRateLimiter) Allow(ctx context.Context) (*redis_rate.Result, error) {
	return d.limiter.Allow(ctx, d.key, redis_rate.PerSecond(int(d.limit)))
}
