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

func TestService_UpdateConfig(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdateConfigRequest
		expect func(t *testing.T, s *service, config *models.Config, err error)
	}{
		{
			name: "config not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.UpdateConfigRequest{Value: "v2"},
			expect: func(t *testing.T, s *service, config *models.Config, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(config)
			},
		},
		{
			name: "updates value and keeps name",
			setup: func(t *testing.T, s *service) uint {
				assert := assert.New(t)
				config, err := s.CreateConfig(context.Background(), types.CreateConfigRequest{Name: models.ConfigGC, Value: "v1", UserID: 1})
				assert.NoError(err)
				return config.ID
			},
			req: types.UpdateConfigRequest{Value: "v2", BIO: "bio"},
			expect: func(t *testing.T, s *service, config *models.Config, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.Config{}
				assert.NoError(s.db.First(&stored, config.ID).Error)
				assert.Equal(models.ConfigGC, stored.Name)
				assert.Equal("v2", stored.Value)
				assert.Equal("bio", stored.BIO)
				assert.Equal(uint(1), stored.UserID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			config, err := s.UpdateConfig(context.Background(), id, tc.req)
			tc.expect(t, s, config, err)
		})
	}
}

func TestService_DestroyConfig(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name: "config not found",
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
				assert := assert.New(t)
				config, err := s.CreateConfig(context.Background(), types.CreateConfigRequest{Name: models.ConfigGC, Value: "v1", UserID: 1})
				assert.NoError(err)
				return config.ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.Config{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.DestroyConfig(context.Background(), id))
		})
	}
}

func TestService_GetConfigs(t *testing.T) {
	s := newTestService(t)
	for _, req := range []types.CreateConfigRequest{
		{Name: models.ConfigGC, Value: "v1", UserID: 1},
		{Name: "feature", Value: "on", UserID: 1},
		{Name: "quota", Value: "on", UserID: 2},
	} {
		_, err := s.CreateConfig(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name   string
		query  types.GetConfigsQuery
		expect func(t *testing.T, configs []models.Config, count int64, err error)
	}{
		{
			name:  "paginates",
			query: types.GetConfigsQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, configs []models.Config, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(configs, 1)
				assert.Equal("quota", configs[0].Name)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by name",
			query: types.GetConfigsQuery{Name: models.ConfigGC, Page: 1, PerPage: 10},
			expect: func(t *testing.T, configs []models.Config, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(configs, 1)
				assert.Equal("v1", configs[0].Value)
				assert.Equal(int64(1), count)
			},
		},
		{
			name:  "filter by value and user",
			query: types.GetConfigsQuery{Value: "on", UserID: 1, Page: 1, PerPage: 10},
			expect: func(t *testing.T, configs []models.Config, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(configs, 1)
				assert.Equal("feature", configs[0].Name)
				assert.Equal(int64(1), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configs, count, err := s.GetConfigs(context.Background(), tc.query)
			tc.expect(t, configs, count, err)
		})
	}
}
