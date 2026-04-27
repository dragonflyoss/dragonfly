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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

func TestService_CreateApplication(t *testing.T) {
	s := newTestApplicationService(t)
	priorityValue := 1

	application, err := s.CreateApplication(context.Background(), types.CreateApplicationRequest{
		Name: "test-application",
		URL:  "https://example.com/file.tar.gz",
		BIO:  "test bio",
		Priority: &types.PriorityConfig{
			Value: &priorityValue,
		},
		UserID: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, "test-application", application.Name)
	assert.Equal(t, "https://example.com/file.tar.gz", application.URL)
	assert.Equal(t, "test bio", application.BIO)
	assert.Equal(t, float64(priorityValue), application.Priority["value"])
	assert.Equal(t, uint(1), application.UserID)

	var persisted models.Application
	err = s.db.First(&persisted, application.ID).Error
	assert.NoError(t, err)
}

func TestService_DestroyApplication(t *testing.T) {
	s := newTestApplicationService(t)

	application := models.Application{
		Name:   "destroy-application",
		URL:    "https://example.com/destroy.tar.gz",
		UserID: 1,
		Priority: map[string]any{
			"value": float64(1),
		},
	}
	assert.NoError(t, s.db.Create(&application).Error)

	assert.NoError(t, s.DestroyApplication(context.Background(), application.ID))

	var count int64
	err := s.db.Unscoped().Model(&models.Application{}).Where("id = ?", application.ID).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestService_UpdateApplication(t *testing.T) {
	s := newTestApplicationService(t)
	oldPriorityValue := 1
	newPriorityValue := 2

	application, err := s.CreateApplication(context.Background(), types.CreateApplicationRequest{
		Name: "update-application",
		URL:  "https://example.com/old.tar.gz",
		BIO:  "old bio",
		Priority: &types.PriorityConfig{
			Value: &oldPriorityValue,
		},
		UserID: 1,
	})
	assert.NoError(t, err)

	_, err = s.UpdateApplication(context.Background(), application.ID, types.UpdateApplicationRequest{
		Name: "updated-application",
		URL:  "https://example.com/new.tar.gz",
		BIO:  "new bio",
		Priority: &types.PriorityConfig{
			Value: &newPriorityValue,
		},
		UserID: 1,
	})
	assert.NoError(t, err)

	updated, err := s.GetApplication(context.Background(), application.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updated-application", updated.Name)
	assert.Equal(t, "https://example.com/new.tar.gz", updated.URL)
	assert.Equal(t, "new bio", updated.BIO)
	assert.Equal(t, float64(newPriorityValue), updated.Priority["value"])
}

func TestService_GetApplication(t *testing.T) {
	s := newTestApplicationService(t)
	priorityValue := 1

	application, err := s.CreateApplication(context.Background(), types.CreateApplicationRequest{
		Name: "get-application",
		URL:  "https://example.com/get.tar.gz",
		BIO:  "get bio",
		Priority: &types.PriorityConfig{
			Value: &priorityValue,
		},
		UserID: 1,
	})
	assert.NoError(t, err)

	got, err := s.GetApplication(context.Background(), application.ID)
	assert.NoError(t, err)
	assert.Equal(t, application.ID, got.ID)
	assert.Equal(t, uint(1), got.User.ID)
	assert.Equal(t, "test-user", got.User.Name)
}

func TestService_GetApplications(t *testing.T) {
	s := newTestApplicationService(t)
	priorityValue := 1

	for i := 1; i <= 2; i++ {
		_, err := s.CreateApplication(context.Background(), types.CreateApplicationRequest{
			Name: fmt.Sprintf("list-application-%d", i),
			URL:  "https://example.com/list.tar.gz",
			Priority: &types.PriorityConfig{
				Value: &priorityValue,
			},
			UserID: 1,
		})
		assert.NoError(t, err)
	}

	applications, count, err := s.GetApplications(context.Background(), types.GetApplicationsQuery{
		Page:    1,
		PerPage: 1,
	})
	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	assert.Equal(t, int64(2), count)
}

func newTestApplicationService(t *testing.T) *service {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&models.User{}, &models.Application{}))
	assert.NoError(t, db.Create(&models.User{
		BaseModel: models.BaseModel{
			ID: 1,
		},
		Email: "test@example.com",
		Name:  "test-user",
		State: models.UserStateEnabled,
	}).Error)

	return &service{db: db}
}
