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
	"slices"
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
	// oapiResourceRegexp is a regular expression to extract the resource type from the path.
	// Example: /oapi/v1/jobs/1 -> jobs.
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
			logger.Errorf("invalid personal access token attempt: %s, error: %v", c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: http.StatusText(http.StatusUnauthorized),
			})

			c.Abort()
			return
		}

		// Check if the token is active.
		if token.State != models.PersonalAccessTokenStateActive {
			logger.Errorf("inactive token used: %s, token name: %s, user_id: %d", c.Request.URL.Path, token.Name, token.UserID)
			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "Token is inactive",
			})

			c.Abort()
			return
		}

		// Check if the token has expired.
		if time.Now().After(token.ExpiredAt) {
			logger.Errorf("expired token used: %s, token name: %s, user_id: %d, expired: %v",
				c.Request.URL.Path, token.Name, token.UserID, token.ExpiredAt)
			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "Token has expired",
			})

			c.Abort()
			return
		}

		// Check if the token's scopes include the required resource type.
		requiredPermission, err := requiredPermission(c.Request.URL.Path)
		if err != nil {
			logger.Errorf("failed to extract resource type from path: %s, error: %v", c.Request.URL.Path, err)
			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: fmt.Sprintf("Failed to extract resource type from path: %s", c.Request.URL.Path),
			})

			c.Abort()
			return
		}

		if !hasPermission(token.Scopes, requiredPermission) {
			logger.Errorf("insufficient scope token used %s. Required permission: %s", token.Name, requiredPermission)
			c.JSON(http.StatusForbidden, ErrorResponse{
				Message: fmt.Sprintf("Token doesn't have permission to access this resource. Required permission: %s", requiredPermission),
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

// hasPermission checks if the required permission exists in the provided permissions list.
// For backward compatibility, an empty permissions list grants all permissions.
// This allows existing systems that don't have explicit permissions set to continue
// working without interruption.
//
// Returns true if:
// 1. The permissions list is empty (backward compatibility mode)
// 2. The requiredPermission is found in the permissions list
func hasPermission(permissions []string, requiredPermission string) bool {
	if len(permissions) == 0 {
		return true
	}

	return slices.Contains(permissions, requiredPermission)
}

// requiredPermission extracts the resource type from the path and returns the required permission.
func requiredPermission(path string) (string, error) {
	matches := oapiResourceRegexp.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to extract resource type from path: %s", path)
	}

	resource := strings.ToLower(matches[1])
	switch resource {
	case "jobs":
		return types.PersonalAccessTokenScopeJob, nil
	case "clusters":
		return types.PersonalAccessTokenScopeCluster, nil
	default:
		return "", fmt.Errorf("unsupported resource type: %s", resource)
	}
}
