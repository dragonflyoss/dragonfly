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

package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/service"
	"d7y.io/dragonfly/v2/manager/types"
)

func Audit(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow the request to proceed through the chain first.
		c.Next()

		actorType := models.ActorTypeUnknown
		actorName := "unknown"
		// Check if user information is available from previous auth middleware.
		// Personal Access Token auth.
		if rawPAT, ok := c.Get("pat"); ok {
			if pat, ok := rawPAT.(*models.PersonalAccessToken); ok {
				actorType = models.ActorTypePat
				actorName = pat.Name
			}
		}

		// JWT auth.
		if rawID, ok := c.Get("id"); ok {
			if id, ok := rawID.(float64); ok {
				if user, err := service.GetUser(c.Request.Context(), uint(id)); err == nil {
					actorType = models.ActorTypeUser
					actorName = user.Name
				} else {
					logger.Errorf("failed to get user by id %d: %v", uint(id), err)
				}
			}
		}

		if err := service.AsyncCreateAudit(c.Request.Context(), &types.CreateAuditRequest{
			ActorType:  actorType,
			ActorName:  actorName,
			EventType:  models.EventTypeAPI,
			Operation:  c.Request.Method,
			OperatedAt: time.Now(),
			State:      httpStatusCodeToState(c.Writer.Status()),
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
		}); err != nil {
			logger.Errorf("failed to create audit: %v", err)
		}
	}
}

// httpStatusCodeToState converts an HTTP status code to an audit state.
func httpStatusCodeToState(statusCode int) string {
	if statusCode >= 200 && statusCode < 300 {
		return models.AuditStateSuccess
	}
	return models.AuditStateFailure
}
