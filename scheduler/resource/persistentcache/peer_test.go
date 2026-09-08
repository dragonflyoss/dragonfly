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

package persistentcache

import (
	"context"
	"testing"
	"time"

	"github.com/bits-and-blooms/bitset"
	"github.com/looplab/fsm"
	"github.com/stretchr/testify/assert"
)

func TestNewPeer(t *testing.T) {
	tests := []struct {
		name       string
		state      string
		persistent bool
		options    []PeerOption
		expect     func(t *testing.T, peer *Peer)
	}{
		{
			name:       "new peer with pending state",
			state:      PeerStatePending,
			persistent: true,
			options:    []PeerOption{WithConcurrentPieceCount(1)},
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				assert.Equal("peer-1", peer.ID)
				assert.True(peer.Persistent)
				assert.Equal(uint32(1), peer.ConcurrentPieceCount)
				assert.Equal(bitset.New(64), peer.FinishedPieces)
				assert.Equal([]string{"parent-1"}, peer.BlockParents)
				assert.Equal("task-1", peer.Task.ID)
				assert.Equal("host-1", peer.Host.ID)
				assert.Equal(time.Second, peer.Cost)
				assert.Equal(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), peer.CreatedAt)
				assert.Equal(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), peer.UpdatedAt)
				assert.Equal(PeerStatePending, peer.FSM.Current())
				assert.NotNil(peer.Log)
			},
		},
		{
			name:       "new peer with running state",
			state:      PeerStateRunning,
			persistent: false,
			options:    []PeerOption{WithConcurrentPieceCount(2)},
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				assert.Equal("peer-1", peer.ID)
				assert.False(peer.Persistent)
				assert.Equal(uint32(2), peer.ConcurrentPieceCount)
				assert.Equal(bitset.New(64), peer.FinishedPieces)
				assert.Equal([]string{"parent-1"}, peer.BlockParents)
				assert.Equal("task-1", peer.Task.ID)
				assert.Equal("host-1", peer.Host.ID)
				assert.Equal(time.Second, peer.Cost)
				assert.Equal(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), peer.CreatedAt)
				assert.Equal(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), peer.UpdatedAt)
				assert.Equal(PeerStateRunning, peer.FSM.Current())
				assert.NotNil(peer.Log)
			},
		},
		{
			name:       "new peer without options uses default concurrent piece count",
			state:      PeerStatePending,
			persistent: true,
			expect: func(t *testing.T, peer *Peer) {
				assert := assert.New(t)
				assert.Equal("peer-1", peer.ID)
				assert.True(peer.Persistent)
				assert.Equal(defaultConcurrentPieceCount, peer.ConcurrentPieceCount)
				assert.Equal(bitset.New(64), peer.FinishedPieces)
				assert.Equal([]string{"parent-1"}, peer.BlockParents)
				assert.Equal("task-1", peer.Task.ID)
				assert.Equal("host-1", peer.Host.ID)
				assert.Equal(time.Second, peer.Cost)
				assert.Equal(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), peer.CreatedAt)
				assert.Equal(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), peer.UpdatedAt)
				assert.Equal(PeerStatePending, peer.FSM.Current())
				assert.NotNil(peer.Log)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, NewPeer(
				"peer-1",
				tc.state,
				tc.persistent,
				bitset.New(64),
				[]string{"parent-1"},
				&Task{ID: "task-1"},
				&Host{ID: "host-1", Hostname: "host-1", IP: "127.0.0.1"},
				time.Second,
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				nil,
				tc.options...,
			))
		})
	}
}

func TestPeer_FSM(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		event  string
		expect func(t *testing.T, peer *Peer, err error)
	}{
		{
			name:  "upload from pending",
			state: PeerStatePending,
			event: PeerEventUpload,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateUploading, peer.FSM.Current())
			},
		},
		{
			name:  "upload from failed",
			state: PeerStateFailed,
			event: PeerEventUpload,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateUploading, peer.FSM.Current())
			},
		},
		{
			name:  "upload from running is rejected",
			state: PeerStateRunning,
			event: PeerEventUpload,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateRunning, peer.FSM.Current())
			},
		},
		{
			name:  "register empty from pending",
			state: PeerStatePending,
			event: PeerEventRegisterEmpty,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateReceivedEmpty, peer.FSM.Current())
			},
		},
		{
			name:  "register empty from succeeded is rejected",
			state: PeerStateSucceeded,
			event: PeerEventRegisterEmpty,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateSucceeded, peer.FSM.Current())
			},
		},
		{
			name:  "register normal from failed",
			state: PeerStateFailed,
			event: PeerEventRegisterNormal,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateReceivedNormal, peer.FSM.Current())
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
			name:  "download from received empty",
			state: PeerStateReceivedEmpty,
			event: PeerEventDownload,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateRunning, peer.FSM.Current())
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
			name:  "succeeded from uploading",
			state: PeerStateUploading,
			event: PeerEventSucceeded,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateSucceeded, peer.FSM.Current())
			},
		},
		{
			name:  "succeeded from running",
			state: PeerStateRunning,
			event: PeerEventSucceeded,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateSucceeded, peer.FSM.Current())
			},
		},
		{
			name:  "succeeded from pending is rejected",
			state: PeerStatePending,
			event: PeerEventSucceeded,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStatePending, peer.FSM.Current())
			},
		},
		{
			name:  "failed from uploading",
			state: PeerStateUploading,
			event: PeerEventFailed,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateFailed, peer.FSM.Current())
			},
		},
		{
			name:  "failed from running",
			state: PeerStateRunning,
			event: PeerEventFailed,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PeerStateFailed, peer.FSM.Current())
			},
		},
		{
			name:  "failed from pending is rejected",
			state: PeerStatePending,
			event: PeerEventFailed,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStatePending, peer.FSM.Current())
			},
		},
		{
			name:  "failed from received normal is rejected",
			state: PeerStateReceivedNormal,
			event: PeerEventFailed,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateReceivedNormal, peer.FSM.Current())
			},
		},
		{
			name:  "failed from succeeded is rejected",
			state: PeerStateSucceeded,
			event: PeerEventFailed,
			expect: func(t *testing.T, peer *Peer, err error) {
				assert := assert.New(t)
				assert.ErrorAs(err, new(fsm.InvalidEventError))
				assert.Equal(PeerStateSucceeded, peer.FSM.Current())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			peer := NewPeer("peer-1", tc.state, true, bitset.New(64), nil, &Task{ID: "task-1"}, &Host{ID: "host-1"}, 0, time.Now(), time.Now(), nil)
			err := peer.FSM.Event(context.Background(), tc.event)
			tc.expect(t, peer, err)
		})
	}
}
