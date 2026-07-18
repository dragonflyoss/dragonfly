/*
 *     Copyright 2022 The Dragonfly Authors
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
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	dfdaemonv2 "d7y.io/api/v2/pkg/apis/dfdaemon/v2"
	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"
	dfdaemonclientmocks "d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/client/mocks"
	"d7y.io/dragonfly/v2/pkg/types"
)

func TestSeedPeer_newSeedPeer(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, s SeedPeer)
	}{
		{
			name: "new seed peer",
			expect: func(t *testing.T, s SeedPeer) {
				assert := assert.New(t)
				assert.Equal(reflect.TypeOf(s).Elem().Name(), "seedPeer")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			hostManager := NewMockHostManager(ctl)
			taskManager := NewMockTaskManager(ctl)
			peerManager := NewMockPeerManager(ctl)
			clientPool := dfdaemonclientmocks.NewMockPool(ctl)

			tc.expect(t, newSeedPeer(peerManager, hostManager, taskManager, clientPool, mockTaskBackToSourceLimit))
		})
	}
}

func TestSeedPeer_TriggerDownloadTask(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, err error)
	}{
		{
			name: "trigger download task failed",
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "no seed peer available")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			hostManager := NewMockHostManager(ctl)
			taskManager := NewMockTaskManager(ctl)
			peerManager := NewMockPeerManager(ctl)
			clientPool := dfdaemonclientmocks.NewMockPool(ctl)

			seedPeer := newSeedPeer(peerManager, hostManager, taskManager, clientPool, mockTaskBackToSourceLimit)
			tc.expect(t, seedPeer.TriggerDownloadTask(context.Background(), mockTaskID, &dfdaemonv2.DownloadTaskRequest{}))
		})
	}
}

func TestSeedPeer_handleDownloadTaskResponse_registersSeedPeer(t *testing.T) {
	assert := assert.New(t)

	hostManager := &hostManager{Map: &sync.Map{}, normals: &sync.Map{}, seeds: &sync.Map{}}
	taskManager := &taskManager{Map: &sync.Map{}}
	peerManager := &peerManager{Map: &sync.Map{}, mu: &sync.Mutex{}}
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	clientPool := dfdaemonclientmocks.NewMockPool(ctl)
	seedPeer := newSeedPeer(peerManager, hostManager, taskManager, clientPool, mockTaskBackToSourceLimit).(*seedPeer)

	seedHost := NewHost("seed-host", "127.0.0.1", "seed", "seed", 4000, 8001, 65001, types.HostTypeSuperSeed)
	hostManager.Store(seedHost)
	task := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithPieceLength(mockTaskPieceLength))
	taskManager.Store(task)

	contentLength := uint64(4096)
	pieceCount := uint64(2)
	download := &commonv2.Download{
		Url:                 mockTaskURL,
		Type:                commonv2.TaskType_STANDARD,
		Priority:            commonv2.Priority_LEVEL6,
		FilteredQueryParams: mockTaskFilteredQueryParams,
		RequestHeader:       mockTaskHeader,
		NeedBackToSource:    true,
		ActualPieceLength:   &mockTaskPieceLength,
		ActualContentLength: &contentLength,
		ActualPieceCount:    &pieceCount,
	}

	started := &dfdaemonv2.DownloadTaskResponse{
		HostId: seedHost.ID,
		TaskId: mockTaskID,
		PeerId: "seed-peer",
		Response: &dfdaemonv2.DownloadTaskResponse_DownloadTaskStartedResponse{
			DownloadTaskStartedResponse: &dfdaemonv2.DownloadTaskStartedResponse{ContentLength: contentLength},
		},
	}
	peer, err := seedPeer.handleDownloadTaskResponse(context.Background(), mockTaskID, download, started, nil)
	assert.NoError(err)
	assert.NotNil(peer)
	assert.Equal("seed-peer", peer.ID)
	assert.True(peer.FSM.Is(PeerStateBackToSource))
	assert.Equal(int64(contentLength), task.ContentLength.Load())
	assert.Equal(int32(pieceCount), task.TotalPieceCount.Load())

	trafficType := commonv2.TrafficType_LOCAL_PEER
	pieceResp := &dfdaemonv2.DownloadTaskResponse{
		HostId: seedHost.ID,
		TaskId: mockTaskID,
		PeerId: "seed-peer",
		Response: &dfdaemonv2.DownloadTaskResponse_DownloadPieceFinishedResponse{
			DownloadPieceFinishedResponse: &dfdaemonv2.DownloadPieceFinishedResponse{Piece: &commonv2.Piece{
				Number:      0,
				Offset:      0,
				Length:      mockTaskPieceLength,
				TrafficType: &trafficType,
			}},
		},
	}
	peer, err = seedPeer.handleDownloadTaskResponse(context.Background(), mockTaskID, download, pieceResp, peer)
	assert.NoError(err)
	assert.True(peer.FinishedPieces.Test(0))
	piece, loaded := task.LoadPiece(0)
	assert.True(loaded)
	assert.Equal(commonv2.TrafficType_LOCAL_PEER, piece.TrafficType)

	assert.NoError(seedPeer.completeSeedPeerDownloadTask(context.Background(), peer))
	assert.True(peer.FSM.Is(PeerStateSucceeded))
	assert.True(task.FSM.Is(TaskStateSucceeded))

	storedPeer, loaded := peerManager.Load("seed-peer")
	assert.True(loaded)
	assert.Equal(peer, storedPeer)
	storedPeer, loaded = task.LoadPeer("seed-peer")
	assert.True(loaded)
	assert.Equal(peer, storedPeer)
}

func TestSeedPeer_TriggerTask(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, peer *Peer, result *schedulerv1.PeerResult, err error)
	}{
		{
			name: "start obtain seed stream failed",
			expect: func(t *testing.T, peer *Peer, result *schedulerv1.PeerResult, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "no seed peer available")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			hostManager := NewMockHostManager(ctl)
			taskManager := NewMockTaskManager(ctl)
			peerManager := NewMockPeerManager(ctl)
			clientPool := dfdaemonclientmocks.NewMockPool(ctl)

			seedPeer := newSeedPeer(peerManager, hostManager, taskManager, clientPool, mockTaskBackToSourceLimit)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, WithDigest(mockTaskDigest))
			peer, result, err := seedPeer.TriggerTask(context.Background(), nil, mockTask)
			tc.expect(t, peer, result, err)
		})
	}
}
