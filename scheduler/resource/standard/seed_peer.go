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

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"stathat.com/c/consistent"

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

	// SeedPeerRefreshInterval is the interval of refreshing seed peers.
	SeedPeerRefreshInterval = 1 * time.Minute
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

	// HasAvailable returns whether there is any available seed peer.
	HasAvailable() bool

	// Serve serves the seed peer service.
	Serve() error

	// Stop seed peer service.
	Stop()
}

// seedPeer contains content for seed peer.
type seedPeer struct {
	// peerManager is PeerManager interface.
	peerManager PeerManager

	// hostManager is HostManager interface.
	hostManager HostManager

	// taskManager is TaskManager interface.
	taskManager TaskManager

	// clientPool is Pool interface.
	clientPool dfdaemonclient.Pool

	// backToSourceLimit is back-to-source limit of task.
	backToSourceLimit int32

	// dialOpts is the options for grpc dial.
	dialOptions []grpc.DialOption

	// hosts is the list of seed peers.
	hosts *sync.Map

	// hashring is the hashring constructed from seed peers.
	hashring *consistent.Consistent

	// done is the channel to stop the seed peer service.
	done chan struct{}
}

// New SeedPeer interface.
func newSeedPeer(peerManager PeerManager, hostManager HostManager, taskManager TaskManager, clientPool dfdaemonclient.Pool, backToSourceLimit int32, dialOptions ...grpc.DialOption) SeedPeer {
	return &seedPeer{
		peerManager:       peerManager,
		hostManager:       hostManager,
		taskManager:       taskManager,
		clientPool:        clientPool,
		backToSourceLimit: backToSourceLimit,
		dialOptions:       dialOptions,
		hosts:             &sync.Map{},
		hashring:          consistent.New(),
		done:              make(chan struct{}),
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
		return err
	}

	// Wait for the download task to complete. Meanwhile, mirror the seed peer
	// progress into scheduler memory. This is important after scheduler restarts:
	// the seed may still have a complete local task, but the scheduler has lost
	// the in-memory task/peer DAG. The DownloadTask stream is the on-demand
	// reconciliation point that makes the seed peer available as a parent again.
	var peer *Peer
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				if peer != nil {
					if err := s.completeSeedPeerDownloadTask(ctx, peer); err != nil {
						return err
					}
				}

				return nil
			}

			if peer != nil && !peer.FSM.Is(PeerStateFailed) {
				if err := peer.FSM.Event(ctx, PeerEventDownloadFailed); err != nil {
					peer.Log.Errorf("seed peer download task failed, but set peer failed state failed: %s", err.Error())
				}
			}

			return err
		}

		peer, err = s.handleDownloadTaskResponse(ctx, taskID, req.GetDownload(), resp, peer)
		if err != nil {
			return err
		}
	}
}

func (s *seedPeer) handleDownloadTaskResponse(ctx context.Context, taskID string, download *commonv2.Download, resp *dfdaemonv2.DownloadTaskResponse, peer *Peer) (*Peer, error) {
	if resp == nil {
		return peer, nil
	}

	if resp.GetTaskId() != "" && resp.GetTaskId() != taskID {
		return nil, fmt.Errorf("seed peer download task response task id %s does not match %s", resp.GetTaskId(), taskID)
	}

	if peer == nil {
		var err error
		peer, err = s.initSeedPeerForDownloadTask(ctx, taskID, download, resp)
		if err != nil {
			return nil, err
		}
	}

	if err := s.ensureSeedPeerDownloading(ctx, peer, download.GetNeedBackToSource()); err != nil {
		return nil, err
	}

	switch r := resp.GetResponse().(type) {
	case *dfdaemonv2.DownloadTaskResponse_DownloadTaskStartedResponse:
		started := r.DownloadTaskStartedResponse
		if started == nil {
			return peer, nil
		}

		peer.Task.ContentLength.Store(int64(started.GetContentLength()))
		if download.GetActualPieceCount() > 0 {
			peer.Task.TotalPieceCount.Store(int32(download.GetActualPieceCount()))
		} else {
			if count := calculatePieceCount(started.GetContentLength(), peer.Task.PieceLength); count > 0 {
				peer.Task.TotalPieceCount.Store(count)
			}
		}

		// When the seed local task is already complete, dfdaemon may return all
		// pieces in the started response and set is_finished. Store those pieces and
		// mark the seed peer completed immediately, so later scheduling retries can
		// select it as parent without waiting for another scheduler announce path.
		if started.GetIsFinished() {
			for _, piece := range started.GetPieces() {
				if err := storeSeedPeerPiece(peer, piece); err != nil {
					return nil, err
				}
			}

			if err := s.completeSeedPeerDownloadTask(ctx, peer); err != nil {
				return nil, err
			}
		}
	case *dfdaemonv2.DownloadTaskResponse_DownloadPieceFinishedResponse:
		finished := r.DownloadPieceFinishedResponse
		if finished == nil {
			return peer, nil
		}

		if err := storeSeedPeerPiece(peer, finished.GetPiece()); err != nil {
			return nil, err
		}
	}

	return peer, nil
}

func (s *seedPeer) initSeedPeerForDownloadTask(ctx context.Context, taskID string, download *commonv2.Download, resp *dfdaemonv2.DownloadTaskResponse) (*Peer, error) {
	if download == nil {
		download = &commonv2.Download{}
	}

	hostID := resp.GetHostId()
	if hostID == "" {
		return nil, fmt.Errorf("seed peer download task response host id is empty")
	}

	peerID := resp.GetPeerId()
	if peerID == "" {
		return nil, fmt.Errorf("seed peer download task response peer id is empty")
	}

	host, loaded := s.hostManager.Load(hostID)
	if !loaded {
		return nil, fmt.Errorf("can not find seed host id: %s", hostID)
	}
	host.UpdatedAt.Store(time.Now())

	task, loaded := s.taskManager.Load(taskID)
	if !loaded {
		options := []TaskOption{WithPieceLength(download.GetActualPieceLength())}
		if download.GetDigest() != "" {
			d, err := digest.Parse(download.GetDigest())
			if err != nil {
				return nil, err
			}

			options = append(options, WithDigest(d))
		}

		newTask := NewTask(taskID, download.GetUrl(), download.GetTag(), download.GetApplication(), download.GetType(),
			download.GetFilteredQueryParams(), download.GetRequestHeader(), s.backToSourceLimit, options...)
		task, loaded = s.taskManager.LoadOrStore(newTask)
	}
	if loaded {
		task.URL = download.GetUrl()
		task.FilteredQueryParams = download.GetFilteredQueryParams()
		task.Header = download.GetRequestHeader()
		if task.PieceLength == 0 && download.GetActualPieceLength() > 0 {
			task.PieceLength = download.GetActualPieceLength()
		}
	}

	if download.GetActualContentLength() > 0 || download.ActualContentLength != nil {
		task.ContentLength.Store(int64(download.GetActualContentLength()))
	}
	if download.GetActualPieceCount() > 0 || download.ActualPieceCount != nil {
		task.TotalPieceCount.Store(int32(download.GetActualPieceCount()))
	}

	if !task.FSM.Is(TaskStateRunning) && !task.FSM.Is(TaskStateSucceeded) {
		if err := task.FSM.Event(ctx, TaskEventDownload); err != nil {
			return nil, err
		}
	}

	peer, loaded := s.peerManager.Load(peerID)
	if loaded {
		peer.UpdatedAt.Store(time.Now())
		return peer, nil
	}

	options := []PeerOption{WithPriority(download.GetPriority())}
	if download.GetRange() != nil {
		options = append(options, WithRange(http.Range{Start: int64(download.GetRange().GetStart()), Length: int64(download.GetRange().GetLength())}))
	}
	if download.ConcurrentPieceCount != nil {
		options = append(options, WithConcurrentPieceCount(download.GetConcurrentPieceCount()))
	}

	peer = NewPeer(peerID, task, host, options...)
	s.peerManager.Store(peer)
	peer.Log.Info("seed peer has been stored from download task response")

	if err := peer.FSM.Event(ctx, PeerEventRegisterNormal); err != nil {
		return nil, err
	}

	return peer, nil
}

func (s *seedPeer) ensureSeedPeerDownloading(ctx context.Context, peer *Peer, needBackToSource bool) error {
	if !peer.Task.FSM.Is(TaskStateRunning) && !peer.Task.FSM.Is(TaskStateSucceeded) {
		if err := peer.Task.FSM.Event(ctx, TaskEventDownload); err != nil {
			return err
		}
	}

	if !peer.FSM.Is(PeerStateReceivedNormal) && !peer.FSM.Is(PeerStateReceivedSmall) && !peer.FSM.Is(PeerStateReceivedTiny) && !peer.FSM.Is(PeerStateReceivedEmpty) {
		return nil
	}

	if needBackToSource {
		return peer.FSM.Event(ctx, PeerEventDownloadBackToSource)
	}

	return peer.FSM.Event(ctx, PeerEventDownload)
}

func (s *seedPeer) completeSeedPeerDownloadTask(ctx context.Context, peer *Peer) error {
	if !peer.FSM.Is(PeerStateSucceeded) {
		if err := peer.FSM.Event(ctx, PeerEventDownloadSucceeded); err != nil {
			return err
		}
	}

	if peer.Range == nil && !peer.Task.FSM.Is(TaskStateSucceeded) {
		if !peer.Task.FSM.Is(TaskStateRunning) {
			if err := peer.Task.FSM.Event(ctx, TaskEventDownload); err != nil {
				return err
			}
		}

		if err := peer.Task.FSM.Event(ctx, TaskEventDownloadSucceeded); err != nil {
			return err
		}
	}

	return nil
}

func storeSeedPeerPiece(peer *Peer, protoPiece *commonv2.Piece) error {
	if protoPiece == nil {
		return nil
	}

	cost := time.Duration(0)
	if protoPiece.GetCost() != nil {
		cost = protoPiece.GetCost().AsDuration()
	}

	createdAt := time.Now()
	if protoPiece.GetCreatedAt() != nil {
		createdAt = protoPiece.GetCreatedAt().AsTime()
	} else if cost > 0 {
		createdAt = createdAt.Add(-cost)
	}

	piece := &Piece{
		Number:      int32(protoPiece.GetNumber()),
		ParentID:    protoPiece.GetParentId(),
		Offset:      protoPiece.GetOffset(),
		Length:      protoPiece.GetLength(),
		TrafficType: protoPiece.GetTrafficType(),
		Cost:        cost,
		CreatedAt:   createdAt,
	}

	if len(protoPiece.GetDigest()) > 0 {
		d, err := digest.Parse(protoPiece.GetDigest())
		if err != nil {
			return err
		}

		piece.Digest = d
	}

	peer.FinishedPieces.Set(uint(piece.Number))
	peer.AppendPieceCost(piece.Cost)
	peer.PieceUpdatedAt.Store(time.Now())
	peer.UpdatedAt.Store(time.Now())
	peer.Host.UpdatedAt.Store(time.Now())
	peer.Task.StorePiece(piece)
	peer.Task.UpdatedAt.Store(time.Now())
	if peer.Task.TotalPieceCount.Load() <= piece.Number {
		peer.Task.TotalPieceCount.Store(piece.Number + 1)
	}

	metrics.DownloadPieceCount.WithLabelValues(piece.TrafficType.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()
	metrics.Traffic.WithLabelValues(piece.TrafficType.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Add(float64(piece.Length))

	return nil
}

func calculatePieceCount(contentLength uint64, pieceLength uint64) int32 {
	if pieceLength == 0 {
		return 0
	}

	return int32((contentLength + pieceLength - 1) / pieceLength)
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
	// The synchronization of the hash ring is handled by the refreshSeedPeers periodically and asynchronously.
	if len(s.hashring.Members()) == 0 {
		return nil, fmt.Errorf("no seed peer available")
	}

	addr, err := s.hashring.Get(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to select seed peer: %w", err)
	}

	host, ok := s.hosts.Load(addr)
	if !ok {
		return nil, fmt.Errorf("failed to load host: %s", addr)
	}

	return host.(*Host), nil
}

// HasAvailable returns whether there is any available seed peer.
func (s *seedPeer) HasAvailable() bool {
	return len(s.hashring.Members()) > 0
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
	s.peerManager.Store(peer)
	peer.Log.Info("seed peer has been stored")

	if err := peer.FSM.Event(ctx, PeerEventRegisterNormal); err != nil {
		return nil, err
	}

	return peer, nil
}

func (s *seedPeer) refresh(ctx context.Context) {
	hosts := s.hostManager.LoadAllSeeds()
	if len(hosts) == 0 {
		logger.Warnf("no seed peer found in host manager")
		return
	}

	healthyHosts := &sync.Map{}
	// Do the health check for each seed peer.
	for _, host := range hosts {
		addr := net.JoinHostPort(host.IP, strconv.Itoa(int(host.Port)))
		if err := healthclient.Check(ctx, addr, s.dialOptions...); err != nil {
			logger.Errorf("failed to check the healthy for seed peer %s: %v", addr, err)
		} else {
			healthyHosts.Store(addr, host)
		}
	}
	s.hosts = healthyHosts

	hashring := consistent.New()
	s.hosts.Range(func(addr, _ any) bool {
		hashring.Add(addr.(string))
		return true
	})

	s.hashring = hashring
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
