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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	managerv1 "d7y.io/api/v2/pkg/apis/manager/v1"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/searcher/mocks"
	"d7y.io/dragonfly/v2/manager/types"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
)

type mockKeepAliveStreamV1 struct {
	grpc.ServerStream
	reqs    []*managerv1.KeepAliveRequest
	err     error
	observe func() string
	states  []string
}

func (m *mockKeepAliveStreamV1) Recv() (*managerv1.KeepAliveRequest, error) {
	if m.observe != nil {
		m.states = append(m.states, m.observe())
	}

	if len(m.reqs) == 0 {
		return nil, m.err
	}

	req := m.reqs[0]
	m.reqs = m.reqs[1:]
	return req, nil
}

func (m *mockKeepAliveStreamV1) SendAndClose(*emptypb.Empty) error {
	return nil
}

func (m *mockKeepAliveStreamV1) Context() context.Context {
	return context.Background()
}

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
				assert.Equal(int64(1), countRows(t, s.db, &models.SeedPeer{}, models.SeedPeer{Hostname: mockNewHostname}))
				seedPeer := findSeedPeer(t, s.db, mockNewHostname)
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
				assert.Equal(int64(1), countRows(t, s.db, &models.SeedPeer{}, models.SeedPeer{Hostname: mockActiveSeedPeerHostname}))
				seedPeer := findSeedPeer(t, s.db, mockActiveSeedPeerHostname)
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
				seedPeer := findSeedPeer(t, s.db, mockInactiveSeedPeerHostname)
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
				assert.Equal(int64(0), countRows(t, s.db, &models.SeedPeer{}, models.SeedPeer{Hostname: mockNewHostname}))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestServerV1(t, nil)
			setCache(t, s.cache, pkgredis.MakeSeedPeerKeyInManager(uint(tc.req.SeedPeerClusterId), tc.req.Hostname, tc.req.Ip), &managerv1.SeedPeer{Hostname: mockCachedHostname})

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
				assert.Equal(int64(1), countRows(t, s.db, &models.Scheduler{}, models.Scheduler{Hostname: mockNewHostname}))
				scheduler := findScheduler(t, s.db, mockNewHostname)
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
				assert.Equal(int64(1), countRows(t, s.db, &models.Scheduler{}, models.Scheduler{Hostname: mockActiveSchedulerHostname}))
				scheduler := findScheduler(t, s.db, mockActiveSchedulerHostname)
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
				scheduler := findScheduler(t, s.db, mockInactiveSchedulerHostname)
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
				assert.Equal(int64(0), countRows(t, s.db, &models.Scheduler{}, models.Scheduler{Hostname: mockNewHostname}))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestServerV1(t, nil)
			setCache(t, s.cache, pkgredis.MakeSchedulerKeyInManager(uint(tc.req.SchedulerClusterId), tc.req.Hostname, tc.req.Ip), &managerv1.Scheduler{Hostname: mockCachedHostname})

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
					DoAndReturn(pickSchedulerCluster(mockSchedulerClusterName))
			},
			expect: func(t *testing.T, s *managerServerV1, resp *managerv1.ListSchedulersResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{mockActiveSchedulerHostname}, hostnames(resp.GetSchedulers()))
				for _, scheduler := range resp.GetSchedulers() {
					assert.Equal(models.SchedulerStateActive, scheduler.GetState())
					assert.Equal(uint64(mockSchedulerClusterID), scheduler.GetSchedulerClusterId())
					assert.JSONEq(mockDefaultFeaturesJSON, string(scheduler.GetFeatures()))
					assert.Equal([]string{mockActiveSeedPeerHostname}, hostnames(scheduler.GetSeedPeers()))
					for _, seedPeer := range scheduler.GetSeedPeers() {
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
				assert.Subset(hostnames(resp.GetSchedulers()), []string{mockActiveSchedulerHostname, mockSecondClusterSchedulerHostname})
				assert.NotContains(hostnames(resp.GetSchedulers()), mockInactiveSchedulerHostname)
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
				assert.Equal([]string{mockCachedHostname}, hostnames(resp.GetSchedulers()))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			ms := mocks.NewMockSearcher(ctl)
			tc.mock(ms.EXPECT())

			s := newTestServerV1(t, ms)
			if tc.cached != nil {
				setCache(t, s.cache, cacheKey, tc.cached)
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
		reqs   []*managerv1.KeepAliveRequest
		err    error
		expect func(t *testing.T, s *managerServerV1, states []string, err error)
	}{
		{
			name: "scheduler turns active on first message and inactive when stream closes",
			reqs: []*managerv1.KeepAliveRequest{schedulerReq, schedulerReq},
			err:  io.EOF,
			expect: func(t *testing.T, s *managerServerV1, states []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{models.SchedulerStateInactive, models.SchedulerStateActive, models.SchedulerStateActive}, states)
				scheduler := findScheduler(t, s.db, mockInactiveSchedulerHostname)
				assert.Equal(models.SchedulerStateInactive, scheduler.State)
				assert.True(scheduler.LastKeepAliveAt.IsZero())
				assert.False(s.cache.Exists(context.Background(), pkgredis.MakeSchedulerKeyInManager(mockSchedulerClusterID, mockInactiveSchedulerHostname, mockInactiveSchedulerIP)))
			},
		},
		{
			name: "seed peer turns active on first message and inactive when stream closes",
			reqs: []*managerv1.KeepAliveRequest{seedPeerReq},
			err:  io.EOF,
			expect: func(t *testing.T, s *managerServerV1, states []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{models.SeedPeerStateInactive, models.SeedPeerStateActive}, states)
				assert.Equal(models.SeedPeerStateInactive, findSeedPeer(t, s.db, mockInactiveSeedPeerHostname).State)
				assert.False(s.cache.Exists(context.Background(), pkgredis.MakeSeedPeerKeyInManager(mockSeedPeerClusterID, mockInactiveSeedPeerHostname, mockInactiveSeedPeerIP)))
			},
		},
		{
			name: "stream failure marks seed peer inactive and returns unknown",
			reqs: []*managerv1.KeepAliveRequest{seedPeerReq},
			err:  errors.New("connection reset"),
			expect: func(t *testing.T, s *managerServerV1, states []string, err error) {
				assert := assert.New(t)
				assert.Equal(codes.Unknown, status.Code(err))
				assert.Equal([]string{models.SeedPeerStateInactive, models.SeedPeerStateActive}, states)
				assert.Equal(models.SeedPeerStateInactive, findSeedPeer(t, s.db, mockInactiveSeedPeerHostname).State)
			},
		},
		{
			name: "unknown seed peer returns internal",
			reqs: []*managerv1.KeepAliveRequest{{
				SourceType: managerv1.SourceType_SEED_PEER_SOURCE,
				Hostname:   mockUnknownHostname,
				Ip:         mockUnknownIP,
				ClusterId:  uint64(mockSeedPeerClusterID),
			}},
			err: io.EOF,
			expect: func(t *testing.T, _ *managerServerV1, _ []string, err error) {
				assert := assert.New(t)
				assert.Equal(codes.Internal, status.Code(err))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestServerV1(t, nil)
			first := tc.reqs[0]
			setCache(t, s.cache, pkgredis.MakeSchedulerKeyInManager(uint(first.ClusterId), first.Hostname, first.Ip), &managerv1.Scheduler{Hostname: mockCachedHostname})
			setCache(t, s.cache, pkgredis.MakeSeedPeerKeyInManager(uint(first.ClusterId), first.Hostname, first.Ip), &managerv1.SeedPeer{Hostname: mockCachedHostname})
			stream := &mockKeepAliveStreamV1{reqs: tc.reqs, err: tc.err, observe: func() string {
				return sourceState(s.db, first.SourceType == managerv1.SourceType_SCHEDULER_SOURCE, first.Hostname, first.Ip, uint(first.ClusterId))
			}}

			err := s.KeepAlive(stream)
			tc.expect(t, s, stream.states, err)
		})
	}
}
