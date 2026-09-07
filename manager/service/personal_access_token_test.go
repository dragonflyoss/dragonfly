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
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

func mockPersonalAccessToken(t *testing.T, db *gorm.DB, name, state string, userID uint) models.PersonalAccessToken {
	assert := assert.New(t)
	personalAccessToken := models.PersonalAccessToken{
		Name:      name,
		Token:     "token-" + name,
		Scopes:    []string{types.PersonalAccessTokenScopeJob},
		State:     state,
		ExpiredAt: time.Now().Add(time.Hour),
		UserID:    userID,
	}
	assert.NoError(db.Create(&personalAccessToken).Error)

	return personalAccessToken
}

func TestService_CreatePersonalAccessToken(t *testing.T) {
	expiredAt := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		req    types.CreatePersonalAccessTokenRequest
		expect func(t *testing.T, s *service, personalAccessToken *models.PersonalAccessToken, err error)
	}{
		{
			name: "generates token and applies default scopes",
			setup: func(t *testing.T, s *service) {
				mockUser(t, s.db, "foo")
			},
			req: types.CreatePersonalAccessTokenRequest{Name: "ci", BIO: "bio", ExpiredAt: expiredAt, UserID: 1},
			expect: func(t *testing.T, s *service, personalAccessToken *models.PersonalAccessToken, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("ci", personalAccessToken.Name)
				assert.NotEmpty(personalAccessToken.Token)
				assert.Equal(models.Array(types.DefaultPersonalAccessTokenScopes), personalAccessToken.Scopes)
				assert.Equal(models.PersonalAccessTokenStateActive, personalAccessToken.State)
				assert.True(expiredAt.Equal(personalAccessToken.ExpiredAt))
				assert.Equal(uint(1), personalAccessToken.UserID)

				another, err := s.CreatePersonalAccessToken(context.Background(), types.CreatePersonalAccessTokenRequest{Name: "other", ExpiredAt: expiredAt, UserID: 1})
				assert.NoError(err)
				assert.NotEqual(personalAccessToken.Token, another.Token)
			},
		},
		{
			name: "keeps explicit scopes",
			setup: func(t *testing.T, s *service) {
				mockUser(t, s.db, "foo")
			},
			req: types.CreatePersonalAccessTokenRequest{Name: "ci", Scopes: []string{types.PersonalAccessTokenScopeJob}, ExpiredAt: expiredAt, UserID: 1},
			expect: func(t *testing.T, s *service, personalAccessToken *models.PersonalAccessToken, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(models.Array{types.PersonalAccessTokenScopeJob}, personalAccessToken.Scopes)
			},
		},
		{
			name: "duplicate name",
			setup: func(t *testing.T, s *service) {
				mockPersonalAccessToken(t, s.db, "ci", models.PersonalAccessTokenStateActive, 1)
			},
			req: types.CreatePersonalAccessTokenRequest{Name: "ci", ExpiredAt: expiredAt, UserID: 1},
			expect: func(t *testing.T, s *service, personalAccessToken *models.PersonalAccessToken, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(personalAccessToken)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			tc.setup(t, s)

			personalAccessToken, err := s.CreatePersonalAccessToken(context.Background(), tc.req)
			tc.expect(t, s, personalAccessToken, err)
		})
	}
}

func TestService_UpdatePersonalAccessToken(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdatePersonalAccessTokenRequest
		expect func(t *testing.T, s *service, personalAccessToken *models.PersonalAccessToken, err error)
	}{
		{
			name: "token not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.UpdatePersonalAccessTokenRequest{BIO: "bio"},
			expect: func(t *testing.T, s *service, personalAccessToken *models.PersonalAccessToken, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(personalAccessToken)
			},
		},
		{
			name: "empty scopes fall back to defaults and state is persisted",
			setup: func(t *testing.T, s *service) uint {
				user := mockUser(t, s.db, "foo")
				return mockPersonalAccessToken(t, s.db, "ci", models.PersonalAccessTokenStateActive, user.ID).ID
			},
			req: types.UpdatePersonalAccessTokenRequest{State: models.PersonalAccessTokenStateInactive},
			expect: func(t *testing.T, s *service, personalAccessToken *models.PersonalAccessToken, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", personalAccessToken.User.Name)

				stored := models.PersonalAccessToken{}
				assert.NoError(s.db.First(&stored, personalAccessToken.ID).Error)
				assert.Equal(models.PersonalAccessTokenStateInactive, stored.State)
				assert.Equal(models.Array(types.DefaultPersonalAccessTokenScopes), stored.Scopes)
				assert.Equal("token-ci", stored.Token)
				assert.Equal("ci", stored.Name)
			},
		},
		{
			name: "explicit scopes replace the stored ones",
			setup: func(t *testing.T, s *service) uint {
				return mockPersonalAccessToken(t, s.db, "ci", models.PersonalAccessTokenStateActive, 0).ID
			},
			req: types.UpdatePersonalAccessTokenRequest{Scopes: []string{types.PersonalAccessTokenScopeCluster}},
			expect: func(t *testing.T, s *service, personalAccessToken *models.PersonalAccessToken, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.PersonalAccessToken{}
				assert.NoError(s.db.First(&stored, personalAccessToken.ID).Error)
				assert.Equal(models.Array{types.PersonalAccessTokenScopeCluster}, stored.Scopes)
				assert.Equal(models.PersonalAccessTokenStateActive, stored.State)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			personalAccessToken, err := s.UpdatePersonalAccessToken(context.Background(), id, tc.req)
			tc.expect(t, s, personalAccessToken, err)
		})
	}
}

func TestService_DestroyPersonalAccessToken(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name: "token not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "hard deletes",
			setup: func(t *testing.T, s *service) uint {
				return mockPersonalAccessToken(t, s.db, "ci", models.PersonalAccessTokenStateActive, 0).ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.PersonalAccessToken{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.DestroyPersonalAccessToken(context.Background(), id))
		})
	}
}

func TestService_GetPersonalAccessTokens(t *testing.T) {
	s := newTestService(t)
	foo := mockUser(t, s.db, "foo")
	bar := mockUser(t, s.db, "bar")
	mockPersonalAccessToken(t, s.db, "foo-active", models.PersonalAccessTokenStateActive, foo.ID)
	mockPersonalAccessToken(t, s.db, "foo-inactive", models.PersonalAccessTokenStateInactive, foo.ID)
	mockPersonalAccessToken(t, s.db, "bar-active", models.PersonalAccessTokenStateActive, bar.ID)

	tests := []struct {
		name   string
		query  types.GetPersonalAccessTokensQuery
		expect func(t *testing.T, personalAccessTokens []models.PersonalAccessToken, count int64, err error)
	}{
		{
			name:  "paginates and preloads users",
			query: types.GetPersonalAccessTokensQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, personalAccessTokens []models.PersonalAccessToken, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(personalAccessTokens, 1)
				assert.Equal(int64(3), count)
				assert.Equal("bar", personalAccessTokens[0].User.Name)
			},
		},
		{
			name:  "filter by user",
			query: types.GetPersonalAccessTokensQuery{UserID: foo.ID, Page: 1, PerPage: 10},
			expect: func(t *testing.T, personalAccessTokens []models.PersonalAccessToken, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(personalAccessTokens, 2)
				assert.Equal(int64(2), count)
			},
		},
		{
			name:  "filter by state and user",
			query: types.GetPersonalAccessTokensQuery{State: models.PersonalAccessTokenStateActive, UserID: foo.ID, Page: 1, PerPage: 10},
			expect: func(t *testing.T, personalAccessTokens []models.PersonalAccessToken, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(personalAccessTokens, 1)
				assert.Equal("foo-active", personalAccessTokens[0].Name)
				assert.Equal(int64(1), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			personalAccessTokens, count, err := s.GetPersonalAccessTokens(context.Background(), tc.query)
			tc.expect(t, personalAccessTokens, count, err)
		})
	}
}
