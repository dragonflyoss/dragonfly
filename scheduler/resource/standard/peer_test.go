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

package standard

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-http-utils/headers"
	"github.com/looplab/fsm"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"
	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"
	v1mocks "d7y.io/api/v2/pkg/apis/scheduler/v1/mocks"
	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"
	schedulerv2mock "d7y.io/api/v2/pkg/apis/scheduler/v2/mocks"

	"d7y.io/dragonfly/v2/pkg/idgen"
	nethttp "d7y.io/dragonfly/v2/pkg/net/http"
	configmocks "d7y.io/dragonfly/v2/scheduler/config/mocks"
)

var (
	mockPeerID     = idgen.PeerID()
	mockSeedPeerID = idgen.PeerID() + "_Seed"
)

func TestPeer_NewPeer(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	stream := schedulerv2mock.NewMockScheduler_AnnouncePeerServer(ctl)

	tests := []struct {
		name    string
		id      string
		options []PeerOption
		expect  func(t *testing.T, peer *Peer, mockTask *Task, mockHost *Host)
	}{
		{
			name:    "new peer",
			id:      mockPeerID,
			options: []PeerOption{},
			expect: func(t *testing.T, peer *Peer, mockTask *Task, mockHost *Host) {
				assert := assert.New(t)
				assert.Equal(mockPeerID, peer.ID)
				assert.Nil(peer.Range)
				assert.Equal(commonv2.Priority_LEVEL0, peer.Priority)
				assert.Equal(defaultConcurrentPieceCount, peer.ConcurrentPieceCount)
				assert.Empty(peer.FinishedPieces)
				assert.Empty(peer.PieceCosts())
				assert.Empty(peer.ReportPieceResultStream)
				assert.Empty(peer.AnnouncePeerStream)
				assert.Equal(PeerStatePending, peer.FSM.Current())
				assert.EqualValues(mockTask, peer.Task)
				assert.EqualValues(mockHost, peer.Host)
				assert.Equal(uint(0), peer.BlockParents.Len())
				assert.False(peer.NeedBackToSource.Load())
				assert.NotEmpty(peer.PieceUpdatedAt.Load())
				assert.NotEmpty(peer.CreatedAt.Load())
				assert.NotEmpty(peer.UpdatedAt.Load())
				assert.NotNil(peer.Log)
			},
		},
		{
			name:    "new peer with priority",
			id:      mockPeerID,
			options: []PeerOption{WithPriority(commonv2.Priority_LEVEL4)},
			expect: func(t *testing.T, peer *Peer, mockTask *Task, mockHost *Host) {
				assert := assert.New(t)
				assert.Equal(mockPeerID, peer.ID)
				assert.Nil(peer.Range)
				assert.Equal(commonv2.Priority_LEVEL4, peer.Priority)
				assert.Equal(defaultConcurrentPieceCount, peer.ConcurrentPieceCount)
				assert.Empty(peer.FinishedPieces)
				assert.Empty(peer.PieceCosts())
				assert.Empty(peer.ReportPieceResultStream)
				assert.Empty(peer.AnnouncePeerStream)
				assert.Equal(PeerStatePending, peer.FSM.Current())
				assert.EqualValues(mockTask, peer.Task)
				assert.EqualValues(mockHost, peer.Host)
				assert.Equal(uint(0), peer.BlockParents.Len())
				assert.False(peer.NeedBackToSource.Load())
				assert.NotEmpty(peer.PieceUpdatedAt.Load())
				assert.NotEmpty(peer.CreatedAt.Load())
				assert.NotEmpty(peer.UpdatedAt.Load())
				assert.NotNil(peer.Log)
			},
		},
		{
			name: "new peer with range",
			id:   mockPeerID,
			options: []PeerOption{WithRange(nethttp.Range{
				Start:  1,
				Length: 10,
			})},
			expect: func(t *testing.T, peer *Peer, mockTask *Task, mockHost *Host) {
				assert := assert.New(t)
				assert.Equal(mockPeerID, peer.ID)
				assert.EqualValues(&nethttp.Range{Start: 1, Length: 10}, peer.Range)
				assert.Equal(commonv2.Priority_LEVEL0, peer.Priority)
				assert.Equal(defaultConcurrentPieceCount, peer.ConcurrentPieceCount)
				assert.Empty(peer.FinishedPieces)
				assert.Empty(peer.PieceCosts())
				assert.Empty(peer.ReportPieceResultStream)
				assert.Empty(peer.AnnouncePeerStream)
				assert.Equal(PeerStatePending, peer.FSM.Current())
				assert.EqualValues(mockTask, peer.Task)
				assert.EqualValues(mockHost, peer.Host)
				assert.Equal(uint(0), peer.BlockParents.Len())
				assert.False(peer.NeedBackToSource.Load())
				assert.NotEmpty(peer.PieceUpdatedAt.Load())
				assert.NotEmpty(peer.CreatedAt.Load())
				assert.NotEmpty(peer.UpdatedAt.Load())
				assert.NotNil(peer.Log)
			},
		},
		{
			name:    "new peer with AnnouncePeerStream",
			id:      mockPeerID,
			options: []PeerOption{WithAnnouncePeerStream(stream)},
			expect: func(t *testing.T, peer *Peer, mockTask *Task, mockHost *Host) {
				assert := assert.New(t)
				assert.Equal(mockPeerID, peer.ID)
				assert.Nil(peer.Range)
				assert.Equal(commonv2.Priority_LEVEL0, peer.Priority)
				assert.Equal(defaultConcurrentPieceCount, peer.ConcurrentPieceCount)
				assert.Empty(peer.FinishedPieces)
				assert.Empty(peer.PieceCosts())
				assert.Empty(peer.ReportPieceResultStream)
				announcePeerStream, loaded := peer.LoadAnnouncePeerStream()
				assert.True(loaded)
				assert.Equal(stream, announcePeerStream)
				assert.Equal(PeerStatePending, peer.FSM.Current())
				assert.EqualValues(mockTask, peer.Task)
				assert.EqualValues(mockHost, peer.Host)
				assert.Equal(uint(0), peer.BlockParents.Len())
				assert.False(peer.NeedBackToSource.Load())
				assert.NotEmpty(peer.PieceUpdatedAt.Load())
				assert.NotEmpty(peer.CreatedAt.Load())
				assert.NotEmpty(peer.UpdatedAt.Load())
				assert.NotNil(peer.Log)
			},
		},
		{
			name:    "new peer with ConcurrentPieceCount",
			id:      mockPeerID,
			options: []PeerOption{WithConcurrentPieceCount(1)},
			expect: func(t *testing.T, peer *Peer, mockTask *Task, mockHost *Host) {
				assert := assert.New(t)
				assert.Equal(mockPeerID, peer.ID)
				assert.Nil(peer.Range)
				assert.Equal(commonv2.Priority_LEVEL0, peer.Priority)
				assert.Equal(uint32(1), peer.ConcurrentPieceCount)
				assert.Empty(peer.FinishedPieces)
				assert.Empty(peer.PieceCosts())
				assert.Empty(peer.ReportPieceResultStream)
				assert.Empty(peer.AnnouncePeerStream)
				assert.Equal(PeerStatePending, peer.FSM.Current())
				assert.EqualValues(mockTask, peer.Task)
				assert.EqualValues(mockHost, peer.Host)
				assert.Equal(uint(0), peer.BlockParents.Len())
				assert.False(peer.NeedBackToSource.Load())
				assert.NotEmpty(peer.PieceUpdatedAt.Load())
				assert.NotEmpty(peer.CreatedAt.Load())
				assert.NotEmpty(peer.UpdatedAt.Load())
				assert.NotNil(peer.Log)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			tc.expect(t, NewPeer(tc.id, mockTask, mockHost, tc.options...), mockTask, mockHost)
		})
	}
}

func TestPeer_FSMEvent(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		event  string
		expect func(t *testing.T, peer *Peer, err error)
	}{
		{
			name:  "register normal from pending",
			state: PeerStatePending,
			event: PeerEventRegisterNormal,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateReceivedNormal, peer.FSM.Current())
			},
		},
		{
			name:  "register tiny from pending",
			state: PeerStatePending,
			event: PeerEventRegisterTiny,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateReceivedTiny, peer.FSM.Current())
			},
		},
		{
			name:  "download from received normal",
			state: PeerStateReceivedNormal,
			event: PeerEventDownload,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateRunning, peer.FSM.Current())
			},
		},
		{
			name:  "download back-to-source from running",
			state: PeerStateRunning,
			event: PeerEventDownloadBackToSource,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateBackToSource, peer.FSM.Current())
			},
		},
		{
			name:  "download succeeded from received small",
			state: PeerStateReceivedSmall,
			event: PeerEventDownloadSucceeded,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateSucceeded, peer.FSM.Current())
			},
		},
		{
			name:  "download failed from pending",
			state: PeerStatePending,
			event: PeerEventDownloadFailed,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateFailed, peer.FSM.Current())
			},
		},
		{
			name:  "download failed from succeeded",
			state: PeerStateSucceeded,
			event: PeerEventDownloadFailed,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateFailed, peer.FSM.Current())
			},
		},
		{
			name:  "leave from failed",
			state: PeerStateFailed,
			event: PeerEventLeave,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateLeave, peer.FSM.Current())
			},
		},
		{
			name:  "download from pending is rejected",
			state: PeerStatePending,
			event: PeerEventDownload,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStatePending, peer.FSM.Current())
			},
		},
		{
			name:  "register normal from running is rejected",
			state: PeerStateRunning,
			event: PeerEventRegisterNormal,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateRunning, peer.FSM.Current())
			},
		},
		{
			name:  "download from succeeded is rejected",
			state: PeerStateSucceeded,
			event: PeerEventDownload,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateSucceeded, peer.FSM.Current())
			},
		},
		{
			name:  "download back-to-source from succeeded is rejected",
			state: PeerStateSucceeded,
			event: PeerEventDownloadBackToSource,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateSucceeded, peer.FSM.Current())
			},
		},
		{
			name:  "download succeeded from pending is rejected",
			state: PeerStatePending,
			event: PeerEventDownloadSucceeded,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStatePending, peer.FSM.Current())
			},
		},
		{
			name:  "download succeeded from failed is rejected",
			state: PeerStateFailed,
			event: PeerEventDownloadSucceeded,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateFailed, peer.FSM.Current())
			},
		},
		{
			name:  "download failed from leave is rejected",
			state: PeerStateLeave,
			event: PeerEventDownloadFailed,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateLeave, peer.FSM.Current())
			},
		},
		{
			name:  "leave from leave is rejected",
			state: PeerStateLeave,
			event: PeerEventLeave,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateLeave, peer.FSM.Current())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			mockTask.StorePeer(peer)
			peer.FSM.SetState(tc.state)

			tc.expect(t, peer, peer.FSM.Event(context.Background(), tc.event))
		})
	}
}

func TestPeer_FSMCallback(t *testing.T) {
	tests := []struct {
		name            string
		state           string
		event           string
		backToSource    bool
		peerFailedCount int32
		expect          func(t *testing.T, peer *Peer, mockParent *Peer, err error)
	}{
		{
			name:  "register normal keeps the parent edge and upload load",
			state: PeerStatePending,
			event: PeerEventRegisterNormal,
			expect: func(t *testing.T, peer *Peer, mockParent *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peer.Parents(), 1)
				assert.Equal(int32(1), mockParent.Host.ConcurrentUploadCount.Load())
			},
		},
		{
			name:  "download back-to-source registers back-to-source peer and releases parent upload load",
			state: PeerStateRunning,
			event: PeerEventDownloadBackToSource,
			expect: func(t *testing.T, peer *Peer, mockParent *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(peer.Task.BackToSourcePeers.Contains(peer.ID))
				assert.Empty(peer.Parents())
				assert.Equal(int32(0), mockParent.Host.ConcurrentUploadCount.Load())
			},
		},
		{
			name:            "download succeeded from back-to-source clears back-to-source peer and resets peer failed count",
			state:           PeerStateBackToSource,
			event:           PeerEventDownloadSucceeded,
			backToSource:    true,
			peerFailedCount: 3,
			expect: func(t *testing.T, peer *Peer, mockParent *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.False(peer.Task.BackToSourcePeers.Contains(peer.ID))
				assert.Equal(int32(0), peer.Task.PeerFailedCount.Load())
				assert.Empty(peer.Parents())
				assert.Equal(int32(0), mockParent.Host.ConcurrentUploadCount.Load())
			},
		},
		{
			name:            "download succeeded from running keeps back-to-source peers and resets peer failed count",
			state:           PeerStateRunning,
			event:           PeerEventDownloadSucceeded,
			backToSource:    true,
			peerFailedCount: 3,
			expect: func(t *testing.T, peer *Peer, mockParent *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(peer.Task.BackToSourcePeers.Contains(peer.ID))
				assert.Equal(int32(0), peer.Task.PeerFailedCount.Load())
				assert.Empty(peer.Parents())
			},
		},
		{
			name:         "download failed from back-to-source counts the failure and clears back-to-source peer",
			state:        PeerStateBackToSource,
			event:        PeerEventDownloadFailed,
			backToSource: true,
			expect: func(t *testing.T, peer *Peer, mockParent *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.False(peer.Task.BackToSourcePeers.Contains(peer.ID))
				assert.Equal(int32(1), peer.Task.PeerFailedCount.Load())
				assert.Empty(peer.Parents())
			},
		},
		{
			name:  "download failed from running does not count as back-to-source failure",
			state: PeerStateRunning,
			event: PeerEventDownloadFailed,
			expect: func(t *testing.T, peer *Peer, mockParent *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(int32(0), peer.Task.PeerFailedCount.Load())
				assert.Empty(peer.Parents())
				assert.Equal(int32(0), mockParent.Host.ConcurrentUploadCount.Load())
			},
		},
		{
			name:         "leave clears back-to-source peer and releases parent upload load",
			state:        PeerStateRunning,
			event:        PeerEventLeave,
			backToSource: true,
			expect: func(t *testing.T, peer *Peer, mockParent *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.False(peer.Task.BackToSourcePeers.Contains(peer.ID))
				assert.Empty(peer.Parents())
				assert.Equal(int32(0), mockParent.Host.ConcurrentUploadCount.Load())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			mockParent := NewPeer(mockSeedPeerID, mockTask, mockHost)
			mockTask.StorePeer(peer)
			mockTask.StorePeer(mockParent)
			if err := mockTask.AddPeerEdge(mockParent, peer); err != nil {
				t.Fatal(err)
			}

			peer.FSM.SetState(tc.state)
			if tc.backToSource {
				mockTask.BackToSourcePeers.Add(peer.ID)
			}

			mockTask.PeerFailedCount.Store(tc.peerFailedCount)

			tc.expect(t, peer, mockParent, peer.FSM.Event(context.Background(), tc.event))
		})
	}
}

func TestPeer_PieceCosts(t *testing.T) {
	tests := []struct {
		name   string
		costs  []time.Duration
		expect func(t *testing.T, peer *Peer)
	}{
		{
			name:  "piece costs slice is empty",
			costs: []time.Duration{},
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				assert.Empty(peer.PieceCosts())
			},
		},
		{
			name:  "append piece cost",
			costs: []time.Duration{1},
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				assert.Equal([]time.Duration{1}, peer.PieceCosts())
			},
		},
		{
			name:  "piece costs are ordered from oldest to newest",
			costs: []time.Duration{3, 1, 2},
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				assert.Equal([]time.Duration{3, 1, 2}, peer.PieceCosts())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			for _, cost := range tc.costs {
				peer.AppendPieceCost(cost)
			}

			tc.expect(t, peer)
		})
	}
}

func TestPeer_PieceCostsStats(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer)
	}{
		{
			name: "stats cover the window excluding the latest cost",
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				peer.AppendPieceCost(time.Duration(100))
				peer.AppendPieceCost(time.Duration(200))
				peer.AppendPieceCost(time.Duration(600))
				stats := peer.PieceCostsStats()
				assert.Equal(3, stats.Count)
				assert.Equal(float64(600), stats.Last)
				assert.Equal(float64(150), stats.MeanExcludingLast())
				assert.Equal(float64(50), stats.StdDevExcludingLast())
			},
		},
		{
			name: "window is bounded by pieceCostsWindowLen",
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				for i := range pieceCostsWindowLen + 100 {
					peer.AppendPieceCost(time.Duration(i))
				}

				costs := peer.PieceCosts()
				assert.Len(costs, pieceCostsWindowLen)
				assert.Equal(time.Duration(100), costs[0])
				assert.Equal(time.Duration(pieceCostsWindowLen+99), costs[len(costs)-1])
				assert.Equal(pieceCostsWindowLen, peer.PieceCostsStats().Count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)

			tc.expect(t, peer)
		})
	}
}

func TestPeer_LoadReportPieceResultStream(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, stream schedulerv1.Scheduler_ReportPieceResultServer)
	}{
		{
			name: "load stream",
			expect: func(t *testing.T, peer *Peer, stream schedulerv1.Scheduler_ReportPieceResultServer) {
				assert := assert.New(t)
				peer.StoreReportPieceResultStream(stream)
				newStream, loaded := peer.LoadReportPieceResultStream()
				assert.True(loaded)
				assert.EqualValues(stream, newStream)
			},
		},
		{
			name: "stream does not exist",
			expect: func(t *testing.T, peer *Peer, stream schedulerv1.Scheduler_ReportPieceResultServer) {
				assert := assert.New(t)
				_, loaded := peer.LoadReportPieceResultStream()
				assert.False(loaded)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			stream := v1mocks.NewMockScheduler_ReportPieceResultServer(ctl)

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			tc.expect(t, peer, stream)
		})
	}
}

func TestPeer_StoreReportPieceResultStream(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, stream schedulerv1.Scheduler_ReportPieceResultServer)
	}{
		{
			name: "store stream",
			expect: func(t *testing.T, peer *Peer, stream schedulerv1.Scheduler_ReportPieceResultServer) {
				assert := assert.New(t)
				peer.StoreReportPieceResultStream(stream)
				newStream, loaded := peer.LoadReportPieceResultStream()
				assert.True(loaded)
				assert.EqualValues(stream, newStream)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			stream := v1mocks.NewMockScheduler_ReportPieceResultServer(ctl)

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			tc.expect(t, peer, stream)
		})
	}
}

func TestPeer_DeleteReportPieceResultStream(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, stream schedulerv1.Scheduler_ReportPieceResultServer)
	}{
		{
			name: "delete stream",
			expect: func(t *testing.T, peer *Peer, stream schedulerv1.Scheduler_ReportPieceResultServer) {
				assert := assert.New(t)
				peer.StoreReportPieceResultStream(stream)
				peer.DeleteReportPieceResultStream()
				_, loaded := peer.LoadReportPieceResultStream()
				assert.False(loaded)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			stream := v1mocks.NewMockScheduler_ReportPieceResultServer(ctl)

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			tc.expect(t, peer, stream)
		})
	}
}

func TestPeer_LoadAnnouncePeerStream(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, stream schedulerv2.Scheduler_AnnouncePeerServer)
	}{
		{
			name: "load stream",
			expect: func(t *testing.T, peer *Peer, stream schedulerv2.Scheduler_AnnouncePeerServer) {
				assert := assert.New(t)
				peer.StoreAnnouncePeerStream(stream)
				newStream, loaded := peer.LoadAnnouncePeerStream()
				assert.True(loaded)
				assert.EqualValues(stream, newStream)
			},
		},
		{
			name: "stream does not exist",
			expect: func(t *testing.T, peer *Peer, stream schedulerv2.Scheduler_AnnouncePeerServer) {
				assert := assert.New(t)
				_, loaded := peer.LoadAnnouncePeerStream()
				assert.False(loaded)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			stream := schedulerv2mock.NewMockScheduler_AnnouncePeerServer(ctl)

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			tc.expect(t, peer, stream)
		})
	}
}

func TestPeer_StoreAnnouncePeerStream(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, stream schedulerv2.Scheduler_AnnouncePeerServer)
	}{
		{
			name: "store stream",
			expect: func(t *testing.T, peer *Peer, stream schedulerv2.Scheduler_AnnouncePeerServer) {
				assert := assert.New(t)
				peer.StoreAnnouncePeerStream(stream)
				newStream, loaded := peer.LoadAnnouncePeerStream()
				assert.True(loaded)
				assert.EqualValues(stream, newStream)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			stream := schedulerv2mock.NewMockScheduler_AnnouncePeerServer(ctl)

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			tc.expect(t, peer, stream)
		})
	}
}

func TestPeer_DeleteAnnouncePeerStream(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, stream schedulerv2.Scheduler_AnnouncePeerServer)
	}{
		{
			name: "delete stream",
			expect: func(t *testing.T, peer *Peer, stream schedulerv2.Scheduler_AnnouncePeerServer) {
				assert := assert.New(t)
				peer.StoreAnnouncePeerStream(stream)
				peer.DeleteAnnouncePeerStream()
				_, loaded := peer.LoadAnnouncePeerStream()
				assert.False(loaded)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			stream := schedulerv2mock.NewMockScheduler_AnnouncePeerServer(ctl)

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			tc.expect(t, peer, stream)
		})
	}
}

func TestPeer_Parents(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, seedPeer *Peer)
	}{
		{
			name: "peer has no parents",
			expect: func(t *testing.T, peer *Peer, seedPeer *Peer) {
				assert := assert.New(t)
				peer.Task.StorePeer(peer)
				assert.Empty(peer.Parents())
			},
		},
		{
			name: "peer has parents",
			expect: func(t *testing.T, peer *Peer, seedPeer *Peer) {
				assert := assert.New(t)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(seedPeer)
				if err := peer.Task.AddPeerEdge(seedPeer, peer); err != nil {
					t.Fatal(err)
				}

				assert.Len(peer.Parents(), 1)
				assert.Equal(mockSeedPeerID, peer.Parents()[0].ID)
			},
		},
		{
			name: "peer is not stored in task",
			expect: func(t *testing.T, peer *Peer, seedPeer *Peer) {
				assert := assert.New(t)
				assert.Nil(peer.Parents())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			seedPeer := NewPeer(mockSeedPeerID, mockTask, mockHost)
			tc.expect(t, peer, seedPeer)
		})
	}
}

func TestPeer_Children(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, seedPeer *Peer)
	}{
		{
			name: "peer has no children",
			expect: func(t *testing.T, peer *Peer, seedPeer *Peer) {
				assert := assert.New(t)
				peer.Task.StorePeer(peer)
				assert.Empty(peer.Children())
			},
		},
		{
			name: "peer has children",
			expect: func(t *testing.T, peer *Peer, seedPeer *Peer) {
				assert := assert.New(t)
				peer.Task.StorePeer(peer)
				peer.Task.StorePeer(seedPeer)
				if err := peer.Task.AddPeerEdge(peer, seedPeer); err != nil {
					t.Fatal(err)
				}

				assert.Len(peer.Children(), 1)
				assert.Equal(mockSeedPeerID, peer.Children()[0].ID)
			},
		},
		{
			name: "peer is not stored in task",
			expect: func(t *testing.T, peer *Peer, seedPeer *Peer) {
				assert := assert.New(t)
				assert.Nil(peer.Children())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)
			seedPeer := NewPeer(mockSeedPeerID, mockTask, mockHost)
			tc.expect(t, peer, seedPeer)
		})
	}
}

func TestPeer_DownloadTinyFile(t *testing.T) {
	testData := []byte("./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" +
		"./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	mockServer := func(t *testing.T, peer *Peer) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert := assert.New(t)
			assert.NotNil(peer)
			assert.Equal(fmt.Sprintf("/download/%s/%s", peer.Task.ID[:3], peer.Task.ID), r.URL.Path)
			assert.Equal(fmt.Sprintf("peerId=%s", peer.ID), r.URL.RawQuery)

			rgs, err := nethttp.ParseRange(r.Header.Get(headers.Range), 128)
			assert.NoError(err)
			assert.Len(rgs, 1)
			rg := rgs[0]

			w.WriteHeader(http.StatusPartialContent)
			n, err := w.Write(testData[rg.Start : rg.Start+rg.Length])
			assert.NoError(err)
			assert.Equal(rg.Length, int64(n))
		}))
	}

	tests := []struct {
		name             string
		mockServer       func(t *testing.T, peer *Peer) *httptest.Server
		useDefaultClient bool
		expect           func(t *testing.T, peer *Peer)
	}{
		{
			name:       "download tiny file",
			mockServer: mockServer,
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				peer.Task.ContentLength.Store(32)
				data, err := peer.DownloadTinyFile()
				assert.NoError(err)
				assert.Equal(testData[:32], data)
			},
		},
		{
			name:       "download tiny file with range",
			mockServer: mockServer,
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				peer.Task.ContentLength.Store(32)
				peer.Range = &nethttp.Range{
					Start:  0,
					Length: 10,
				}
				data, err := peer.DownloadTinyFile()
				assert.NoError(err)
				assert.Equal(testData[:32], data)
			},
		},
		{
			name: "download tiny file failed because of http status code",
			mockServer: func(t *testing.T, peer *Peer) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				peer.Task.ID = "foobar"
				_, err := peer.DownloadTinyFile()
				assert.Error(err)
			},
		},
		{
			name:       "download tiny file failed because of invalid task id",
			mockServer: mockServer,
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				peer.Task.ID = "foo"
				_, err := peer.DownloadTinyFile()
				assert.Error(err)
			},
		},
		{
			name:             "download tiny file failed because the default client rejects loopback address",
			mockServer:       mockServer,
			useDefaultClient: true,
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				peer.Task.ContentLength.Store(32)
				_, err := peer.DownloadTinyFile()
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost)

			s := tc.mockServer(t, peer)
			defer s.Close()

			if !tc.useDefaultClient {
				WithTinyFileHTTPClient(s.Client())(peer)
			}

			u, err := url.Parse(s.URL)
			if err != nil {
				t.Fatal(err)
			}

			ip, rawPort, err := net.SplitHostPort(u.Host)
			if err != nil {
				t.Fatal(err)
			}

			port, err := strconv.ParseInt(rawPort, 10, 32)
			if err != nil {
				t.Fatal(err)
			}

			mockHost.IP = ip
			mockHost.DownloadPort = int32(port)
			tc.expect(t, peer)
		})
	}
}

func TestPeer_CalculatePriority(t *testing.T) {
	tests := []struct {
		name        string
		priority    commonv2.Priority
		application string
		url         string
		mock        func(md *configmocks.MockDynconfigInterfaceMockRecorder)
		expect      func(t *testing.T, priority commonv2.Priority)
	}{
		{
			name:        "peer has priority",
			priority:    commonv2.Priority_LEVEL4,
			application: mockTaskApplication,
			url:         mockTaskURL,
			mock:        func(md *configmocks.MockDynconfigInterfaceMockRecorder) {},
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL4, priority)
			},
		},
		{
			name:        "get applications failed",
			priority:    commonv2.Priority_LEVEL0,
			application: mockTaskApplication,
			url:         mockTaskURL,
			mock: func(md *configmocks.MockDynconfigInterfaceMockRecorder) {
				md.GetApplications().Return(nil, errors.New("bas")).Times(1)
			},
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL0, priority)
			},
		},
		{
			name:        "can not found applications",
			priority:    commonv2.Priority_LEVEL0,
			application: mockTaskApplication,
			url:         mockTaskURL,
			mock: func(md *configmocks.MockDynconfigInterfaceMockRecorder) {
				md.GetApplications().Return([]*managerv2.Application{}, nil).Times(1)
			},
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL0, priority)
			},
		},
		{
			name:        "can not found matching application",
			priority:    commonv2.Priority_LEVEL0,
			application: mockTaskApplication,
			url:         mockTaskURL,
			mock: func(md *configmocks.MockDynconfigInterfaceMockRecorder) {
				md.GetApplications().Return([]*managerv2.Application{
					{
						Name: "baw",
					},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL0, priority)
			},
		},
		{
			name:        "can not found priority",
			priority:    commonv2.Priority_LEVEL0,
			application: "bae",
			url:         mockTaskURL,
			mock: func(md *configmocks.MockDynconfigInterfaceMockRecorder) {
				md.GetApplications().Return([]*managerv2.Application{
					{
						Name: "bae",
					},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL0, priority)
			},
		},
		{
			name:        "match the priority of application",
			priority:    commonv2.Priority_LEVEL0,
			application: "baz",
			url:         mockTaskURL,
			mock: func(md *configmocks.MockDynconfigInterfaceMockRecorder) {
				md.GetApplications().Return([]*managerv2.Application{
					{
						Name: "baz",
						Priority: &managerv2.ApplicationPriority{
							Value: commonv2.Priority_LEVEL1,
						},
					},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL1, priority)
			},
		},
		{
			name:        "match the priority of url",
			priority:    commonv2.Priority_LEVEL0,
			application: "bak",
			url:         "example.com",
			mock: func(md *configmocks.MockDynconfigInterfaceMockRecorder) {
				md.GetApplications().Return([]*managerv2.Application{
					{
						Name: "bak",
						Priority: &managerv2.ApplicationPriority{
							Value: commonv2.Priority_LEVEL1,
							Urls: []*managerv2.URLPriority{
								{
									Regex: "am",
									Value: commonv2.Priority_LEVEL2,
								},
							},
						},
					},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL2, priority)
			},
		},
		{
			name:        "url does not match and falls back to the priority of application",
			priority:    commonv2.Priority_LEVEL0,
			application: "bak",
			url:         "example.com",
			mock: func(md *configmocks.MockDynconfigInterfaceMockRecorder) {
				md.GetApplications().Return([]*managerv2.Application{
					{
						Name: "bak",
						Priority: &managerv2.ApplicationPriority{
							Value: commonv2.Priority_LEVEL1,
							Urls: []*managerv2.URLPriority{
								{
									Regex: "zzz",
									Value: commonv2.Priority_LEVEL2,
								},
							},
						},
					},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL1, priority)
			},
		},
		{
			name:        "url regex is invalid and falls back to the priority of application",
			priority:    commonv2.Priority_LEVEL0,
			application: "bak",
			url:         "example.com",
			mock: func(md *configmocks.MockDynconfigInterfaceMockRecorder) {
				md.GetApplications().Return([]*managerv2.Application{
					{
						Name: "bak",
						Priority: &managerv2.ApplicationPriority{
							Value: commonv2.Priority_LEVEL1,
							Urls: []*managerv2.URLPriority{
								{
									Regex: "(",
									Value: commonv2.Priority_LEVEL2,
								},
							},
						},
					},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL1, priority)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			dynconfig := configmocks.NewMockDynconfigInterface(ctl)

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, tc.url, mockTaskTag, tc.application, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer := NewPeer(mockPeerID, mockTask, mockHost, WithPriority(tc.priority))
			tc.mock(dynconfig.EXPECT())
			tc.expect(t, peer.CalculatePriority(dynconfig))
		})
	}
}
