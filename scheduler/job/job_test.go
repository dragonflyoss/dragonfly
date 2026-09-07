/*
 *     Copyright 2024 The Dragonfly Authors
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

package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	dfdaemonv2 "d7y.io/api/v2/pkg/apis/dfdaemon/v2"
	dfdaemonv2mocks "d7y.io/api/v2/pkg/apis/dfdaemon/v2/mocks"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	internaljob "d7y.io/dragonfly/v2/internal/job"
	managertypes "d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/pkg/idgen"
	dfdaemonclientmocks "d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/client/mocks"
	pkgtypes "d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
	resource "d7y.io/dragonfly/v2/scheduler/resource/standard"
)

var (
	mockAdvertiseIP = "127.0.0.1"
	mockTaskURL     = "http://example.com/foo"
	mockTaskID      = idgen.TaskIDV2ByURLBased(mockTaskURL, nil, "", "", idgen.ParseFilteredQueryParams(""), "")
	mockGroupUUID   = "group"
	mockTaskUUID    = "task"

	mockConfig = &config.Config{
		Server:  config.ServerConfig{AdvertiseIP: net.ParseIP(mockAdvertiseIP)},
		Manager: config.ManagerConfig{SchedulerClusterID: 1},
	}
)

func newMockHost(hostname, ip string, typ pkgtypes.HostType) *resource.Host {
	return resource.NewHost(idgen.HostID(ip, hostname, typ != pkgtypes.HostTypeNormal), ip, hostname, hostname, 8003, 8001, 8004, typ)
}

func newMockHosts(n int, typ pkgtypes.HostType) []*resource.Host {
	hosts := make([]*resource.Host, 0, n)
	for i := range n {
		hosts = append(hosts, newMockHost(fmt.Sprintf("host%d", i+1), fmt.Sprintf("127.0.0.%d", i+1), typ))
	}

	return hosts
}

func ptr[T any](v T) *T {
	return &v
}

func TestJob_preheat(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient)
		expect func(t *testing.T, result string, err error)
	}{
		{
			name: "invalid json request",
			data: "foo",
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(result)
			},
		},
		{
			name: "invalid url fails validation",
			data: `{"url":"foo","urls":["foo"],"timeout":1000000000}`,
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "percentage over 100 fails validation",
			data: `{"urls":["http://example.com/foo"],"percentage":101,"timeout":1000000000}`,
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "count over 200 fails validation",
			data: `{"urls":["http://example.com/foo"],"count":201,"timeout":1000000000}`,
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "single seed peer scope without available seed peer",
			data: `{"urls":["http://example.com/foo"],"scope":"single_seed_peer","timeout":1000000000}`,
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				seedPeer.EXPECT().HasAvailable().Return(false).Times(1)
			},
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "all seed peers scope without available seed peer",
			data: `{"urls":["http://example.com/foo"],"scope":"all_seed_peers","timeout":1000000000}`,
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				seedPeer.EXPECT().HasAvailable().Return(false).Times(1)
			},
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "unknown scope falls back to single seed peer",
			data: `{"urls":["http://example.com/foo"],"scope":"foo","timeout":1000000000}`,
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				seedPeer.EXPECT().HasAvailable().Return(false).Times(1)
			},
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "all peers scope without available peer",
			data: `{"urls":["http://example.com/foo"],"scope":"all_peers","timeout":1000000000}`,
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				hostManager.EXPECT().LoadAllNormals().Return(nil).Times(1)
			},
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "all peers scope marshals preheat response",
			data: `{"urls":["http://example.com/foo"],"scope":"all_peers","timeout":1000000000,"concurrent_task_count":1,"concurrent_peer_count":1}`,
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				hostManager.EXPECT().LoadAllNormals().Return(newMockHosts(1, pkgtypes.HostTypeNormal)).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
			},
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var resp internaljob.PreheatResponse
				assert.NoError(json.Unmarshal([]byte(result), &resp))
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Equal([]*internaljob.PreheatSuccessTask{{URL: mockTaskURL, Hostname: "host1", IP: "127.0.0.1"}}, resp.SuccessTasks)
				assert.Empty(resp.FailureTasks)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			seedPeer := resource.NewMockSeedPeer(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			pool := dfdaemonclientmocks.NewMockPool(ctl)
			client := dfdaemonclientmocks.NewMockClient(ctl)
			stream := dfdaemonv2mocks.NewMockDfdaemonUpload_DownloadTaskClient(ctl)
			res.EXPECT().SeedPeer().Return(seedPeer).AnyTimes()
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			res.EXPECT().PeerClientPool().Return(pool).AnyTimes()
			if tc.mock != nil {
				tc.mock(seedPeer, hostManager, pool, client, stream)
			}

			j := &job{resource: res, config: mockConfig}
			result, err := j.preheat(context.Background(), tc.data)
			tc.expect(t, result, err)
		})
	}
}

func TestJob_PreheatSingleSeedPeer(t *testing.T) {
	tests := []struct {
		name   string
		urls   []string
		mock   func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient, seedHost *resource.Host)
		expect func(t *testing.T, seedHost *resource.Host, resp *internaljob.PreheatResponse, err error)
	}{
		{
			name: "no available seed peer",
			urls: []string{mockTaskURL},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient, seedHost *resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(false).Times(1)
			},
			expect: func(t *testing.T, seedHost *resource.Host, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "preheat succeeded and host is resolved from the last piece",
			urls: []string{mockTaskURL},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient, seedHost *resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				seedPeer.EXPECT().Select(gomock.Any(), mockTaskID).Return(seedHost, nil).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(1)
				gomock.InOrder(
					stream.EXPECT().Recv().Return(&dfdaemonv2.DownloadTaskResponse{HostId: seedHost.ID}, nil).Times(1),
					stream.EXPECT().Recv().Return(nil, io.EOF).Times(1),
				)
				hostManager.EXPECT().Load(seedHost.ID).Return(seedHost, true).Times(1)
			},
			expect: func(t *testing.T, seedHost *resource.Host, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Equal([]*internaljob.PreheatSuccessTask{{URL: mockTaskURL, Hostname: seedHost.Hostname, IP: seedHost.IP}}, resp.SuccessTasks)
				assert.Empty(resp.FailureTasks)
			},
		},
		{
			name: "preheat succeeded but host is unknown",
			urls: []string{mockTaskURL},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient, seedHost *resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				seedPeer.EXPECT().Select(gomock.Any(), mockTaskID).Return(seedHost, nil).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
				hostManager.EXPECT().Load("").Return(nil, false).Times(1)
			},
			expect: func(t *testing.T, seedHost *resource.Host, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]*internaljob.PreheatSuccessTask{{URL: mockTaskURL, Hostname: "unknown", IP: "unknown"}}, resp.SuccessTasks)
			},
		},
		{
			name: "success tasks of multiple urls are aggregated",
			urls: []string{mockTaskURL, "http://example.com/bar"},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient, seedHost *resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				seedPeer.EXPECT().Select(gomock.Any(), gomock.Any()).Return(seedHost, nil).Times(2)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(2)
				client.EXPECT().DownloadTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(stream, nil).Times(2)
				stream.EXPECT().Recv().Return(nil, io.EOF).Times(2)
				hostManager.EXPECT().Load("").Return(seedHost, true).Times(2)
			},
			expect: func(t *testing.T, seedHost *resource.Host, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.ElementsMatch([]*internaljob.PreheatSuccessTask{
					{URL: mockTaskURL, Hostname: seedHost.Hostname, IP: seedHost.IP},
					{URL: "http://example.com/bar", Hostname: seedHost.Hostname, IP: seedHost.IP},
				}, resp.SuccessTasks)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			seedPeer := resource.NewMockSeedPeer(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			pool := dfdaemonclientmocks.NewMockPool(ctl)
			client := dfdaemonclientmocks.NewMockClient(ctl)
			stream := dfdaemonv2mocks.NewMockDfdaemonUpload_DownloadTaskClient(ctl)
			seedHost := newMockHost("bar", "127.0.0.1", pkgtypes.HostTypeSuperSeed)
			res.EXPECT().SeedPeer().Return(seedPeer).AnyTimes()
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			res.EXPECT().PeerClientPool().Return(pool).AnyTimes()
			tc.mock(seedPeer, hostManager, pool, client, stream, seedHost)

			j := &job{resource: res, config: mockConfig}
			resp, err := j.PreheatSingleSeedPeer(context.Background(), &internaljob.PreheatRequest{
				URLs:                tc.urls,
				ConcurrentTaskCount: 1,
				GroupUUID:           mockGroupUUID,
				TaskUUID:            mockTaskUUID,
			}, logger.WithPreheatJob(mockGroupUUID, mockTaskUUID, tc.urls))
			tc.expect(t, seedHost, resp, err)
		})
	}
}

func TestJob_PreheatAllSeedPeers(t *testing.T) {
	tests := []struct {
		name   string
		urls   []string
		mock   func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient)
		expect func(t *testing.T, resp *internaljob.PreheatResponse, err error)
	}{
		{
			name: "no available seed peer",
			urls: []string{mockTaskURL},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				seedPeer.EXPECT().HasAvailable().Return(false).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "all seed peers succeed",
			urls: []string{mockTaskURL},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(newMockHosts(2, pkgtypes.HostTypeSuperSeed)).Times(1)
				pool.EXPECT().Get(gomock.Any()).Return(client, nil).Times(2)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(2)
				stream.EXPECT().Recv().Return(nil, io.EOF).Times(2)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.ElementsMatch([]*internaljob.PreheatSuccessTask{
					{URL: mockTaskURL, Hostname: "host1", IP: "127.0.0.1"},
					{URL: mockTaskURL, Hostname: "host2", IP: "127.0.0.2"},
				}, resp.SuccessTasks)
				assert.Empty(resp.FailureTasks)
			},
		},
		{
			name: "seed peer whose client cannot be obtained is reported as failure",
			urls: []string{mockTaskURL},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(newMockHosts(2, pkgtypes.HostTypeSuperSeed)).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(nil, errors.New("foo")).Times(1)
				pool.EXPECT().Get("127.0.0.2:8003").Return(client, nil).Times(1)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]*internaljob.PreheatSuccessTask{{URL: mockTaskURL, Hostname: "host2", IP: "127.0.0.2"}}, resp.SuccessTasks)
				assert.Equal([]*internaljob.PreheatFailureTask{{URL: mockTaskURL, Hostname: "host1", IP: "127.0.0.1", Description: "group uuid group failed: foo"}}, resp.FailureTasks)
			},
		},
		{
			name: "download task request failure makes the whole preheat fail",
			urls: []string{mockTaskURL},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(newMockHosts(1, pkgtypes.HostTypeSuperSeed)).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.Nil(resp)
				assert.Error(err)
			},
		},
		{
			name: "stream receive failure makes the whole preheat fail",
			urls: []string{mockTaskURL},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(newMockHosts(1, pkgtypes.HostTypeSuperSeed)).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Recv().Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "success on a seed peer that also failed is dropped",
			urls: []string{mockTaskURL, "http://example.com/bar"},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(newMockHosts(1, pkgtypes.HostTypeSuperSeed)).Times(1)
				gomock.InOrder(
					pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1),
					pool.EXPECT().Get("127.0.0.1:8003").Return(nil, errors.New("foo")).Times(1),
				)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			seedPeer := resource.NewMockSeedPeer(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			pool := dfdaemonclientmocks.NewMockPool(ctl)
			client := dfdaemonclientmocks.NewMockClient(ctl)
			stream := dfdaemonv2mocks.NewMockDfdaemonUpload_DownloadTaskClient(ctl)
			res.EXPECT().SeedPeer().Return(seedPeer).AnyTimes()
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			res.EXPECT().PeerClientPool().Return(pool).AnyTimes()
			tc.mock(seedPeer, hostManager, pool, client, stream)

			j := &job{resource: res, config: mockConfig}
			resp, err := j.PreheatAllSeedPeers(context.Background(), &internaljob.PreheatRequest{
				URLs:                tc.urls,
				ConcurrentTaskCount: 1,
				ConcurrentPeerCount: 2,
				GroupUUID:           mockGroupUUID,
				TaskUUID:            mockTaskUUID,
			}, logger.WithPreheatJob(mockGroupUUID, mockTaskUUID, tc.urls))
			tc.expect(t, resp, err)
		})
	}
}

func TestJob_PreheatAllPeers(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient)
		expect func(t *testing.T, resp *internaljob.PreheatResponse, err error)
	}{
		{
			name: "no available peer",
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				hostManager.EXPECT().LoadAllNormals().Return(nil).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "all peers succeed",
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				hostManager.EXPECT().LoadAllNormals().Return(newMockHosts(2, pkgtypes.HostTypeNormal)).Times(1)
				pool.EXPECT().Get(gomock.Any()).Return(client, nil).Times(2)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(2)
				stream.EXPECT().Recv().Return(nil, io.EOF).Times(2)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.ElementsMatch([]*internaljob.PreheatSuccessTask{
					{URL: mockTaskURL, Hostname: "host1", IP: "127.0.0.1"},
					{URL: mockTaskURL, Hostname: "host2", IP: "127.0.0.2"},
				}, resp.SuccessTasks)
				assert.Empty(resp.FailureTasks)
			},
		},
		{
			name: "peer whose client cannot be obtained is reported as failure",
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				hostManager.EXPECT().LoadAllNormals().Return(newMockHosts(2, pkgtypes.HostTypeNormal)).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(nil, errors.New("foo")).Times(1)
				pool.EXPECT().Get("127.0.0.2:8003").Return(client, nil).Times(1)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]*internaljob.PreheatSuccessTask{{URL: mockTaskURL, Hostname: "host2", IP: "127.0.0.2"}}, resp.SuccessTasks)
				assert.Equal([]*internaljob.PreheatFailureTask{{URL: mockTaskURL, Hostname: "host1", IP: "127.0.0.1", Description: fmt.Sprintf("task %s failed: foo", mockTaskID)}}, resp.FailureTasks)
			},
		},
		{
			name: "all peers failed returns error",
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, stream *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient) {
				hostManager.EXPECT().LoadAllNormals().Return(newMockHosts(1, pkgtypes.HostTypeNormal)).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().DownloadTask(gomock.Any(), mockTaskID, gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Recv().Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.PreheatResponse, err error) {
				assert := assert.New(t)
				assert.Nil(resp)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			pool := dfdaemonclientmocks.NewMockPool(ctl)
			client := dfdaemonclientmocks.NewMockClient(ctl)
			stream := dfdaemonv2mocks.NewMockDfdaemonUpload_DownloadTaskClient(ctl)
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			res.EXPECT().PeerClientPool().Return(pool).AnyTimes()
			tc.mock(hostManager, pool, client, stream)

			j := &job{resource: res, config: mockConfig}
			resp, err := j.PreheatAllPeers(context.Background(), &internaljob.PreheatRequest{
				URLs:                []string{mockTaskURL},
				ConcurrentTaskCount: 1,
				ConcurrentPeerCount: 2,
				GroupUUID:           mockGroupUUID,
				TaskUUID:            mockTaskUUID,
			}, logger.WithPreheatJob(mockGroupUUID, mockTaskUUID, []string{mockTaskURL}))
			tc.expect(t, resp, err)
		})
	}
}

func TestJob_selectSeedPeers(t *testing.T) {
	tests := []struct {
		name       string
		hosts      []*resource.Host
		ips        []string
		count      *uint32
		percentage *uint32
		mock       func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host)
		expect     func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error)
	}{
		{
			name: "no available seed peer",
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(false).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "seed peer list is empty",
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:  "ips select only matching seed peers",
			hosts: newMockHosts(3, pkgtypes.HostTypeSuperSeed),
			ips:   []string{"127.0.0.2", "10.0.0.1"},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[1:2], seedPeers)
			},
		},
		{
			name:       "ips take priority over count and percentage",
			hosts:      newMockHosts(3, pkgtypes.HostTypeSuperSeed),
			ips:        []string{"127.0.0.3"},
			count:      ptr[uint32](3),
			percentage: ptr[uint32](100),
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[2:], seedPeers)
			},
		},
		{
			name:  "no seed peer matches ips",
			hosts: newMockHosts(3, pkgtypes.HostTypeSuperSeed),
			ips:   []string{"10.0.0.1"},
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:  "count selects the first n seed peers",
			hosts: newMockHosts(3, pkgtypes.HostTypeSuperSeed),
			count: ptr[uint32](2),
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[:2], seedPeers)
			},
		},
		{
			name:  "count over available seed peers is clamped",
			hosts: newMockHosts(3, pkgtypes.HostTypeSuperSeed),
			count: ptr[uint32](5),
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts, seedPeers)
			},
		},
		{
			name:       "count takes priority over percentage",
			hosts:      newMockHosts(3, pkgtypes.HostTypeSuperSeed),
			count:      ptr[uint32](1),
			percentage: ptr[uint32](100),
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[:1], seedPeers)
			},
		},
		{
			name:       "percentage over 100 is clamped to available seed peers",
			hosts:      newMockHosts(2, pkgtypes.HostTypeSuperSeed),
			percentage: ptr[uint32](200),
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts, seedPeers)
			},
		},
		{
			name:       "valid percentage selects a proportional number of seed peers",
			hosts:      newMockHosts(4, pkgtypes.HostTypeSuperSeed),
			percentage: ptr[uint32](50),
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[:2], seedPeers)
			},
		},
		{
			name:       "percentage rounding down to zero still selects one seed peer",
			hosts:      newMockHosts(3, pkgtypes.HostTypeSuperSeed),
			percentage: ptr[uint32](10),
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[:1], seedPeers)
			},
		},
		{
			name:       "zero percentage selects no seed peer",
			hosts:      newMockHosts(3, pkgtypes.HostTypeSuperSeed),
			percentage: ptr[uint32](0),
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(seedPeers)
			},
		},
		{
			name:  "nil count and percentage select all seed peers",
			hosts: newMockHosts(3, pkgtypes.HostTypeSuperSeed),
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, hosts []*resource.Host) {
				seedPeer.EXPECT().HasAvailable().Return(true).Times(1)
				hostManager.EXPECT().LoadAllSeeds().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, seedPeers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts, seedPeers)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			seedPeer := resource.NewMockSeedPeer(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			res.EXPECT().SeedPeer().Return(seedPeer).AnyTimes()
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			tc.mock(seedPeer, hostManager, tc.hosts)

			j := &job{resource: res, config: mockConfig}
			seedPeers, err := j.selectSeedPeers(tc.ips, tc.count, tc.percentage, logger.WithPeerID("test"))
			tc.expect(t, tc.hosts, seedPeers, err)
		})
	}
}

func TestJob_selectPeers(t *testing.T) {
	tests := []struct {
		name       string
		hosts      []*resource.Host
		ips        []string
		count      *uint32
		percentage *uint32
		mock       func(hostManager *resource.MockHostManager, hosts []*resource.Host)
		expect     func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error)
	}{
		{
			name: "peer list is empty",
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:  "ips select only matching peers",
			hosts: newMockHosts(3, pkgtypes.HostTypeNormal),
			ips:   []string{"127.0.0.2", "10.0.0.1"},
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[1:2], peers)
			},
		},
		{
			name:       "ips take priority over count and percentage",
			hosts:      newMockHosts(3, pkgtypes.HostTypeNormal),
			ips:        []string{"127.0.0.3"},
			count:      ptr[uint32](3),
			percentage: ptr[uint32](100),
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[2:], peers)
			},
		},
		{
			name:  "no peer matches ips",
			hosts: newMockHosts(3, pkgtypes.HostTypeNormal),
			ips:   []string{"10.0.0.1"},
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:  "count selects the first n peers",
			hosts: newMockHosts(3, pkgtypes.HostTypeNormal),
			count: ptr[uint32](2),
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[:2], peers)
			},
		},
		{
			name:  "count over available peers is clamped",
			hosts: newMockHosts(3, pkgtypes.HostTypeNormal),
			count: ptr[uint32](5),
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts, peers)
			},
		},
		{
			name:       "count takes priority over percentage",
			hosts:      newMockHosts(3, pkgtypes.HostTypeNormal),
			count:      ptr[uint32](1),
			percentage: ptr[uint32](100),
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[:1], peers)
			},
		},
		{
			name:       "percentage over 100 is clamped to available peers",
			hosts:      newMockHosts(3, pkgtypes.HostTypeNormal),
			percentage: ptr[uint32](200),
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts, peers)
			},
		},
		{
			name:       "valid percentage selects a proportional number of peers",
			hosts:      newMockHosts(4, pkgtypes.HostTypeNormal),
			percentage: ptr[uint32](50),
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[:2], peers)
			},
		},
		{
			name:       "percentage rounding down to zero still selects one peer",
			hosts:      newMockHosts(3, pkgtypes.HostTypeNormal),
			percentage: ptr[uint32](10),
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts[:1], peers)
			},
		},
		{
			name:       "zero percentage selects no peer",
			hosts:      newMockHosts(3, pkgtypes.HostTypeNormal),
			percentage: ptr[uint32](0),
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(peers)
			},
		},
		{
			name:  "nil count and percentage select all peers",
			hosts: newMockHosts(3, pkgtypes.HostTypeNormal),
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().LoadAllNormals().Return(hosts).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, peers []*resource.Host, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(hosts, peers)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			tc.mock(hostManager, tc.hosts)

			j := &job{resource: res, config: mockConfig}
			peers, err := j.selectPeers(tc.ips, tc.count, tc.percentage, logger.WithPeerID("test"))
			tc.expect(t, tc.hosts, peers, err)
		})
	}
}

func TestJob_syncPeers(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(hostManager *resource.MockHostManager, hosts []*resource.Host)
		expect func(t *testing.T, hosts []*resource.Host, result string, err error)
	}{
		{
			name: "all hosts are marshalled",
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().Len().Return(len(hosts)).Times(1)
				hostManager.EXPECT().Range(gomock.Any()).Do(func(f func(any, any) bool) {
					for _, host := range hosts {
						f(host.ID, host)
					}
				}).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, result string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var got []struct{ ID string }
				assert.NoError(json.Unmarshal([]byte(result), &got))
				assert.Equal([]struct{ ID string }{{hosts[0].ID}, {hosts[1].ID}}, got)
			},
		},
		{
			name: "values that are not hosts are skipped",
			mock: func(hostManager *resource.MockHostManager, hosts []*resource.Host) {
				hostManager.EXPECT().Len().Return(len(hosts) + 1).Times(1)
				hostManager.EXPECT().Range(gomock.Any()).Do(func(f func(any, any) bool) {
					f(hosts[0].ID, hosts[0])
					f("foo", "bar")
					f(hosts[1].ID, hosts[1])
				}).Times(1)
			},
			expect: func(t *testing.T, hosts []*resource.Host, result string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var got []struct{ ID string }
				assert.NoError(json.Unmarshal([]byte(result), &got))
				assert.Equal([]struct{ ID string }{{hosts[0].ID}, {hosts[1].ID}}, got)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			hosts := newMockHosts(2, pkgtypes.HostTypeNormal)
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			tc.mock(hostManager, hosts)

			j := &job{resource: res, config: mockConfig}
			result, err := j.syncPeers()
			tc.expect(t, hosts, result, err)
		})
	}
}

func TestJob_getTask(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(hostManager *resource.MockHostManager)
		expect func(t *testing.T, result string, err error)
	}{
		{
			name: "invalid json request",
			data: "foo",
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(result)
			},
		},
		{
			name: "missing task id fails validation",
			data: `{"timeout":1000000000}`,
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "response is marshalled",
			data: `{"task_id":"foo","timeout":1000000000,"concurrent_peer_count":1}`,
			mock: func(hostManager *resource.MockHostManager) {
				hostManager.EXPECT().LoadAll().Return(nil).Times(1)
			},
			expect: func(t *testing.T, result string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var resp internaljob.GetTaskResponse
				assert.NoError(json.Unmarshal([]byte(result), &resp))
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Empty(resp.Peers)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			if tc.mock != nil {
				tc.mock(hostManager)
			}

			j := &job{resource: res, config: mockConfig}
			result, err := j.getTask(context.Background(), tc.data)
			tc.expect(t, result, err)
		})
	}
}

func TestJob_GetTask(t *testing.T) {
	tests := []struct {
		name   string
		scope  string
		mock   func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, host *resource.Host)
		expect func(t *testing.T, host *resource.Host, resp *internaljob.GetTaskResponse, err error)
	}{
		{
			name:  "all_seed_peers scope queries seed peers only",
			scope: managertypes.AllSeedPeersScope,
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, host *resource.Host) {
				hostManager.EXPECT().LoadAllSeeds().Return(nil).Times(1)
			},
			expect: func(t *testing.T, host *resource.Host, resp *internaljob.GetTaskResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Empty(resp.Peers)
			},
		},
		{
			name:  "all_peers scope queries all peers",
			scope: managertypes.AllPeersScope,
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, host *resource.Host) {
				hostManager.EXPECT().LoadAll().Return(nil).Times(1)
			},
			expect: func(t *testing.T, host *resource.Host, resp *internaljob.GetTaskResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Empty(resp.Peers)
			},
		},
		{
			name:  "empty scope queries all peers",
			scope: "",
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, host *resource.Host) {
				hostManager.EXPECT().LoadAll().Return(nil).Times(1)
			},
			expect: func(t *testing.T, host *resource.Host, resp *internaljob.GetTaskResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Empty(resp.Peers)
			},
		},
		{
			name:  "peer whose client cannot be obtained is skipped",
			scope: managertypes.AllPeersScope,
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, host *resource.Host) {
				hostManager.EXPECT().LoadAll().Return([]*resource.Host{host}).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, host *resource.Host, resp *internaljob.GetTaskResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Empty(resp.Peers)
			},
		},
		{
			name:  "peer whose stat fails is skipped",
			scope: managertypes.AllPeersScope,
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, host *resource.Host) {
				hostManager.EXPECT().LoadAll().Return([]*resource.Host{host}).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().StatLocalTask(gomock.Any(), gomock.Eq(&dfdaemonv2.StatLocalTaskRequest{TaskId: mockTaskID, RemoteIp: &mockAdvertiseIP})).Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, host *resource.Host, resp *internaljob.GetTaskResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Empty(resp.Peers)
			},
		},
		{
			name:  "finished peer is reported",
			scope: managertypes.AllPeersScope,
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, host *resource.Host) {
				hostManager.EXPECT().LoadAll().Return([]*resource.Host{host}).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().StatLocalTask(gomock.Any(), gomock.Eq(&dfdaemonv2.StatLocalTaskRequest{TaskId: mockTaskID, RemoteIp: &mockAdvertiseIP})).Return(&dfdaemonv2.StatLocalTaskResponse{FinishedAt: timestamppb.Now()}, nil).Times(1)
			},
			expect: func(t *testing.T, host *resource.Host, resp *internaljob.GetTaskResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]*internaljob.Peer{{
					ID:         host.ID,
					Hostname:   host.Hostname,
					IP:         host.IP,
					HostType:   pkgtypes.HostTypeNormalName,
					CreatedAt:  host.CreatedAt.Load(),
					UpdatedAt:  host.UpdatedAt.Load(),
					IsFinished: true,
				}}, resp.Peers)
			},
		},
		{
			name:  "unfinished peer is reported",
			scope: managertypes.AllPeersScope,
			mock: func(hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, host *resource.Host) {
				hostManager.EXPECT().LoadAll().Return([]*resource.Host{host}).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().StatLocalTask(gomock.Any(), gomock.Any()).Return(&dfdaemonv2.StatLocalTaskResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, host *resource.Host, resp *internaljob.GetTaskResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]*internaljob.Peer{{
					ID:        host.ID,
					Hostname:  host.Hostname,
					IP:        host.IP,
					HostType:  pkgtypes.HostTypeNormalName,
					CreatedAt: host.CreatedAt.Load(),
					UpdatedAt: host.UpdatedAt.Load(),
				}}, resp.Peers)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			pool := dfdaemonclientmocks.NewMockPool(ctl)
			client := dfdaemonclientmocks.NewMockClient(ctl)
			host := newMockHost("foo", "127.0.0.1", pkgtypes.HostTypeNormal)
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			res.EXPECT().PeerClientPool().Return(pool).AnyTimes()
			tc.mock(hostManager, pool, client, host)

			j := &job{resource: res, config: mockConfig}
			resp, err := j.GetTask(context.Background(), &internaljob.GetTaskRequest{TaskID: mockTaskID, Scope: tc.scope, ConcurrentPeerCount: 1}, logger.WithPeerID("test"))
			tc.expect(t, host, resp, err)
		})
	}
}

func TestJob_deleteTask(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(taskManager *resource.MockTaskManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, task *resource.Task, peer *resource.Peer)
		expect func(t *testing.T, task *resource.Task, peer *resource.Peer, result string, err error)
	}{
		{
			name: "invalid json request",
			data: "foo",
			expect: func(t *testing.T, task *resource.Task, peer *resource.Peer, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(result)
			},
		},
		{
			name: "missing task id fails validation",
			data: `{"timeout":1000000000}`,
			expect: func(t *testing.T, task *resource.Task, peer *resource.Peer, result string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "task not found returns empty response",
			data: fmt.Sprintf(`{"task_id":"%s","timeout":1000000000}`, mockTaskID),
			mock: func(taskManager *resource.MockTaskManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, task *resource.Task, peer *resource.Peer) {
				taskManager.EXPECT().Load(mockTaskID).Return(nil, false).Times(1)
			},
			expect: func(t *testing.T, task *resource.Task, peer *resource.Peer, result string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var resp internaljob.DeleteTaskResponse
				assert.NoError(json.Unmarshal([]byte(result), &resp))
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Empty(resp.SuccessTasks)
				assert.Empty(resp.FailureTasks)
			},
		},
		{
			name: "unfinished peers are ignored",
			data: fmt.Sprintf(`{"task_id":"%s","timeout":1000000000}`, mockTaskID),
			mock: func(taskManager *resource.MockTaskManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, task *resource.Task, peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateRunning)
				task.StorePeer(peer)
				taskManager.EXPECT().Load(mockTaskID).Return(task, true).Times(1)
			},
			expect: func(t *testing.T, task *resource.Task, peer *resource.Peer, result string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var resp internaljob.DeleteTaskResponse
				assert.NoError(json.Unmarshal([]byte(result), &resp))
				assert.Empty(resp.SuccessTasks)
				assert.Empty(resp.FailureTasks)
				_, loaded := task.LoadPeer(peer.ID)
				assert.True(loaded)
			},
		},
		{
			name: "peer whose client cannot be obtained is reported as failure",
			data: fmt.Sprintf(`{"task_id":"%s","timeout":1000000000}`, mockTaskID),
			mock: func(taskManager *resource.MockTaskManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, task *resource.Task, peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateSucceeded)
				task.StorePeer(peer)
				taskManager.EXPECT().Load(mockTaskID).Return(task, true).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, task *resource.Task, peer *resource.Peer, result string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var resp internaljob.DeleteTaskResponse
				assert.NoError(json.Unmarshal([]byte(result), &resp))
				assert.Empty(resp.SuccessTasks)
				assert.Equal([]*internaljob.DeleteFailureTask{{Hostname: "foo", IP: "127.0.0.1", HostType: pkgtypes.HostTypeNormalName, Description: fmt.Sprintf("task %s failed: foo", mockTaskID)}}, resp.FailureTasks)
				_, loaded := task.LoadPeer(peer.ID)
				assert.True(loaded)
			},
		},
		{
			name: "delete task rpc failure is reported as failure",
			data: fmt.Sprintf(`{"task_id":"%s","timeout":1000000000}`, mockTaskID),
			mock: func(taskManager *resource.MockTaskManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, task *resource.Task, peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateFailed)
				task.StorePeer(peer)
				taskManager.EXPECT().Load(mockTaskID).Return(task, true).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().DeleteTask(gomock.Any(), gomock.Eq(&dfdaemonv2.DeleteTaskRequest{TaskId: mockTaskID, RemoteIp: &mockAdvertiseIP})).Return(errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, task *resource.Task, peer *resource.Peer, result string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var resp internaljob.DeleteTaskResponse
				assert.NoError(json.Unmarshal([]byte(result), &resp))
				assert.Empty(resp.SuccessTasks)
				assert.Equal([]*internaljob.DeleteFailureTask{{Hostname: "foo", IP: "127.0.0.1", HostType: pkgtypes.HostTypeNormalName, Description: fmt.Sprintf("task %s failed: foo", mockTaskID)}}, resp.FailureTasks)
				_, loaded := task.LoadPeer(peer.ID)
				assert.True(loaded)
			},
		},
		{
			name: "deleted peer is removed from task",
			data: fmt.Sprintf(`{"task_id":"%s","timeout":1000000000}`, mockTaskID),
			mock: func(taskManager *resource.MockTaskManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient, task *resource.Task, peer *resource.Peer) {
				peer.FSM.SetState(resource.PeerStateLeave)
				task.StorePeer(peer)
				taskManager.EXPECT().Load(mockTaskID).Return(task, true).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().DeleteTask(gomock.Any(), gomock.Eq(&dfdaemonv2.DeleteTaskRequest{TaskId: mockTaskID, RemoteIp: &mockAdvertiseIP})).Return(nil).Times(1)
			},
			expect: func(t *testing.T, task *resource.Task, peer *resource.Peer, result string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var resp internaljob.DeleteTaskResponse
				assert.NoError(json.Unmarshal([]byte(result), &resp))
				assert.Equal(uint(1), resp.SchedulerClusterID)
				assert.Equal([]*internaljob.DeleteSuccessTask{{Hostname: "foo", IP: "127.0.0.1", HostType: pkgtypes.HostTypeNormalName}}, resp.SuccessTasks)
				assert.Empty(resp.FailureTasks)
				_, loaded := task.LoadPeer(peer.ID)
				assert.False(loaded)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			taskManager := resource.NewMockTaskManager(ctl)
			pool := dfdaemonclientmocks.NewMockPool(ctl)
			client := dfdaemonclientmocks.NewMockClient(ctl)
			task := resource.NewTask(mockTaskID, mockTaskURL, "", "", commonv2.TaskType_STANDARD, nil, nil, 200)
			peer := resource.NewPeer(idgen.PeerID(), task, newMockHost("foo", "127.0.0.1", pkgtypes.HostTypeNormal))
			res.EXPECT().TaskManager().Return(taskManager).AnyTimes()
			res.EXPECT().PeerClientPool().Return(pool).AnyTimes()
			if tc.mock != nil {
				tc.mock(taskManager, pool, client, task, peer)
			}

			j := &job{resource: res, config: mockConfig}
			result, err := j.deleteTask(context.Background(), tc.data)
			tc.expect(t, task, peer, result, err)
		})
	}
}

func TestJob_ListTaskEntries(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient)
		expect func(t *testing.T, resp *internaljob.ListTaskEntriesResponse, err error)
	}{
		{
			name: "normal peer is preferred and multiple entries are recursive",
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient) {
				hostManager.EXPECT().LoadAllNormals().Return(newMockHosts(1, pkgtypes.HostTypeNormal)).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().ListTaskEntries(gomock.Any(), gomock.Any()).Return(&dfdaemonv2.ListTaskEntriesResponse{
					Entries: []*dfdaemonv2.Entry{{Url: "http://example.com/foo/a"}, {Url: "http://example.com/foo/b"}},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.ListTaskEntriesResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(resp.Recursive)
				assert.Len(resp.Entries, 2)
				assert.Equal(uint(1), resp.SchedulerID)
			},
		},
		{
			name: "seed peer is used when no normal peer is available and a single entry is not recursive",
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient) {
				hostManager.EXPECT().LoadAllNormals().Return(nil).Times(1)
				seedPeer.EXPECT().Select(gomock.Any(), mockTaskID).Return(newMockHost("bar", "127.0.0.2", pkgtypes.HostTypeSuperSeed), nil).Times(1)
				pool.EXPECT().Get("127.0.0.2:8003").Return(client, nil).Times(1)
				client.EXPECT().ListTaskEntries(gomock.Any(), gomock.Any()).Return(&dfdaemonv2.ListTaskEntriesResponse{
					Entries: []*dfdaemonv2.Entry{{Url: mockTaskURL}},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.ListTaskEntriesResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.False(resp.Recursive)
				assert.Len(resp.Entries, 1)
			},
		},
		{
			name: "seed peer selection failed",
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient) {
				hostManager.EXPECT().LoadAllNormals().Return(nil).Times(1)
				seedPeer.EXPECT().Select(gomock.Any(), mockTaskID).Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.ListTaskEntriesResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "client cannot be obtained",
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient) {
				hostManager.EXPECT().LoadAllNormals().Return(newMockHosts(1, pkgtypes.HostTypeNormal)).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.ListTaskEntriesResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "list task entries rpc failed",
			mock: func(seedPeer *resource.MockSeedPeer, hostManager *resource.MockHostManager, pool *dfdaemonclientmocks.MockPool, client *dfdaemonclientmocks.MockClient) {
				hostManager.EXPECT().LoadAllNormals().Return(newMockHosts(1, pkgtypes.HostTypeNormal)).Times(1)
				pool.EXPECT().Get("127.0.0.1:8003").Return(client, nil).Times(1)
				client.EXPECT().ListTaskEntries(gomock.Any(), gomock.Any()).Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, resp *internaljob.ListTaskEntriesResponse, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			res := resource.NewMockResource(ctl)
			seedPeer := resource.NewMockSeedPeer(ctl)
			hostManager := resource.NewMockHostManager(ctl)
			pool := dfdaemonclientmocks.NewMockPool(ctl)
			client := dfdaemonclientmocks.NewMockClient(ctl)
			res.EXPECT().SeedPeer().Return(seedPeer).AnyTimes()
			res.EXPECT().HostManager().Return(hostManager).AnyTimes()
			res.EXPECT().PeerClientPool().Return(pool).AnyTimes()
			tc.mock(seedPeer, hostManager, pool, client)

			j := &job{resource: res, config: mockConfig}
			resp, err := j.ListTaskEntries(context.Background(), &internaljob.ListTaskEntriesRequest{TaskID: mockTaskID, Url: mockTaskURL}, logger.WithTaskID(mockTaskID))
			tc.expect(t, resp, err)
		})
	}
}
