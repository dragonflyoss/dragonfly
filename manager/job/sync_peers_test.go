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

package job

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/config"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/pkg/idgen"
	pkgtypes "d7y.io/dragonfly/v2/pkg/types"
	resource "d7y.io/dragonfly/v2/scheduler/resource/standard"
)

func mockDB(t *testing.T, peers []*models.Peer) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "sync_peers.db")), &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormlogger.Discard,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&models.Peer{}); err != nil {
		t.Fatal(err)
	}

	if len(peers) > 0 {
		if err := db.Create(peers).Error; err != nil {
			t.Fatal(err)
		}
	}

	return db
}

func TestSyncPeers_mergePeers(t *testing.T) {
	tests := []struct {
		name    string
		peers   []*models.Peer
		results []*resource.Host
		expect  func(t *testing.T, peers []models.Peer)
	}{
		{
			name: "peers missing from the results are hard deleted across batches",
			peers: []*models.Peer{
				{Hostname: "peer-1", IP: "10.0.0.1", Type: pkgtypes.HostTypeNormalName, State: models.PeerStateActive, SchedulerClusterID: 1},
				{Hostname: "peer-2", IP: "10.0.0.2", Type: pkgtypes.HostTypeNormalName, State: models.PeerStateActive, SchedulerClusterID: 1},
				{Hostname: "peer-3", IP: "10.0.0.3", Type: pkgtypes.HostTypeNormalName, State: models.PeerStateActive, SchedulerClusterID: 1},
				{Hostname: "peer-4", IP: "10.0.0.4", Type: pkgtypes.HostTypeNormalName, State: models.PeerStateActive, SchedulerClusterID: 1},
				{Hostname: "peer-5", IP: "10.0.0.5", Type: pkgtypes.HostTypeNormalName, State: models.PeerStateActive, SchedulerClusterID: 1},
			},
			results: nil,
			expect: func(t *testing.T, peers []models.Peer) {
				assert := assert.New(t)
				assert.Empty(peers)
			},
		},
		{
			name:  "results without a scheduler cluster id are ignored",
			peers: nil,
			results: []*resource.Host{
				{ID: idgen.HostID("10.0.0.1", "peer-1", false), Hostname: "peer-1", IP: "10.0.0.1", Type: pkgtypes.HostTypeNormal, SchedulerClusterID: 0},
			},
			expect: func(t *testing.T, peers []models.Peer) {
				assert := assert.New(t)
				assert.Empty(peers)
			},
		},
		{
			name:  "peers only present in the results are inserted as active",
			peers: nil,
			results: []*resource.Host{
				{
					ID:                 idgen.HostID("10.0.0.1", "peer-1", false),
					Type:               pkgtypes.HostTypeNormal,
					Hostname:           "peer-1",
					IP:                 "10.0.0.1",
					Port:               8000,
					DownloadPort:       8001,
					ProxyPort:          8002,
					ObjectStoragePort:  8003,
					OS:                 "linux",
					Platform:           "ubuntu",
					PlatformFamily:     "debian",
					PlatformVersion:    "22.04",
					KernelVersion:      "6.1",
					Network:            resource.Network{IDC: "idc-1", Location: "loc-1"},
					Build:              resource.Build{GitVersion: "v2.3.0", GitCommit: "abc123", Platform: "linux/amd64"},
					SchedulerClusterID: 1,
				},
				{
					ID:                 idgen.HostID("10.0.0.2", "seed-1", true),
					Type:               pkgtypes.HostTypeSuperSeed,
					Hostname:           "seed-1",
					IP:                 "10.0.0.2",
					SchedulerClusterID: 1,
				},
			},
			expect: func(t *testing.T, peers []models.Peer) {
				assert := assert.New(t)
				assert.Len(peers, 2)

				byHostname := make(map[string]models.Peer, len(peers))
				for _, peer := range peers {
					byHostname[peer.Hostname] = peer
				}

				peer := byHostname["peer-1"]
				assert.Equal(models.PeerStateActive, peer.State)
				assert.Equal(pkgtypes.HostTypeNormalName, peer.Type)
				assert.Equal("10.0.0.1", peer.IP)
				assert.Equal(int32(8000), peer.Port)
				assert.Equal(int32(8001), peer.DownloadPort)
				assert.Equal(int32(8002), peer.ProxyPort)
				assert.Equal(int32(8003), peer.ObjectStoragePort)
				assert.Equal("linux", peer.OS)
				assert.Equal("ubuntu", peer.Platform)
				assert.Equal("debian", peer.PlatformFamily)
				assert.Equal("22.04", peer.PlatformVersion)
				assert.Equal("6.1", peer.KernelVersion)
				assert.Equal("idc-1", peer.IDC)
				assert.Equal("loc-1", peer.Location)
				assert.Equal("v2.3.0", peer.GitVersion)
				assert.Equal("abc123", peer.GitCommit)
				assert.Equal("linux/amd64", peer.BuildPlatform)
				assert.Equal(uint(1), peer.SchedulerClusterID)

				seed := byHostname["seed-1"]
				assert.Equal(models.PeerStateActive, seed.State)
				assert.Equal(pkgtypes.HostTypeSuperSeedName, seed.Type)
				assert.Equal(uint(1), seed.SchedulerClusterID)
			},
		},
		{
			name: "existing peers found in the results are kept and the others replaced",
			peers: []*models.Peer{
				{Hostname: "peer-1", IP: "10.0.0.1", Type: pkgtypes.HostTypeNormalName, State: models.PeerStateInactive, SchedulerClusterID: 1},
				{Hostname: "peer-2", IP: "10.0.0.2", Type: pkgtypes.HostTypeNormalName, State: models.PeerStateActive, SchedulerClusterID: 1},
				{Hostname: "peer-3", IP: "10.0.0.3", Type: pkgtypes.HostTypeNormalName, State: models.PeerStateActive, SchedulerClusterID: 1},
				{Hostname: "other-1", IP: "10.0.1.1", Type: pkgtypes.HostTypeNormalName, State: models.PeerStateActive, SchedulerClusterID: 2},
			},
			results: []*resource.Host{
				{ID: idgen.HostID("10.0.0.1", "peer-1", false), Hostname: "peer-1", IP: "10.0.0.1", Type: pkgtypes.HostTypeNormal, SchedulerClusterID: 1},
				{ID: idgen.HostID("10.0.0.4", "peer-4", false), Hostname: "peer-4", IP: "10.0.0.4", Type: pkgtypes.HostTypeNormal, Network: resource.Network{IDC: "idc-4"}, SchedulerClusterID: 1},
				{ID: idgen.HostID("10.0.0.5", "peer-5", false), Hostname: "peer-5", IP: "10.0.0.5", Type: pkgtypes.HostTypeNormal, SchedulerClusterID: 0},
			},
			expect: func(t *testing.T, peers []models.Peer) {
				assert := assert.New(t)
				hostnames := make([]string, 0, len(peers))
				byHostname := make(map[string]models.Peer, len(peers))
				for _, peer := range peers {
					hostnames = append(hostnames, peer.Hostname)
					byHostname[peer.Hostname] = peer
				}

				assert.ElementsMatch([]string{"peer-1", "peer-4", "other-1"}, hostnames)

				assert.Equal(uint(1), byHostname["peer-1"].SchedulerClusterID)
				assert.Equal(models.PeerStateActive, byHostname["peer-4"].State)
				assert.Equal("idc-4", byHostname["peer-4"].IDC)
				assert.Equal(uint(1), byHostname["peer-4"].SchedulerClusterID)
				assert.Equal(models.PeerStateActive, byHostname["other-1"].State)
				assert.Equal(uint(2), byHostname["other-1"].SchedulerClusterID)
			},
		},
		{
			name: "seed peer is matched by its seed host id",
			peers: []*models.Peer{
				{Hostname: "seed-1", IP: "10.0.0.1", Type: pkgtypes.HostTypeSuperSeedName, State: models.PeerStateActive, SchedulerClusterID: 1},
			},
			results: []*resource.Host{
				{ID: idgen.HostID("10.0.0.1", "seed-1", true), Hostname: "seed-1", IP: "10.0.0.1", Type: pkgtypes.HostTypeSuperSeed, SchedulerClusterID: 1},
			},
			expect: func(t *testing.T, peers []models.Peer) {
				assert := assert.New(t)
				assert.Len(peers, 1)
				assert.Equal("seed-1", peers[0].Hostname)
				assert.Equal(pkgtypes.HostTypeSuperSeedName, peers[0].Type)
			},
		},
		{
			name: "seed peer is not matched by a normal host id and is replaced",
			peers: []*models.Peer{
				{Hostname: "seed-1", IP: "10.0.0.1", Type: pkgtypes.HostTypeSuperSeedName, State: models.PeerStateInactive, SchedulerClusterID: 1},
			},
			results: []*resource.Host{
				{ID: idgen.HostID("10.0.0.1", "seed-1", false), Hostname: "seed-1", IP: "10.0.0.1", Type: pkgtypes.HostTypeNormal, SchedulerClusterID: 1},
			},
			expect: func(t *testing.T, peers []models.Peer) {
				assert := assert.New(t)
				assert.Len(peers, 1)
				assert.Equal("seed-1", peers[0].Hostname)
				assert.Equal(pkgtypes.HostTypeNormalName, peers[0].Type)
				assert.Equal(models.PeerStateActive, peers[0].State)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := mockDB(t, tc.peers)
			s := &syncPeers{
				config: &config.Config{Job: config.JobConfig{SyncPeers: config.SyncPeersConfig{BatchSize: 2}}},
				db:     db,
				mu:     &sync.Mutex{},
			}

			scheduler := models.Scheduler{Hostname: "scheduler", IP: "127.0.0.1", SchedulerClusterID: 1}
			s.mergePeers(context.Background(), scheduler, tc.results, logger.WithScheduler(scheduler.Hostname, scheduler.IP, uint64(scheduler.SchedulerClusterID)))

			var peers []models.Peer
			if err := db.Order("id").Find(&peers).Error; err != nil {
				t.Fatal(err)
			}

			tc.expect(t, peers)
		})
	}
}
