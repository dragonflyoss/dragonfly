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

	// Serve serves the seed peer service.
	Serve() error

	// Stop seed peer service.
	Stop() error
}

// seedPeer contains content for seed peer.
type seedPeer struct {
	// peerManager is PeerManager interface.
	peerManager PeerManager

	// hostManager is HostManager interface.
	hostManager HostManager

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
func newSeedPeer(peerManager PeerManager, hostManager HostManager, dialOptions ...grpc.DialOption) SeedPeer {
	return &seedPeer{
		peerManager: peerManager,
		hostManager: hostManager,
		dialOptions: dialOptions,
		hosts:       &sync.Map{},
		hashring:    consistent.New(),
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

	addr := fmt.Sprintf("%s:%d", selected.IP, selected.Port)
	logger.Infof("selected seed peer %s for task %s", addr, taskID)

	// TODO(chlins): reuse the client if we encounter the performance issue in future.
	client, err := dfdaemonclient.GetV2ByAddr(ctx, addr, s.dialOptions...)
	if err != nil {
		return err
	}
	defer client.Close()

	stream, err := client.DownloadTask(ctx, taskID, req)
	if err != nil {
		return err
	}

	// Wait for the download task to complete.
	for {
		_, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
	}
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

	addr := fmt.Sprintf("%s:%d", selected.IP, selected.Port)
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

			peer.StorePiece(piece)
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

func (s *seedPeer) refresh(ctx context.Context) error {
	hosts := s.hostManager.LoadAllSeeds()
	if len(hosts) == 0 {
		return fmt.Errorf("no seed peer found in host manager")
	}

	healthyHosts := &sync.Map{}
	// Do the health check for each seed peer.
	for _, host := range hosts {
		addr := fmt.Sprintf("%s:%d", host.IP, host.Port)
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
	return nil
}

// Serve serves the seed peer service.
func (s *seedPeer) Serve() error {
	ticker := time.NewTicker(SeedPeerRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.refresh(context.Background()); err != nil {
				logger.Errorf("failed to refresh seed peers: %v", err)
			}
		case <-s.done:
			return nil
		}
	}
}

// Stop seed peer service.
func (s *seedPeer) Stop() error {
	close(s.done)
	return nil
}
