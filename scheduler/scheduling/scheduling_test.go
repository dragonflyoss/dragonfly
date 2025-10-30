/*
 *     Copyright 2020 The Dragonfly Authors
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

package scheduling

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"
	schedulerv1mocks "d7y.io/api/v2/pkg/apis/scheduler/v1/mocks"
	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"
	schedulerv2mocks "d7y.io/api/v2/pkg/apis/scheduler/v2/mocks"

	"d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/pkg/container/set"
	"d7y.io/dragonfly/v2/pkg/digest"
	"d7y.io/dragonfly/v2/pkg/idgen"
	nethttp "d7y.io/dragonfly/v2/pkg/net/http"
	pkgtypes "d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
	configmocks "d7y.io/dragonfly/v2/scheduler/config/mocks"
	"d7y.io/dragonfly/v2/scheduler/resource/persistentcache"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
	"d7y.io/dragonfly/v2/scheduler/scheduling/evaluator"
)

var (
	mockPluginDir       = "bas"
	mockSchedulerConfig = &config.SchedulerConfig{
		RetryLimit:             2,
		RetryBackToSourceLimit: 1,
		RetryInterval:          10 * time.Millisecond,
		BackToSourceCount:      int(mockTaskBackToSourceLimit),
		Algorithm:              evaluator.DefaultAlgorithm,
	}

	mockRawHost = standard.Host{
		ID:              mockHostID,
		Type:            pkgtypes.HostTypeNormal,
		Hostname:        "foo",
		IP:              "127.0.0.1",
		Port:            8003,
		DownloadPort:    8001,
		ProxyPort:       8004,
		OS:              "darwin",
		Platform:        "darwin",
		PlatformFamily:  "Standalone Workstation",
		PlatformVersion: "11.1",
		KernelVersion:   "20.2.0",
		CPU:             mockCPU,
		Memory:          mockMemory,
		Network:         mockNetwork,
		Disk:            mockDisk,
		Build:           mockBuild,
		CreatedAt:       atomic.NewTime(time.Now()),
		UpdatedAt:       atomic.NewTime(time.Now()),
	}

	mockRawSeedHost = standard.Host{
		ID:              mockSeedHostID,
		Type:            pkgtypes.HostTypeSuperSeed,
		Hostname:        "bar",
		IP:              "127.0.0.1",
		Port:            8003,
		DownloadPort:    8001,
		ProxyPort:       8004,
		OS:              "darwin",
		Platform:        "darwin",
		PlatformFamily:  "Standalone Workstation",
		PlatformVersion: "11.1",
		KernelVersion:   "20.2.0",
		CPU:             mockCPU,
		Memory:          mockMemory,
		Network:         mockNetwork,
		Disk:            mockDisk,
		Build:           mockBuild,
		CreatedAt:       atomic.NewTime(time.Now()),
		UpdatedAt:       atomic.NewTime(time.Now()),
	}

	mockCPU = standard.CPU{
		LogicalCount:   4,
		PhysicalCount:  2,
		Percent:        1,
		ProcessPercent: 0.5,
		Times: standard.CPUTimes{
			User:      240662.2,
			System:    317950.1,
			Idle:      3393691.3,
			Nice:      0,
			Iowait:    0,
			Irq:       0,
			Softirq:   0,
			Steal:     0,
			Guest:     0,
			GuestNice: 0,
		},
	}

	mockMemory = standard.Memory{
		Total:              17179869184,
		Available:          5962813440,
		Used:               11217055744,
		UsedPercent:        65.291858,
		ProcessUsedPercent: 41.525125,
		Free:               2749598908,
	}

	mockNetwork = standard.Network{
		TCPConnectionCount:       10,
		UploadTCPConnectionCount: 1,
		Location:                 mockHostLocation,
		IDC:                      mockHostIDC,
		RxBandwidth:              100,
		MaxRxBandwidth:           200,
		TxBandwidth:              100,
		MaxTxBandwidth:           200,
	}

	mockDisk = standard.Disk{
		Total:             499963174912,
		Free:              37226479616,
		Used:              423809622016,
		UsedPercent:       91.92547406065952,
		InodesTotal:       4882452880,
		InodesUsed:        7835772,
		InodesFree:        4874617108,
		InodesUsedPercent: 0.1604884305611568,
	}

	mockBuild = standard.Build{
		GitVersion: "v1.0.0",
		GitCommit:  "221176b117c6d59366d68f2b34d38be50c935883",
		GoVersion:  "1.18",
		Platform:   "darwin",
	}

	mockTaskBackToSourceLimit   int32  = 200
	mockTaskURL                        = "http://example.com/foo"
	mockTaskPieceLength         uint64 = 2048
	mockTaskID                         = idgen.TaskIDV2ByURLBased(mockTaskURL, &mockTaskPieceLength, mockTaskTag, mockTaskApplication, mockTaskFilteredQueryParams)
	mockTaskDigest                     = digest.New(digest.AlgorithmSHA256, "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4")
	mockTaskTag                        = "d7y"
	mockTaskApplication                = "foo"
	mockTaskFilteredQueryParams        = []string{"bar"}
	mockTaskHeader                     = map[string]string{"content-length": "100"}
	mockHostID                         = idgen.HostIDV2("127.0.0.1", "foo", false)
	mockSeedHostID                     = idgen.HostIDV2("127.0.0.1", "bar", true)
	mockHostLocation                   = "baz"
	mockHostIDC                        = "bas"
	mockPeerID                         = idgen.PeerIDV2()
	mockSeedPeerID                     = idgen.PeerIDV2()
	mockPiece                          = standard.Piece{
		Number:      1,
		ParentID:    "foo",
		Offset:      2,
		Length:      10,
		Digest:      digest.New(digest.AlgorithmMD5, "1f70f5a1630d608a71442c54ab706638"),
		TrafficType: commonv2.TrafficType_REMOTE_PEER,
		Cost:        1 * time.Minute,
		CreatedAt:   time.Now(),
	}
)

func TestScheduling_New(t *testing.T) {
	tests := []struct {
		name      string
		pluginDir string
		expect    func(t *testing.T, s any)
	}{
		{
			name:      "new scheduling",
			pluginDir: "bar",
			expect: func(t *testing.T, s any) {
				assert := assert.New(t)
				assert.Equal(reflect.TypeOf(s).Elem().Name(), "scheduling")
			},
		},
		{
			name:      "new scheduling with empty pluginDir",
			pluginDir: "",
			expect: func(t *testing.T, s any) {
				assert := assert.New(t)
				assert.Equal(reflect.TypeOf(s).Elem().Name(), "scheduling")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			dynconfig := configmocks.NewMockDynconfigInterface(ctl)
			persistentCacheResource := persistentcache.NewMockResource(ctl)

			tc.expect(t, New(mockSchedulerConfig, persistentCacheResource, dynconfig, tc.pluginDir))
		})
	}
}

func TestScheduling_ScheduleCandidateParents(t *testing.T) {
	needBackToSourceDescription := "peer's NeedBackToSource is true"
	exceededLimitDescription := "scheduling exceeded RetryBackToSourceLimit"

	tests := []struct {
		name   string
		mock   func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder)
		expect func(t *testing.T, peer *standard.Peer, err error)
	}{
		{
			name: "context was done",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				cancel()
			},
			expect: func(t *testing.T, peer *standard.Peer, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, context.Canceled)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "peer needs back-to-source and peer stream load failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.NeedBackToSource.Store(true)
				peer.FSM.SetState(standard.PeerStateRunning)
			},
			expect: func(t *testing.T, peer *standard.Peer, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, status.Error(codes.FailedPrecondition, "load stream failed"))
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "peer needs back-to-source and send NeedBackToSourceResponse failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.NeedBackToSource.Store(true)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreAnnouncePeerStream(stream)

				ma.Send(gomock.Eq(&schedulerv2.AnnouncePeerResponse{
					Response: &schedulerv2.AnnouncePeerResponse_NeedBackToSourceResponse{
						NeedBackToSourceResponse: &schedulerv2.NeedBackToSourceResponse{
							Description: &needBackToSourceDescription,
						},
					},
				})).Return(errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, status.Error(codes.FailedPrecondition, "foo"))
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "peer needs back-to-source and send NeedBackToSourceResponse success",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.NeedBackToSource.Store(true)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreAnnouncePeerStream(stream)

				ma.Send(gomock.Eq(&schedulerv2.AnnouncePeerResponse{
					Response: &schedulerv2.AnnouncePeerResponse_NeedBackToSourceResponse{
						NeedBackToSourceResponse: &schedulerv2.NeedBackToSourceResponse{
							Description: &needBackToSourceDescription,
						},
					},
				})).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule exceeds RetryBackToSourceLimit and peer stream load failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, status.Error(codes.FailedPrecondition, "load stream failed"))
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule exceeds RetryBackToSourceLimit and send NeedBackToSourceResponse failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreAnnouncePeerStream(stream)

				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1),
					ma.Send(gomock.Eq(&schedulerv2.AnnouncePeerResponse{
						Response: &schedulerv2.AnnouncePeerResponse_NeedBackToSourceResponse{
							NeedBackToSourceResponse: &schedulerv2.NeedBackToSourceResponse{
								Description: &exceededLimitDescription,
							},
						},
					})).Return(errors.New("foo")).Times(1),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, status.Error(codes.FailedPrecondition, "foo"))
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule exceeds RetryBackToSourceLimit and send NeedBackToSourceResponse success",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreAnnouncePeerStream(stream)

				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1),
					ma.Send(gomock.Eq(&schedulerv2.AnnouncePeerResponse{
						Response: &schedulerv2.AnnouncePeerResponse_NeedBackToSourceResponse{
							NeedBackToSourceResponse: &schedulerv2.NeedBackToSourceResponse{
								Description: &exceededLimitDescription,
							},
						},
					})).Return(nil).Times(1),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule exceeds RetryLimit",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.Task.BackToSourceLimit.Store(-1)
				peer.StoreAnnouncePeerStream(stream)

				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, status.Error(codes.FailedPrecondition, "scheduling exceeded RetryLimit"))
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule succeeded",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv2.Scheduler_AnnouncePeerServer, ma *schedulerv2mocks.MockScheduler_AnnouncePeerServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				task.StorePeer(seedPeer)
				peer.FSM.SetState(standard.PeerStateRunning)
				seedPeer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreAnnouncePeerStream(stream)
				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2),
					ma.Send(gomock.Any()).Return(nil).Times(1),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(len(peer.Parents()), 1)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			stream := schedulerv2mocks.NewMockScheduler_AnnouncePeerServer(ctl)
			persistentCacheResource := persistentcache.NewMockResource(ctl)
			dynconfig := configmocks.NewMockDynconfigInterface(ctl)
			ctx, cancel := context.WithCancel(context.Background())
			mockHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, mockHost)
			mockSeedHost := standard.NewHost(
				mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
				mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type)
			seedPeer := standard.NewPeer(mockSeedPeerID, mockTask, mockSeedHost)
			blocklist := set.NewSafeSet[string]()

			tc.mock(cancel, peer, seedPeer, blocklist, stream, stream.EXPECT(), dynconfig.EXPECT())
			scheduling := New(mockSchedulerConfig, persistentCacheResource, dynconfig, mockPluginDir)
			tc.expect(t, peer, scheduling.ScheduleCandidateParents(ctx, peer, blocklist))
		})
	}
}

func TestScheduling_ScheduleParentAndCandidateParents(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder)
		expect func(t *testing.T, peer *standard.Peer)
	}{
		{
			name: "context was done",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				cancel()
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "peer needs back-to-source and peer stream load failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.NeedBackToSource.Store(true)
				peer.FSM.SetState(standard.PeerStateRunning)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "peer needs back-to-source and send Code_SchedNeedBackSource failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.NeedBackToSource.Store(true)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreReportPieceResultStream(stream)

				mr.Send(gomock.Eq(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedNeedBackSource})).Return(errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "peer needs back-to-source and send Code_SchedNeedBackSource success",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.NeedBackToSource.Store(true)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreReportPieceResultStream(stream)

				mr.Send(gomock.Eq(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedNeedBackSource})).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateBackToSource))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "peer needs back-to-source and task state is TaskStateFailed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.NeedBackToSource.Store(true)
				peer.FSM.SetState(standard.PeerStateRunning)
				task.FSM.SetState(standard.TaskStateFailed)
				peer.StoreReportPieceResultStream(stream)

				mr.Send(gomock.Eq(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedNeedBackSource})).Return(nil).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateBackToSource))
				assert.True(peer.Task.FSM.Is(standard.TaskStateRunning))
			},
		},
		{
			name: "schedule exceeds RetryBackToSourceLimit and peer stream load failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule exceeds RetryBackToSourceLimit and send Code_SchedNeedBackSource failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreReportPieceResultStream(stream)

				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1),
					mr.Send(gomock.Eq(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedNeedBackSource})).Return(errors.New("foo")).Times(1),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule exceeds RetryBackToSourceLimit and send Code_SchedNeedBackSource success",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreReportPieceResultStream(stream)

				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1),
					mr.Send(gomock.Eq(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedNeedBackSource})).Return(nil).Times(1),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateBackToSource))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule exceeds RetryBackToSourceLimit and  task state is TaskStateFailed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				task.FSM.SetState(standard.TaskStateFailed)
				peer.StoreReportPieceResultStream(stream)

				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1),
					mr.Send(gomock.Eq(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedNeedBackSource})).Return(nil).Times(1),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateBackToSource))
				assert.True(peer.Task.FSM.Is(standard.TaskStateRunning))
			},
		},
		{
			name: "schedule exceeds RetryLimit and peer stream load failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.Task.BackToSourceLimit.Store(-1)
				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule exceeds RetryLimit and send Code_SchedTaskStatusError failed",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.Task.BackToSourceLimit.Store(-1)
				peer.StoreReportPieceResultStream(stream)

				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2),
					mr.Send(gomock.Eq(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedTaskStatusError})).Return(errors.New("foo")).Times(1),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule exceeds RetryLimit and send Code_SchedTaskStatusError success",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.Task.BackToSourceLimit.Store(-1)
				peer.StoreReportPieceResultStream(stream)

				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2),
					mr.Send(gomock.Eq(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedTaskStatusError})).Return(nil).Times(1),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 0)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
		{
			name: "schedule succeeded",
			mock: func(cancel context.CancelFunc, peer *standard.Peer, seedPeer *standard.Peer, blocklist set.SafeSet[string], stream schedulerv1.Scheduler_ReportPieceResultServer, mr *schedulerv1mocks.MockScheduler_ReportPieceResultServerMockRecorder, md *configmocks.MockDynconfigInterfaceMockRecorder) {
				task := peer.Task
				task.StorePeer(peer)
				task.StorePeer(seedPeer)
				peer.FSM.SetState(standard.PeerStateRunning)
				seedPeer.FSM.SetState(standard.PeerStateRunning)
				peer.StoreReportPieceResultStream(stream)
				gomock.InOrder(
					md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2),
					mr.Send(gomock.Any()).Return(nil).Times(1),
				)
			},
			expect: func(t *testing.T, peer *standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(peer.Parents()), 1)
				assert.True(peer.FSM.Is(standard.PeerStateRunning))
				assert.True(peer.Task.FSM.Is(standard.TaskStatePending))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			stream := schedulerv1mocks.NewMockScheduler_ReportPieceResultServer(ctl)
			persistentCacheResource := persistentcache.NewMockResource(ctl)
			dynconfig := configmocks.NewMockDynconfigInterface(ctl)
			ctx, cancel := context.WithCancel(context.Background())
			mockHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, mockHost)
			mockSeedHost := standard.NewHost(
				mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
				mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type)
			seedPeer := standard.NewPeer(mockSeedPeerID, mockTask, mockSeedHost)
			blocklist := set.NewSafeSet[string]()

			tc.mock(cancel, peer, seedPeer, blocklist, stream, stream.EXPECT(), dynconfig.EXPECT())
			scheduling := New(mockSchedulerConfig, persistentCacheResource, dynconfig, mockPluginDir)
			scheduling.ScheduleParentAndCandidateParents(ctx, peer, blocklist)
			tc.expect(t, peer)
		})
	}
}

func TestScheduling_FindCandidateParents(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder)
		expect func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool)
	}{
		{
			name: "task peers state is failed",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateFailed)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.Equal(len(parents), 0)
				assert.False(ok)
			},
		},
		{
			name: "task peers is empty",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "task contains only one peer and peer is itself",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateReceivedNormal)
				peer.Task.StorePeer(peer)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "peer is in blocklist",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				blocklist.Add(mockPeers[0].ID)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "peer is bad node",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateReceivedNormal)
				mockPeers[0].FSM.SetState(standard.PeerStateFailed)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "parent is peer's descendant",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				if err := peer.Task.AddPeerEdge(peer, mockPeers[0]); err != nil {
					t.Fatal(err)
				}

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "parent free upload load is zero",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateReceivedNormal)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				mockPeers[0].Host.ConcurrentUploadLimit.Store(0)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "parent is disabled share data with other peers",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateReceivedNormal)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				mockPeers[0].Host.DisableShared = true

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "find back-to-source parent",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				peer.Task.BackToSourcePeers.Add(mockPeers[0].ID)
				mockPeers[0].FSM.SetState(standard.PeerStateBackToSource)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(mockPeers[0].ID, parents[0].ID)
			},
		},
		{
			name: "find normal parent",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				peer.Task.StorePeer(mockPeers[2])
				mockPeers[0].Host.Type = pkgtypes.HostTypeNormal
				mockPeers[1].Host.Type = pkgtypes.HostTypeSuperSeed

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(mockPeers[1].ID, parents[0].ID)
			},
		},
		{
			name: "parent state is PeerStateSucceeded",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateSucceeded)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(mockPeers[0].ID, parents[0].ID)
			},
		},
		{
			name: "find parent with ancestor",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				mockPeers[2].FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].Host.Type = pkgtypes.HostTypeNormal
				mockPeers[1].Host.Type = pkgtypes.HostTypeNormal
				mockPeers[2].Host.Type = pkgtypes.HostTypeSuperSeed
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				peer.Task.StorePeer(mockPeers[2])
				if err := peer.Task.AddPeerEdge(mockPeers[2], mockPeers[0]); err != nil {
					t.Fatal(err)
				}

				if err := peer.Task.AddPeerEdge(mockPeers[2], mockPeers[1]); err != nil {
					t.Fatal(err)
				}

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(mockPeers[2].ID, parents[0].ID)
			},
		},
		{
			name: "find parent and fetch candidateParentLimit from manager dynconfig",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				peer.Task.BackToSourcePeers.Add(mockPeers[0].ID)
				peer.Task.BackToSourcePeers.Add(mockPeers[1].ID)
				mockPeers[0].FSM.SetState(standard.PeerStateBackToSource)
				mockPeers[1].FSM.SetState(standard.PeerStateBackToSource)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{
					CandidateParentLimit: 3,
				}, nil).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Contains([]string{mockPeers[0].ID, mockPeers[1].ID, peer.ID}, parents[0].ID)
			},
		},
		{
			name: "candidateParents is longer than candidateParentLimit",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				peer.Task.BackToSourcePeers.Add(mockPeers[0].ID)
				mockPeers[0].FSM.SetState(standard.PeerStateBackToSource)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{
					CandidateParentLimit: 1,
				}, nil).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(len(parents), 1)
				assert.Equal(mockPeers[0].ID, parents[0].ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			dynconfig := configmocks.NewMockDynconfigInterface(ctl)
			persistentCacheResource := persistentcache.NewMockResource(ctl)
			mockHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, mockHost)

			var mockPeers []*standard.Peer
			for i := range 11 {
				mockHost := standard.NewHost(
					idgen.HostIDV2("127.0.0.1", uuid.New().String(), false), mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
				peer := standard.NewPeer(idgen.PeerIDV1(fmt.Sprintf("127.0.0.%d", i)), mockTask, mockHost)
				mockPeers = append(mockPeers, peer)
			}

			blocklist := set.NewSafeSet[string]()
			tc.mock(peer, mockPeers, blocklist, dynconfig.EXPECT())
			scheduling := New(mockSchedulerConfig, persistentCacheResource, dynconfig, mockPluginDir)
			parents, found := scheduling.FindCandidateParents(context.Background(), peer, blocklist)
			tc.expect(t, peer, mockPeers, parents, found)
		})
	}
}

func TestScheduling_FindParentAndCandidateParents(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder)
		expect func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool)
	}{
		{
			name: "task peers state is failed",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateFailed)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.Equal(len(parents), 0)
				assert.False(ok)
			},
		},
		{
			name: "task peers is empty",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "task contains only one peer and peer is itself",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "peer is in blocklist",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				blocklist.Add(mockPeers[0].ID)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "peer is bad node",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateFailed)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "parent is peer's descendant",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				if err := peer.Task.AddPeerEdge(peer, mockPeers[0]); err != nil {
					t.Fatal(err)
				}

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "find back-to-source parent",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				peer.Task.BackToSourcePeers.Add(mockPeers[0].ID)
				mockPeers[0].FSM.SetState(standard.PeerStateBackToSource)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(mockPeers[0].ID, parents[0].ID)
			},
		},
		{
			name: "find normal parent",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				mockPeers[0].Host.Type = pkgtypes.HostTypeSuperSeed
				mockPeers[1].Host.Type = pkgtypes.HostTypeNormal

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(mockPeers[0].ID, parents[0].ID)
			},
		},
		{
			name: "parent state is PeerStateSucceeded",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateSucceeded)
				mockPeers[1].FSM.SetState(standard.PeerStateSucceeded)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				mockPeers[0].Host.Type = pkgtypes.HostTypeSuperSeed
				mockPeers[1].Host.Type = pkgtypes.HostTypeNormal

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(mockPeers[1].ID, parents[0].ID)
			},
		},
		{
			name: "find parent with ancestor",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				mockPeers[2].FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].Host.Type = pkgtypes.HostTypeNormal
				mockPeers[1].Host.Type = pkgtypes.HostTypeNormal
				mockPeers[2].Host.Type = pkgtypes.HostTypeSuperSeed
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				peer.Task.StorePeer(mockPeers[2])
				if err := peer.Task.AddPeerEdge(mockPeers[2], mockPeers[0]); err != nil {
					t.Fatal(err)
				}

				if err := peer.Task.AddPeerEdge(mockPeers[2], mockPeers[1]); err != nil {
					t.Fatal(err)
				}

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{}, errors.New("foo")).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(mockPeers[2].ID, parents[0].ID)
			},
		},
		{
			name: "find parent and fetch candidateParentLimit from manager dynconfig",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				peer.Task.BackToSourcePeers.Add(mockPeers[0].ID)
				peer.Task.BackToSourcePeers.Add(mockPeers[1].ID)
				mockPeers[0].FSM.SetState(standard.PeerStateBackToSource)
				mockPeers[1].FSM.SetState(standard.PeerStateBackToSource)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{
					CandidateParentLimit: 3,
				}, nil).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Contains([]string{mockPeers[0].ID, mockPeers[1].ID, peer.ID}, parents[0].ID)
			},
		},
		{
			name: "candidateParents is longer than candidateParentLimit",
			mock: func(peer *standard.Peer, mockPeers []*standard.Peer, blocklist set.SafeSet[string], md *configmocks.MockDynconfigInterfaceMockRecorder) {
				peer.FSM.SetState(standard.PeerStateRunning)
				mockPeers[0].FSM.SetState(standard.PeerStateRunning)
				mockPeers[1].FSM.SetState(standard.PeerStateRunning)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(mockPeers[0])
				peer.Task.StorePeer(mockPeers[1])
				peer.Task.BackToSourcePeers.Add(mockPeers[0].ID)
				mockPeers[0].FSM.SetState(standard.PeerStateBackToSource)

				md.GetSchedulerClusterConfig().Return(types.SchedulerClusterConfig{
					CandidateParentLimit: 1,
				}, nil).Times(2)
			},
			expect: func(t *testing.T, peer *standard.Peer, mockPeers []*standard.Peer, parents []*standard.Peer, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(len(parents), 1)
				assert.Equal(mockPeers[0].ID, parents[0].ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			dynconfig := configmocks.NewMockDynconfigInterface(ctl)
			persistentCacheResource := persistentcache.NewMockResource(ctl)
			mockHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, mockHost)

			var mockPeers []*standard.Peer
			for i := range 11 {
				mockHost := standard.NewHost(
					idgen.HostIDV2("127.0.0.1", uuid.New().String(), false), mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
				peer := standard.NewPeer(idgen.PeerIDV1(fmt.Sprintf("127.0.0.%d", i)), mockTask, mockHost)
				mockPeers = append(mockPeers, peer)
			}

			blocklist := set.NewSafeSet[string]()
			tc.mock(peer, mockPeers, blocklist, dynconfig.EXPECT())
			scheduling := New(mockSchedulerConfig, persistentCacheResource, dynconfig, mockPluginDir)
			parents, found := scheduling.FindParentAndCandidateParents(context.Background(), peer, blocklist)
			tc.expect(t, peer, mockPeers, parents, found)
		})
	}
}

func TestScheduling_constructSuccessNormalTaskResponse(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, resp *schedulerv2.AnnouncePeerResponse_NormalTaskResponse, candidateParents []*standard.Peer)
	}{
		{
			name: "construct success normal task response",
			expect: func(t *testing.T, resp *schedulerv2.AnnouncePeerResponse_NormalTaskResponse, candidateParents []*standard.Peer) {
				dgst := candidateParents[0].Task.Digest.String()

				assert := assert.New(t)
				assert.EqualValues(resp, &schedulerv2.AnnouncePeerResponse_NormalTaskResponse{
					NormalTaskResponse: &schedulerv2.NormalTaskResponse{
						CandidateParents: []*commonv2.Peer{
							{
								Id: candidateParents[0].ID,
								Range: &commonv2.Range{
									Start:  uint64(candidateParents[0].Range.Start),
									Length: uint64(candidateParents[0].Range.Length),
								},
								Priority:             candidateParents[0].Priority,
								ConcurrentPieceCount: candidateParents[0].ConcurrentPieceCount,
								Cost:                 durationpb.New(candidateParents[0].Cost.Load()),
								State:                candidateParents[0].FSM.Current(),
								Task: &commonv2.Task{
									Id:                  candidateParents[0].Task.ID,
									Type:                candidateParents[0].Task.Type,
									Url:                 candidateParents[0].Task.URL,
									Digest:              &dgst,
									Tag:                 &candidateParents[0].Task.Tag,
									Application:         &candidateParents[0].Task.Application,
									FilteredQueryParams: candidateParents[0].Task.FilteredQueryParams,
									RequestHeader:       candidateParents[0].Task.Header,
									ContentLength:       uint64(candidateParents[0].Task.ContentLength.Load()),
									PieceCount:          uint32(candidateParents[0].Task.TotalPieceCount.Load()),
									SizeScope:           candidateParents[0].Task.SizeScope(),
									State:               candidateParents[0].Task.FSM.Current(),
									PeerCount:           uint32(candidateParents[0].Task.PeerCount()),
									CreatedAt:           timestamppb.New(candidateParents[0].Task.CreatedAt.Load()),
									UpdatedAt:           timestamppb.New(candidateParents[0].Task.UpdatedAt.Load()),
								},
								Host: &commonv2.Host{
									Id:              candidateParents[0].Host.ID,
									Type:            uint32(candidateParents[0].Host.Type),
									Hostname:        candidateParents[0].Host.Hostname,
									Ip:              candidateParents[0].Host.IP,
									Port:            candidateParents[0].Host.Port,
									DownloadPort:    candidateParents[0].Host.DownloadPort,
									ProxyPort:       candidateParents[0].Host.ProxyPort,
									Os:              candidateParents[0].Host.OS,
									Platform:        candidateParents[0].Host.Platform,
									PlatformFamily:  candidateParents[0].Host.PlatformFamily,
									PlatformVersion: candidateParents[0].Host.PlatformVersion,
									KernelVersion:   candidateParents[0].Host.KernelVersion,
									Cpu: &commonv2.CPU{
										LogicalCount:   candidateParents[0].Host.CPU.LogicalCount,
										PhysicalCount:  candidateParents[0].Host.CPU.PhysicalCount,
										Percent:        candidateParents[0].Host.CPU.Percent,
										ProcessPercent: candidateParents[0].Host.CPU.ProcessPercent,
										Times: &commonv2.CPUTimes{
											User:      candidateParents[0].Host.CPU.Times.User,
											System:    candidateParents[0].Host.CPU.Times.System,
											Idle:      candidateParents[0].Host.CPU.Times.Idle,
											Nice:      candidateParents[0].Host.CPU.Times.Nice,
											Iowait:    candidateParents[0].Host.CPU.Times.Iowait,
											Irq:       candidateParents[0].Host.CPU.Times.Irq,
											Softirq:   candidateParents[0].Host.CPU.Times.Softirq,
											Steal:     candidateParents[0].Host.CPU.Times.Steal,
											Guest:     candidateParents[0].Host.CPU.Times.Guest,
											GuestNice: candidateParents[0].Host.CPU.Times.GuestNice,
										},
									},
									Memory: &commonv2.Memory{
										Total:              candidateParents[0].Host.Memory.Total,
										Available:          candidateParents[0].Host.Memory.Available,
										Used:               candidateParents[0].Host.Memory.Used,
										UsedPercent:        candidateParents[0].Host.Memory.UsedPercent,
										ProcessUsedPercent: candidateParents[0].Host.Memory.ProcessUsedPercent,
										Free:               candidateParents[0].Host.Memory.Free,
									},
									Network: &commonv2.Network{
										TcpConnectionCount:       candidateParents[0].Host.Network.TCPConnectionCount,
										UploadTcpConnectionCount: candidateParents[0].Host.Network.UploadTCPConnectionCount,
										Location:                 &candidateParents[0].Host.Network.Location,
										Idc:                      &candidateParents[0].Host.Network.IDC,
										RxBandwidth:              &candidateParents[0].Host.Network.RxBandwidth,
										MaxRxBandwidth:           candidateParents[0].Host.Network.MaxRxBandwidth,
										TxBandwidth:              &candidateParents[0].Host.Network.TxBandwidth,
										MaxTxBandwidth:           candidateParents[0].Host.Network.MaxTxBandwidth,
									},
									Disk: &commonv2.Disk{
										Total:             candidateParents[0].Host.Disk.Total,
										Free:              candidateParents[0].Host.Disk.Free,
										Used:              candidateParents[0].Host.Disk.Used,
										UsedPercent:       candidateParents[0].Host.Disk.UsedPercent,
										InodesTotal:       candidateParents[0].Host.Disk.InodesTotal,
										InodesUsed:        candidateParents[0].Host.Disk.InodesUsed,
										InodesFree:        candidateParents[0].Host.Disk.InodesFree,
										InodesUsedPercent: candidateParents[0].Host.Disk.InodesUsedPercent,
									},
									Build: &commonv2.Build{
										GitVersion: candidateParents[0].Host.Build.GitVersion,
										GitCommit:  &candidateParents[0].Host.Build.GitCommit,
										GoVersion:  &candidateParents[0].Host.Build.GoVersion,
										Platform:   &candidateParents[0].Host.Build.Platform,
									},
								},
								NeedBackToSource: candidateParents[0].NeedBackToSource.Load(),
								CreatedAt:        timestamppb.New(candidateParents[0].CreatedAt.Load()),
								UpdatedAt:        timestamppb.New(candidateParents[0].UpdatedAt.Load()),
							},
						},
					},
				})
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			candidateParents := []*standard.Peer{standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost, standard.WithRange(nethttp.Range{
				Start:  1,
				Length: 10,
			}))}
			candidateParents[0].Task.StorePiece(&mockPiece)

			tc.expect(t, constructSuccessNormalTaskResponse(candidateParents), candidateParents)
		})
	}
}

func TestScheduling_constructSuccessPeerPacket(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, packet *schedulerv1.PeerPacket, parent *standard.Peer, candidateParents []*standard.Peer)
	}{
		{
			name: "construct success peer packet",
			expect: func(t *testing.T, packet *schedulerv1.PeerPacket, parent *standard.Peer, candidateParents []*standard.Peer) {
				assert := assert.New(t)
				assert.EqualValues(packet, &schedulerv1.PeerPacket{
					TaskId: mockTaskID,
					SrcPid: mockPeerID,
					MainPeer: &schedulerv1.PeerPacket_DestPeer{
						Ip:      parent.Host.IP,
						RpcPort: parent.Host.Port,
						PeerId:  parent.ID,
					},
					CandidatePeers: []*schedulerv1.PeerPacket_DestPeer{
						{
							Ip:      candidateParents[0].Host.IP,
							RpcPort: candidateParents[0].Host.Port,
							PeerId:  candidateParents[0].ID,
						},
					},
					Code: commonv1.Code_Success,
				})
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))

			peer := standard.NewPeer(mockPeerID, mockTask, mockHost)
			parent := standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost)
			candidateParents := []*standard.Peer{standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost)}

			tc.expect(t, constructSuccessPeerPacket(peer, parent, candidateParents), parent, candidateParents)
		})
	}
}
