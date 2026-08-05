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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				assert.EqualError(err, "DRAGONFLY_INITIAL_ROOT_PASSWORD must be between 8 and 20 characters")
			},
		},
		{
			name: "environment variable is too long",
			env:  strings.Repeat("a", MaxRootPasswordLength+1),
			set:  true,
			expect: func(t *testing.T, password string, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "DRAGONFLY_INITIAL_ROOT_PASSWORD must be between 8 and 20 characters")
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
				assert.Equal(data, "users")
			},
		},
		{
			name: `path is /api/v1/users/`,
			path: "/api/v1/users/",
			expect: func(t *testing.T, data string, err error) {
				assert := assert.New(t)
				assert.Equal(data, "users")
			},
		},
		{
			name: `path is /api/v1/users/name`,
			path: "/api/v1/users/name",
			expect: func(t *testing.T, data string, err error) {
				assert := assert.New(t)
				assert.Equal(data, "users")
			},
		},
		{
			name: `path is /api/user`,
			path: "/api/user",
			expect: func(t *testing.T, data string, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "cannot find group name")
			},
		},
		{
			name: "path is empty",
			path: "",
			expect: func(t *testing.T, data string, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "cannot find group name")
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
		method         string
		expectedAction string
	}{
		{
			method:         "GET",
			expectedAction: ReadAction,
		},
		{
			method:         "POST",
			expectedAction: AllAction,
		},
		{
			method:         "UNKNOWN",
			expectedAction: ReadAction,
		},
	}

	for _, tt := range tests {
		action := HTTPMethodToAction(tt.method)
		if action != tt.expectedAction {
			t.Errorf("HttpMethodToAction(%v) = %v, want %v", tt.method, action, tt.expectedAction)
		}
	}
}

// newTestDB returns a sqlite database with the tables InitRBAC touches.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// SingularTable matches the production configuration, so that the table the
	// casbin gorm adapter uses is the one migrated here.
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "rbac.db")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         gormlogger.Discard,
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&managermodels.User{}, &managermodels.CasbinRule{}))

	return db
}

// newTestRouter returns a router whose paths match apiGroupRegexp, so that
// GetPermissions derives the users and clusters objects from it.
func newTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.GET("/api/v1/users", func(c *gin.Context) { c.Status(http.StatusOK) })
	g.POST("/api/v1/clusters", func(c *gin.Context) { c.Status(http.StatusOK) })

	return g
}

// TestInitRBACConcurrentReplicas covers replicas that build their enforcer before
// another replica seeds the root user. Every enforcer must hold the root grant,
// not only the one that did the seeding.
func TestInitRBACConcurrentReplicas(t *testing.T) {
	assert := assert.New(t)
	db := newTestDB(t)
	router := newTestRouter(t)

	// Both replicas load an empty casbin_rule before either seeds the root user.
	enforcerA, err := NewEnforcer(db)
	require.NoError(t, err)
	enforcerB, err := NewEnforcer(db)
	require.NoError(t, err)

	// Replica A seeds the root user, replica B finds it already there.
	require.NoError(t, InitRBAC(enforcerA, router, db))
	require.NoError(t, InitRBAC(enforcerB, router, db))

	var rootUser managermodels.User
	require.NoError(t, db.Where(&managermodels.User{Name: RootUserName}).First(&rootUser).Error)
	sub := fmt.Sprint(rootUser.ID)

	okA, err := enforcerA.Enforce(sub, "clusters", AllAction)
	require.NoError(t, err)
	assert.True(okA, "enforcer of the replica that seeded the root user denies root")

	okB, err := enforcerB.Enforce(sub, "clusters", AllAction)
	require.NoError(t, err)
	assert.True(okB, "enforcer of the replica that did not seed the root user denies root")
}
