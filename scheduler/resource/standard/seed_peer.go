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

//go:generate mockgen -destination seed_peer_mock.go -source seed_peer.go -package standard

package standard

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"d7y.io/hashring"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cdnsystemv1 "d7y.io/api/v2/pkg/apis/cdnsystem/v1"
	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	dfdaemonv2 "d7y.io/api/v2/pkg/apis/dfdaemon/v2"
	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/dfnet"
	"d7y.io/dragonfly/v2/pkg/digest"
	"d7y.io/dragonfly/v2/pkg/idgen"
	"d7y.io/dragonfly/v2/pkg/net/http"
	cndsystemclient "d7y.io/dragonfly/v2/pkg/rpc/cdnsystem/client"
	"d7y.io/dragonfly/v2/pkg/rpc/common"
	dfdaemonclient "d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/client"
	healthclient "d7y.io/dragonfly/v2/pkg/rpc/health/client"
	"d7y.io/dragonfly/v2/scheduler/metrics"
)

const (
	// Default value of seed peer failed timeout.
	SeedPeerFailedTimeout = 30 * time.Minute

	// SeedPeerRefreshInterval is the interval of refreshing seed peers, which
	// drops the seed peers whose announcements have stopped.
	SeedPeerRefreshInterval = 30 * time.Second

	// SeedPeerHealthCheckConcurrencyLimit is the number of seed peers checked
	// for health concurrently by refresh.
	SeedPeerHealthCheckConcurrencyLimit = 32

	// Bound the local cache lookup after the seed's download completes.
	seedPeerCacheStatTimeout = 5 * time.Second
)

// SeedPeer is the interface used for seed peer.
type SeedPeer interface {
	// TriggerDownloadTask triggers the seed peer to download task.
	// Used only in v2 version of the grpc.
	TriggerDownloadTask(context.Context, string, *dfdaemonv2.DownloadTaskRequest) error

	// TriggerTask triggers the seed peer to download task.
	// Used only in v1 version of the grpc.
	TriggerTask(context.Context, *http.Range, *Task) (*Peer, *schedulerv1.PeerResult, error)

	// Select selects a seed peer target by the task id.
	Select(context.Context, string) (*Host, error)

	// Store stores the announced seed peer in the hashring.
	Store(*Host)

	// Delete deletes the seed peer from the hashring.
	Delete(*Host)

	// HasAvailable returns whether there is any available seed peer.
	HasAvailable() bool

	// Serve serves the seed peer service.
	Serve() error

	// Stop seed peer service.
	Stop()
}

// seedPeer contains content for seed peer.
type seedPeer struct {
	// taskManager is TaskManager interface.
	taskManager TaskManager

	// peerManager is PeerManager interface.
	peerManager PeerManager

	// hostManager is HostManager interface.
	hostManager HostManager

	// clientPool is Pool interface.
	clientPool dfdaemonclient.Pool

	// dialOpts is the options for grpc dial.
	dialOptions []grpc.DialOption

	// mu protects hosts and hashring.
	mu sync.RWMutex

	// hosts is the seed peers keyed by name.
	hosts map[string]*Host

	// hashring is the hashring constructed from the names of seed peers.
	hashring *hashring.VNodeHashRing

	// done is the channel to stop the seed peer service.
	done chan struct{}
}

// New SeedPeer interface.
func newSeedPeer(taskManager TaskManager, peerManager PeerManager, hostManager HostManager, clientPool dfdaemonclient.Pool, dialOptions ...grpc.DialOption) SeedPeer {
	return &seedPeer{
		taskManager: taskManager,
		peerManager: peerManager,
		hostManager: hostManager,
		clientPool:  clientPool,
		dialOptions: dialOptions,
		hosts:       map[string]*Host{},
		hashring:    hashring.New(),
		done:        make(chan struct{}),
	}
}

// TriggerDownloadTask triggers the seed peer to download task.
// Used only in v2 version of the grpc.
func (s *seedPeer) TriggerDownloadTask(ctx context.Context, taskID string, req *dfdaemonv2.DownloadTaskRequest) error {
	ctx, cancel := context.WithCancel(trace.ContextWithSpan(ctx, trace.SpanFromContext(ctx)))
	defer cancel()

	selected, err := s.Select(ctx, taskID)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(selected.IP, strconv.Itoa(int(selected.Port)))
	logger.Infof("selected seed peer %s for task %s", addr, taskID)

	client, err := s.clientPool.Get(addr, s.dialOptions...)
	if err != nil {
		return err
	}

	stream, err := client.DownloadTask(ctx, taskID, req)
	if err != nil {
		s.evict(selected, err)
		return err
	}

	// Wait for the download task to complete.
	for {
		_, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return s.restoreCachedPeer(ctx, selected, client, taskID, req)
			}

			s.evict(selected, err)
			return err
		}
	}
}

// restoreCachedPeer makes a completed seed cache available to the scheduler that
// triggered the download. Different client versions can hash the same task to
// different schedulers, so the seed's normal announcement may reach another
// scheduler. A cache hit can also return an empty DownloadTask stream.
func (s *seedPeer) restoreCachedPeer(ctx context.Context, selected *Host, client dfdaemonclient.Client, taskID string, req *dfdaemonv2.DownloadTaskRequest) error {
	if req.GetDownload() == nil || req.GetDownload().GetRange() != nil {
		return nil
	}

	task, loaded := s.taskManager.Load(taskID)
	if !loaded || task.FSM.Is(TaskStateLeave) || task.PieceLength == 0 ||
		task.ContentLength.Load() <= 0 || task.TotalPieceCount.Load() <= 0 {
		return nil
	}
	if host, loaded := s.hostManager.Load(selected.ID); !loaded || host != selected ||
		!host.IsSeedPeer() || host.DisableShared {
		return nil
	}
	if hasAnnouncedSeedCache(task, selected) {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, seedPeerCacheStatTimeout)
	defer cancel()
	// StatTask would consult the seed's own scheduler, reproducing the routing
	// mismatch. StatLocalTask checks the selected seed's local storage directly.
	cache, err := client.StatLocalTask(ctx, &dfdaemonv2.StatLocalTaskRequest{TaskId: taskID})
	if err != nil {
		return fmt.Errorf("stat completed seed cache for task %s: %w", taskID, err)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validateSeedCache(task, cache); err != nil {
		return err
	}

	// The lookup may have raced with host replacement, task GC, or a regular
	// announcement. Do not attach a cache peer to stale resource objects.
	if current, loaded := s.taskManager.Load(taskID); !loaded || current != task || task.FSM.Is(TaskStateLeave) {
		return nil
	}
	if current, loaded := s.hostManager.Load(selected.ID); !loaded || current != selected ||
		!selected.IsSeedPeer() || selected.DisableShared {
		return nil
	}
	if hasAnnouncedSeedCache(task, selected) {
		return nil
	}
	if err := validateSeedCache(task, cache); err != nil {
		return err
	}

	// Keep this identity separate from the seed's announced peers. In particular,
	// never advance an announcement's FSM or mutate its shared piece bitmap.
	peer := NewPeer("seed-cache-"+selected.ID+"-"+taskID, task, selected)
	peer.SeedCache = true
	pieceCount := (cache.GetContentLength()-1)/cache.GetPieceLength() + 1
	peer.FinishedPieces.FlipRange(0, uint(pieceCount))
	peer.FSM.SetState(PeerStateSucceeded)
	stored, loaded := s.peerManager.LoadOrStore(peer)
	if stored.Task != task || stored.Host != selected || !stored.SeedCache || !stored.FSM.Is(PeerStateSucceeded) {
		return fmt.Errorf("completed seed cache peer for task %s is no longer available", taskID)
	}
	stored.UpdatedAt.Store(time.Now())
	stored.PieceUpdatedAt.Store(time.Now())
	task.UpdatedAt.Store(time.Now())
	if !loaded {
		stored.Log.Info("restored completed seed cache from local task metadata")
	}
	return nil
}

// hasAnnouncedSeedCache reports only full, completed peers announced to this
// scheduler. A previously imported cache must be checked against the seed again.
func hasAnnouncedSeedCache(task *Task, selected *Host) bool {
	for _, peer := range task.LoadPeers() {
		if peer.Host == selected && !peer.SeedCache && peer.Range == nil && peer.FSM.Is(PeerStateSucceeded) {
			return true
		}
	}
	return false
}

func validateSeedCache(task *Task, cache *dfdaemonv2.StatLocalTaskResponse) error {
	if cache == nil || cache.GetTaskId() != task.ID {
		return fmt.Errorf("seed cache task ID does not match %s", task.ID)
	}
	if cache.FinishedAt == nil || cache.FinishedAt.CheckValid() != nil {
		return fmt.Errorf("seed cache task %s is not complete", task.ID)
	}
	if cache.FailedAt != nil {
		return fmt.Errorf("seed cache task %s has failed", task.ID)
	}
	if cache.PieceLength == nil || cache.ContentLength == nil || cache.GetPieceLength() == 0 ||
		cache.GetContentLength() == 0 || cache.GetPieceLength() != task.PieceLength ||
		task.ContentLength.Load() <= 0 || cache.GetContentLength() != uint64(task.ContentLength.Load()) {
		return fmt.Errorf("seed cache metadata does not match task %s", task.ID)
	}
	pieceCount := (cache.GetContentLength()-1)/cache.GetPieceLength() + 1
	if task.TotalPieceCount.Load() <= 0 || pieceCount != uint64(task.TotalPieceCount.Load()) {
		return fmt.Errorf("seed cache piece count does not match task %s", task.ID)
	}
	return nil
}

// evict deletes the seed peer from the hashring when the error reports it
// unreachable, the seed peer is stored again by its next announcement, refer to
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/outlier.
func (s *seedPeer) evict(host *Host, err error) {
	if status.Code(err) != codes.Unavailable {
		return
	}

	host.Log.Warnf("evict unavailable seed peer from hashring: %s", err)
	s.Delete(host)
}

// TriggerTask triggers the seed peer to download task.
// Used only in v1 version of the grpc.
func (s *seedPeer) TriggerTask(ctx context.Context, rg *http.Range, task *Task) (*Peer, *schedulerv1.PeerResult, error) {
	urlMeta := &commonv1.UrlMeta{
		Tag:         task.Tag,
		Filter:      idgen.FormatFilteredQueryParams(task.FilteredQueryParams),
		Header:      task.Header,
		Application: task.Application,
		Priority:    commonv1.Priority_LEVEL0,
	}

	if task.Digest != nil {
		urlMeta.Digest = task.Digest.String()
	}

	if rg != nil {
		urlMeta.Range = rg.URLMetaString()
	}

	selected, err := s.Select(ctx, task.ID)
	if err != nil {
		return nil, nil, err
	}

	addr := net.JoinHostPort(selected.IP, strconv.Itoa(int(selected.Port)))
	logger.Infof("selected seed peer %s for task %s", addr, task.ID)

	// TODO(chlins): reuse the client if we encounter the performance issue in future.
	client, err := cndsystemclient.GetClientByAddr(ctx, dfnet.NetAddr{Type: dfnet.TCP, Addr: addr}, s.dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	defer client.Close()

	stream, err := client.ObtainSeeds(ctx, &cdnsystemv1.SeedRequest{
		TaskId:  task.ID,
		Url:     task.URL,
		UrlMeta: urlMeta,
	})
	if err != nil {
		return nil, nil, err
	}

	var (
		peer        *Peer
		initialized bool
	)

	for {
		pieceSeed, err := stream.Recv()
		if err != nil {
			// If the peer initialization succeeds and the download fails,
			// set peer status is PeerStateFailed.
			if peer != nil {
				if err := peer.FSM.Event(ctx, PeerEventDownloadFailed); err != nil {
					return nil, nil, err
				}
			}

			return nil, nil, err
		}

		if !initialized {
			initialized = true

			// Initialize seed peer.
			peer, err = s.initSeedPeer(ctx, rg, task, pieceSeed.HostId, pieceSeed.PeerId)
			if err != nil {
				return nil, nil, err
			}
		}

		if pieceSeed.PieceInfo != nil {
			// Handle begin of piece.
			if pieceSeed.PieceInfo.PieceNum == common.BeginOfPiece {
				peer.Log.Infof("receive begin of piece from seed peer: %#v %#v", pieceSeed, pieceSeed.PieceInfo)
				if err := peer.FSM.Event(ctx, PeerEventDownload); err != nil {
					return nil, nil, err
				}

				continue
			}

			// Handle piece download successfully.
			peer.Log.Infof("receive piece from seed peer: %#v %#v", pieceSeed, pieceSeed.PieceInfo)
			cost := time.Duration(int64(pieceSeed.PieceInfo.DownloadCost) * int64(time.Millisecond))
			piece := &Piece{
				Number:      pieceSeed.PieceInfo.PieceNum,
				Offset:      pieceSeed.PieceInfo.RangeStart,
				Length:      uint64(pieceSeed.PieceInfo.RangeSize),
				TrafficType: commonv2.TrafficType_BACK_TO_SOURCE,
				Cost:        cost,
				CreatedAt:   time.Now().Add(-cost),
			}

			if len(pieceSeed.PieceInfo.PieceMd5) > 0 {
				piece.Digest = digest.New(digest.AlgorithmMD5, pieceSeed.PieceInfo.PieceMd5)
			}

			peer.FinishedPieces.Set(uint(pieceSeed.PieceInfo.PieceNum))
			peer.AppendPieceCost(piece.Cost)

			// When the piece is downloaded successfully,
			// peer.UpdatedAt needs to be updated to prevent
			// the peer from being GC during the download process.
			peer.UpdatedAt.Store(time.Now())
			peer.PieceUpdatedAt.Store(time.Now())
			task.StorePiece(piece)

			// Collect Traffic metrics.
			trafficType := commonv2.TrafficType_BACK_TO_SOURCE
			if pieceSeed.Reuse {
				trafficType = commonv2.TrafficType_LOCAL_PEER
			}
			metrics.Traffic.WithLabelValues(trafficType.String(), peer.Task.Type.String(),
				peer.Host.Type.Name()).Add(float64(pieceSeed.PieceInfo.RangeSize))
		}

		// Handle end of piece.
		if pieceSeed.Done {
			peer.Log.Infof("receive done piece")
			return peer, &schedulerv1.PeerResult{
				TotalPieceCount: pieceSeed.TotalPieceCount,
				ContentLength:   pieceSeed.ContentLength,
			}, nil
		}
	}
}

// Select selects a seed peer by the task id.
func (s *seedPeer) Select(ctx context.Context, taskID string) (*Host, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	vnode, ok := s.hashring.Get(taskID)
	if !ok {
		return nil, fmt.Errorf("no available seed peer")
	}

	host, ok := s.hosts[vnode.Name()]
	if !ok {
		return nil, fmt.Errorf("failed to load host: %s", vnode.Name())
	}

	return host, nil
}

// Store stores the announced seed peer in the hashring. The hashring is keyed
// by the name of the seed peer instead of the address, so that the seed peer
// keeps its position in the hashring when its ip changes, refer to
// https://github.com/grpc/proposal/blob/master/A76-ring-hash-improvements.md.
func (s *seedPeer) Store(host *Host) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.hosts[host.Name]; !ok {
		s.hashring.Add(host.Name)
	}

	// The most recently announced host of the seed peer wins, such as
	// the new pod after a restart.
	s.hosts[host.Name] = host
}

// Delete deletes the seed peer from the hashring.
func (s *seedPeer) Delete(host *Host) {
	s.deleteHosts([]*Host{host})
}

// deleteHosts deletes the hosts from the hashring at once, and keeps the
// seed peer when a more recently announced host of the same name has
// replaced the host.
func (s *seedPeer) deleteHosts(hosts []*Host) {
	s.mu.Lock()
	defer s.mu.Unlock()

	changed := false
	for _, host := range hosts {
		if current, ok := s.hosts[host.Name]; !ok || current != host {
			continue
		}

		delete(s.hosts, host.Name)
		changed = true
	}

	if changed {
		ring := hashring.New()
		for name := range s.hosts {
			ring.Add(name)
		}

		s.hashring = ring
	}
}

// HasAvailable returns whether there is any available seed peer.
func (s *seedPeer) HasAvailable() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return !s.hashring.IsEmpty()
}

// Initialize seed peer.
func (s *seedPeer) initSeedPeer(ctx context.Context, rg *http.Range, task *Task, hostID string, peerID string) (*Peer, error) {
	// Load host from manager.
	host, loaded := s.hostManager.Load(hostID)
	if !loaded {
		task.Log.Errorf("can not find seed host id: %s", hostID)
		return nil, fmt.Errorf("can not find host id: %s", hostID)
	}
	host.UpdatedAt.Store(time.Now())

	// Load peer from manager.
	peer, loaded := s.peerManager.Load(peerID)
	if loaded {
		return peer, nil
	}
	task.Log.Infof("can not find seed peer: %s", peerID)

	options := []PeerOption{}
	if rg != nil {
		options = append(options, WithRange(*rg))
	}

	// New and store seed peer without range.
	peer = NewPeer(peerID, task, host, options...)
	if peer, loaded = s.peerManager.LoadOrStore(peer); loaded {
		return peer, nil
	}
	peer.Log.Info("seed peer has been stored")

	if err := peer.FSM.Event(ctx, PeerEventRegisterNormal); err != nil {
		return nil, err
	}

	return peer, nil
}

// refresh drops the seed peers whose announcements have stopped, the one
// missed an announcement is dropped when it fails the health check, refer to
// https://www.cs.cornell.edu/projects/Quicksilver/public_pdfs/SWIM.pdf.
func (s *seedPeer) refresh(ctx context.Context) {
	now := time.Now()

	var expiredHosts, suspectedHosts []*Host
	s.mu.RLock()
	for _, host := range s.hosts {
		if host.AnnounceInterval <= 0 {
			continue
		}

		// If the host's elapsed exceeds twice the announcing interval,
		// then delete the seed peer without health check.
		elapsed := now.Sub(host.UpdatedAt.Load())
		if elapsed > host.AnnounceInterval*2 {
			expiredHosts = append(expiredHosts, host)
			continue
		}

		// If the host's elapsed exceeds the announcing interval,
		// then check the health of the seed peer.
		if elapsed > host.AnnounceInterval {
			suspectedHosts = append(suspectedHosts, host)
		}
	}
	s.mu.RUnlock()

	// Check the health of the suspected seed peers concurrently, an unhealthy
	// seed peer costs the health check timeout and must not delay the others.
	unhealthyHosts := make([]*Host, len(suspectedHosts))
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(SeedPeerHealthCheckConcurrencyLimit)
	for i, host := range suspectedHosts {
		eg.Go(func() error {
			addr := net.JoinHostPort(host.IP, strconv.Itoa(int(host.Port)))
			if err := healthclient.Check(ctx, addr, s.dialOptions...); err != nil {
				host.Log.Infof("host missed an announcement and failed to check the healthy: %s", err)
				unhealthyHosts[i] = host
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		logger.Errorf("failed to check the healthy for seed peers: %v", err)
	}

	// Delete the seed peers, and keep the host which announced during the refresh.
	var hosts []*Host
	for _, host := range append(expiredHosts, unhealthyHosts...) {
		if host == nil || host.UpdatedAt.Load().After(now) {
			continue
		}

		host.Log.Info("seed peer has been deleted from hashring")
		hosts = append(hosts, host)
	}

	if len(hosts) > 0 {
		s.deleteHosts(hosts)
	}
}

// Serve serves the seed peer service.
func (s *seedPeer) Serve() error {
	go s.clientPool.Serve()

	ticker := time.NewTicker(SeedPeerRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.refresh(context.Background())
		case <-s.done:
			return nil
		}
	}
}

// Stop seed peer service.
func (s *seedPeer) Stop() {
	close(s.done)
}
