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
	"testing"
	"time"

	"github.com/bits-and-blooms/bitset"
	"github.com/stretchr/testify/assert"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

func TestNewPeer(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		state          string
		persistent     bool
		finishedPieces *bitset.BitSet
		blockParents   []string
		task           *Task
		host           *Host
		cost           time.Duration
		createdAt      time.Time
		updatedAt      time.Time
		log            *logger.SugaredLoggerOnWith
		expectedState  string
	}{
		{
			name:           "new peer with pending state",
			id:             "peer-1",
			state:          PeerStatePending,
			persistent:     true,
			finishedPieces: bitset.New(64),
			blockParents:   []string{"parent-1"},
			task: &Task{
				ID: "task-1",
			},
			host: &Host{
				ID:       "host-1",
				Hostname: "host-1",
				IP:       "127.0.0.1",
			},
			cost:          time.Second,
			createdAt:     time.Now(),
			updatedAt:     time.Now(),
			expectedState: PeerStatePending,
		},
		{
			name:           "new peer with running state",
			id:             "peer-2",
			state:          PeerStateRunning,
			persistent:     false,
			finishedPieces: bitset.New(128),
			blockParents:   []string{"parent-2", "parent-3"},
			task: &Task{
				ID: "task-2",
			},
			host: &Host{
				ID:       "host-2",
				Hostname: "host-2",
				IP:       "127.0.0.2",
			},
			cost:          2 * time.Second,
			createdAt:     time.Now(),
			updatedAt:     time.Now(),
			expectedState: PeerStateRunning,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			peer := NewPeer(
				tc.id,
				tc.state,
				tc.persistent,
				tc.finishedPieces,
				tc.blockParents,
				tc.task,
				tc.host,
				tc.cost,
				tc.createdAt,
				tc.updatedAt,
				tc.log,
			)

			assert.Equal(t, tc.id, peer.ID)
			assert.Equal(t, tc.persistent, peer.Persistent)
			assert.Equal(t, tc.finishedPieces, peer.FinishedPieces)
			assert.Equal(t, tc.blockParents, peer.BlockParents)
			assert.Equal(t, tc.task, peer.Task)
			assert.Equal(t, tc.host, peer.Host)
			assert.Equal(t, tc.cost, peer.Cost)
			assert.Equal(t, tc.createdAt, peer.CreatedAt)
			assert.Equal(t, tc.updatedAt, peer.UpdatedAt)
			assert.Equal(t, tc.expectedState, peer.FSM.Current())
			assert.NotNil(t, peer.Log)
		})
	}
}
