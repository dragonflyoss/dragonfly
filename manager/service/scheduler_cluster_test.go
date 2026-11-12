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

func TestService_CreateSchedulerCluster(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		req    types.CreateSchedulerClusterRequest
		expect func(t *testing.T, s *service, schedulerCluster *models.SchedulerCluster, err error)
	}{
		{
			name:  "creates cluster with json configs",
			setup: func(t *testing.T, s *service) {},
			req: types.CreateSchedulerClusterRequest{
				Name:             "foo",
				BIO:              "bio",
				Config:           &types.SchedulerClusterConfig{CandidateParentLimit: 4, FilterParentLimit: 40},
				ClientConfig:     &types.SchedulerClusterClientConfig{LoadLimit: 50},
				SeedClientConfig: &types.SchedulerClusterSeedClientConfig{},
				Scopes:           &types.SchedulerClusterScopes{IDC: "idc"},
				IsDefault:        true,
			},
			expect: func(t *testing.T, s *service, schedulerCluster *models.SchedulerCluster, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", schedulerCluster.Name)
				assert.True(schedulerCluster.IsDefault)
				assert.Equal(float64(4), schedulerCluster.Config["candidate_parent_limit"])
				assert.Equal(float64(50), schedulerCluster.ClientConfig["load_limit"])
				assert.Equal("idc", schedulerCluster.Scopes["idc"])

				stored := models.SchedulerCluster{}
				assert.NoError(s.db.Preload("SeedPeerClusters").First(&stored, schedulerCluster.ID).Error)
				assert.Empty(stored.SeedPeerClusters)
			},
		},
		{
			name: "links to the given seed peer cluster",
			setup: func(t *testing.T, s *service) {
				mockSeedPeerCluster(t, s.db, "seed")
			},
			req: types.CreateSchedulerClusterRequest{
				Name:              "foo",
				Config:            &types.SchedulerClusterConfig{},
				ClientConfig:      &types.SchedulerClusterClientConfig{},
				SeedClientConfig:  &types.SchedulerClusterSeedClientConfig{},
				SeedPeerClusterID: 1,
			},
			expect: func(t *testing.T, s *service, schedulerCluster *models.SchedulerCluster, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.SchedulerCluster{}
				assert.NoError(s.db.Preload("SeedPeerClusters").First(&stored, schedulerCluster.ID).Error)
				assert.Len(stored.SeedPeerClusters, 1)
				assert.Equal("seed", stored.SeedPeerClusters[0].Name)
			},
		},
		{
			name:  "seed peer cluster not found",
			setup: func(t *testing.T, s *service) {},
			req: types.CreateSchedulerClusterRequest{
				Name:              "foo",
				Config:            &types.SchedulerClusterConfig{},
				ClientConfig:      &types.SchedulerClusterClientConfig{},
				SeedClientConfig:  &types.SchedulerClusterSeedClientConfig{},
				SeedPeerClusterID: 99,
			},
			expect: func(t *testing.T, s *service, schedulerCluster *models.SchedulerCluster, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(schedulerCluster)
			},
		},
		{
			name: "duplicate name",
			setup: func(t *testing.T, s *service) {
				mockSchedulerCluster(t, s.db, "foo")
			},
			req: types.CreateSchedulerClusterRequest{
				Name:             "foo",
				Config:           &types.SchedulerClusterConfig{},
				ClientConfig:     &types.SchedulerClusterClientConfig{},
				SeedClientConfig: &types.SchedulerClusterSeedClientConfig{},
			},
			expect: func(t *testing.T, s *service, schedulerCluster *models.SchedulerCluster, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(schedulerCluster)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			tc.setup(t, s)

			schedulerCluster, err := s.CreateSchedulerCluster(context.Background(), tc.req)
			tc.expect(t, s, schedulerCluster, err)
		})
	}
}

func TestService_DestroySchedulerCluster(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name: "cluster not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "cluster still has schedulers",
			setup: func(t *testing.T, s *service) uint {
				schedulerCluster := mockSchedulerCluster(t, s.db, "foo")
				mockScheduler(t, s.db, schedulerCluster.ID, "scheduler-1", models.SchedulerStateActive, nil)
				return schedulerCluster.ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.Error(err)

				var count int64
				assert.NoError(s.db.Model(&models.SchedulerCluster{}).Count(&count).Error)
				assert.Equal(int64(1), count)
			},
		},
		{
			name: "hard deletes empty cluster",
			setup: func(t *testing.T, s *service) uint {
				return mockSchedulerCluster(t, s.db, "foo").ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.SchedulerCluster{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.DestroySchedulerCluster(context.Background(), id))
		})
	}
}

func TestService_UpdateSchedulerCluster(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdateSchedulerClusterRequest
		expect func(t *testing.T, s *service, id uint, schedulerCluster *models.SchedulerCluster, err error)
	}{
		{
			name: "cluster not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.UpdateSchedulerClusterRequest{Name: "bar"},
			expect: func(t *testing.T, s *service, id uint, schedulerCluster *models.SchedulerCluster, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(schedulerCluster)
			},
		},
		{
			name: "unsets is_default",
			setup: func(t *testing.T, s *service) uint {
				assert := assert.New(t)
				schedulerCluster := mockSchedulerCluster(t, s.db, "foo")
				assert.NoError(s.db.Model(&schedulerCluster).Update("is_default", true).Error)
				return schedulerCluster.ID
			},
			req: types.UpdateSchedulerClusterRequest{IsDefault: false},
			expect: func(t *testing.T, s *service, id uint, schedulerCluster *models.SchedulerCluster, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.SchedulerCluster{}
				assert.NoError(s.db.First(&stored, id).Error)
				assert.False(stored.IsDefault)
				assert.Equal("foo", stored.Name)
			},
		},
		{
			name: "sets is_default",
			setup: func(t *testing.T, s *service) uint {
				return mockSchedulerCluster(t, s.db, "foo").ID
			},
			req: types.UpdateSchedulerClusterRequest{IsDefault: true},
			expect: func(t *testing.T, s *service, id uint, schedulerCluster *models.SchedulerCluster, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.SchedulerCluster{}
				assert.NoError(s.db.First(&stored, id).Error)
				assert.True(stored.IsDefault)
			},
		},
		{
			name: "nil config keeps stored config",
			setup: func(t *testing.T, s *service) uint {
				return mockSchedulerCluster(t, s.db, "foo").ID
			},
			req: types.UpdateSchedulerClusterRequest{Name: "bar", ClientConfig: &types.SchedulerClusterClientConfig{LoadLimit: 80}},
			expect: func(t *testing.T, s *service, id uint, schedulerCluster *models.SchedulerCluster, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.SchedulerCluster{}
				assert.NoError(s.db.First(&stored, id).Error)
				assert.Equal("bar", stored.Name)
				assert.Equal(float64(4), stored.Config["candidate_parent_limit"])
				assert.Equal(float64(80), stored.ClientConfig["load_limit"])
			},
		},
		{
			name: "relinks to another seed peer cluster",
			setup: func(t *testing.T, s *service) uint {
				assert := assert.New(t)
				schedulerCluster := mockSchedulerCluster(t, s.db, "foo")
				seedPeerCluster := mockSeedPeerCluster(t, s.db, "seed-1")
				mockSeedPeerCluster(t, s.db, "seed-2")
				assert.NoError(s.db.Model(&seedPeerCluster).Association("SchedulerClusters").Append(&schedulerCluster))
				return schedulerCluster.ID
			},
			req: types.UpdateSchedulerClusterRequest{SeedPeerClusterID: 2},
			expect: func(t *testing.T, s *service, id uint, schedulerCluster *models.SchedulerCluster, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.SchedulerCluster{}
				assert.NoError(s.db.Preload("SeedPeerClusters").First(&stored, id).Error)
				assert.Len(stored.SeedPeerClusters, 1)
				assert.Equal("seed-2", stored.SeedPeerClusters[0].Name)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			schedulerCluster, err := s.UpdateSchedulerCluster(context.Background(), id, tc.req)
			tc.expect(t, s, id, schedulerCluster, err)
		})
	}
}

func TestService_GetSchedulerClusters(t *testing.T) {
	s := mockService(t)
	seedPeerCluster := mockSeedPeerCluster(t, s.db, "seed")
	for _, name := range []string{"foo", "bar", "baz"} {
		schedulerCluster := mockSchedulerCluster(t, s.db, name)
		if err := s.db.Model(&seedPeerCluster).Association("SchedulerClusters").Append(&schedulerCluster); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name   string
		query  types.GetSchedulerClustersQuery
		expect func(t *testing.T, schedulerClusters []models.SchedulerCluster, count int64, err error)
	}{
		{
			name:  "first page preloads seed peer clusters",
			query: types.GetSchedulerClustersQuery{Page: 1, PerPage: 2},
			expect: func(t *testing.T, schedulerClusters []models.SchedulerCluster, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulerClusters, 2)
				assert.Equal(int64(3), count)
				assert.Len(schedulerClusters[0].SeedPeerClusters, 1)
				assert.Equal("seed", schedulerClusters[0].SeedPeerClusters[0].Name)
			},
		},
		{
			name:  "last page",
			query: types.GetSchedulerClustersQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, schedulerClusters []models.SchedulerCluster, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulerClusters, 1)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by name",
			query: types.GetSchedulerClustersQuery{Name: "baz", Page: 1, PerPage: 10},
			expect: func(t *testing.T, schedulerClusters []models.SchedulerCluster, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulerClusters, 1)
				assert.Equal("baz", schedulerClusters[0].Name)
				assert.Equal(int64(1), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			schedulerClusters, count, err := s.GetSchedulerClusters(context.Background(), tc.query)
			tc.expect(t, schedulerClusters, count, err)
		})
	}
}

func TestService_AddSchedulerToSchedulerCluster(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T, s *service)
		id          uint
		schedulerID uint
		expect      func(t *testing.T, s *service, err error)
	}{
		{
			name: "scheduler cluster not found",
			setup: func(t *testing.T, s *service) {
				schedulerCluster := mockSchedulerCluster(t, s.db, "foo")
				mockScheduler(t, s.db, schedulerCluster.ID, "scheduler-1", models.SchedulerStateActive, nil)
			},
			id:          99,
			schedulerID: 1,
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "scheduler not found",
			setup: func(t *testing.T, s *service) {
				mockSchedulerCluster(t, s.db, "foo")
			},
			id:          1,
			schedulerID: 99,
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "moves scheduler into the cluster",
			setup: func(t *testing.T, s *service) {
				schedulerCluster := mockSchedulerCluster(t, s.db, "foo")
				mockSchedulerCluster(t, s.db, "bar")
				mockScheduler(t, s.db, schedulerCluster.ID, "scheduler-1", models.SchedulerStateActive, nil)
			},
			id:          2,
			schedulerID: 1,
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				scheduler := models.Scheduler{}
				assert.NoError(s.db.First(&scheduler, 1).Error)
				assert.Equal(uint(2), scheduler.SchedulerClusterID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			tc.setup(t, s)

			tc.expect(t, s, s.AddSchedulerToSchedulerCluster(context.Background(), tc.id, tc.schedulerID))
		})
	}
}
