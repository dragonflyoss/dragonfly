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
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootPassword(t *testing.T) {
	tests := []struct {
		name   string
		env    string
		set    bool
		expect func(t *testing.T, password string, explicit bool, err error)
	}{
		{
			name: "environment variable is not set",
			set:  false,
			expect: func(t *testing.T, password string, explicit bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(DefaultRootPassword, password)
				assert.False(explicit)
			},
		},
		{
			name: "environment variable is empty",
			env:  "",
			set:  true,
			expect: func(t *testing.T, password string, explicit bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(DefaultRootPassword, password)
				assert.False(explicit)
			},
		},
		{
			name: "environment variable is set",
			env:  "dragonfly-root",
			set:  true,
			expect: func(t *testing.T, password string, explicit bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("dragonfly-root", password)
				assert.True(explicit)
			},
		},
		{
			name: "environment variable is too short",
			env:  strings.Repeat("a", MinRootPasswordLength-1),
			set:  true,
			expect: func(t *testing.T, password string, explicit bool, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "DRAGONFLY_ROOT_PASSWORD must be between 8 and 20 characters")
			},
		},
		{
			name: "environment variable is too long",
			env:  strings.Repeat("a", MaxRootPasswordLength+1),
			set:  true,
			expect: func(t *testing.T, password string, explicit bool, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "DRAGONFLY_ROOT_PASSWORD must be between 8 and 20 characters")
			},
		},
		{
			name: "environment variable is of the minimum length",
			env:  strings.Repeat("a", MinRootPasswordLength),
			set:  true,
			expect: func(t *testing.T, password string, explicit bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(password, MinRootPasswordLength)
				assert.True(explicit)
			},
		},
		{
			name: "environment variable is of the maximum length",
			env:  strings.Repeat("a", MaxRootPasswordLength),
			set:  true,
			expect: func(t *testing.T, password string, explicit bool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(password, MaxRootPasswordLength)
				assert.True(explicit)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.set {
				t.Setenv(DragonflyRootPasswordEnvName, tc.env)
			} else {
				os.Unsetenv(DragonflyRootPasswordEnvName)
			}

			password, explicit, err := rootPassword()
			tc.expect(t, password, explicit, err)
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
