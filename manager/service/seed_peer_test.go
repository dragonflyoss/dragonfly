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

func TestService_UpdateSeedPeer(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdateSeedPeerRequest
		expect func(t *testing.T, s *service, seedPeer *models.SeedPeer, err error)
	}{
		{
			name: "seed peer not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.UpdateSeedPeerRequest{IDC: "idc"},
			expect: func(t *testing.T, s *service, seedPeer *models.SeedPeer, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(seedPeer)
			},
		},
		{
			name: "zero values leave stored fields untouched",
			setup: func(t *testing.T, s *service) uint {
				seedPeerCluster := mockSeedPeerCluster(t, s.db, "foo")
				return mockSeedPeer(t, s.db, seedPeerCluster.ID, "seed-peer-1").ID
			},
			req: types.UpdateSeedPeerRequest{IDC: "idc", DownloadPort: 4002},
			expect: func(t *testing.T, s *service, seedPeer *models.SeedPeer, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.SeedPeer{}
				assert.NoError(s.db.First(&stored, seedPeer.ID).Error)
				assert.Equal("seed-peer-1", stored.Hostname)
				assert.Equal("idc", stored.IDC)
				assert.Equal("127.0.0.1", stored.IP)
				assert.Equal(int32(4000), stored.Port)
				assert.Equal(int32(4002), stored.DownloadPort)
				assert.Equal(uint(1), stored.SeedPeerClusterID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := mockService(t)
			id := tc.setup(t, s)

			seedPeer, err := s.UpdateSeedPeer(context.Background(), id, tc.req)
			tc.expect(t, s, seedPeer, err)
		})
	}
}

func TestService_GetSeedPeers(t *testing.T) {
	s := mockService(t)
	foo := mockSeedPeerCluster(t, s.db, "foo")
	bar := mockSeedPeerCluster(t, s.db, "bar")
	mockSeedPeer(t, s.db, foo.ID, "seed-peer-1")
	mockSeedPeer(t, s.db, foo.ID, "seed-peer-2")
	mockSeedPeer(t, s.db, bar.ID, "seed-peer-3")
	if err := s.db.Model(&models.SeedPeer{}).Where("host_name = ?", "seed-peer-2").Update("idc", "idc-2").Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		query  types.GetSeedPeersQuery
		expect func(t *testing.T, seedPeers []models.SeedPeer, count int64, err error)
	}{
		{
			name:  "paginates",
			query: types.GetSeedPeersQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, seedPeers []models.SeedPeer, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(seedPeers, 1)
				assert.Equal("seed-peer-3", seedPeers[0].Hostname)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by cluster",
			query: types.GetSeedPeersQuery{SeedPeerClusterID: foo.ID, Page: 1, PerPage: 10},
			expect: func(t *testing.T, seedPeers []models.SeedPeer, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(seedPeers, 2)
				assert.Equal(int64(2), count)
			},
		},
		{
			name:  "filter by cluster and idc",
			query: types.GetSeedPeersQuery{SeedPeerClusterID: foo.ID, IDC: "idc-2", Page: 1, PerPage: 10},
			expect: func(t *testing.T, seedPeers []models.SeedPeer, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(seedPeers, 1)
				assert.Equal("seed-peer-2", seedPeers[0].Hostname)
				assert.Equal(int64(1), count)
			},
		},
		{
			name:  "filter by hostname",
			query: types.GetSeedPeersQuery{Hostname: "seed-peer-3", Page: 1, PerPage: 10},
			expect: func(t *testing.T, seedPeers []models.SeedPeer, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(seedPeers, 1)
				assert.Equal(bar.ID, seedPeers[0].SeedPeerClusterID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			seedPeers, count, err := s.GetSeedPeers(context.Background(), tc.query)
			tc.expect(t, seedPeers, count, err)
		})
	}
}
