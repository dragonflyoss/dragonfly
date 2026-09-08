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

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

func mockPeer(t *testing.T, db *gorm.DB, schedulerClusterID uint, hostname, peerType, state string) models.Peer {
	assert := assert.New(t)
	peer := models.Peer{
		Hostname:           hostname,
		Type:               peerType,
		IP:                 "127.0.0.1",
		Port:               4000,
		DownloadPort:       4001,
		ProxyPort:          4002,
		State:              state,
		SchedulerClusterID: schedulerClusterID,
	}
	assert.NoError(db.Create(&peer).Error)

	return peer
}

func TestService_DestroyPeer(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name: "peer not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "hard deletes",
			setup: func(t *testing.T, s *service) uint {
				schedulerCluster := mockSchedulerCluster(t, s.db, "foo")
				return mockPeer(t, s.db, schedulerCluster.ID, "peer-1", "normal", models.PeerStateActive).ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.Peer{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.DestroyPeer(context.Background(), id))
		})
	}
}

func TestService_GetPeers(t *testing.T) {
	s := mockService(t)
	foo := mockSchedulerCluster(t, s.db, "foo")
	bar := mockSchedulerCluster(t, s.db, "bar")
	mockPeer(t, s.db, foo.ID, "peer-1", "normal", models.PeerStateActive)
	mockPeer(t, s.db, foo.ID, "peer-2", "super", models.PeerStateInactive)
	mockPeer(t, s.db, bar.ID, "peer-3", "normal", models.PeerStateActive)

	tests := []struct {
		name   string
		query  types.GetPeersQuery
		expect func(t *testing.T, peers []models.Peer, count int64, err error)
	}{
		{
			name:  "paginates and preloads scheduler cluster",
			query: types.GetPeersQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, peers []models.Peer, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peers, 1)
				assert.Equal("peer-3", peers[0].Hostname)
				assert.Equal("bar", peers[0].SchedulerCluster.Name)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by cluster",
			query: types.GetPeersQuery{SchedulerClusterID: foo.ID, Page: 1, PerPage: 10},
			expect: func(t *testing.T, peers []models.Peer, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peers, 2)
				assert.Equal(int64(2), count)
			},
		},
		{
			name:  "filter by type and state",
			query: types.GetPeersQuery{Type: "normal", State: models.PeerStateActive, Page: 1, PerPage: 10},
			expect: func(t *testing.T, peers []models.Peer, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peers, 2)
				assert.Equal("peer-1", peers[0].Hostname)
				assert.Equal("peer-3", peers[1].Hostname)
			},
		},
		{
			name:  "filter by cluster and state",
			query: types.GetPeersQuery{SchedulerClusterID: foo.ID, State: models.PeerStateInactive, Page: 1, PerPage: 10},
			expect: func(t *testing.T, peers []models.Peer, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(peers, 1)
				assert.Equal("peer-2", peers[0].Hostname)
				assert.Equal(int64(1), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			peers, count, err := s.GetPeers(context.Background(), tc.query)
			tc.expect(t, peers, count, err)
		})
	}
}
