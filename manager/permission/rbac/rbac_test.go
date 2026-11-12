/*
 *     Copyright 2020 The Dragonfly Authors
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

package rbac

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	managermodels "d7y.io/dragonfly/v2/manager/models"
)

func TestInitialRootPassword(t *testing.T) {
	tests := []struct {
		name   string
		env    string
		set    bool
		expect func(t *testing.T, password string, err error)
	}{
		{
			name: "environment variable is not set",
			set:  false,
			expect: func(t *testing.T, password string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(DefaultRootPassword, password)
			},
		},
		{
			name: "environment variable is empty",
			env:  "",
			set:  true,
			expect: func(t *testing.T, password string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(DefaultRootPassword, password)
			},
		},
		{
			name: "environment variable is set",
			env:  "dragonfly-root",
			set:  true,
			expect: func(t *testing.T, password string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("dragonfly-root", password)
			},
		},
		{
			name: "environment variable is too short",
			env:  strings.Repeat("a", MinRootPasswordLength-1),
			set:  true,
			expect: func(t *testing.T, password string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "environment variable is too long",
			env:  strings.Repeat("a", MaxRootPasswordLength+1),
			set:  true,
			expect: func(t *testing.T, password string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "environment variable is of the minimum length",
			env:  strings.Repeat("a", MinRootPasswordLength),
			set:  true,
			expect: func(t *testing.T, password string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(password, MinRootPasswordLength)
			},
		},
		{
			name: "environment variable is of the maximum length",
			env:  strings.Repeat("a", MaxRootPasswordLength),
			set:  true,
			expect: func(t *testing.T, password string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(password, MaxRootPasswordLength)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.set {
				t.Setenv(DragonflyInitialRootPasswordEnvName, tc.env)
			} else {
				os.Unsetenv(DragonflyInitialRootPasswordEnvName)
			}

			password, err := initialRootPassword()
			tc.expect(t, password, err)
		})
	}
}

func TestGetApiGroupName(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		expect func(t *testing.T, data string, err error)
	}{
		{
			name: `path is /api/v1/users`,
			path: "/api/v1/users",
			expect: func(t *testing.T, data string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("users", data)
			},
		},
		{
			name: `path is /api/v1/users/`,
			path: "/api/v1/users/",
			expect: func(t *testing.T, data string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("users", data)
			},
		},
		{
			name: `path is /api/v1/users/name`,
			path: "/api/v1/users/name",
			expect: func(t *testing.T, data string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("users", data)
			},
		},
		{
			name: `path is /api/user`,
			path: "/api/user",
			expect: func(t *testing.T, data string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "path is empty",
			path: "",
			expect: func(t *testing.T, data string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			name, err := GetAPIGroupName(tc.path)
			tc.expect(t, name, err)
		})
	}
}

func TestHTTPMethodToAction(t *testing.T) {
	tests := []struct {
		name   string
		method string
		expect func(t *testing.T, action string)
	}{
		{
			name:   "GET",
			method: http.MethodGet,
			expect: func(t *testing.T, action string) {
				assert := assert.New(t)
				assert.Equal(ReadAction, action)
			},
		},
		{
			name:   "HEAD",
			method: http.MethodHead,
			expect: func(t *testing.T, action string) {
				assert := assert.New(t)
				assert.Equal(ReadAction, action)
			},
		},
		{
			name:   "OPTIONS",
			method: http.MethodOptions,
			expect: func(t *testing.T, action string) {
				assert := assert.New(t)
				assert.Equal(ReadAction, action)
			},
		},
		{
			name:   "POST",
			method: http.MethodPost,
			expect: func(t *testing.T, action string) {
				assert := assert.New(t)
				assert.Equal(AllAction, action)
			},
		},
		{
			name:   "PUT",
			method: http.MethodPut,
			expect: func(t *testing.T, action string) {
				assert := assert.New(t)
				assert.Equal(AllAction, action)
			},
		},
		{
			name:   "PATCH",
			method: http.MethodPatch,
			expect: func(t *testing.T, action string) {
				assert := assert.New(t)
				assert.Equal(AllAction, action)
			},
		},
		{
			name:   "DELETE",
			method: http.MethodDelete,
			expect: func(t *testing.T, action string) {
				assert := assert.New(t)
				assert.Equal(AllAction, action)
			},
		},
		{
			name:   "UNKNOWN",
			method: "UNKNOWN",
			expect: func(t *testing.T, action string) {
				assert := assert.New(t)
				assert.Equal(ReadAction, action)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, HTTPMethodToAction(tc.method))
		})
	}
}

func TestInitRBAC(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "rbac.db")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         gormlogger.Discard,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&managermodels.User{}, &managermodels.CasbinRule{}); err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.GET("/api/v1/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.POST("/api/v1/clusters", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	enforcers := make([]*casbin.Enforcer, 2)
	for i := range enforcers {
		enforcer, err := NewEnforcer(db)
		if err != nil {
			t.Fatal(err)
		}

		enforcers[i] = enforcer
	}

	for _, enforcer := range enforcers {
		if err := InitRBAC(enforcer, router, db); err != nil {
			t.Fatal(err)
		}
	}

	var rootUser managermodels.User
	if err := db.Where(&managermodels.User{Name: RootUserName}).First(&rootUser).Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		enforcer *casbin.Enforcer
	}{
		{
			name:     "enforcer that seeded the root user",
			enforcer: enforcers[0],
		},
		{
			name:     "enforcer built before the root user was seeded",
			enforcer: enforcers[1],
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			ok, err := tc.enforcer.Enforce(fmt.Sprint(rootUser.ID), "clusters", AllAction)
			assert.NoError(err)
			assert.True(ok)
		})
	}
}
