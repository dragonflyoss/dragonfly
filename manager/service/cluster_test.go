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

func mockCreateClusterRequest(name string) types.CreateClusterRequest {
	return types.CreateClusterRequest{
		Name:                   name,
		BIO:                    "bio",
		Scopes:                 &types.SchedulerClusterScopes{IDC: "idc", Location: "location"},
		SchedulerClusterConfig: &types.SchedulerClusterConfig{CandidateParentLimit: 4, FilterParentLimit: 40, JobRateLimit: 10},
		SeedPeerClusterConfig:  &types.SeedPeerClusterConfig{LoadLimit: 300},
		PeerClusterConfig:      &types.SchedulerClusterClientConfig{LoadLimit: 50},
		IsDefault:              true,
	}
}

func mockCluster(t *testing.T, s *service, name string) *types.CreateClusterResponse {
	assert := assert.New(t)
	resp, err := s.CreateCluster(context.Background(), mockCreateClusterRequest(name))
	assert.NoError(err)

	return resp
}

func TestService_CreateCluster(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service)
		req    types.CreateClusterRequest
		expect func(t *testing.T, s *service, resp *types.CreateClusterResponse, err error)
	}{
		{
			name:  "creates linked scheduler cluster and seed peer cluster",
			setup: func(t *testing.T, s *service) {},
			req:   mockCreateClusterRequest("foo"),
			expect: func(t *testing.T, s *service, resp *types.CreateClusterResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", resp.Name)
				assert.Equal("bio", resp.BIO)
				assert.Equal(resp.ID, resp.SchedulerClusterID)
				assert.True(resp.IsDefault)

				seedPeerCluster := models.SeedPeerCluster{}
				assert.NoError(s.db.Preload("SchedulerClusters").First(&seedPeerCluster, resp.SeedPeerClusterID).Error)
				assert.Equal("foo", seedPeerCluster.Name)
				assert.Equal(float64(300), seedPeerCluster.Config["load_limit"])
				assert.Len(seedPeerCluster.SchedulerClusters, 1)
				assert.Equal(resp.SchedulerClusterID, seedPeerCluster.SchedulerClusters[0].ID)

				schedulerCluster := models.SchedulerCluster{}
				assert.NoError(s.db.First(&schedulerCluster, resp.SchedulerClusterID).Error)
				assert.Equal(float64(4), schedulerCluster.Config["candidate_parent_limit"])
				assert.Equal(float64(50), schedulerCluster.ClientConfig["load_limit"])
				assert.Equal(float64(300), schedulerCluster.SeedClientConfig["load_limit"])
				assert.Equal("idc", schedulerCluster.Scopes["idc"])
				assert.True(schedulerCluster.IsDefault)
			},
		},
		{
			name: "duplicate scheduler cluster name creates nothing",
			setup: func(t *testing.T, s *service) {
				mockSchedulerCluster(t, s.db, "foo")
			},
			req: mockCreateClusterRequest("foo"),
			expect: func(t *testing.T, s *service, resp *types.CreateClusterResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(resp)

				var count int64
				assert.NoError(s.db.Model(&models.SeedPeerCluster{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
		{
			name: "duplicate seed peer cluster name rolls back scheduler cluster",
			setup: func(t *testing.T, s *service) {
				mockSeedPeerCluster(t, s.db, "foo")
			},
			req: mockCreateClusterRequest("foo"),
			expect: func(t *testing.T, s *service, resp *types.CreateClusterResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(resp)

				var count int64
				assert.NoError(s.db.Model(&models.SchedulerCluster{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			tc.setup(t, s)

			resp, err := s.CreateCluster(context.Background(), tc.req)
			tc.expect(t, s, resp, err)
		})
	}
}

func TestService_DestroyCluster(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, s *service, id uint, err error)
	}{
		{
			name: "cluster not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			expect: func(t *testing.T, s *service, id uint, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "cluster still has schedulers",
			setup: func(t *testing.T, s *service) uint {
				resp := mockCluster(t, s, "foo")
				mockScheduler(t, s.db, resp.SchedulerClusterID, "scheduler-1", models.SchedulerStateActive, nil)
				return resp.SchedulerClusterID
			},
			expect: func(t *testing.T, s *service, id uint, err error) {
				assert := assert.New(t)
				assert.Error(err)

				var count int64
				assert.NoError(s.db.Model(&models.SchedulerCluster{}).Count(&count).Error)
				assert.Equal(int64(1), count)
			},
		},
		{
			name: "hard deletes scheduler cluster together with its seed peer cluster",
			setup: func(t *testing.T, s *service) uint {
				return mockCluster(t, s, "foo").SchedulerClusterID
			},
			expect: func(t *testing.T, s *service, id uint, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var schedulerClusterCount int64
				assert.NoError(s.db.Unscoped().Model(&models.SchedulerCluster{}).Count(&schedulerClusterCount).Error)
				assert.Equal(int64(0), schedulerClusterCount)

				var seedPeerClusterCount int64
				assert.NoError(s.db.Unscoped().Model(&models.SeedPeerCluster{}).Count(&seedPeerClusterCount).Error)
				assert.Equal(int64(0), seedPeerClusterCount)

				var joinCount int64
				assert.NoError(s.db.Table("seed_peer_cluster_scheduler_cluster").Where("scheduler_cluster_id = ?", id).Count(&joinCount).Error)
				assert.Equal(int64(0), joinCount)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, id, s.DestroyCluster(context.Background(), id))
		})
	}
}

func TestService_UpdateCluster(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdateClusterRequest
		expect func(t *testing.T, s *service, id uint, resp *types.UpdateClusterResponse, err error)
	}{
		{
			name: "cluster not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.UpdateClusterRequest{Name: "bar"},
			expect: func(t *testing.T, s *service, id uint, resp *types.UpdateClusterResponse, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(resp)
			},
		},
		{
			name: "unsets is_default without touching other fields",
			setup: func(t *testing.T, s *service) uint {
				return mockCluster(t, s, "foo").SchedulerClusterID
			},
			req: types.UpdateClusterRequest{IsDefault: false},
			expect: func(t *testing.T, s *service, id uint, resp *types.UpdateClusterResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.False(resp.IsDefault)

				schedulerCluster := models.SchedulerCluster{}
				assert.NoError(s.db.First(&schedulerCluster, id).Error)
				assert.False(schedulerCluster.IsDefault)
				assert.Equal("foo", schedulerCluster.Name)
				assert.Equal("bio", schedulerCluster.BIO)
				assert.Equal(float64(4), schedulerCluster.Config["candidate_parent_limit"])
			},
		},
		{
			name: "renames both clusters and updates seed peer cluster config",
			setup: func(t *testing.T, s *service) uint {
				return mockCluster(t, s, "foo").SchedulerClusterID
			},
			req: types.UpdateClusterRequest{
				Name:                  "bar",
				BIO:                   "new bio",
				SeedPeerClusterConfig: &types.SeedPeerClusterConfig{LoadLimit: 100},
				IsDefault:             true,
			},
			expect: func(t *testing.T, s *service, id uint, resp *types.UpdateClusterResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(id, resp.SchedulerClusterID)
				assert.Equal("bar", resp.Name)

				schedulerCluster := models.SchedulerCluster{}
				assert.NoError(s.db.First(&schedulerCluster, id).Error)
				assert.Equal("bar", schedulerCluster.Name)
				assert.Equal("new bio", schedulerCluster.BIO)
				assert.Equal(float64(100), schedulerCluster.SeedClientConfig["load_limit"])
				assert.Equal(float64(4), schedulerCluster.Config["candidate_parent_limit"])
				assert.True(schedulerCluster.IsDefault)

				seedPeerCluster := models.SeedPeerCluster{}
				assert.NoError(s.db.First(&seedPeerCluster, resp.SeedPeerClusterID).Error)
				assert.Equal("bar", seedPeerCluster.Name)
				assert.Equal("new bio", seedPeerCluster.BIO)
				assert.Equal(float64(100), seedPeerCluster.Config["load_limit"])
			},
		},
		{
			name: "scheduler cluster without seed peer cluster rolls back the rename",
			setup: func(t *testing.T, s *service) uint {
				return mockSchedulerCluster(t, s.db, "baz").ID
			},
			req: types.UpdateClusterRequest{Name: "renamed"},
			expect: func(t *testing.T, s *service, id uint, resp *types.UpdateClusterResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(resp)

				schedulerCluster := models.SchedulerCluster{}
				assert.NoError(s.db.First(&schedulerCluster, id).Error)
				assert.Equal("baz", schedulerCluster.Name)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			resp, err := s.UpdateCluster(context.Background(), id, tc.req)
			tc.expect(t, s, id, resp, err)
		})
	}
}

func TestService_GetCluster(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, id uint, resp *types.GetClusterResponse, err error)
	}{
		{
			name: "cluster not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			expect: func(t *testing.T, id uint, resp *types.GetClusterResponse, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(resp)
			},
		},
		{
			name: "scheduler cluster without seed peer cluster",
			setup: func(t *testing.T, s *service) uint {
				return mockSchedulerCluster(t, s.db, "foo").ID
			},
			expect: func(t *testing.T, id uint, resp *types.GetClusterResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(resp)
			},
		},
		{
			name: "maps stored json configs back to typed configs",
			setup: func(t *testing.T, s *service) uint {
				return mockCluster(t, s, "foo").SchedulerClusterID
			},
			expect: func(t *testing.T, id uint, resp *types.GetClusterResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(id, resp.ID)
				assert.Equal(id, resp.SchedulerClusterID)
				assert.NotZero(resp.SeedPeerClusterID)
				assert.Equal("foo", resp.Name)
				assert.True(resp.IsDefault)
				assert.Equal("idc", resp.Scopes.IDC)
				assert.Equal("location", resp.Scopes.Location)
				assert.Equal(uint32(4), resp.SchedulerClusterConfig.CandidateParentLimit)
				assert.Equal(uint32(40), resp.SchedulerClusterConfig.FilterParentLimit)
				assert.Equal(uint32(10), resp.SchedulerClusterConfig.JobRateLimit)
				assert.Equal(uint32(300), resp.SeedPeerClusterConfig.LoadLimit)
				assert.Equal(uint32(50), resp.PeerClusterConfig.LoadLimit)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			resp, err := s.GetCluster(context.Background(), id)
			tc.expect(t, id, resp, err)
		})
	}
}

func TestService_GetClusters(t *testing.T) {
	s := newTestService(t)
	for _, name := range []string{"foo", "bar", "baz"} {
		mockCluster(t, s, name)
	}

	tests := []struct {
		name   string
		query  types.GetClustersQuery
		expect func(t *testing.T, resp []types.GetClusterResponse, count int64, err error)
	}{
		{
			name:  "first page",
			query: types.GetClustersQuery{Page: 1, PerPage: 2},
			expect: func(t *testing.T, resp []types.GetClusterResponse, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(resp, 2)
				assert.Equal(int64(3), count)
				assert.Equal(uint32(300), resp[0].SeedPeerClusterConfig.LoadLimit)
				assert.Equal("idc", resp[0].Scopes.IDC)
			},
		},
		{
			name:  "last page",
			query: types.GetClustersQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, resp []types.GetClusterResponse, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(resp, 1)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by name",
			query: types.GetClustersQuery{Name: "bar", Page: 1, PerPage: 10},
			expect: func(t *testing.T, resp []types.GetClusterResponse, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(resp, 1)
				assert.Equal("bar", resp[0].Name)
				assert.Equal(int64(1), count)
			},
		},
		{
			name:  "no match",
			query: types.GetClustersQuery{Name: "missing", Page: 1, PerPage: 10},
			expect: func(t *testing.T, resp []types.GetClusterResponse, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(resp)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, count, err := s.GetClusters(context.Background(), tc.query)
			tc.expect(t, resp, count, err)
		})
	}
}
