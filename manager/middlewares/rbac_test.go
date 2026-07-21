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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"d7y.io/dragonfly/v2/manager/permission/rbac"
)

const rbacModelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && (r.act == p.act || p.act == "*")
`

func mockRBACRouter(e *casbin.Enforcer, id float64) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("id", id)
		c.Next()
	})

	r.Use(RBAC(e))
	r.GET("/api/v1/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r
}

func newRBACEnforcer(t *testing.T, id uint) *casbin.Enforcer {
	m, err := model.NewModelFromString(rbacModelText)
	assert.NoError(t, err)

	e, err := casbin.NewEnforcer(m)
	assert.NoError(t, err)

	_, err = e.AddPermissionForUser(rbac.RootRole, "users", rbac.AllAction)
	assert.NoError(t, err)

	_, err = e.AddRoleForUser(fmt.Sprint(id), rbac.RootRole)
	assert.NoError(t, err)
	return e
}

func TestRBAC(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		id     uint
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "small id is allowed",
			id:   1,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
			},
		},
		{
			name: "id just below the scientific notation boundary is allowed",
			id:   999999,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
			},
		},
		{
			name: "id at the scientific notation boundary is allowed",
			id:   1000000,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
			},
		},
		{
			name: "large non-round id is allowed",
			id:   1234567,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
			},
		},
		{
			name: "max uint32 id is allowed",
			id:   4294967295,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newRBACEnforcer(t, tc.id)
			router := mockRBACRouter(e, float64(tc.id))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			router.ServeHTTP(w, req)

			tc.expect(t, w)
		})
	}
}
