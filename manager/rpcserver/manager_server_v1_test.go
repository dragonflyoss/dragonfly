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

package rpcserver

import (
	"context"
	"errors"
	"io"
	"testing"

	cachev9 "github.com/go-redis/cache/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	managerv1 "d7y.io/api/v2/pkg/apis/manager/v1"
	managerv1mocks "d7y.io/api/v2/pkg/apis/manager/v1/mocks"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/searcher/mocks"
	"d7y.io/dragonfly/v2/manager/types"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
)

func TestManagerServerV1_UpdateSeedPeer(t *testing.T) {
	tests := []struct {
		name   string
		req    *managerv1.UpdateSeedPeerRequest
		expect func(t *testing.T, s *managerServerV1, resp *managerv1.SeedPeer, err error)
	}{
		{
			name: "creates seed peer with default inactive state when no row matches",
			req: &managerv1.UpdateSeedPeerRequest{
				Hostname:          mockNewHostname,
				Ip:                mockNewIP,
				Type:              mockSeedPeerType,
				Idc:               mockIDC,
				Location:          mockLocation,
				Port:              mockSeedPeerPort,
				DownloadPort:      mockSeedPeerDownloadPort,
				SeedPeerClusterId: uint64(mockSeedPeerClusterID),
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.SeedPeer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.SeedPeer{}).Where(models.SeedPeer{Hostname: mockNewHostname}).Count(&count).Error)
				assert.Equal(int64(1), count)
				var seedPeer models.SeedPeer
				assert.NoError(s.db.First(&seedPeer, models.SeedPeer{Hostname: mockNewHostname}).Error)
				assert.Equal(uint64(seedPeer.ID), resp.GetId())
				assert.Equal(mockNewIP, seedPeer.IP)
				assert.Equal(mockSeedPeerType, seedPeer.Type)
				assert.Equal(mockIDC, seedPeer.IDC)
				assert.Equal(mockLocation, seedPeer.Location)
				assert.Equal(mockSeedPeerPort, seedPeer.Port)
				assert.Equal(mockSeedPeerDownloadPort, seedPeer.DownloadPort)
				assert.Equal(mockSeedPeerClusterID, seedPeer.SeedPeerClusterID)
				assert.Equal(models.SeedPeerStateInactive, seedPeer.State)
			},
		},
		{
			name: "changed port updates the existing row in place",
			req: &managerv1.UpdateSeedPeerRequest{
				Hostname:          mockActiveSeedPeerHostname,
				Ip:                mockActiveSeedPeerIP,
				Type:              mockSeedPeerType,
				Port:              mockUpdatedSeedPeerPort,
				DownloadPort:      mockSeedPeerDownloadPort,
				SeedPeerClusterId: uint64(mockSeedPeerClusterID),
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.SeedPeer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.SeedPeer{}).Where(models.SeedPeer{Hostname: mockActiveSeedPeerHostname}).Count(&count).Error)
				assert.Equal(int64(1), count)
				var seedPeer models.SeedPeer
				assert.NoError(s.db.First(&seedPeer, models.SeedPeer{Hostname: mockActiveSeedPeerHostname}).Error)
				assert.Equal(uint64(seedPeer.ID), resp.GetId())
				assert.Equal(mockUpdatedSeedPeerPort, seedPeer.Port)
			},
		},
		{
			name: "updates matching seed peer to active and invalidates cache",
			req: &managerv1.UpdateSeedPeerRequest{
				Hostname:          mockInactiveSeedPeerHostname,
				Ip:                mockInactiveSeedPeerIP,
				Type:              mockSeedPeerType,
				Idc:               mockUpdatedIDC,
				Location:          mockUpdatedLocation,
				Port:              mockSeedPeerPort,
				DownloadPort:      mockSeedPeerDownloadPort,
				SeedPeerClusterId: uint64(mockSeedPeerClusterID),
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.SeedPeer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				var seedPeer models.SeedPeer
				assert.NoError(s.db.First(&seedPeer, models.SeedPeer{Hostname: mockInactiveSeedPeerHostname}).Error)
				assert.Equal(uint64(seedPeer.ID), resp.GetId())
				assert.Equal(mockUpdatedIDC, seedPeer.IDC)
				assert.Equal(mockUpdatedLocation, seedPeer.Location)
				assert.Equal(models.SeedPeerStateActive, seedPeer.State)
				assert.False(s.cache.Exists(context.Background(), pkgredis.MakeSeedPeerKeyInManager(mockSeedPeerClusterID, mockInactiveSeedPeerHostname, mockInactiveSeedPeerIP)))
			},
		},
		{
			name: "unknown cluster returns internal and creates nothing",
			req: &managerv1.UpdateSeedPeerRequest{
				Hostname:          mockNewHostname,
				Ip:                mockNewIP,
				Port:              mockSeedPeerPort,
				DownloadPort:      mockSeedPeerDownloadPort,
				SeedPeerClusterId: uint64(mockUnknownClusterID),
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.SeedPeer, err error) {
				assert := assert.New(t)
				assert.Nil(resp)
				assert.Equal(codes.Internal, status.Code(err))
				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.SeedPeer{}).Where(models.SeedPeer{Hostname: mockNewHostname}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &managerServerV1{db: mockDB(t), cache: mockCache()}
			if err := s.cache.Set(&cachev9.Item{Key: pkgredis.MakeSeedPeerKeyInManager(uint(tc.req.SeedPeerClusterId), tc.req.Hostname, tc.req.Ip), Value: &managerv1.SeedPeer{Hostname: mockCachedHostname}}); err != nil {
				t.Fatal(err)
			}

			resp, err := s.UpdateSeedPeer(context.Background(), tc.req)
			tc.expect(t, s, resp, err)
		})
	}
}

func TestManagerServerV1_UpdateScheduler(t *testing.T) {
	tests := []struct {
		name   string
		req    *managerv1.UpdateSchedulerRequest
		expect func(t *testing.T, s *managerServerV1, resp *managerv1.Scheduler, err error)
	}{
		{
			name: "creates scheduler with default features when no row matches",
			req: &managerv1.UpdateSchedulerRequest{
				Hostname:           mockNewHostname,
				Ip:                 mockNewIP,
				Port:               mockSchedulerPort,
				Idc:                mockIDC,
				Location:           mockLocation,
				SchedulerClusterId: uint64(mockSchedulerClusterID),
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.Scheduler{}).Where(models.Scheduler{Hostname: mockNewHostname}).Count(&count).Error)
				assert.Equal(int64(1), count)
				var scheduler models.Scheduler
				assert.NoError(s.db.First(&scheduler, models.Scheduler{Hostname: mockNewHostname}).Error)
				assert.Equal(uint64(scheduler.ID), resp.GetId())
				assert.Equal(mockNewIP, scheduler.IP)
				assert.Equal(mockSchedulerPort, scheduler.Port)
				assert.Equal(mockIDC, scheduler.IDC)
				assert.Equal(mockLocation, scheduler.Location)
				assert.Equal(mockSchedulerClusterID, scheduler.SchedulerClusterID)
				assert.Equal(models.Array(types.DefaultSchedulerFeatures), scheduler.Features)
				assert.Equal(models.SchedulerStateInactive, scheduler.State)
				assert.JSONEq(mockDefaultFeaturesJSON, string(resp.GetFeatures()))
			},
		},
		{
			name: "changed port updates the existing row in place",
			req: &managerv1.UpdateSchedulerRequest{
				Hostname:           mockActiveSchedulerHostname,
				Ip:                 mockActiveSchedulerIP,
				Port:               mockUpdatedSchedulerPort,
				SchedulerClusterId: uint64(mockSchedulerClusterID),
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.Scheduler{}).Where(models.Scheduler{Hostname: mockActiveSchedulerHostname}).Count(&count).Error)
				assert.Equal(int64(1), count)
				var scheduler models.Scheduler
				assert.NoError(s.db.First(&scheduler, models.Scheduler{Hostname: mockActiveSchedulerHostname}).Error)
				assert.Equal(uint64(scheduler.ID), resp.GetId())
				assert.Equal(mockUpdatedSchedulerPort, scheduler.Port)
				assert.Equal(models.SchedulerStateActive, scheduler.State)
			},
		},
		{
			name: "updates matching scheduler keeping features and invalidates cache",
			req: &managerv1.UpdateSchedulerRequest{
				Hostname:           mockInactiveSchedulerHostname,
				Ip:                 mockInactiveSchedulerIP,
				Port:               mockSchedulerPort,
				Idc:                mockUpdatedIDC,
				Location:           mockUpdatedLocation,
				SchedulerClusterId: uint64(mockSchedulerClusterID),
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				var scheduler models.Scheduler
				assert.NoError(s.db.First(&scheduler, models.Scheduler{Hostname: mockInactiveSchedulerHostname}).Error)
				assert.Equal(uint64(scheduler.ID), resp.GetId())
				assert.Equal(mockUpdatedIDC, scheduler.IDC)
				assert.Equal(mockUpdatedLocation, scheduler.Location)
				assert.Equal(models.Array(types.DefaultSchedulerFeatures), scheduler.Features)
				assert.JSONEq(mockDefaultFeaturesJSON, string(resp.GetFeatures()))
				assert.False(s.cache.Exists(context.Background(), pkgredis.MakeSchedulerKeyInManager(mockSchedulerClusterID, mockInactiveSchedulerHostname, mockInactiveSchedulerIP)))
			},
		},
		{
			name: "unknown cluster returns internal and creates nothing",
			req: &managerv1.UpdateSchedulerRequest{
				Hostname:           mockNewHostname,
				Ip:                 mockNewIP,
				Port:               mockSchedulerPort,
				SchedulerClusterId: uint64(mockUnknownClusterID),
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.Scheduler, err error) {
				assert := assert.New(t)
				assert.Nil(resp)
				assert.Equal(codes.Internal, status.Code(err))
				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.Scheduler{}).Where(models.Scheduler{Hostname: mockNewHostname}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &managerServerV1{db: mockDB(t), cache: mockCache()}
			if err := s.cache.Set(&cachev9.Item{Key: pkgredis.MakeSchedulerKeyInManager(uint(tc.req.SchedulerClusterId), tc.req.Hostname, tc.req.Ip), Value: &managerv1.Scheduler{Hostname: mockCachedHostname}}); err != nil {
				t.Fatal(err)
			}

			resp, err := s.UpdateScheduler(context.Background(), tc.req)
			tc.expect(t, s, resp, err)
		})
	}
}

func TestManagerServerV1_ListSchedulers(t *testing.T) {
	hostInfo := map[string]string{"idc": mockIDC, "location": mockLocation, "security_domain": "domain-1"}
	req := &managerv1.ListSchedulersRequest{
		Hostname: mockPeerHostname,
		Ip:       mockPeerIP,
		HostInfo: hostInfo,
		Version:  mockVersion,
		Commit:   mockCommit,
	}
	cacheKey := pkgredis.MakeSchedulersKeyForPeerInManager(mockPeerHostname, mockPeerIP, mockVersion)
	tests := []struct {
		name   string
		cached *managerv1.ListSchedulersResponse
		mock   func(m *mocks.MockSearcherMockRecorder)
		expect func(t *testing.T, s *managerServerV1, resp *managerv1.ListSchedulersResponse, err error)
	}{
		{
			name: "forwards host info to searcher and returns active schedulers of the matched cluster",
			mock: func(m *mocks.MockSearcherMockRecorder) {
				m.FindSchedulerClusters(gomock.Any(), gomock.Any(), mockPeerIP, mockPeerHostname, hostInfo, gomock.Any()).
					DoAndReturn(func(_ context.Context, clusters []models.SchedulerCluster, _, _ string, _ map[string]string, _ *zap.SugaredLogger) ([]models.SchedulerCluster, error) {
						for _, cluster := range clusters {
							if cluster.Name == mockSchedulerClusterName {
								return []models.SchedulerCluster{cluster}, nil
							}
						}

						return nil, errors.New("scheduler cluster not found")
					})
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.ListSchedulersResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(resp.GetSchedulers(), 1)
				for _, scheduler := range resp.GetSchedulers() {
					assert.Equal(mockActiveSchedulerHostname, scheduler.GetHostname())
					assert.Equal(models.SchedulerStateActive, scheduler.GetState())
					assert.Equal(uint64(mockSchedulerClusterID), scheduler.GetSchedulerClusterId())
					assert.JSONEq(mockDefaultFeaturesJSON, string(scheduler.GetFeatures()))
					assert.Len(scheduler.GetSeedPeers(), 1)
					for _, seedPeer := range scheduler.GetSeedPeers() {
						assert.Equal(mockActiveSeedPeerHostname, seedPeer.GetHostname())
						assert.Equal(models.SeedPeerStateActive, seedPeer.GetState())
						assert.Equal(uint64(mockSeedPeerClusterID), seedPeer.GetSeedPeerClusterId())
					}
				}

				assert.True(s.cache.Exists(context.Background(), cacheKey))
			},
		},
		{
			name: "searcher failure falls back to active schedulers of every cluster",
			mock: func(m *mocks.MockSearcherMockRecorder) {
				m.FindSchedulerClusters(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("no matching scheduler cluster"))
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.ListSchedulersResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				var hostnames []string
				for _, scheduler := range resp.GetSchedulers() {
					hostnames = append(hostnames, scheduler.GetHostname())
				}

				assert.Subset(hostnames, []string{mockActiveSchedulerHostname, mockSecondClusterSchedulerHostname})
				assert.NotContains(hostnames, mockInactiveSchedulerHostname)
				assert.True(s.cache.Exists(context.Background(), cacheKey))
			},
		},
		{
			name: "no candidate cluster returns empty response without caching",
			mock: func(m *mocks.MockSearcherMockRecorder) {
				m.FindSchedulerClusters(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.SchedulerCluster{}, nil)
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.ListSchedulersResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(resp.GetSchedulers())
				assert.False(s.cache.Exists(context.Background(), cacheKey))
			},
		},
		{
			name:   "cache hit skips searcher and db",
			cached: &managerv1.ListSchedulersResponse{Schedulers: []*managerv1.Scheduler{{Hostname: mockCachedHostname}}},
			mock:   func(m *mocks.MockSearcherMockRecorder) {},
			expect: func(t *testing.T, _ *managerServerV1, resp *managerv1.ListSchedulersResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(resp.GetSchedulers(), 1)
				assert.Equal(mockCachedHostname, resp.GetSchedulers()[0].GetHostname())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			ms := mocks.NewMockSearcher(ctl)
			tc.mock(ms.EXPECT())

			s := &managerServerV1{db: mockDB(t), cache: mockCache(), searcher: ms}
			if tc.cached != nil {
				if err := s.cache.Set(&cachev9.Item{Key: cacheKey, Value: tc.cached}); err != nil {
					t.Fatal(err)
				}
			}

			resp, err := s.ListSchedulers(context.Background(), req)
			tc.expect(t, s, resp, err)
		})
	}
}

func TestManagerServerV1_KeepAlive(t *testing.T) {
	schedulerReq := &managerv1.KeepAliveRequest{
		SourceType: managerv1.SourceType_SCHEDULER_SOURCE,
		Hostname:   mockInactiveSchedulerHostname,
		Ip:         mockInactiveSchedulerIP,
		ClusterId:  uint64(mockSchedulerClusterID),
	}
	seedPeerReq := &managerv1.KeepAliveRequest{
		SourceType: managerv1.SourceType_SEED_PEER_SOURCE,
		Hostname:   mockInactiveSeedPeerHostname,
		Ip:         mockInactiveSeedPeerIP,
		ClusterId:  uint64(mockSeedPeerClusterID),
	}
	tests := []struct {
		name   string
		mock   func(t *testing.T, s *managerServerV1, ms *managerv1mocks.MockManager_KeepAliveServerMockRecorder, states *[]string)
		expect func(t *testing.T, s *managerServerV1, states []string, err error)
	}{
		{
			name: "scheduler turns active on first message and inactive when stream closes",
			mock: func(t *testing.T, s *managerServerV1, ms *managerv1mocks.MockManager_KeepAliveServerMockRecorder, states *[]string) {
				if err := s.cache.Set(&cachev9.Item{Key: pkgredis.MakeSchedulerKeyInManager(mockSchedulerClusterID, mockInactiveSchedulerHostname, mockInactiveSchedulerIP), Value: &managerv1.Scheduler{Hostname: mockCachedHostname}}); err != nil {
					t.Fatal(err)
				}

				ms.Context().Return(context.Background()).AnyTimes()
				gomock.InOrder(
					ms.Recv().DoAndReturn(func() (*managerv1.KeepAliveRequest, error) {
						assert := assert.New(t)
						var scheduler models.Scheduler
						assert.NoError(s.db.First(&scheduler, models.Scheduler{Hostname: mockInactiveSchedulerHostname, IP: mockInactiveSchedulerIP, SchedulerClusterID: mockSchedulerClusterID}).Error)
						*states = append(*states, scheduler.State)
						return schedulerReq, nil
					}),
					ms.Recv().DoAndReturn(func() (*managerv1.KeepAliveRequest, error) {
						assert := assert.New(t)
						var scheduler models.Scheduler
						assert.NoError(s.db.First(&scheduler, models.Scheduler{Hostname: mockInactiveSchedulerHostname, IP: mockInactiveSchedulerIP, SchedulerClusterID: mockSchedulerClusterID}).Error)
						*states = append(*states, scheduler.State)
						return schedulerReq, nil
					}),
					ms.Recv().DoAndReturn(func() (*managerv1.KeepAliveRequest, error) {
						assert := assert.New(t)
						var scheduler models.Scheduler
						assert.NoError(s.db.First(&scheduler, models.Scheduler{Hostname: mockInactiveSchedulerHostname, IP: mockInactiveSchedulerIP, SchedulerClusterID: mockSchedulerClusterID}).Error)
						*states = append(*states, scheduler.State)
						return nil, io.EOF
					}),
				)
			},
			expect: func(t *testing.T, s *managerServerV1, states []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{models.SchedulerStateInactive, models.SchedulerStateActive, models.SchedulerStateActive}, states)
				var scheduler models.Scheduler
				assert.NoError(s.db.First(&scheduler, models.Scheduler{Hostname: mockInactiveSchedulerHostname}).Error)
				assert.Equal(models.SchedulerStateInactive, scheduler.State)
				assert.True(scheduler.LastKeepAliveAt.IsZero())
				assert.False(s.cache.Exists(context.Background(), pkgredis.MakeSchedulerKeyInManager(mockSchedulerClusterID, mockInactiveSchedulerHostname, mockInactiveSchedulerIP)))
			},
		},
		{
			name: "seed peer turns active on first message and inactive when stream closes",
			mock: func(t *testing.T, s *managerServerV1, ms *managerv1mocks.MockManager_KeepAliveServerMockRecorder, states *[]string) {
				if err := s.cache.Set(&cachev9.Item{Key: pkgredis.MakeSeedPeerKeyInManager(mockSeedPeerClusterID, mockInactiveSeedPeerHostname, mockInactiveSeedPeerIP), Value: &managerv1.SeedPeer{Hostname: mockCachedHostname}}); err != nil {
					t.Fatal(err)
				}

				ms.Context().Return(context.Background()).AnyTimes()
				gomock.InOrder(
					ms.Recv().DoAndReturn(func() (*managerv1.KeepAliveRequest, error) {
						assert := assert.New(t)
						var seedPeer models.SeedPeer
						assert.NoError(s.db.First(&seedPeer, models.SeedPeer{Hostname: mockInactiveSeedPeerHostname, IP: mockInactiveSeedPeerIP, SeedPeerClusterID: mockSeedPeerClusterID}).Error)
						*states = append(*states, seedPeer.State)
						return seedPeerReq, nil
					}),
					ms.Recv().DoAndReturn(func() (*managerv1.KeepAliveRequest, error) {
						assert := assert.New(t)
						var seedPeer models.SeedPeer
						assert.NoError(s.db.First(&seedPeer, models.SeedPeer{Hostname: mockInactiveSeedPeerHostname, IP: mockInactiveSeedPeerIP, SeedPeerClusterID: mockSeedPeerClusterID}).Error)
						*states = append(*states, seedPeer.State)
						return nil, io.EOF
					}),
				)
			},
			expect: func(t *testing.T, s *managerServerV1, states []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{models.SeedPeerStateInactive, models.SeedPeerStateActive}, states)
				var seedPeer models.SeedPeer
				assert.NoError(s.db.First(&seedPeer, models.SeedPeer{Hostname: mockInactiveSeedPeerHostname}).Error)
				assert.Equal(models.SeedPeerStateInactive, seedPeer.State)
				assert.False(s.cache.Exists(context.Background(), pkgredis.MakeSeedPeerKeyInManager(mockSeedPeerClusterID, mockInactiveSeedPeerHostname, mockInactiveSeedPeerIP)))
			},
		},
		{
			name: "stream failure marks seed peer inactive and returns unknown",
			mock: func(t *testing.T, s *managerServerV1, ms *managerv1mocks.MockManager_KeepAliveServerMockRecorder, states *[]string) {
				ms.Context().Return(context.Background()).AnyTimes()
				gomock.InOrder(
					ms.Recv().DoAndReturn(func() (*managerv1.KeepAliveRequest, error) {
						assert := assert.New(t)
						var seedPeer models.SeedPeer
						assert.NoError(s.db.First(&seedPeer, models.SeedPeer{Hostname: mockInactiveSeedPeerHostname, IP: mockInactiveSeedPeerIP, SeedPeerClusterID: mockSeedPeerClusterID}).Error)
						*states = append(*states, seedPeer.State)
						return seedPeerReq, nil
					}),
					ms.Recv().DoAndReturn(func() (*managerv1.KeepAliveRequest, error) {
						assert := assert.New(t)
						var seedPeer models.SeedPeer
						assert.NoError(s.db.First(&seedPeer, models.SeedPeer{Hostname: mockInactiveSeedPeerHostname, IP: mockInactiveSeedPeerIP, SeedPeerClusterID: mockSeedPeerClusterID}).Error)
						*states = append(*states, seedPeer.State)
						return nil, errors.New("connection reset")
					}),
				)
			},
			expect: func(t *testing.T, s *managerServerV1, states []string, err error) {
				assert := assert.New(t)
				assert.Equal(codes.Unknown, status.Code(err))
				assert.Equal([]string{models.SeedPeerStateInactive, models.SeedPeerStateActive}, states)
				var seedPeer models.SeedPeer
				assert.NoError(s.db.First(&seedPeer, models.SeedPeer{Hostname: mockInactiveSeedPeerHostname}).Error)
				assert.Equal(models.SeedPeerStateInactive, seedPeer.State)
			},
		},
		{
			name: "unknown seed peer returns internal",
			mock: func(_ *testing.T, _ *managerServerV1, ms *managerv1mocks.MockManager_KeepAliveServerMockRecorder, _ *[]string) {
				ms.Recv().Return(&managerv1.KeepAliveRequest{
					SourceType: managerv1.SourceType_SEED_PEER_SOURCE,
					Hostname:   mockUnknownHostname,
					Ip:         mockUnknownIP,
					ClusterId:  uint64(mockSeedPeerClusterID),
				}, nil)
			},
			expect: func(t *testing.T, _ *managerServerV1, _ []string, err error) {
				assert := assert.New(t)
				assert.Equal(codes.Internal, status.Code(err))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			stream := managerv1mocks.NewMockManager_KeepAliveServer(ctl)
			s := &managerServerV1{db: mockDB(t), cache: mockCache()}
			var states []string
			tc.mock(t, s, stream.EXPECT(), &states)

			err := s.KeepAlive(stream)
			tc.expect(t, s, states, err)
		})
	}
}
