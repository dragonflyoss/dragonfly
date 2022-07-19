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

package resource

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/looplab/fsm"
	"go.uber.org/atomic"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/container/set"
	"d7y.io/dragonfly/v2/pkg/dag"
	"d7y.io/dragonfly/v2/pkg/rpc/base"
	rpcscheduler "d7y.io/dragonfly/v2/pkg/rpc/scheduler"
)

const (
	// Tiny file size is 128 bytes.
	TinyFileSize = 128

	// Peer failure limit in task.
	FailedPeerCountLimit = 200

	// Peer count limit for task.
	PeerCountLimitForTask = 10 * 1000
)

const (
	// Task has been created but did not start running.
	TaskStatePending = "Pending"

	// Task is downloading resources from seed peer or back-to-source.
	TaskStateRunning = "Running"

	// Task has been downloaded successfully.
	TaskStateSucceeded = "Succeeded"

	// Task has been downloaded failed.
	TaskStateFailed = "Failed"
)

const (
	// Task is downloading.
	TaskEventDownload = "Download"

	// Task downloaded successfully.
	TaskEventDownloadSucceeded = "DownloadSucceeded"

	// Task downloaded failed.
	TaskEventDownloadFailed = "DownloadFailed"
)

// Option is a functional option for task.
type Option func(task *Task)

// WithBackToSourceLimit set BackToSourceLimit for task.
func WithBackToSourceLimit(limit int32) Option {
	return func(task *Task) {
		task.BackToSourceLimit.Add(limit)
	}
}

type Task struct {
	// ID is task id.
	ID string

	// URL is task download url.
	URL string

	// Type is task type.
	Type base.TaskType

	// URLMeta is task download url meta.
	URLMeta *base.UrlMeta

	// DirectPiece is tiny piece data.
	DirectPiece []byte

	// ContentLength is task total content length.
	ContentLength *atomic.Int64

	// TotalPieceCount is total piece count.
	TotalPieceCount *atomic.Int32

	// BackToSourceLimit is back-to-source limit.
	BackToSourceLimit *atomic.Int32

	// BackToSourcePeers is back-to-source sync map.
	BackToSourcePeers set.SafeSet

	// Task state machine.
	FSM *fsm.FSM

	// Piece sync map.
	Pieces *sync.Map

	// DAG is directed acyclic graph of peers.
	DAG dag.DAG

	// PeerFailedCount is peer failed count,
	// if one peer succeeds, the value is reset to zero.
	PeerFailedCount *atomic.Int32

	// CreateAt is task create time.
	CreateAt *atomic.Time

	// UpdateAt is task update time.
	UpdateAt *atomic.Time

	// Task log.
	Log *logger.SugaredLoggerOnWith
}

// New task instance.
func NewTask(id, url string, taskType base.TaskType, meta *base.UrlMeta, options ...Option) *Task {
	t := &Task{
		ID:                id,
		URL:               url,
		Type:              taskType,
		URLMeta:           meta,
		ContentLength:     atomic.NewInt64(0),
		TotalPieceCount:   atomic.NewInt32(0),
		BackToSourceLimit: atomic.NewInt32(0),
		BackToSourcePeers: set.NewSafeSet(),
		Pieces:            &sync.Map{},
		DAG:               dag.NewDAG(),
		PeerFailedCount:   atomic.NewInt32(0),
		CreateAt:          atomic.NewTime(time.Now()),
		UpdateAt:          atomic.NewTime(time.Now()),
		Log:               logger.WithTaskIDAndURL(id, url),
	}

	// Initialize state machine.
	t.FSM = fsm.NewFSM(
		TaskStatePending,
		fsm.Events{
			{Name: TaskEventDownload, Src: []string{TaskStatePending, TaskStateSucceeded, TaskStateFailed}, Dst: TaskStateRunning},
			{Name: TaskEventDownloadSucceeded, Src: []string{TaskStateRunning, TaskStateFailed}, Dst: TaskStateSucceeded},
			{Name: TaskEventDownloadFailed, Src: []string{TaskStateRunning}, Dst: TaskStateFailed},
		},
		fsm.Callbacks{
			TaskEventDownload: func(e *fsm.Event) {
				t.UpdateAt.Store(time.Now())
				t.Log.Infof("task state is %s", e.FSM.Current())
			},
			TaskEventDownloadSucceeded: func(e *fsm.Event) {
				t.UpdateAt.Store(time.Now())
				t.Log.Infof("task state is %s", e.FSM.Current())
			},
			TaskEventDownloadFailed: func(e *fsm.Event) {
				t.UpdateAt.Store(time.Now())
				t.Log.Infof("task state is %s", e.FSM.Current())
			},
		},
	)

	for _, opt := range options {
		opt(t)
	}

	return t
}

// LoadPeer return peer for a key.
func (t *Task) LoadPeer(key string) (*Peer, bool) {
	vertex, err := t.DAG.GetVertex(key)
	if err != nil {
		t.Log.Error(err)
		return nil, false
	}

	value := vertex.Value
	if value == nil {
		return nil, false
	}

	return value.(*Peer), true
}

// StorePeer set peer.
func (t *Task) StorePeer(peer *Peer) {
	t.DAG.AddVertex(peer.ID, peer) // nolint: errcheck
}

// DeletePeer deletes peer for a key.
func (t *Task) DeletePeer(key string) {
	if err := t.DeletePeerInEdges(key); err != nil {
		t.Log.Error(err)
	}

	if err := t.DeletePeerOutEdges(key); err != nil {
		t.Log.Error(err)
	}

	t.DAG.DeleteVertex(key)
}

// AddPeerEdge adds inedges between two peers.
func (t *Task) AddPeerEdge(fromPeer *Peer, toPeer *Peer) error {
	if err := t.DAG.AddEdge(fromPeer.ID, toPeer.ID); err != nil {
		return err
	}

	fromPeer.Host.UploadPeerCount.Inc()
	return nil
}

// DeleteInEdges deletes inedges of peer.
func (t *Task) DeletePeerInEdges(key string) error {
	vertex, err := t.DAG.GetVertex(key)
	if err != nil {
		return err
	}

	for _, value := range vertex.Parents.Values() {
		vertex, ok := value.(*dag.Vertex)
		if !ok {
			continue
		}

		vertexVal := vertex.Value
		if vertexVal == nil {
			continue
		}

		parent, ok := vertexVal.(*Peer)
		if !ok {
			continue
		}

		parent.Host.UploadPeerCount.Dec()
	}

	vertex.DeleteInEdges()
	return nil
}

// DeletePeerOutEdges deletes outedges of peer.
func (t *Task) DeletePeerOutEdges(key string) error {
	vertex, err := t.DAG.GetVertex(key)
	if err != nil {
		return err
	}

	value := vertex.Value
	if value == nil {
		return errors.New("vertex value is nil")
	}

	peer, ok := value.(*Peer)
	if !ok {
		return errors.New("vertex value is not peer")
	}

	peer.Host.UploadPeerCount.Sub(int32(vertex.Children.Len()))
	vertex.DeleteOutEdges()
	return nil
}

// PeerCount returns count of peer.
func (t *Task) PeerCount() int {
	return t.DAG.VertexCount()
}

// HasAvailablePeer returns whether there is an available peer.
func (t *Task) HasAvailablePeer() bool {
	var hasAvailablePeer bool
	for _, vertex := range t.DAG.Vertices() {
		value := vertex.Value
		if value == nil {
			continue
		}

		peer, ok := value.(*Peer)
		if !ok {
			continue
		}

		if peer.FSM.Is(PeerStateSucceeded) {
			hasAvailablePeer = true
			break
		}
	}

	return hasAvailablePeer
}

// LoadSeedPeer return latest seed peer in peers sync map.
func (t *Task) LoadSeedPeer() (*Peer, bool) {
	var peers []*Peer
	for _, vertex := range t.DAG.Vertices() {
		value := vertex.Value
		if value == nil {
			continue
		}

		peer, ok := value.(*Peer)
		if !ok {
			continue
		}

		if peer.Host.Type != HostTypeNormal {
			peers = append(peers, peer)
		}
	}

	sort.Slice(
		peers,
		func(i, j int) bool {
			return peers[i].UpdateAt.Load().After(peers[j].UpdateAt.Load())
		},
	)

	if len(peers) > 0 {
		return peers[0], true
	}

	return nil, false
}

// IsSeedPeerFailed returns whether the seed peer in the task failed.
func (t *Task) IsSeedPeerFailed() bool {
	seedPeer, ok := t.LoadSeedPeer()
	return ok && seedPeer.FSM.Is(PeerStateFailed) && time.Since(seedPeer.CreateAt.Load()) < SeedPeerFailedTimeout
}

// LoadPiece return piece for a key.
func (t *Task) LoadPiece(key int32) (*base.PieceInfo, bool) {
	rawPiece, ok := t.Pieces.Load(key)
	if !ok {
		return nil, false
	}

	return rawPiece.(*base.PieceInfo), ok
}

// StorePiece set piece.
func (t *Task) StorePiece(piece *base.PieceInfo) {
	t.Pieces.Store(piece.PieceNum, piece)
}

// LoadOrStorePiece returns piece the key if present.
// Otherwise, it stores and returns the given piece.
// The loaded result is true if the piece was loaded, false if stored.
func (t *Task) LoadOrStorePiece(piece *base.PieceInfo) (*base.PieceInfo, bool) {
	rawPiece, loaded := t.Pieces.LoadOrStore(piece.PieceNum, piece)
	return rawPiece.(*base.PieceInfo), loaded
}

// DeletePiece deletes piece for a key.
func (t *Task) DeletePiece(key int32) {
	t.Pieces.Delete(key)
}

// SizeScope return task size scope type.
func (t *Task) SizeScope() (base.SizeScope, error) {
	if t.ContentLength.Load() < 0 {
		return -1, errors.New("invalid content length")
	}

	if t.TotalPieceCount.Load() <= 0 {
		return -1, errors.New("invalid total piece count")
	}

	if t.ContentLength.Load() <= TinyFileSize {
		return base.SizeScope_TINY, nil
	}

	if t.TotalPieceCount.Load() == 1 {
		return base.SizeScope_SMALL, nil
	}

	return base.SizeScope_NORMAL, nil
}

// CanBackToSource represents whether peer can back-to-source.
func (t *Task) CanBackToSource() bool {
	return int32(t.BackToSourcePeers.Len()) < t.BackToSourceLimit.Load() && (t.Type == base.TaskType_Normal || t.Type == base.TaskType_DfStore)
}

// NotifyPeers notify all peers in the task with the state code.
func (t *Task) NotifyPeers(peerPacket *rpcscheduler.PeerPacket, event string) {
	for _, vertex := range t.DAG.Vertices() {
		value := vertex.Value
		if value == nil {
			continue
		}

		peer := value.(*Peer)
		if peer.FSM.Is(PeerStateRunning) {
			stream, ok := peer.LoadStream()
			if !ok {
				continue
			}

			if err := stream.Send(peerPacket); err != nil {
				t.Log.Errorf("send packet to peer %s failed: %s", peer.ID, err.Error())
				continue
			}
			t.Log.Infof("task notify peer %s code %s", peer.ID, peerPacket.Code)

			if err := peer.FSM.Event(event); err != nil {
				peer.Log.Errorf("peer fsm event failed: %s", err.Error())
				continue
			}
		}
	}
}
