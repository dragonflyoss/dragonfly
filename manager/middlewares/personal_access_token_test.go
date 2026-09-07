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
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

func newPersonalAccessTokenDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "pat.db")), &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormlogger.Discard,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&models.PersonalAccessToken{}); err != nil {
		t.Fatal(err)
	}

	if err := db.Create([]*models.PersonalAccessToken{
		{Name: "active", Token: "active-token", Scopes: models.Array{types.PersonalAccessTokenScopeJob}, State: models.PersonalAccessTokenStateActive, ExpiredAt: time.Now().Add(time.Hour)},
		{Name: "unscoped", Token: "unscoped-token", Scopes: models.Array{}, State: models.PersonalAccessTokenStateActive, ExpiredAt: time.Now().Add(time.Hour)},
		{Name: "cluster", Token: "cluster-token", Scopes: models.Array{types.PersonalAccessTokenScopeCluster}, State: models.PersonalAccessTokenStateActive, ExpiredAt: time.Now().Add(time.Hour)},
		{Name: "inactive", Token: "inactive-token", Scopes: models.Array{types.PersonalAccessTokenScopeJob}, State: models.PersonalAccessTokenStateInactive, ExpiredAt: time.Now().Add(time.Hour)},
		{Name: "expired", Token: "expired-token", Scopes: models.Array{types.PersonalAccessTokenScopeJob}, State: models.PersonalAccessTokenStateActive, ExpiredAt: time.Now().Add(-time.Hour)},
	}).Error; err != nil {
		t.Fatal(err)
	}

	return db
}

func mockPersonalAccessTokenRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(PersonalAccessToken(db))

	for _, path := range []string{"/oapi/v1/jobs", "/oapi/v1/users"} {
		r.GET(path, func(c *gin.Context) {
			if pat, ok := c.Get("pat"); ok {
				c.String(http.StatusOK, pat.(*models.PersonalAccessToken).Name)
				return
			}

			c.Status(http.StatusOK)
		})
	}

	return r
}

func TestPersonalAccessToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := mockPersonalAccessTokenRouter(newPersonalAccessTokenDB(t))

	tests := []struct {
		name          string
		target        string
		authorization string
		expect        func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:   "missing authorization header",
			target: "/oapi/v1/jobs",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
			},
		},
		{
			name:          "authorization header without bearer scheme",
			target:        "/oapi/v1/jobs",
			authorization: "Basic active-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
			},
		},
		{
			name:          "authorization header with too many fields",
			target:        "/oapi/v1/jobs",
			authorization: "Bearer active-token extra",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
			},
		},
		{
			name:          "unknown token",
			target:        "/oapi/v1/jobs",
			authorization: "Bearer unknown-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
			},
		},
		{
			name:          "inactive token",
			target:        "/oapi/v1/jobs",
			authorization: "Bearer inactive-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusForbidden, w.Code)
				assert.JSONEq(`{"message":"Token is inactive"}`, w.Body.String())
			},
		},
		{
			name:          "expired token",
			target:        "/oapi/v1/jobs",
			authorization: "Bearer expired-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusForbidden, w.Code)
				assert.JSONEq(`{"message":"Token has expired"}`, w.Body.String())
			},
		},
		{
			name:          "token without required scope",
			target:        "/oapi/v1/jobs",
			authorization: "Bearer cluster-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusForbidden, w.Code)
				assert.JSONEq(`{"message":"Token doesn't have permission to access this resource. Required permission: job"}`, w.Body.String())
			},
		},
		{
			name:          "path with unsupported resource",
			target:        "/oapi/v1/users",
			authorization: "Bearer active-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusForbidden, w.Code)
				assert.JSONEq(`{"message":"Failed to extract resource type from path: /oapi/v1/users"}`, w.Body.String())
			},
		},
		{
			name:          "valid bearer token",
			target:        "/oapi/v1/jobs",
			authorization: "Bearer active-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.Equal("active", w.Body.String())
			},
		},
		{
			name:          "bearer scheme is case insensitive",
			target:        "/oapi/v1/jobs",
			authorization: "bearer active-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.Equal("active", w.Body.String())
			},
		},
		{
			name:   "valid token from access_token query",
			target: "/oapi/v1/jobs?access_token=active-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.Equal("active", w.Body.String())
			},
		},
		{
			name:          "access_token query takes precedence over header",
			target:        "/oapi/v1/jobs?access_token=active-token",
			authorization: "Bearer inactive-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.Equal("active", w.Body.String())
			},
		},
		{
			name:          "token without scopes grants all resources",
			target:        "/oapi/v1/jobs",
			authorization: "Bearer unscoped-token",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.Equal("unscoped", w.Body.String())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.target, nil)
			if tc.authorization != "" {
				req.Header.Set(headers.Authorization, tc.authorization)
			}

			router.ServeHTTP(w, req)

			tc.expect(t, w)
		})
	}
}

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name               string
		permissions        []string
		requiredPermission string
		expect             func(t *testing.T, ok bool)
	}{
		{
			name:               "nil permissions grant all",
			permissions:        nil,
			requiredPermission: types.PersonalAccessTokenScopeJob,
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name:               "empty permissions grant all",
			permissions:        []string{},
			requiredPermission: types.PersonalAccessTokenScopeCluster,
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name:               "required permission is contained",
			permissions:        []string{types.PersonalAccessTokenScopeJob, types.PersonalAccessTokenScopeCluster},
			requiredPermission: types.PersonalAccessTokenScopeCluster,
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name:               "required permission is missing",
			permissions:        []string{types.PersonalAccessTokenScopeJob},
			requiredPermission: types.PersonalAccessTokenScopeCluster,
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, hasPermission(tc.permissions, tc.requiredPermission))
		})
	}
}

func TestRequiredPermission(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		expect func(t *testing.T, permission string, err error)
	}{
		{
			name: "jobs collection",
			path: "/oapi/v1/jobs",
			expect: func(t *testing.T, permission string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.PersonalAccessTokenScopeJob, permission)
			},
		},
		{
			name: "jobs resource with id",
			path: "/oapi/v1/jobs/1",
			expect: func(t *testing.T, permission string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.PersonalAccessTokenScopeJob, permission)
			},
		},
		{
			name: "clusters under another api version",
			path: "/oapi/v2/clusters/1/schedulers",
			expect: func(t *testing.T, permission string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.PersonalAccessTokenScopeCluster, permission)
			},
		},
		{
			name: "resource name is case insensitive",
			path: "/oapi/v1/JOBS",
			expect: func(t *testing.T, permission string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.PersonalAccessTokenScopeJob, permission)
			},
		},
		{
			name: "unsupported resource",
			path: "/oapi/v1/users",
			expect: func(t *testing.T, permission string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(permission)
			},
		},
		{
			name: "path outside of the oapi prefix",
			path: "/api/v1/jobs",
			expect: func(t *testing.T, permission string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(permission)
			},
		},
		{
			name: "empty path",
			path: "",
			expect: func(t *testing.T, permission string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(permission)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			permission, err := requiredPermission(tc.path)
			tc.expect(t, permission, err)
		})
	}
}
