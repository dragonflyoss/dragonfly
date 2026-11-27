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

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"

	"d7y.io/dragonfly/v2/pkg/container/set"
	"d7y.io/dragonfly/v2/pkg/gc"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
	schedulingmocks "d7y.io/dragonfly/v2/scheduler/scheduling/mocks"
)

type seedCacheServiceFixture struct {
	service    *V2
	resource   standard.Resource
	downloader *standard.Peer
	cached     *standard.Peer
	announced  *standard.Peer
	foreign    *standard.Peer
}

func newSeedCacheServiceFixture(t *testing.T) seedCacheServiceFixture {
	t.Helper()
	resource, err := standard.New(config.New(), gc.New(), insecure.NewCredentials())
	require.NoError(t, err)

	task := standard.NewTask("seed-cache-task", "https://example.com/cache", "", "",
		commonv2.TaskType_STANDARD, nil, nil, 10, standard.WithPieceLength(4*1024*1024))
	task.ContentLength.Store(32 * 1024 * 1024)
	task.TotalPieceCount.Store(8)
	task.FSM.SetState(standard.TaskStateRunning)
	resource.TaskManager().LoadOrStore(task)

	newPeer := func(id, hostID string, typ types.HostType, task *standard.Task, cached bool) *standard.Peer {
		host := standard.NewHost(hostID, "127.0.0.1", hostID, hostID, 4000, 4003, 4001, typ)
		resource.HostManager().LoadOrStore(host)
		peer := standard.NewPeer(id, task, host)
		peer.SeedCache = cached
		peer.FSM.SetState(standard.PeerStateSucceeded)
		peer, loaded := resource.PeerManager().LoadOrStore(peer)
		require.False(t, loaded)
		return peer
	}

	downloader := newPeer("downloader", "client-host", types.HostTypeNormal, task, false)
	downloader.FSM.SetState(standard.PeerStateRunning)
	cached := newPeer("cached-seed-parent", "cached-seed-host", types.HostTypeSuperSeed, task, true)
	announced := newPeer("announced-seed-parent", "announced-seed-host", types.HostTypeSuperSeed, task, false)
	require.NoError(t, task.AddPeerEdge(cached, downloader))

	foreignTask := standard.NewTask("other-task", "https://example.com/other", "", "",
		commonv2.TaskType_STANDARD, nil, nil, 10)
	resource.TaskManager().LoadOrStore(foreignTask)
	foreign := newPeer("other-task-cached-parent", "other-seed-host", types.HostTypeSuperSeed, foreignTask, true)

	return seedCacheServiceFixture{
		service: &V2{resource: resource}, resource: resource, downloader: downloader,
		cached: cached, announced: announced, foreign: foreign,
	}
}

func requireSeedCachePeerRemoved(t *testing.T, fixture seedCacheServiceFixture) {
	t.Helper()
	_, loaded := fixture.resource.PeerManager().Load(fixture.cached.ID)
	require.False(t, loaded, "the failed cache advertisement must leave the global peer index")
	_, loaded = fixture.downloader.Task.LoadPeer(fixture.cached.ID)
	require.False(t, loaded, "the failed cache advertisement must leave the task DAG")
	_, loaded = fixture.cached.Host.LoadPeer(fixture.cached.ID)
	require.False(t, loaded, "the failed cache advertisement must leave the host peer index")
	require.Equal(t, standard.PeerStateLeave, fixture.cached.FSM.Current())
	require.Zero(t, fixture.cached.Host.ConcurrentUploadCount.Load(), "deleting the parent must release its upload edge")
	require.Equal(t, 2, fixture.downloader.Task.PeerCount())
	_, loaded = fixture.resource.PeerManager().Load(fixture.announced.ID)
	require.True(t, loaded, "a normally announced seed parent must remain available")
	require.Equal(t, standard.PeerStateSucceeded, fixture.announced.FSM.Current())
}

func TestServiceV2_RescheduleInvalidatesSeedCacheParents(t *testing.T) {
	fixture := newSeedCacheServiceFixture(t)
	scheduling := schedulingmocks.NewMockScheduling(gomock.NewController(t))
	fixture.service.scheduling = scheduling

	// A missing task in SyncPieces produces a reschedule with the original
	// candidate list, rather than a DownloadPieceFailed request. Invalidate the
	// cache advertisement before scheduling another parent.
	scheduling.EXPECT().ScheduleCandidateParents(gomock.Any(), fixture.downloader, gomock.Any()).
		DoAndReturn(func(_ context.Context, peer *standard.Peer, blocked set.SafeSet[string]) error {
			requireSeedCachePeerRemoved(t, fixture)
			require.Same(t, fixture.downloader, peer)
			require.True(t, blocked.Contains(fixture.cached.ID))
			return nil
		}).Times(2)

	parents := []*commonv2.Peer{{Id: fixture.cached.ID}, {Id: fixture.announced.ID}, {Id: "missing-parent"}}
	require.NoError(t, fixture.service.handleReschedulePeerRequest(context.Background(), fixture.downloader.ID, parents))
	requireSeedCachePeerRemoved(t, fixture)
	for _, parent := range parents {
		require.True(t, fixture.downloader.BlockParents.Contains(parent.Id))
	}

	// An already invalidated advertisement can appear in a retried report.
	require.NoError(t, fixture.service.handleReschedulePeerRequest(context.Background(), fixture.downloader.ID, parents))
	requireSeedCachePeerRemoved(t, fixture)
}

func TestServiceV2_TemporaryPieceFailureInvalidatesOnlySeedCacheParent(t *testing.T) {
	for _, test := range []struct {
		name      string
		parent    string
		temporary bool
		removed   bool
		canceled  bool
	}{
		{name: "cached seed parent", parent: "cached-seed-parent", temporary: true, removed: true},
		{name: "cached seed parent after request cancellation", parent: "cached-seed-parent", temporary: true, removed: true, canceled: true},
		{name: "normally announced seed parent", parent: "announced-seed-parent", temporary: true},
		{name: "cached parent of a different task", parent: "other-task-cached-parent", temporary: true},
		{name: "missing parent", parent: "missing-parent", temporary: true},
		{name: "permanent failure", parent: "cached-seed-parent"},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newSeedCacheServiceFixture(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if test.canceled {
				cancel()
			}
			err := fixture.service.handleDownloadPieceFailedRequest(ctx, fixture.downloader.ID,
				&schedulerv2.DownloadPieceFailedRequest{ParentId: test.parent, Temporary: test.temporary})
			if test.temporary {
				require.NoError(t, err)
				require.True(t, fixture.downloader.BlockParents.Contains(test.parent))
			} else {
				require.Equal(t, codes.FailedPrecondition, status.Code(err))
				require.False(t, fixture.downloader.BlockParents.Contains(test.parent))
			}

			if test.removed {
				requireSeedCachePeerRemoved(t, fixture)
			} else {
				_, loaded := fixture.resource.PeerManager().Load(fixture.cached.ID)
				require.True(t, loaded)
				require.Equal(t, standard.PeerStateSucceeded, fixture.cached.FSM.Current())
			}
			_, loaded := fixture.resource.PeerManager().Load(fixture.announced.ID)
			require.True(t, loaded)
			require.Equal(t, standard.PeerStateSucceeded, fixture.announced.FSM.Current())
			_, loaded = fixture.resource.PeerManager().Load(fixture.foreign.ID)
			require.True(t, loaded, "one task's failure must not invalidate another task's cache")
			require.Equal(t, standard.PeerStateSucceeded, fixture.foreign.FSM.Current())
		})
	}
}

func TestServiceV2_InvalidateSeedCachePeerRejectsOtherTask(t *testing.T) {
	fixture := newSeedCacheServiceFixture(t)
	fixture.service.invalidateSeedCachePeer(context.Background(), fixture.downloader.Task, fixture.foreign)
	peer, loaded := fixture.resource.PeerManager().Load(fixture.foreign.ID)
	require.True(t, loaded)
	require.Same(t, fixture.foreign, peer)
	require.Equal(t, standard.PeerStateSucceeded, peer.FSM.Current())
	require.Equal(t, 1, peer.Task.PeerCount())
}
