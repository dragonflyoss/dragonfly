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

func TestService_CreateScheduler(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		req    types.CreateSchedulerRequest
		expect func(t *testing.T, s *service, scheduler *models.Scheduler, err error)
	}{
		{
			name: "nil features fall back to defaults",
			setup: func(t *testing.T, s *service) {
				mockSchedulerCluster(t, s.db, "foo")
			},
			req: types.CreateSchedulerRequest{Hostname: "scheduler-1", IDC: "idc", IP: "127.0.0.1", Port: 8002, SchedulerClusterID: 1},
			expect: func(t *testing.T, s *service, scheduler *models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(models.Array(types.DefaultSchedulerFeatures), scheduler.Features)
				assert.Nil(scheduler.Config)

				stored := models.Scheduler{}
				assert.NoError(s.db.First(&stored, scheduler.ID).Error)
				assert.Equal("scheduler-1", stored.Hostname)
				assert.Equal(models.SchedulerStateInactive, stored.State)
				assert.Equal(models.Array(types.DefaultSchedulerFeatures), stored.Features)
				assert.Equal(uint(1), stored.SchedulerClusterID)
			},
		},
		{
			name: "explicit features and config are stored",
			setup: func(t *testing.T, s *service) {
				mockSchedulerCluster(t, s.db, "foo")
			},
			req: types.CreateSchedulerRequest{
				Hostname:           "scheduler-1",
				IP:                 "127.0.0.1",
				Port:               8002,
				Features:           []string{types.SchedulerFeatureSchedule},
				Config:             &types.SchedulerConfig{ManagerKeepAliveInterval: 5 * time.Second},
				SchedulerClusterID: 1,
			},
			expect: func(t *testing.T, s *service, scheduler *models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(models.Array{types.SchedulerFeatureSchedule}, scheduler.Features)
				assert.Equal(float64(5*time.Second), scheduler.Config["manager_keep_alive_interval"])
			},
		},
		{
			name: "empty features are kept empty",
			setup: func(t *testing.T, s *service) {
				mockSchedulerCluster(t, s.db, "foo")
			},
			req: types.CreateSchedulerRequest{Hostname: "scheduler-1", IP: "127.0.0.1", Port: 8002, Features: []string{}, SchedulerClusterID: 1},
			expect: func(t *testing.T, s *service, scheduler *models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(scheduler.Features)
			},
		},
		{
			name: "duplicate hostname ip and cluster",
			setup: func(t *testing.T, s *service) {
				schedulerCluster := mockSchedulerCluster(t, s.db, "foo")
				mockScheduler(t, s.db, schedulerCluster.ID, "scheduler-1", models.SchedulerStateActive, nil)
			},
			req: types.CreateSchedulerRequest{Hostname: "scheduler-1", IP: "127.0.0.1", Port: 8002, SchedulerClusterID: 1},
			expect: func(t *testing.T, s *service, scheduler *models.Scheduler, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(scheduler)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			tc.setup(t, s)

			scheduler, err := s.CreateScheduler(context.Background(), tc.req)
			tc.expect(t, s, scheduler, err)
		})
	}
}

func TestService_DestroyScheduler(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name: "scheduler not found",
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
				schedulerCluster := mockSchedulerCluster(t, s.db, "foo")
				return mockScheduler(t, s.db, schedulerCluster.ID, "scheduler-1", models.SchedulerStateActive, nil).ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.Scheduler{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.DestroyScheduler(context.Background(), id))
		})
	}
}

func TestService_GetSchedulers(t *testing.T) {
	s := mockService(t)
	foo := mockSchedulerCluster(t, s.db, "foo")
	bar := mockSchedulerCluster(t, s.db, "bar")
	mockScheduler(t, s.db, foo.ID, "scheduler-1", models.SchedulerStateActive, nil)
	mockScheduler(t, s.db, foo.ID, "scheduler-2", models.SchedulerStateInactive, nil)
	mockScheduler(t, s.db, bar.ID, "scheduler-3", models.SchedulerStateActive, nil)

	tests := []struct {
		name   string
		query  types.GetSchedulersQuery
		expect func(t *testing.T, schedulers []models.Scheduler, count int64, err error)
	}{
		{
			name:  "paginates",
			query: types.GetSchedulersQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, schedulers []models.Scheduler, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 1)
				assert.Equal("scheduler-3", schedulers[0].Hostname)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by cluster",
			query: types.GetSchedulersQuery{SchedulerClusterID: foo.ID, Page: 1, PerPage: 10},
			expect: func(t *testing.T, schedulers []models.Scheduler, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 2)
				assert.Equal(int64(2), count)
			},
		},
		{
			name:  "filter by cluster and state",
			query: types.GetSchedulersQuery{SchedulerClusterID: foo.ID, State: models.SchedulerStateActive, Page: 1, PerPage: 10},
			expect: func(t *testing.T, schedulers []models.Scheduler, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 1)
				assert.Equal("scheduler-1", schedulers[0].Hostname)
				assert.Equal(int64(1), count)
			},
		},
		{
			name:  "filter by idc",
			query: types.GetSchedulersQuery{IDC: "idc-scheduler-2", Page: 1, PerPage: 10},
			expect: func(t *testing.T, schedulers []models.Scheduler, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 1)
				assert.Equal("scheduler-2", schedulers[0].Hostname)
			},
		},
		{
			name:  "no match",
			query: types.GetSchedulersQuery{Hostname: "missing", Page: 1, PerPage: 10},
			expect: func(t *testing.T, schedulers []models.Scheduler, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(schedulers)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			schedulers, count, err := s.GetSchedulers(context.Background(), tc.query)
			tc.expect(t, schedulers, count, err)
		})
	}
}
