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
	"gorm.io/gorm"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

func mockOauth(t *testing.T, db *gorm.DB, name, clientID string) models.Oauth {
	assert := assert.New(t)
	oauth := models.Oauth{
		Name:         name,
		ClientID:     clientID,
		ClientSecret: "secret-" + clientID,
		RedirectURL:  "https://example.com/callback",
	}
	assert.NoError(db.Create(&oauth).Error)

	return oauth
}

func TestService_OauthSignin(t *testing.T) {
	s := newTestService(t)
	mockOauth(t, s.db, "github", "github-client")
	mockOauth(t, s.db, "google", "google-client")
	mockOauth(t, s.db, "gitlab", "gitlab-client")

	tests := []struct {
		name   string
		oauth  string
		expect func(t *testing.T, url string, err error)
	}{
		{
			name:  "oauth not found",
			oauth: "missing",
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Empty(url)
			},
		},
		{
			name:  "unsupported provider",
			oauth: "gitlab",
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(url)
			},
		},
		{
			name:  "github",
			oauth: "github",
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Contains(url, "https://github.com/login/oauth/authorize")
				assert.Contains(url, "client_id=github-client")
				assert.Contains(url, "redirect_uri=https%3A%2F%2Fexample.com%2Fcallback")
				assert.Contains(url, "state=")
			},
		},
		{
			name:  "google",
			oauth: "google",
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Contains(url, "https://accounts.google.com/o/oauth2/auth")
				assert.Contains(url, "client_id=google-client")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url, err := s.OauthSignin(context.Background(), tc.oauth)
			tc.expect(t, url, err)
		})
	}
}

func TestService_UpdateOauth(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdateOauthRequest
		expect func(t *testing.T, s *service, oauth *models.Oauth, err error)
	}{
		{
			name: "oauth not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.UpdateOauthRequest{ClientSecret: "new-secret"},
			expect: func(t *testing.T, s *service, oauth *models.Oauth, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(oauth)
			},
		},
		{
			name: "rotates secret and keeps other fields",
			setup: func(t *testing.T, s *service) uint {
				return mockOauth(t, s.db, "github", "github-client").ID
			},
			req: types.UpdateOauthRequest{ClientSecret: "new-secret", BIO: "bio"},
			expect: func(t *testing.T, s *service, oauth *models.Oauth, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.Oauth{}
				assert.NoError(s.db.First(&stored, oauth.ID).Error)
				assert.Equal("github", stored.Name)
				assert.Equal("github-client", stored.ClientID)
				assert.Equal("new-secret", stored.ClientSecret)
				assert.Equal("bio", stored.BIO)
				assert.Equal("https://example.com/callback", stored.RedirectURL)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			oauth, err := s.UpdateOauth(context.Background(), id, tc.req)
			tc.expect(t, s, oauth, err)
		})
	}
}

func TestService_DestroyOauth(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name: "oauth not found",
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
				return mockOauth(t, s.db, "github", "github-client").ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.Oauth{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.DestroyOauth(context.Background(), id))
		})
	}
}

func TestService_GetOauths(t *testing.T) {
	s := newTestService(t)
	mockOauth(t, s.db, "github", "github-client")
	mockOauth(t, s.db, "google", "google-client")

	tests := []struct {
		name   string
		query  types.GetOauthsQuery
		expect func(t *testing.T, oauths []models.Oauth, count int64, err error)
	}{
		{
			name:  "paginates",
			query: types.GetOauthsQuery{Page: 1, PerPage: 1},
			expect: func(t *testing.T, oauths []models.Oauth, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(oauths, 1)
				assert.Equal(int64(2), count)
			},
		},
		{
			name:  "filter by name",
			query: types.GetOauthsQuery{Name: "google", Page: 1, PerPage: 10},
			expect: func(t *testing.T, oauths []models.Oauth, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(oauths, 1)
				assert.Equal("google-client", oauths[0].ClientID)
				assert.Equal(int64(1), count)
			},
		},
		{
			name:  "filter by client id",
			query: types.GetOauthsQuery{ClientID: "github-client", Page: 1, PerPage: 10},
			expect: func(t *testing.T, oauths []models.Oauth, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(oauths, 1)
				assert.Equal("github", oauths[0].Name)
			},
		},
		{
			name:  "mismatched name and client id",
			query: types.GetOauthsQuery{Name: "github", ClientID: "google-client", Page: 1, PerPage: 10},
			expect: func(t *testing.T, oauths []models.Oauth, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(oauths)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			oauths, count, err := s.GetOauths(context.Background(), tc.query)
			tc.expect(t, oauths, count, err)
		})
	}
}
