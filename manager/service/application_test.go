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

var mockApplicationPriorityValue = 3

func mockApplication(t *testing.T, db *gorm.DB, name string, userID uint) models.Application {
	assert := assert.New(t)
	application := models.Application{
		Name:     name,
		URL:      "https://example.com/" + name,
		BIO:      "bio",
		Priority: models.JSONMap{"value": float64(mockApplicationPriorityValue)},
		UserID:   userID,
	}
	assert.NoError(db.Create(&application).Error)

	return application
}

func TestService_CreateApplication(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		req    types.CreateApplicationRequest
		expect func(t *testing.T, s *service, application *models.Application, err error)
	}{
		{
			name:  "creates application with priority converted to a map",
			setup: func(t *testing.T, s *service) {},
			req: types.CreateApplicationRequest{
				Name: "foo",
				URL:  "https://example.com/foo",
				BIO:  "bio",
				Priority: &types.PriorityConfig{
					Value: &mockApplicationPriorityValue,
					URLs:  []types.URLPriorityConfig{{Regex: "blobs/sha256.*", Value: 5}},
				},
				UserID: 1,
			},
			expect: func(t *testing.T, s *service, application *models.Application, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", application.Name)
				assert.Equal("https://example.com/foo", application.URL)
				assert.Equal("bio", application.BIO)
				assert.Equal(float64(3), application.Priority["value"])
				assert.Equal([]any{map[string]any{"regex": "blobs/sha256.*", "value": float64(5)}}, application.Priority["urls"])
				assert.Equal(uint(1), application.UserID)

				var count int64
				assert.NoError(s.db.Model(&models.Application{}).Count(&count).Error)
				assert.Equal(int64(1), count)
			},
		},
		{
			name: "duplicate name is rejected",
			setup: func(t *testing.T, s *service) {
				mockApplication(t, s.db, "foo", 1)
			},
			req: types.CreateApplicationRequest{
				Name:     "foo",
				URL:      "https://example.com/other",
				Priority: &types.PriorityConfig{Value: &mockApplicationPriorityValue},
				UserID:   1,
			},
			expect: func(t *testing.T, s *service, application *models.Application, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(application)

				var count int64
				assert.NoError(s.db.Model(&models.Application{}).Count(&count).Error)
				assert.Equal(int64(1), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			tc.setup(t, s)

			application, err := s.CreateApplication(context.Background(), tc.req)
			tc.expect(t, s, application, err)
		})
	}
}

func TestService_GetApplication(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, application *models.Application, err error)
	}{
		{
			name: "not found",
			setup: func(t *testing.T, s *service) uint {
				return 1
			},
			expect: func(t *testing.T, application *models.Application, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(application)
			},
		},
		{
			name: "populates the user from the user id",
			setup: func(t *testing.T, s *service) uint {
				user := mockUser(t, s.db, "foo")
				return mockApplication(t, s.db, "foo", user.ID).ID
			},
			expect: func(t *testing.T, application *models.Application, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", application.Name)
				assert.Equal(float64(3), application.Priority["value"])
				assert.Equal(application.UserID, application.User.ID)
				assert.Equal("foo", application.User.Name)
				assert.Equal("foo@example.com", application.User.Email)
			},
		},
		{
			name: "application without user is returned without user",
			setup: func(t *testing.T, s *service) uint {
				return mockApplication(t, s.db, "foo", 0).ID
			},
			expect: func(t *testing.T, application *models.Application, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", application.Name)
				assert.Equal(uint(0), application.User.ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			application, err := s.GetApplication(context.Background(), id)
			tc.expect(t, application, err)
		})
	}
}

func TestService_UpdateApplication(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdateApplicationRequest
		expect func(t *testing.T, s *service, id uint, application *models.Application, err error)
	}{
		{
			name: "not found",
			setup: func(t *testing.T, s *service) uint {
				return 1
			},
			req: types.UpdateApplicationRequest{Name: "bar", UserID: 1},
			expect: func(t *testing.T, s *service, id uint, application *models.Application, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(application)
			},
		},
		{
			name: "updates fields and priority and populates the user",
			setup: func(t *testing.T, s *service) uint {
				user := mockUser(t, s.db, "foo")
				return mockApplication(t, s.db, "foo", user.ID).ID
			},
			req: types.UpdateApplicationRequest{
				Name: "bar",
				URL:  "https://example.com/bar",
				BIO:  "new bio",
				Priority: &types.PriorityConfig{
					Value: &mockApplicationPriorityValue,
					URLs:  []types.URLPriorityConfig{{Regex: "manifests/.*", Value: 1}},
				},
				UserID: 1,
			},
			expect: func(t *testing.T, s *service, id uint, application *models.Application, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", application.User.Name)

				var stored models.Application
				assert.NoError(s.db.First(&stored, id).Error)
				assert.Equal("bar", stored.Name)
				assert.Equal("https://example.com/bar", stored.URL)
				assert.Equal("new bio", stored.BIO)
				assert.Equal(float64(3), stored.Priority["value"])
				assert.Equal([]any{map[string]any{"regex": "manifests/.*", "value": float64(1)}}, stored.Priority["urls"])
				assert.Equal("foo", stored.User.Name)
			},
		},
		{
			name: "nil priority keeps the stored priority",
			setup: func(t *testing.T, s *service) uint {
				user := mockUser(t, s.db, "foo")
				return mockApplication(t, s.db, "foo", user.ID).ID
			},
			req: types.UpdateApplicationRequest{BIO: "new bio", UserID: 1},
			expect: func(t *testing.T, s *service, id uint, application *models.Application, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var stored models.Application
				assert.NoError(s.db.First(&stored, id).Error)
				assert.Equal("foo", stored.Name)
				assert.Equal("new bio", stored.BIO)
				assert.Equal(float64(3), stored.Priority["value"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			application, err := s.UpdateApplication(context.Background(), id, tc.req)
			tc.expect(t, s, id, application, err)
		})
	}
}

func TestService_GetApplications(t *testing.T) {
	tests := []struct {
		name   string
		query  types.GetApplicationsQuery
		expect func(t *testing.T, applications []models.Application, count int64, err error)
	}{
		{
			name:  "first page populates the user of every application",
			query: types.GetApplicationsQuery{Page: 1, PerPage: 2},
			expect: func(t *testing.T, applications []models.Application, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(int64(3), count)
				assert.Len(applications, 2)
				assert.Equal("app-1", applications[0].Name)
				assert.Equal("app-2", applications[1].Name)
				for _, application := range applications {
					assert.Equal(application.UserID, application.User.ID)
					assert.Equal("foo", application.User.Name)
				}
			},
		},
		{
			name:  "last page",
			query: types.GetApplicationsQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, applications []models.Application, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(int64(3), count)
				assert.Len(applications, 1)
				assert.Equal("app-3", applications[0].Name)
				assert.Equal("foo", applications[0].User.Name)
			},
		},
		{
			name:  "page beyond the last is empty",
			query: types.GetApplicationsQuery{Page: 3, PerPage: 2},
			expect: func(t *testing.T, applications []models.Application, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(int64(3), count)
				assert.Empty(applications)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			user := mockUser(t, s.db, "foo")
			for _, name := range []string{"app-1", "app-2", "app-3"} {
				mockApplication(t, s.db, name, user.ID)
			}

			applications, count, err := s.GetApplications(context.Background(), tc.query)
			tc.expect(t, applications, count, err)
		})
	}
}

func TestService_DestroyApplication(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name: "not found",
			setup: func(t *testing.T, s *service) uint {
				return 1
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "hard deletes the application",
			setup: func(t *testing.T, s *service) uint {
				mockApplication(t, s.db, "other", 0)
				return mockApplication(t, s.db, "foo", 0).ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.Application{}).Count(&count).Error)
				assert.Equal(int64(1), count)

				var remaining models.Application
				assert.NoError(s.db.First(&remaining).Error)
				assert.Equal("other", remaining.Name)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.DestroyApplication(context.Background(), id))
		})
	}
}
