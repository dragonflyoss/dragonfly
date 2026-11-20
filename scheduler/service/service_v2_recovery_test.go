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
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"

	"d7y.io/dragonfly/v2/pkg/container/set"
	"d7y.io/dragonfly/v2/pkg/net/http"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
	schedulingmocks "d7y.io/dragonfly/v2/scheduler/scheduling/mocks"
)

func TestServiceV2_RegisterFailedTaskWithSucceededParent(t *testing.T) {
	ctrl := gomock.NewController(t)
	resource := standard.NewMockResource(ctrl)
	hostManager := standard.NewMockHostManager(ctrl)
	taskManager := standard.NewMockTaskManager(ctrl)
	peerManager := standard.NewMockPeerManager(ctrl)
	scheduling := schedulingmocks.NewMockScheduling(ctrl)
	task := standard.NewTask(mockTaskID, mockTaskURL, "", "", commonv2.TaskType_STANDARD, nil, nil, 3)
	task.FSM.SetState(standard.TaskStateFailed)
	task.ContentLength.Store(1024)
	task.TotalPieceCount.Store(2)
	host := standard.NewHost("new-host", "127.0.0.1", "", "", 8003, 8001, 8004, types.HostTypeNormal)
	parentHost := standard.NewHost("cache-host", "127.0.0.2", "", "", 8003, 8001, 8004, types.HostTypeNormal)
	parent := standard.NewPeer("cached-peer", task, parentHost)
	parent.FSM.SetState(standard.PeerStateSucceeded)
	task.StorePeer(parent)
	peer := standard.NewPeer("new-peer", task, host, standard.WithPriority(commonv2.Priority_LEVEL6))
	task.StorePeer(peer)

	resource.EXPECT().HostManager().Return(hostManager)
	hostManager.EXPECT().Load(host.ID).Return(host, true)
	resource.EXPECT().TaskManager().Return(taskManager)
	taskManager.EXPECT().Load(task.ID).Return(task, true)
	resource.EXPECT().PeerManager().Return(peerManager)
	peerManager.EXPECT().Load(peer.ID).Return(peer, true)
	// 不触发 SeedPeer 回源，进入原有父节点调度流程；后续仍使用原有回源重试策略。
	scheduling.EXPECT().ScheduleCandidateParents(gomock.Any(), peer, gomock.Any()).DoAndReturn(
		func(context.Context, *standard.Peer, set.SafeSet[string]) error {
			require.False(t, peer.NeedBackToSource.Load())
			return nil
		})

	svc := &V2{resource: resource, scheduling: scheduling, config: &config.Config{Scheduler: mockSchedulerConfig}}
	require.NoError(t, svc.handleRegisterPeerRequest(context.Background(), nil, host.ID, task.ID, peer.ID,
		&schedulerv2.RegisterPeerRequest{Download: &commonv2.Download{}}))
	require.Equal(t, standard.PeerStateReceivedNormal, peer.FSM.Current())
	require.Equal(t, standard.TaskStateRunning, task.FSM.Current())
}

func TestServiceV2_ConcurrentTaskRecovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	resource := standard.NewMockResource(ctrl)
	peerManager := standard.NewMockPeerManager(ctrl)
	resource.EXPECT().PeerManager().Return(peerManager).AnyTimes()
	task := standard.NewTask(mockTaskID, mockTaskURL, "", "", commonv2.TaskType_STANDARD, nil, nil, 3)
	task.FSM.SetState(standard.TaskStateRunning)
	task.ContentLength.Store(1024)
	task.TotalPieceCount.Store(2)
	host := standard.NewHost(mockHostID, "127.0.0.1", "", "", 8003, 8001, 8004, types.HostTypeNormal)
	svc := &V2{resource: resource}

	var requests []func() error
	for _, id := range []string{"source-1", "source-2", "p2p-1", "p2p-2", "p2p-3"} {
		peer := standard.NewPeer(id, task, host, standard.WithPriority(commonv2.Priority_LEVEL6))
		task.StorePeer(peer)
		peerManager.EXPECT().Load(id).Return(peer, true)
		if id == "source-1" || id == "source-2" {
			peer.FSM.SetState(standard.PeerStateBackToSource)
			task.BackToSourcePeers.Add(id)
			requests = append(requests, func() error {
				return svc.handleDownloadPeerBackToSourceFailedRequest(context.Background(), peer.ID)
			})
		} else {
			peer.FSM.SetState(standard.PeerStateRunning)
			requests = append(requests, func() error {
				return svc.handleDownloadPeerFinishedRequest(context.Background(), peer.ID)
			})
		}
	}

	start := make(chan struct{})
	results := make(chan error, len(requests))
	var wg sync.WaitGroup
	for _, request := range requests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			results <- request()
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	for err := range results {
		require.NoError(t, err)
	}
	require.Equal(t, standard.TaskStateSucceeded, task.FSM.Current())
	require.Equal(t, int64(1024), task.ContentLength.Load())
	require.Equal(t, int32(2), task.TotalPieceCount.Load())
	require.Zero(t, task.BackToSourcePeers.Len())
}

func TestServiceV2_RecoverTaskAfterSourceFailure(t *testing.T) {
	for _, ranged := range []bool{false, true} {
		name := "complete P2P"
		if ranged {
			name = "range P2P"
		}
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			resource := standard.NewMockResource(ctrl)
			peerManager := standard.NewMockPeerManager(ctrl)
			resource.EXPECT().PeerManager().Return(peerManager).AnyTimes()
			task := standard.NewTask(mockTaskID, mockTaskURL, "", "", commonv2.TaskType_STANDARD, nil, nil, 3)
			task.FSM.SetState(standard.TaskStateRunning)
			task.ContentLength.Store(1)
			task.TotalPieceCount.Store(1)
			task.DirectPiece = []byte{1}
			task.StorePiece(&mockPiece)
			host := standard.NewHost(mockHostID, "127.0.0.1", "", "", 8003, 8001, 8004, types.HostTypeNormal)
			source := standard.NewPeer("source", task, host, standard.WithPriority(commonv2.Priority_LEVEL6))
			source.FSM.SetState(standard.PeerStateBackToSource)
			task.BackToSourcePeers.Add(source.ID)
			task.StorePeer(source)
			peer := standard.NewPeer("p2p", task, host, standard.WithPriority(commonv2.Priority_LEVEL6))
			peer.FSM.SetState(standard.PeerStateRunning)
			if ranged {
				peer.Range = &http.Range{Start: 0, Length: 1}
			}
			task.StorePeer(peer)
			peerManager.EXPECT().Load(source.ID).Return(source, true)
			peerManager.EXPECT().Load(peer.ID).Return(peer, true)
			svc := &V2{resource: resource}

			require.NoError(t, svc.handleDownloadPeerBackToSourceFailedRequest(context.Background(), source.ID))
			require.Equal(t, standard.TaskStateFailed, task.FSM.Current())
			require.Equal(t, int64(1), task.ContentLength.Load())
			require.Equal(t, int32(1), task.TotalPieceCount.Load())
			require.Equal(t, []byte{1}, task.DirectPiece)
			piece, loaded := task.LoadPiece(mockPiece.Number)
			require.True(t, loaded)
			require.Same(t, &mockPiece, piece)
			require.False(t, task.BackToSourcePeers.Contains(source.ID))

			require.NoError(t, svc.handleDownloadPeerFinishedRequest(context.Background(), peer.ID))
			require.Equal(t, standard.PeerStateSucceeded, peer.FSM.Current())
			wantState := standard.TaskStateSucceeded
			if ranged {
				wantState = standard.TaskStateFailed
			}
			require.Equal(t, wantState, task.FSM.Current())

			// 另一个回源 Peer 的失败晚于完整 P2P 成功，不能再使任务退化或清空缓存。
			late := standard.NewPeer("late-source", task, host, standard.WithPriority(commonv2.Priority_LEVEL6))
			late.FSM.SetState(standard.PeerStateBackToSource)
			task.StorePeer(late)
			peerManager.EXPECT().Load(late.ID).Return(late, true)
			require.NoError(t, svc.handleDownloadPeerBackToSourceFailedRequest(context.Background(), late.ID))
			require.Equal(t, wantState, task.FSM.Current())
			require.Equal(t, int64(1), task.ContentLength.Load())
			require.Equal(t, int32(1), task.TotalPieceCount.Load())
			require.Equal(t, []byte{1}, task.DirectPiece)
		})
	}
}
