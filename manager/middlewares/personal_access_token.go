/*
 *     Copyright 2023 The Dragonfly Authors
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
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"gorm.io/gorm"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

var (
	oapiResourceRegexp = regexp.MustCompile(`^/oapi/v[0-9]+/([-_a-zA-Z]*)[/.*]*`)
)

func PersonalAccessToken(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get bearer token from Authorization header.
		authorization := c.GetHeader(headers.Authorization)
		tokenFields := strings.Fields(authorization)
		if len(tokenFields) != 2 || tokenFields[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: http.StatusText(http.StatusUnauthorized),
			})
			c.Abort()
			return
		}

		// Check if the personal access token is valid.
		personalAccessToken := tokenFields[1]
		var token models.PersonalAccessToken
		if err := gdb.WithContext(c).Where("token = ?", personalAccessToken).First(&token).Error; err != nil {
			logger.Errorf("Invalid personal access token attempt: %s, error: %v", c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: http.StatusText(http.StatusUnauthorized),
			})
			c.Abort()
			return
		}

		// Check if the token is active.
		if token.State != models.PersonalAccessTokenStateActive {
			logger.Errorf("Inactive token used: %s, token name: %s, user_id: %d", c.Request.URL.Path, token.Name, token.UserID)
			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "Token is inactive",
			})
			c.Abort()
			return
		}

		// Check if the token has expired.
		if time.Now().After(token.ExpiredAt) {
			logger.Errorf("Expired token used: %s, token name: %s, user_id: %d, expired: %v",
				c.Request.URL.Path, token.Name, token.UserID, token.ExpiredAt)
			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "Token has expired",
			})
			c.Abort()
			return
		}

		// Check if the token's scopes include the required resource type.
		hasScope := false
		resourceType := getAPIResourceType(c.Request.URL.Path)
		for _, scope := range token.Scopes {
			if scope == resourceType {
				hasScope = true
				break
			}
		}

		if !hasScope {
			logger.Errorf("Insufficient scope token used: %s, token name: %s, user_id: %d, required: %s, available: %v",
				c.Request.URL.Path, token.Name, token.UserID, resourceType, token.Scopes)
			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: fmt.Sprintf("Token doesn't have permission to access this resource. Required scope: %s", resourceType),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getAPIResourceType extracts the resource type from the path.
// For example: /oapi/v1/jobs -> job, /oapi/v1/clusters -> cluster.
func getAPIResourceType(path string) string {
	matches := oapiResourceRegexp.FindStringSubmatch(path)
	if len(matches) != 2 {
		return ""
	}

	resource := strings.ToLower(matches[1])
	switch resource {
	case "jobs":
		return types.PersonalAccessTokenScopeJob
	case "clusters":
		return types.PersonalAccessTokenScopeCluster
	case "preheats":
		return types.PersonalAccessTokenScopePreheat
	default:
		return resource
	}
}
