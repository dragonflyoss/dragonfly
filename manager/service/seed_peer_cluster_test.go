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

func TestService_DestroySeedPeerCluster(t *testing.T) {
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
			name: "cluster still has seed peers",
			setup: func(t *testing.T, s *service) uint {
				seedPeerCluster := mockSeedPeerCluster(t, s.db, "foo")
				mockSeedPeer(t, s.db, seedPeerCluster.ID, "seed-peer-1")
				return seedPeerCluster.ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.Error(err)

				var count int64
				assert.NoError(s.db.Model(&models.SeedPeerCluster{}).Count(&count).Error)
				assert.Equal(int64(1), count)
			},
		},
		{
			name: "hard deletes empty cluster",
			setup: func(t *testing.T, s *service) uint {
				return mockSeedPeerCluster(t, s.db, "foo").ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.SeedPeerCluster{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.DestroySeedPeerCluster(context.Background(), id))
		})
	}
}

func TestService_UpdateSeedPeerCluster(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdateSeedPeerClusterRequest
		expect func(t *testing.T, s *service, id uint, seedPeerCluster *models.SeedPeerCluster, err error)
	}{
		{
			name: "cluster not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.UpdateSeedPeerClusterRequest{Name: "bar"},
			expect: func(t *testing.T, s *service, id uint, seedPeerCluster *models.SeedPeerCluster, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(seedPeerCluster)
			},
		},
		{
			name: "nil config keeps stored config",
			setup: func(t *testing.T, s *service) uint {
				return mockSeedPeerCluster(t, s.db, "foo").ID
			},
			req: types.UpdateSeedPeerClusterRequest{Name: "bar", BIO: "bio"},
			expect: func(t *testing.T, s *service, id uint, seedPeerCluster *models.SeedPeerCluster, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.SeedPeerCluster{}
				assert.NoError(s.db.First(&stored, id).Error)
				assert.Equal("bar", stored.Name)
				assert.Equal("bio", stored.BIO)
				assert.Equal(float64(300), stored.Config["load_limit"])
			},
		},
		{
			name: "replaces config",
			setup: func(t *testing.T, s *service) uint {
				return mockSeedPeerCluster(t, s.db, "foo").ID
			},
			req: types.UpdateSeedPeerClusterRequest{Config: &types.SeedPeerClusterConfig{LoadLimit: 100}},
			expect: func(t *testing.T, s *service, id uint, seedPeerCluster *models.SeedPeerCluster, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(float64(100), seedPeerCluster.Config["load_limit"])

				stored := models.SeedPeerCluster{}
				assert.NoError(s.db.First(&stored, id).Error)
				assert.Equal("foo", stored.Name)
				assert.Equal(float64(100), stored.Config["load_limit"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			seedPeerCluster, err := s.UpdateSeedPeerCluster(context.Background(), id, tc.req)
			tc.expect(t, s, id, seedPeerCluster, err)
		})
	}
}

func TestService_GetSeedPeerClusters(t *testing.T) {
	s := newTestService(t)
	for _, name := range []string{"foo", "bar", "baz"} {
		mockSeedPeerCluster(t, s.db, name)
	}

	tests := []struct {
		name   string
		query  types.GetSeedPeerClustersQuery
		expect func(t *testing.T, seedPeerClusters []models.SeedPeerCluster, count int64, err error)
	}{
		{
			name:  "first page",
			query: types.GetSeedPeerClustersQuery{Page: 1, PerPage: 2},
			expect: func(t *testing.T, seedPeerClusters []models.SeedPeerCluster, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(seedPeerClusters, 2)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "last page",
			query: types.GetSeedPeerClustersQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, seedPeerClusters []models.SeedPeerCluster, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(seedPeerClusters, 1)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by name",
			query: types.GetSeedPeerClustersQuery{Name: "bar", Page: 1, PerPage: 10},
			expect: func(t *testing.T, seedPeerClusters []models.SeedPeerCluster, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(seedPeerClusters, 1)
				assert.Equal("bar", seedPeerClusters[0].Name)
				assert.Equal(int64(1), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			seedPeerClusters, count, err := s.GetSeedPeerClusters(context.Background(), tc.query)
			tc.expect(t, seedPeerClusters, count, err)
		})
	}
}

func TestService_AddSeedPeerToSeedPeerCluster(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, s *service)
		id         uint
		seedPeerID uint
		expect     func(t *testing.T, s *service, err error)
	}{
		{
			name: "seed peer cluster not found",
			setup: func(t *testing.T, s *service) {
				seedPeerCluster := mockSeedPeerCluster(t, s.db, "foo")
				mockSeedPeer(t, s.db, seedPeerCluster.ID, "seed-peer-1")
			},
			id:         99,
			seedPeerID: 1,
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "seed peer not found",
			setup: func(t *testing.T, s *service) {
				mockSeedPeerCluster(t, s.db, "foo")
			},
			id:         1,
			seedPeerID: 99,
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "moves seed peer into the cluster",
			setup: func(t *testing.T, s *service) {
				seedPeerCluster := mockSeedPeerCluster(t, s.db, "foo")
				mockSeedPeerCluster(t, s.db, "bar")
				mockSeedPeer(t, s.db, seedPeerCluster.ID, "seed-peer-1")
			},
			id:         2,
			seedPeerID: 1,
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				seedPeer := models.SeedPeer{}
				assert.NoError(s.db.First(&seedPeer, 1).Error)
				assert.Equal(uint(2), seedPeer.SeedPeerClusterID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			tc.setup(t, s)

			tc.expect(t, s, s.AddSeedPeerToSeedPeerCluster(context.Background(), tc.id, tc.seedPeerID))
		})
	}
}

func TestService_AddSchedulerClusterToSeedPeerCluster(t *testing.T) {
	tests := []struct {
		name               string
		setup              func(t *testing.T, s *service)
		id                 uint
		schedulerClusterID uint
		expect             func(t *testing.T, s *service, err error)
	}{
		{
			name: "seed peer cluster not found",
			setup: func(t *testing.T, s *service) {
				mockSchedulerCluster(t, s.db, "foo")
			},
			id:                 99,
			schedulerClusterID: 1,
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "scheduler cluster not found",
			setup: func(t *testing.T, s *service) {
				mockSeedPeerCluster(t, s.db, "seed")
			},
			id:                 1,
			schedulerClusterID: 99,
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "replaces the existing link",
			setup: func(t *testing.T, s *service) {
				assert := assert.New(t)
				schedulerCluster := mockSchedulerCluster(t, s.db, "foo")
				seedPeerCluster := mockSeedPeerCluster(t, s.db, "seed-1")
				mockSeedPeerCluster(t, s.db, "seed-2")
				assert.NoError(s.db.Model(&seedPeerCluster).Association("SchedulerClusters").Append(&schedulerCluster))
			},
			id:                 2,
			schedulerClusterID: 1,
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				schedulerCluster := models.SchedulerCluster{}
				assert.NoError(s.db.Preload("SeedPeerClusters").First(&schedulerCluster, 1).Error)
				assert.Len(schedulerCluster.SeedPeerClusters, 1)
				assert.Equal("seed-2", schedulerCluster.SeedPeerClusters[0].Name)

				previous := models.SeedPeerCluster{}
				assert.NoError(s.db.Preload("SchedulerClusters").First(&previous, 1).Error)
				assert.Empty(previous.SchedulerClusters)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			tc.setup(t, s)

			tc.expect(t, s, s.AddSchedulerClusterToSeedPeerCluster(context.Background(), tc.id, tc.schedulerClusterID))
		})
	}
}
