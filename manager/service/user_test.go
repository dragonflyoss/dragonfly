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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/permission/rbac"
	"d7y.io/dragonfly/v2/manager/types"
)

func mockSignUpRequest(name, email string) types.SignUpRequest {
	return types.SignUpRequest{
		SignInRequest: types.SignInRequest{Name: name, Password: "dragonfly"},
		Email:         email,
		Location:      "hangzhou",
	}
}

func TestService_SignUp(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		req    types.SignUpRequest
		expect func(t *testing.T, s *service, user *models.User, err error)
	}{
		{
			name:  "hashes password and grants guest role",
			setup: func(t *testing.T, s *service) {},
			req:   mockSignUpRequest("foo", "foo@example.com"),
			expect: func(t *testing.T, s *service, user *models.User, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", user.Name)
				assert.Equal("foo@example.com", user.Email)
				assert.Equal("hangzhou", user.Location)
				assert.Equal(models.UserStateEnabled, user.State)
				assert.NotEqual("dragonfly", user.EncryptedPassword)
				assert.NoError(bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte("dragonfly")))

				roles, err := s.enforcer.GetRolesForUser(fmt.Sprint(user.ID))
				assert.NoError(err)
				assert.Equal([]string{rbac.GuestRole}, roles)
			},
		},
		{
			name: "duplicate name",
			setup: func(t *testing.T, s *service) {
				mockUser(t, s.db, "foo")
			},
			req: mockSignUpRequest("foo", "other@example.com"),
			expect: func(t *testing.T, s *service, user *models.User, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(user)
			},
		},
		{
			name: "duplicate email",
			setup: func(t *testing.T, s *service) {
				mockUser(t, s.db, "foo")
			},
			req: mockSignUpRequest("bar", "foo@example.com"),
			expect: func(t *testing.T, s *service, user *models.User, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(user)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			tc.setup(t, s)

			user, err := s.SignUp(context.Background(), tc.req)
			tc.expect(t, s, user, err)
		})
	}
}

func TestService_SignIn(t *testing.T) {
	s := mockService(t)
	_, err := s.SignUp(context.Background(), mockSignUpRequest("foo", "foo@example.com"))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		req    types.SignInRequest
		expect func(t *testing.T, user *models.User, err error)
	}{
		{
			name: "correct password",
			req:  types.SignInRequest{Name: "foo", Password: "dragonfly"},
			expect: func(t *testing.T, user *models.User, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", user.Name)
			},
		},
		{
			name: "wrong password",
			req:  types.SignInRequest{Name: "foo", Password: "wrong-password"},
			expect: func(t *testing.T, user *models.User, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, bcrypt.ErrMismatchedHashAndPassword)
				assert.Nil(user)
			},
		},
		{
			name: "unknown user",
			req:  types.SignInRequest{Name: "bar", Password: "dragonfly"},
			expect: func(t *testing.T, user *models.User, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(user)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user, err := s.SignIn(context.Background(), tc.req)
			tc.expect(t, user, err)
		})
	}
}

func TestService_ResetPassword(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.ResetPasswordRequest
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name: "user not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.ResetPasswordRequest{OldPassword: "dragonfly", NewPassword: "new-password"},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "wrong old password keeps the current one",
			setup: func(t *testing.T, s *service) uint {
				assert := assert.New(t)
				user, err := s.SignUp(context.Background(), mockSignUpRequest("foo", "foo@example.com"))
				assert.NoError(err)
				return user.ID
			},
			req: types.ResetPasswordRequest{OldPassword: "wrong-password", NewPassword: "new-password"},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, bcrypt.ErrMismatchedHashAndPassword)

				_, err = s.SignIn(context.Background(), types.SignInRequest{Name: "foo", Password: "dragonfly"})
				assert.NoError(err)
			},
		},
		{
			name: "correct old password stores the new one",
			setup: func(t *testing.T, s *service) uint {
				assert := assert.New(t)
				user, err := s.SignUp(context.Background(), mockSignUpRequest("foo", "foo@example.com"))
				assert.NoError(err)
				return user.ID
			},
			req: types.ResetPasswordRequest{OldPassword: "dragonfly", NewPassword: "new-password"},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				_, err = s.SignIn(context.Background(), types.SignInRequest{Name: "foo", Password: "new-password"})
				assert.NoError(err)

				_, err = s.SignIn(context.Background(), types.SignInRequest{Name: "foo", Password: "dragonfly"})
				assert.ErrorIs(err, bcrypt.ErrMismatchedHashAndPassword)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.ResetPassword(context.Background(), id, tc.req))
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdateUserRequest
		expect func(t *testing.T, s *service, user *models.User, err error)
	}{
		{
			name: "user not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.UpdateUserRequest{BIO: "bio"},
			expect: func(t *testing.T, s *service, user *models.User, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(user)
			},
		},
		{
			name: "updates profile fields only",
			setup: func(t *testing.T, s *service) uint {
				return mockUser(t, s.db, "foo").ID
			},
			req: types.UpdateUserRequest{Email: "new@example.com", Phone: "123", Location: "beijing", BIO: "bio"},
			expect: func(t *testing.T, s *service, user *models.User, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.User{}
				assert.NoError(s.db.First(&stored, user.ID).Error)
				assert.Equal("foo", stored.Name)
				assert.Equal("new@example.com", stored.Email)
				assert.Equal("123", stored.Phone)
				assert.Equal("beijing", stored.Location)
				assert.Equal("bio", stored.BIO)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			user, err := s.UpdateUser(context.Background(), id, tc.req)
			tc.expect(t, s, user, err)
		})
	}
}

func TestService_GetUsers(t *testing.T) {
	s := mockService(t)
	for _, user := range []models.User{
		{Name: "foo", Email: "foo@example.com", Location: "hangzhou", State: models.UserStateEnabled},
		{Name: "bar", Email: "bar@example.com", Location: "hangzhou", State: models.UserStateDisabled},
		{Name: "baz", Email: "baz@example.com", Location: "beijing", State: models.UserStateEnabled},
	} {
		if err := s.db.Create(&user).Error; err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name   string
		query  types.GetUsersQuery
		expect func(t *testing.T, users []models.User, count int64, err error)
	}{
		{
			name:  "paginates",
			query: types.GetUsersQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, users []models.User, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(users, 1)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by location",
			query: types.GetUsersQuery{Location: "hangzhou", Page: 1, PerPage: 10},
			expect: func(t *testing.T, users []models.User, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(users, 2)
				assert.Equal(int64(2), count)
			},
		},
		{
			name:  "filter by state and location",
			query: types.GetUsersQuery{State: models.UserStateEnabled, Location: "hangzhou", Page: 1, PerPage: 10},
			expect: func(t *testing.T, users []models.User, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(users, 1)
				assert.Equal("foo", users[0].Name)
				assert.Equal(int64(1), count)
			},
		},
		{
			name:  "filter by email",
			query: types.GetUsersQuery{Email: "baz@example.com", Page: 1, PerPage: 10},
			expect: func(t *testing.T, users []models.User, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(users, 1)
				assert.Equal("baz", users[0].Name)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			users, count, err := s.GetUsers(context.Background(), tc.query)
			tc.expect(t, users, count, err)
		})
	}
}
