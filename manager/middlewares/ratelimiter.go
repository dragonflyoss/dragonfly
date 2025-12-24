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

package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"d7y.io/dragonfly/v2/internal/ratelimiter"
	"d7y.io/dragonfly/v2/manager/types"
)

// CreateJobRateLimiter create job rate limiter middleware
func CreateJobRateLimiter(limiter ratelimiter.JobRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json types.CreateJobRequest
		if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
			return
		}

		if allowed := limiter.AllowByClusterIDs(c, json.SchedulerClusterIDs); !allowed {
			c.String(http.StatusTooManyRequests, "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
