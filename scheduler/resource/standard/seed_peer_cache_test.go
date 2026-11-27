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

package standard

import (
	"context"
	"errors"
	"io"
	"net"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/utils/keymutex"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	dfdaemonv2 "d7y.io/api/v2/pkg/apis/dfdaemon/v2"
	dfdaemonv2mocks "d7y.io/api/v2/pkg/apis/dfdaemon/v2/mocks"

	"d7y.io/dragonfly/v2/pkg/container/set"
	dfdaemonclientmocks "d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/client/mocks"
	"d7y.io/dragonfly/v2/pkg/types"
)

const (
	seedCacheTestContentLength = uint64(32 << 20)
	seedCacheTestPieceLength   = uint64(4 << 20)
	seedCacheTestPieceCount    = uint64(8)
)

// Use the real managers so these tests verify that a recovered parent is
// published consistently to the peer manager, host, and task DAG.
type seedCacheFixture struct {
	seed    *seedPeer
	tasks   *taskManager
	hosts   *hostManager
	peers   *peerManager
	task    *Task
	host    *Host
	client  *dfdaemonclientmocks.MockClient
	stream  *dfdaemonv2mocks.MockDfdaemonUpload_DownloadTaskClient
	eof     *gomock.Call
	request *dfdaemonv2.DownloadTaskRequest
}

func newSeedCacheFixture(t *testing.T) *seedCacheFixture {
	t.Helper()
	ctl := gomock.NewController(t)
	tasks := &taskManager{Map: &sync.Map{}}
	hosts := &hostManager{Map: &sync.Map{}, normals: &sync.Map{}, seeds: &sync.Map{}}
	peers := &peerManager{Map: &sync.Map{}, mu: keymutex.NewHashed(peerMutexShardCount)}
	pool := dfdaemonclientmocks.NewMockPool(ctl)
	client := dfdaemonclientmocks.NewMockClient(ctl)
	stream := dfdaemonv2mocks.NewMockDfdaemonUpload_DownloadTaskClient(ctl)
	task := NewTask("cached-task", "http://example.com/cached-file", "", "", commonv2.TaskType_STANDARD,
		nil, nil, 1, WithPieceLength(seedCacheTestPieceLength))
	task.ContentLength.Store(int64(seedCacheTestContentLength))
	task.TotalPieceCount.Store(int32(seedCacheTestPieceCount))
	require.NoError(t, task.FSM.Event(context.Background(), TaskEventDownload))
	tasks.LoadOrStore(task)
	host := NewHost("seed-host", "192.0.2.2", "seed-client-0", "seed-client-0", 4000, 4001, 4002, types.HostTypeSuperSeed)
	hosts.LoadOrStore(host)
	seed := newSeedPeer(tasks, peers, hosts, pool).(*seedPeer)
	seed.Store(host)

	addr := net.JoinHostPort(host.IP, strconv.Itoa(int(host.Port)))
	pool.EXPECT().Get(addr).Return(client, nil).AnyTimes()
	client.EXPECT().DownloadTask(gomock.Any(), task.ID, gomock.Any()).Return(stream, nil).AnyTimes()
	eof := stream.EXPECT().Recv().Return(nil, io.EOF).AnyTimes()

	contentLength, pieceLength, pieceCount := seedCacheTestContentLength, seedCacheTestPieceLength, seedCacheTestPieceCount
	return &seedCacheFixture{
		seed: seed, tasks: tasks, hosts: hosts, peers: peers, task: task, host: host, client: client, stream: stream, eof: eof,
		request: &dfdaemonv2.DownloadTaskRequest{Download: &commonv2.Download{
			Url: task.URL, ActualContentLength: &contentLength, ActualPieceLength: &pieceLength, ActualPieceCount: &pieceCount,
		}},
	}
}

func (f *seedCacheFixture) localTask() *dfdaemonv2.StatLocalTaskResponse {
	contentLength, pieceLength := seedCacheTestContentLength, seedCacheTestPieceLength
	return &dfdaemonv2.StatLocalTaskResponse{
		TaskId: f.task.ID, ContentLength: &contentLength, PieceLength: &pieceLength,
		FinishedAt: timestamppb.New(time.Now().Add(-time.Minute)),
	}
}

func (f *seedCacheFixture) expectLocalTask(t *testing.T, response *dfdaemonv2.StatLocalTaskResponse, err error, duringStat func()) {
	t.Helper()
	f.client.EXPECT().StatLocalTask(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req *dfdaemonv2.StatLocalTaskRequest, _ ...grpc.CallOption) (*dfdaemonv2.StatLocalTaskResponse, error) {
			require.Equal(t, f.task.ID, req.GetTaskId())
			deadline, bounded := ctx.Deadline()
			require.True(t, bounded, "the cache verification RPC must have a deadline")
			require.LessOrEqual(t, time.Until(deadline), 5*time.Second)
			if duringStat != nil {
				duringStat()
			}
			return response, err
		}).Times(1)
}

func (f *seedCacheFixture) requireCachedParent(t *testing.T) *Peer {
	t.Helper()
	parents := f.task.LoadPeers()
	require.Len(t, parents, 1)
	parent := parents[0]
	require.Same(t, f.host, parent.Host)
	require.Same(t, f.task, parent.Task)
	require.True(t, parent.FSM.Is(PeerStateSucceeded))
	require.True(t, parent.SeedCache)
	require.Nil(t, parent.Range)
	require.False(t, parent.NeedBackToSource.Load())
	require.Equal(t, uint(seedCacheTestPieceCount), parent.FinishedPieces.Count())
	for i := uint(0); i < uint(seedCacheTestPieceCount); i++ {
		require.True(t, parent.FinishedPieces.Test(i), "piece %d must be available", i)
	}
	stored, loaded := f.peers.Load(parent.ID)
	require.True(t, loaded)
	require.Same(t, parent, stored)
	stored, loaded = f.host.LoadPeer(parent.ID)
	require.True(t, loaded)
	require.Same(t, parent, stored)
	require.True(t, f.task.HasAvailablePeer("ordinary-client-host", set.NewSafeSet[string]()))
	require.Equal(t, int64(seedCacheTestContentLength), f.task.ContentLength.Load())
	require.Equal(t, seedCacheTestPieceLength, f.task.PieceLength)
	require.Equal(t, int32(seedCacheTestPieceCount), f.task.TotalPieceCount.Load())
	require.Empty(t, f.task.BackToSourcePeers.Values())
	f.task.Pieces.Range(func(_, _ any) bool {
		t.Error("cache recovery must not invent piece download records")
		return false
	})
	return parent
}

func TestSeedPeer_TriggerDownloadTaskRestoresCachedParent(t *testing.T) {
	f := newSeedCacheFixture(t)
	f.expectLocalTask(t, f.localTask(), nil, nil)
	// The seed has already finished the task and reports its MetadataOnly
	// announcement to a different scheduler. Its upload stream is empty.
	require.NoError(t, f.seed.TriggerDownloadTask(context.Background(), f.task.ID, f.request))
	parent := f.requireCachedParent(t)

	// Revalidate the local cache but reuse the same completed parent, without
	// resetting its state or creating another peer in the DAG.
	f.expectLocalTask(t, f.localTask(), nil, nil)
	require.NoError(t, f.seed.TriggerDownloadTask(context.Background(), f.task.ID, f.request))
	require.Same(t, parent, f.requireCachedParent(t))
}

func TestSeedPeer_TriggerDownloadTaskDoesNotImportAfterCancellation(t *testing.T) {
	f := newSeedCacheFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	f.expectLocalTask(t, f.localTask(), nil, cancel)
	_ = f.seed.TriggerDownloadTask(ctx, f.task.ID, f.request)
	require.Empty(t, f.task.LoadPeers())
}

func TestSeedPeer_TriggerDownloadTaskDoesNotImportAfterStreamError(t *testing.T) {
	f := newSeedCacheFixture(t)
	f.eof.Times(0)
	f.stream.EXPECT().Recv().Return(nil, errors.New("seed download failed")).Times(1)
	// No StatLocalTask expectation: a failed stream cannot prove a completed
	// cache, even if the seed might still hold some pieces locally.
	require.Error(t, f.seed.TriggerDownloadTask(context.Background(), f.task.ID, f.request))
	require.Empty(t, f.task.LoadPeers())
}

func TestSeedPeer_TriggerDownloadTaskKeepsAnnouncedParent(t *testing.T) {
	f := newSeedCacheFixture(t)
	parent := NewPeer("actual-announced-peer", f.task, f.host)
	parent.FSM.SetState(PeerStateSucceeded)
	f.peers.LoadOrStore(parent)
	// No StatLocalTask expectation: a completed real announcement already
	// gives this scheduler all the state it needs.
	require.NoError(t, f.seed.TriggerDownloadTask(context.Background(), f.task.ID, f.request))
	require.Equal(t, []*Peer{parent}, f.task.LoadPeers())
	require.Zero(t, parent.FinishedPieces.Count(), "the real peer's bitset must not be rewritten")
}

func TestSeedPeer_TriggerDownloadTaskKeepsAnnouncementReceivedDuringStat(t *testing.T) {
	f := newSeedCacheFixture(t)
	parent := NewPeer("actual-announced-peer", f.task, f.host)
	parent.FSM.SetState(PeerStateSucceeded)
	f.expectLocalTask(t, f.localTask(), nil, func() { f.peers.LoadOrStore(parent) })
	require.NoError(t, f.seed.TriggerDownloadTask(context.Background(), f.task.ID, f.request))
	require.Equal(t, []*Peer{parent}, f.task.LoadPeers())
	require.False(t, parent.SeedCache)
	require.Zero(t, parent.FinishedPieces.Count())
}

func TestSeedPeer_TriggerDownloadTaskDoesNotReviveUnavailableCachePeer(t *testing.T) {
	for _, state := range []string{PeerStateFailed, PeerStateLeave} {
		t.Run(state, func(t *testing.T) {
			f := newSeedCacheFixture(t)
			f.expectLocalTask(t, f.localTask(), nil, nil)
			require.NoError(t, f.seed.TriggerDownloadTask(context.Background(), f.task.ID, f.request))
			parent := f.requireCachedParent(t)
			parent.FSM.SetState(state)
			f.expectLocalTask(t, f.localTask(), nil, nil)
			require.Error(t, f.seed.TriggerDownloadTask(context.Background(), f.task.ID, f.request))
			require.Equal(t, []*Peer{parent}, f.task.LoadPeers())
			require.Equal(t, state, parent.FSM.Current())
			require.False(t, f.task.HasAvailablePeer("ordinary-client-host", set.NewSafeSet[string]()))
		})
	}
}

func TestSeedPeer_TriggerDownloadTaskRejectsUnverifiedCache(t *testing.T) {
	tests := []struct {
		name       string
		before     func(*seedCacheFixture)
		change     func(*dfdaemonv2.StatLocalTaskResponse)
		duringStat func(*seedCacheFixture)
		statErr    error
		skipStat   bool
		nilReply   bool
	}{
		{name: "missing local task", nilReply: true},
		{name: "unfinished task", change: func(r *dfdaemonv2.StatLocalTaskResponse) { r.FinishedAt = nil }},
		{name: "invalid finish timestamp", change: func(r *dfdaemonv2.StatLocalTaskResponse) { r.FinishedAt.Nanos = 1_000_000_000 }},
		{name: "task failed after completion", change: func(r *dfdaemonv2.StatLocalTaskResponse) { r.FailedAt = timestamppb.New(time.Now()) }},
		{name: "task has an older failure marker", change: func(r *dfdaemonv2.StatLocalTaskResponse) {
			r.FailedAt = timestamppb.New(r.FinishedAt.AsTime().Add(-time.Hour))
		}},
		{name: "wrong task", change: func(r *dfdaemonv2.StatLocalTaskResponse) { r.TaskId = "another-task" }},
		{name: "missing content length", change: func(r *dfdaemonv2.StatLocalTaskResponse) { r.ContentLength = nil }},
		{name: "wrong content length", change: func(r *dfdaemonv2.StatLocalTaskResponse) { *r.ContentLength-- }},
		{name: "missing piece length", change: func(r *dfdaemonv2.StatLocalTaskResponse) { r.PieceLength = nil }},
		{name: "zero piece length", change: func(r *dfdaemonv2.StatLocalTaskResponse) { *r.PieceLength = 0 }},
		{name: "different piece length", change: func(r *dfdaemonv2.StatLocalTaskResponse) { *r.PieceLength *= 2 }},
		{name: "different piece count", before: func(f *seedCacheFixture) { f.task.TotalPieceCount.Store(7) }},
		{name: "local stat failed", statErr: errors.New("seed unavailable")},
		{name: "partial download", before: func(f *seedCacheFixture) {
			f.request.Download.Range = &commonv2.Range{Length: seedCacheTestPieceLength}
		}, skipStat: true},
		{name: "unknown task metadata", before: func(f *seedCacheFixture) { f.task.ContentLength.Store(-1) }, skipStat: true},
		{name: "task was deleted", before: func(f *seedCacheFixture) { f.tasks.Delete(f.task.ID) }, skipStat: true},
		{name: "host was deleted", before: func(f *seedCacheFixture) { f.hosts.Delete(f.host.ID) }, skipStat: true},
		{name: "task deleted while stat runs", duringStat: func(f *seedCacheFixture) { f.tasks.Delete(f.task.ID) }},
		{name: "host deleted while stat runs", duringStat: func(f *seedCacheFixture) { f.hosts.Delete(f.host.ID) }},
		{name: "task leaves while stat runs", duringStat: func(f *seedCacheFixture) { f.task.FSM.SetState(TaskStateLeave) }},
		{name: "task length reset while stat runs", duringStat: func(f *seedCacheFixture) { f.task.ContentLength.Store(-1) }},
		{name: "task piece count reset while stat runs", duringStat: func(f *seedCacheFixture) { f.task.TotalPieceCount.Store(0) }},
		{name: "task replaced while stat runs", duringStat: func(f *seedCacheFixture) {
			f.tasks.Delete(f.task.ID)
			replacement := NewTask(f.task.ID, f.task.URL, "", "", commonv2.TaskType_STANDARD, nil, nil, 1)
			f.tasks.LoadOrStore(replacement)
		}},
		{name: "host replaced while stat runs", duringStat: func(f *seedCacheFixture) {
			f.hosts.Delete(f.host.ID)
			replacement := NewHost(f.host.ID, f.host.IP, f.host.Name, f.host.Hostname, 4000, 4001, 4002, types.HostTypeSuperSeed)
			f.hosts.LoadOrStore(replacement)
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newSeedCacheFixture(t)
			if tc.before != nil {
				tc.before(f)
			}
			if !tc.skipStat {
				response := f.localTask()
				if tc.change != nil {
					tc.change(response)
				}
				if tc.nilReply {
					response = nil
				}
				f.expectLocalTask(t, response, tc.statErr, func() {
					if tc.duringStat != nil {
						tc.duringStat(f)
					}
				})
			}
			_ = f.seed.TriggerDownloadTask(context.Background(), f.task.ID, f.request)
			require.Empty(t, f.task.LoadPeers(), "an unverified cache must never become a parent")
			f.peers.Range(func(_, _ any) bool {
				t.Error("an unverified cache must not enter the peer manager")
				return false
			})
		})
	}
}

func TestSeedPeer_TriggerDownloadTaskConcurrentCacheRecovery(t *testing.T) {
	f := newSeedCacheFixture(t)
	response := f.localTask()
	// Multiple successful empty streams may race with the first import. The
	// response is immutable, as are cached peers once they become visible.
	f.client.EXPECT().StatLocalTask(gomock.Any(), gomock.Any()).Return(response, nil).AnyTimes()
	const callers = 24
	start := make(chan struct{})
	results := make(chan error, callers)
	stopReading, readerDone := make(chan struct{}), make(chan struct{})
	observedParent := make(chan struct{})
	readFailure := make(chan string, 1)
	go func() {
		defer close(readerDone)
		observed := false
		for {
			select {
			case <-stopReading:
				return
			default:
			}
			for _, parent := range f.task.LoadPeers() {
				// Real scheduling readers may observe the peer as soon as the
				// DAG publishes it. They must always see its complete bitmap.
				if !parent.FSM.Is(PeerStateSucceeded) || parent.FinishedPieces.Count() != uint(seedCacheTestPieceCount) {
					select {
					case readFailure <- "a partially initialized cached parent became visible":
					default:
					}
				}
				if !observed {
					close(observedParent)
					observed = true
				}
			}
			runtime.Gosched()
		}
	}()
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			results <- f.seed.TriggerDownloadTask(context.Background(), f.task.ID, f.request)
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	select {
	case <-observedParent:
	case <-time.After(time.Second):
		t.Error("the concurrent reader did not observe the recovered parent")
	}
	close(stopReading)
	<-readerDone
	for err := range results {
		require.NoError(t, err)
	}
	select {
	case message := <-readFailure:
		t.Fatal(message)
	default:
	}
	f.requireCachedParent(t)
}
