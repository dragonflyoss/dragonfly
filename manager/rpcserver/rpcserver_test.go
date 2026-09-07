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

package rpcserver

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	cachev9 "github.com/go-redis/cache/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"

	"d7y.io/dragonfly/v2/manager/cache"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/searcher"
	"d7y.io/dragonfly/v2/manager/types"
)

const (
	mockSchedulerClusterID       uint = 1
	mockSecondSchedulerClusterID uint = 2
	mockEmptySchedulerClusterID  uint = 3
	mockSeedPeerClusterID        uint = 1
	mockSecondSeedPeerClusterID  uint = 2
	mockUnknownClusterID         uint = 99

	mockSchedulerClusterName       = "scheduler-cluster-1"
	mockSecondSchedulerClusterName = "scheduler-cluster-2"
	mockEmptySchedulerClusterName  = "scheduler-cluster-3"
	mockSeedPeerClusterName        = "seed-peer-cluster-1"
	mockSecondSeedPeerClusterName  = "seed-peer-cluster-2"

	mockSchedulerClusterConfigJSON           = `{"candidate_parent_limit":4}`
	mockSchedulerClusterClientConfigJSON     = `{"load_limit":100}`
	mockSchedulerClusterSeedClientConfigJSON = `{"load_limit":200}`
	mockSchedulerClusterScopesJSON           = `{"idc":"idc-1"}`
	mockSeedPeerClusterConfigJSON            = `{"load_limit":300}`
	mockDefaultFeaturesJSON                  = `["schedule","preheat"]`

	mockActiveSchedulerHostname              = "scheduler-1"
	mockActiveSchedulerIP                    = "10.0.0.1"
	mockInactiveSchedulerHostname            = "scheduler-2"
	mockInactiveSchedulerIP                  = "10.0.0.2"
	mockPreheatOnlySchedulerHostname         = "scheduler-3"
	mockPreheatOnlySchedulerIP               = "10.0.0.3"
	mockSecondClusterSchedulerHostname       = "scheduler-4"
	mockSecondClusterSchedulerIP             = "10.0.0.4"
	mockSchedulerPort                  int32 = 8002
	mockUpdatedSchedulerPort           int32 = 8003

	mockActiveSeedPeerHostname         = "seed-peer-1"
	mockActiveSeedPeerIP               = "10.0.0.11"
	mockInactiveSeedPeerHostname       = "seed-peer-2"
	mockInactiveSeedPeerIP             = "10.0.0.12"
	mockLonelySeedPeerHostname         = "seed-peer-3"
	mockLonelySeedPeerIP               = "10.0.0.13"
	mockSeedPeerType                   = "super"
	mockSeedPeerPort             int32 = 65002
	mockSeedPeerDownloadPort     int32 = 65001
	mockUpdatedSeedPeerPort      int32 = 65003

	mockNewHostname     = "new-host"
	mockNewIP           = "10.0.0.20"
	mockCachedHostname  = "cached"
	mockPeerHostname    = "dfdaemon"
	mockPeerIP          = "10.0.1.1"
	mockUnknownHostname = "unknown"
	mockUnknownIP       = "10.0.0.99"
	mockIDC             = "idc-1"
	mockLocation        = "location-1"
	mockUpdatedIDC      = "idc-2"
	mockUpdatedLocation = "location-2"
	mockVersion         = "v2.0.0"
	mockCommit          = "abc123"
)

func newTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_foreign_keys=1", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         gormlogger.Discard,
	})
	if err != nil {
		t.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}

	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() {
		assert := assert.New(t)
		assert.NoError(sqlDB.Close())
	})

	if err := db.AutoMigrate(
		&models.SchedulerCluster{},
		&models.Scheduler{},
		&models.SeedPeerCluster{},
		&models.SeedPeer{},
		&models.Application{},
	); err != nil {
		t.Fatal(err)
	}

	seedTestDB(t, db)
	return db
}

func seedTestDB(t *testing.T, db *gorm.DB) {
	schedulerClusters := []models.SchedulerCluster{
		{
			Name:             mockSchedulerClusterName,
			Config:           models.JSONMap{"candidate_parent_limit": 4},
			ClientConfig:     models.JSONMap{"load_limit": 100},
			SeedClientConfig: models.JSONMap{"load_limit": 200},
			Scopes:           models.JSONMap{"idc": mockIDC},
			IsDefault:        true,
			Schedulers: []models.Scheduler{
				{
					Hostname: mockActiveSchedulerHostname,
					IP:       mockActiveSchedulerIP,
					Port:     mockSchedulerPort,
					IDC:      mockIDC,
					Location: mockLocation,
					State:    models.SchedulerStateActive,
					Features: types.DefaultSchedulerFeatures,
				},
				{
					Hostname: mockInactiveSchedulerHostname,
					IP:       mockInactiveSchedulerIP,
					Port:     mockSchedulerPort,
					State:    models.SchedulerStateInactive,
					Features: types.DefaultSchedulerFeatures,
				},
			},
			SeedPeerClusters: []models.SeedPeerCluster{
				{
					Name:   mockSeedPeerClusterName,
					Config: models.JSONMap{"load_limit": 300},
					SeedPeers: []models.SeedPeer{
						{
							Hostname:     mockActiveSeedPeerHostname,
							IP:           mockActiveSeedPeerIP,
							Type:         mockSeedPeerType,
							Port:         mockSeedPeerPort,
							DownloadPort: mockSeedPeerDownloadPort,
							State:        models.SeedPeerStateActive,
						},
						{
							Hostname:     mockInactiveSeedPeerHostname,
							IP:           mockInactiveSeedPeerIP,
							Type:         mockSeedPeerType,
							Port:         mockSeedPeerPort,
							DownloadPort: mockSeedPeerDownloadPort,
							State:        models.SeedPeerStateInactive,
						},
					},
				},
			},
		},
		{
			Name: mockSecondSchedulerClusterName,
			Schedulers: []models.Scheduler{
				{
					Hostname: mockPreheatOnlySchedulerHostname,
					IP:       mockPreheatOnlySchedulerIP,
					Port:     mockSchedulerPort,
					State:    models.SchedulerStateActive,
					Features: models.Array{types.SchedulerFeaturePreheat},
				},
				{
					Hostname: mockSecondClusterSchedulerHostname,
					IP:       mockSecondClusterSchedulerIP,
					Port:     mockSchedulerPort,
					State:    models.SchedulerStateActive,
					Features: models.Array{types.SchedulerFeatureSchedule},
				},
			},
		},
		{
			Name: mockEmptySchedulerClusterName,
		},
	}

	for i := range schedulerClusters {
		if err := db.Create(&schedulerClusters[i]).Error; err != nil {
			t.Fatal(err)
		}
	}

	if err := db.Create(&models.SeedPeerCluster{
		Name: mockSecondSeedPeerClusterName,
		SeedPeers: []models.SeedPeer{
			{
				Hostname:     mockLonelySeedPeerHostname,
				IP:           mockLonelySeedPeerIP,
				Type:         mockSeedPeerType,
				Port:         mockSeedPeerPort,
				DownloadPort: mockSeedPeerDownloadPort,
				State:        models.SeedPeerStateInactive,
			},
		},
	}).Error; err != nil {
		t.Fatal(err)
	}
}

func newTestCache() *cache.Cache {
	return &cache.Cache{
		Cache: cachev9.New(&cachev9.Options{LocalCache: cachev9.NewTinyLFU(1000, time.Minute)}),
		TTL:   time.Minute,
	}
}

func newTestServerV1(t *testing.T, sr searcher.Searcher) *managerServerV1 {
	return &managerServerV1{db: newTestDB(t), cache: newTestCache(), searcher: sr}
}

func newTestServerV2(t *testing.T, sr searcher.Searcher) *managerServerV2 {
	return &managerServerV2{db: newTestDB(t), cache: newTestCache(), searcher: sr}
}

func setCache(t *testing.T, c *cache.Cache, key string, value any) {
	assert := assert.New(t)
	assert.NoError(c.Set(&cachev9.Item{Key: key, Value: value}))
}

func countRows(t *testing.T, db *gorm.DB, model any, query any) int64 {
	assert := assert.New(t)
	var n int64
	assert.NoError(db.Unscoped().Model(model).Where(query).Count(&n).Error)
	return n
}

func countAllRows(t *testing.T, db *gorm.DB, model any) int64 {
	assert := assert.New(t)
	var n int64
	assert.NoError(db.Unscoped().Model(model).Count(&n).Error)
	return n
}

func findScheduler(t *testing.T, db *gorm.DB, hostname string) models.Scheduler {
	assert := assert.New(t)
	var scheduler models.Scheduler
	assert.NoError(db.First(&scheduler, models.Scheduler{Hostname: hostname}).Error)
	return scheduler
}

func findSeedPeer(t *testing.T, db *gorm.DB, hostname string) models.SeedPeer {
	assert := assert.New(t)
	var seedPeer models.SeedPeer
	assert.NoError(db.First(&seedPeer, models.SeedPeer{Hostname: hostname}).Error)
	return seedPeer
}

func hostnames[T interface{ GetHostname() string }](items []T) []string {
	names := []string{}
	for _, item := range items {
		names = append(names, item.GetHostname())
	}

	return names
}

func schedulerHostnamesByCluster(clusters []models.SchedulerCluster) map[string][]string {
	names := map[string][]string{}
	for _, cluster := range clusters {
		for _, scheduler := range cluster.Schedulers {
			names[cluster.Name] = append(names[cluster.Name], scheduler.Hostname)
		}
	}

	return names
}

func pickSchedulerCluster(name string) func(context.Context, []models.SchedulerCluster, string, string, map[string]string, *zap.SugaredLogger) ([]models.SchedulerCluster, error) {
	return func(_ context.Context, clusters []models.SchedulerCluster, _, _ string, _ map[string]string, _ *zap.SugaredLogger) ([]models.SchedulerCluster, error) {
		for _, cluster := range clusters {
			if cluster.Name == name {
				return []models.SchedulerCluster{cluster}, nil
			}
		}

		return nil, fmt.Errorf("scheduler cluster %s not found", name)
	}
}

func sourceState(db *gorm.DB, scheduler bool, hostname, ip string, clusterID uint) string {
	if scheduler {
		var scheduler models.Scheduler
		if err := db.First(&scheduler, models.Scheduler{Hostname: hostname, IP: ip, SchedulerClusterID: clusterID}).Error; err != nil {
			return err.Error()
		}

		return scheduler.State
	}

	var seedPeer models.SeedPeer
	if err := db.First(&seedPeer, models.SeedPeer{Hostname: hostname, IP: ip, SeedPeerClusterID: clusterID}).Error; err != nil {
		return err.Error()
	}

	return seedPeer.State
}

func ptr[T any](v T) *T {
	return &v
}

type mockKeepAliveStream struct {
	grpc.ServerStream
	reqs    []*managerv2.KeepAliveRequest
	err     error
	observe func() string
	states  []string
}

func (m *mockKeepAliveStream) Recv() (*managerv2.KeepAliveRequest, error) {
	if m.observe != nil {
		m.states = append(m.states, m.observe())
	}

	if len(m.reqs) == 0 {
		return nil, m.err
	}

	req := m.reqs[0]
	m.reqs = m.reqs[1:]
	return req, nil
}

func (m *mockKeepAliveStream) SendAndClose(*emptypb.Empty) error {
	return nil
}

func (m *mockKeepAliveStream) Context() context.Context {
	return context.Background()
}
