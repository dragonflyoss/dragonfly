/*
 *     Copyright 2026 The Dragonfly Authors
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

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"d7y.io/dragonfly/v2/manager/permission/rbac"
	"d7y.io/dragonfly/v2/manager/types"
)

func TestService_CreateRole(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		req    types.CreateRoleRequest
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name:  "adds one policy per permission",
			setup: func(t *testing.T, s *service) {},
			req: types.CreateRoleRequest{
				Role: "developer",
				Permissions: []rbac.Permission{
					{Object: "clusters", Action: rbac.ReadAction},
					{Object: "jobs", Action: rbac.AllAction},
				},
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.ElementsMatch([][]string{
					{"developer", "clusters", rbac.ReadAction},
					{"developer", "jobs", rbac.AllAction},
				}, s.GetRole(context.Background(), "developer"))
				assert.Contains(s.GetRoles(context.Background()), "developer")
			},
		},
		{
			name: "re-creating the role does not duplicate policies",
			setup: func(t *testing.T, s *service) {
				assert := assert.New(t)
				assert.NoError(s.CreateRole(context.Background(), types.CreateRoleRequest{
					Role:        "developer",
					Permissions: []rbac.Permission{{Object: "clusters", Action: rbac.ReadAction}},
				}))
			},
			req: types.CreateRoleRequest{
				Role:        "developer",
				Permissions: []rbac.Permission{{Object: "clusters", Action: rbac.ReadAction}},
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([][]string{{"developer", "clusters", rbac.ReadAction}}, s.GetRole(context.Background(), "developer"))
			},
		},
		{
			name:  "role without permissions is not registered",
			setup: func(t *testing.T, s *service) {},
			req:   types.CreateRoleRequest{Role: "empty"},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(s.GetRole(context.Background(), "empty"))
				assert.NotContains(s.GetRoles(context.Background()), "empty")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			tc.setup(t, s)

			tc.expect(t, s, s.CreateRole(context.Background(), tc.req))
		})
	}
}

func TestService_DestroyRole(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		role   string
		expect func(t *testing.T, s *service, ok bool, err error)
	}{
		{
			name: "removes every policy of the role",
			setup: func(t *testing.T, s *service) {
				assert := assert.New(t)
				assert.NoError(s.CreateRole(context.Background(), types.CreateRoleRequest{
					Role: "developer",
					Permissions: []rbac.Permission{
						{Object: "clusters", Action: rbac.ReadAction},
						{Object: "jobs", Action: rbac.AllAction},
					},
				}))
			},
			role: "developer",
			expect: func(t *testing.T, s *service, ok bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(ok)
				assert.Empty(s.GetRole(context.Background(), "developer"))
				assert.NotContains(s.GetRoles(context.Background()), "developer")
			},
		},
		{
			name:  "unknown role",
			setup: func(t *testing.T, s *service) {},
			role:  "unknown",
			expect: func(t *testing.T, s *service, ok bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.False(ok)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			tc.setup(t, s)

			ok, err := s.DestroyRole(context.Background(), tc.role)
			tc.expect(t, s, ok, err)
		})
	}
}

func TestService_AddPermissionForRole(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		req    types.AddPermissionForRoleRequest
		expect func(t *testing.T, s *service, ok bool, err error)
	}{
		{
			name: "new permission",
			setup: func(t *testing.T, s *service) {
				assert := assert.New(t)
				assert.NoError(s.CreateRole(context.Background(), types.CreateRoleRequest{
					Role:        "developer",
					Permissions: []rbac.Permission{{Object: "clusters", Action: rbac.ReadAction}},
				}))
			},
			req: types.AddPermissionForRoleRequest{Permission: rbac.Permission{Object: "jobs", Action: rbac.AllAction}},
			expect: func(t *testing.T, s *service, ok bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(ok)
				assert.ElementsMatch([][]string{
					{"developer", "clusters", rbac.ReadAction},
					{"developer", "jobs", rbac.AllAction},
				}, s.GetRole(context.Background(), "developer"))
			},
		},
		{
			name: "existing permission",
			setup: func(t *testing.T, s *service) {
				assert := assert.New(t)
				assert.NoError(s.CreateRole(context.Background(), types.CreateRoleRequest{
					Role:        "developer",
					Permissions: []rbac.Permission{{Object: "clusters", Action: rbac.ReadAction}},
				}))
			},
			req: types.AddPermissionForRoleRequest{Permission: rbac.Permission{Object: "clusters", Action: rbac.ReadAction}},
			expect: func(t *testing.T, s *service, ok bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.False(ok)
				assert.Len(s.GetRole(context.Background(), "developer"), 1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			tc.setup(t, s)

			ok, err := s.AddPermissionForRole(context.Background(), "developer", tc.req)
			tc.expect(t, s, ok, err)
		})
	}
}

func TestService_DeletePermissionForRole(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		req    types.DeletePermissionForRoleRequest
		expect func(t *testing.T, s *service, ok bool, err error)
	}{
		{
			name: "existing permission",
			setup: func(t *testing.T, s *service) {
				assert := assert.New(t)
				assert.NoError(s.CreateRole(context.Background(), types.CreateRoleRequest{
					Role: "developer",
					Permissions: []rbac.Permission{
						{Object: "clusters", Action: rbac.ReadAction},
						{Object: "jobs", Action: rbac.AllAction},
					},
				}))
			},
			req: types.DeletePermissionForRoleRequest{Permission: rbac.Permission{Object: "jobs", Action: rbac.AllAction}},
			expect: func(t *testing.T, s *service, ok bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(ok)
				assert.Equal([][]string{{"developer", "clusters", rbac.ReadAction}}, s.GetRole(context.Background(), "developer"))
			},
		},
		{
			name: "missing permission",
			setup: func(t *testing.T, s *service) {
				assert := assert.New(t)
				assert.NoError(s.CreateRole(context.Background(), types.CreateRoleRequest{
					Role:        "developer",
					Permissions: []rbac.Permission{{Object: "clusters", Action: rbac.ReadAction}},
				}))
			},
			req: types.DeletePermissionForRoleRequest{Permission: rbac.Permission{Object: "jobs", Action: rbac.AllAction}},
			expect: func(t *testing.T, s *service, ok bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.False(ok)
				assert.Len(s.GetRole(context.Background(), "developer"), 1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			tc.setup(t, s)

			ok, err := s.DeletePermissionForRole(context.Background(), "developer", tc.req)
			tc.expect(t, s, ok, err)
		})
	}
}
