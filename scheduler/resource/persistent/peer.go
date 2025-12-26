/*
 *     Copyright 2025 The Dragonfly Authors
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

package persistent

import (
	"context"
	"time"

	"github.com/bits-and-blooms/bitset"
	"github.com/looplab/fsm"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

const (
	// defaultConcurrentPieceCount is the fallback value for concurrent pieces per peer
	// when ConcurrentPieceCount is not reported by the client.
	defaultConcurrentPieceCount uint32 = 8
)

const (
	// Peer has been created but did not start running.
	PeerStatePending = "Pending"

	// Peer is uploading resources for p2p cluster.
	PeerStateUploading = "Uploading"

	// Peer successfully registered as empty scope size.
	PeerStateReceivedEmpty = "ReceivedEmpty"

	// Peer successfully registered as normal scope size.
	PeerStateReceivedNormal = "ReceivedNormal"

	// Peer is downloading resources from peer.
	PeerStateRunning = "Running"

	// Peer has been downloaded successfully.
	PeerStateSucceeded = "Succeeded"

	// Peer has been downloaded failed.
	PeerStateFailed = "Failed"
)

const (
	// Peer is uploding.
	PeerEventUpload = "Upload"

	// Peer is registered as empty scope size.
	PeerEventRegisterEmpty = "RegisterEmpty"

	// Peer is registered as normal scope size.
	PeerEventRegisterNormal = "RegisterNormal"

	// Peer is downloading.
	PeerEventDownload = "Download"

	// Peer downloaded or uploaded successfully.
	PeerEventSucceeded = "Succeeded"

	// Peer downloaded or uploaded failed.
	PeerEventFailed = "Failed"
)

// PeerOption is a functional option for persistent peer.
type PeerOption func(peer *Peer)

// WithConcurrentPieceCount set ConcurrentPieceCount for peer.
func WithConcurrentPieceCount(count uint32) PeerOption {
	return func(p *Peer) {
		p.ConcurrentPieceCount = count
	}
}

// Peer contains content for persistent peer.
type Peer struct {
	// ID is persistent peer id.
	ID string

	// Persistent is whether the peer is persistent.
	Persistent bool

	// ConcurrentPieceCount is the count of pieces that can be downloaded concurrently.
	ConcurrentPieceCount uint32

	// Pieces is finished pieces bitset.
	FinishedPieces *bitset.BitSet

	// Persistent peer state machine.
	FSM *fsm.FSM

	// Task is persistent task.
	Task *Task

	// Host is the peer host.
	Host *Host

	// BlockParents is bad parents ids.
	BlockParents []string

	// Cost is the cost of downloading.
	Cost time.Duration

	// CreatedAt is persistent peer create time.
	CreatedAt time.Time

	// UpdatedAt is persistent peer update time.
	UpdatedAt time.Time

	// Persistent peer log.
	Log *logger.SugaredLoggerOnWith
}

// New persistent peer instance.
func NewPeer(id, state string, isPersistent bool, finishedPieces *bitset.BitSet, blockParents []string, task *Task, host *Host,
	cost time.Duration, createdAt, updatedAt time.Time, log *logger.SugaredLoggerOnWith, options ...PeerOption) *Peer {
	p := &Peer{
		ID:                   id,
		Persistent:           isPersistent,
		ConcurrentPieceCount: defaultConcurrentPieceCount,
		FinishedPieces:       finishedPieces,
		Task:                 task,
		Host:                 host,
		BlockParents:         blockParents,
		Cost:                 cost,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
		Log:                  logger.WithPeer(host.ID, task.ID, id),
	}

	// Initialize state machine.
	p.FSM = fsm.NewFSM(
		PeerStatePending,
		fsm.Events{
			fsm.EventDesc{Name: PeerEventUpload, Src: []string{PeerStatePending, PeerStateFailed}, Dst: PeerStateUploading},
			fsm.EventDesc{Name: PeerEventRegisterEmpty, Src: []string{PeerStatePending, PeerStateFailed}, Dst: PeerStateReceivedEmpty},
			fsm.EventDesc{Name: PeerEventRegisterNormal, Src: []string{PeerStatePending, PeerStateFailed}, Dst: PeerStateReceivedNormal},
			fsm.EventDesc{Name: PeerEventDownload, Src: []string{PeerStateReceivedEmpty, PeerStateReceivedNormal}, Dst: PeerStateRunning},
			fsm.EventDesc{Name: PeerEventSucceeded, Src: []string{PeerStateUploading, PeerStateRunning}, Dst: PeerStateSucceeded},
			fsm.EventDesc{Name: PeerEventFailed, Src: []string{PeerStatePending, PeerStateReceivedEmpty, PeerStateReceivedNormal, PeerStateUploading, PeerStateRunning}, Dst: PeerStateFailed},
		},
		fsm.Callbacks{
			PeerEventUpload: func(ctx context.Context, e *fsm.Event) {
				p.Log.Infof("peer state is %s", e.FSM.Current())
			},
			PeerEventRegisterEmpty: func(ctx context.Context, e *fsm.Event) {
				p.Log.Infof("peer state is %s", e.FSM.Current())
			},
			PeerEventRegisterNormal: func(ctx context.Context, e *fsm.Event) {
				p.Log.Infof("peer state is %s", e.FSM.Current())
			},
			PeerEventDownload: func(ctx context.Context, e *fsm.Event) {
				p.Log.Infof("peer state is %s", e.FSM.Current())
			},
			PeerEventSucceeded: func(ctx context.Context, e *fsm.Event) {
				p.Log.Infof("peer state is %s", e.FSM.Current())
			},
			PeerEventFailed: func(ctx context.Context, e *fsm.Event) {
				p.Log.Infof("peer state is %s", e.FSM.Current())
			},
		},
	)
	p.FSM.SetState(state)

	for _, opt := range options {
		opt(p)
	}

	return p
}
