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

package gc

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"d7y.io/dragonfly/v2/manager/models"
	pkggc "d7y.io/dragonfly/v2/pkg/gc"
)

func mockDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "gc.db")), &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormlogger.Discard,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&models.Config{}, &models.Audit{}, &models.Job{}, &models.Scheduler{}, &models.SeedPeer{}, &models.SchedulerCluster{}, &models.SeedPeerCluster{}); err != nil {
		t.Fatal(err)
	}

	return db
}

func mockGCConfig(t *testing.T, db *gorm.DB, gcConfig *models.GCConfig) {
	value, err := json.Marshal(gcConfig)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Create(&models.Config{Name: models.ConfigGC, Value: string(value)}).Error; err != nil {
		t.Fatal(err)
	}
}

func TestAudit_RunGC(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		ctx    context.Context
		seed   func(t *testing.T, db *gorm.DB)
		expect func(t *testing.T, db *gorm.DB, err error)
	}{
		{
			name: "gc config is missing",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Empty(jobs)
			},
		},
		{
			name: "gc config is not valid json",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {
				if err := db.Create(&models.Config{Name: models.ConfigGC, Value: "{"}).Error; err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				var syntaxErr *json.SyntaxError
				assert.ErrorAs(err, &syntaxErr)
			},
		},
		{
			name: "audits older than the ttl are deleted and newer ones kept",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {
				mockGCConfig(t, db, &models.GCConfig{Audit: &models.GCAuditConfig{TTL: time.Hour}})
				if err := db.Create([]*models.Audit{
					{BaseModel: models.BaseModel{CreatedAt: now.Add(-2 * time.Hour)}, ActorType: models.ActorTypeUser, ActorName: "old-1", EventType: models.EventTypeAPI, Operation: "GET", State: models.AuditStateSuccess},
					{BaseModel: models.BaseModel{CreatedAt: now.Add(-90 * time.Minute)}, ActorType: models.ActorTypeUser, ActorName: "old-2", EventType: models.EventTypeAPI, Operation: "GET", State: models.AuditStateSuccess},
					{BaseModel: models.BaseModel{CreatedAt: now.Add(-30 * time.Minute)}, ActorType: models.ActorTypeUser, ActorName: "new-1", EventType: models.EventTypeAPI, Operation: "GET", State: models.AuditStateSuccess},
					{BaseModel: models.BaseModel{CreatedAt: now}, ActorType: models.ActorTypeUser, ActorName: "new-2", EventType: models.EventTypeAPI, Operation: "GET", State: models.AuditStateSuccess},
				}).Error; err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var audits []models.Audit
				assert.NoError(db.Order("id").Find(&audits).Error)
				assert.Len(audits, 2)
				assert.Equal("new-1", audits[0].ActorName)
				assert.Equal("new-2", audits[1].ActorName)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal(AuditGCTaskID, jobs[0].TaskID)
				assert.Equal(GCStateSuccess, jobs[0].State)
				assert.Equal(uint(0), jobs[0].UserID)
				assert.Equal(AuditGCTaskID, jobs[0].Args["type"])
				assert.EqualValues(2, jobs[0].Result["purged"])
			},
		},
		{
			name: "gc config without audit section falls back to the default ttl",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {
				mockGCConfig(t, db, &models.GCConfig{Job: &models.GCJobConfig{TTL: time.Hour}})
				if err := db.Create([]*models.Audit{
					{BaseModel: models.BaseModel{CreatedAt: now.Add(-models.DefaultGCAuditTTL - time.Hour)}, ActorType: models.ActorTypeUser, ActorName: "old", EventType: models.EventTypeAPI, Operation: "GET", State: models.AuditStateSuccess},
					{BaseModel: models.BaseModel{CreatedAt: now.Add(-models.DefaultGCAuditTTL + time.Hour)}, ActorType: models.ActorTypeUser, ActorName: "new", EventType: models.EventTypeAPI, Operation: "GET", State: models.AuditStateSuccess},
				}).Error; err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var audits []models.Audit
				assert.NoError(db.Find(&audits).Error)
				assert.Len(audits, 1)
				assert.Equal("new", audits[0].ActorName)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.EqualValues(1, jobs[0].Result["purged"])
			},
		},
		{
			name: "task id and user id from context are recorded",
			ctx:  context.WithValue(context.WithValue(context.Background(), pkggc.ContextKeyUserID, uint(7)), pkggc.ContextKeyTaskID, "task-audit"),
			seed: func(t *testing.T, db *gorm.DB) {
				mockGCConfig(t, db, &models.GCConfig{Audit: &models.GCAuditConfig{TTL: time.Hour}})
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal("task-audit", jobs[0].TaskID)
				assert.Equal(uint(7), jobs[0].UserID)
				assert.Equal(GCStateSuccess, jobs[0].State)
				assert.EqualValues(0, jobs[0].Result["purged"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := mockDB(t)
			tc.seed(t, db)

			err := NewAuditGCTask(db).Runner.RunGC(tc.ctx)
			tc.expect(t, db, err)
		})
	}
}

func TestJob_RunGC(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		ctx    context.Context
		seed   func(t *testing.T, db *gorm.DB)
		expect func(t *testing.T, db *gorm.DB, err error)
	}{
		{
			name: "gc config is missing",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Empty(jobs)
			},
		},
		{
			name: "jobs older than the ttl are deleted with their cluster associations",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {
				mockGCConfig(t, db, &models.GCConfig{Job: &models.GCJobConfig{TTL: time.Hour}})
				if err := db.Create([]*models.SchedulerCluster{{Name: "sc-1", Config: models.JSONMap{}, ClientConfig: models.JSONMap{}, SeedClientConfig: models.JSONMap{}}}).Error; err != nil {
					t.Fatal(err)
				}

				if err := db.Create([]*models.SeedPeerCluster{{Name: "spc-1", Config: models.JSONMap{}}}).Error; err != nil {
					t.Fatal(err)
				}

				var schedulerCluster models.SchedulerCluster
				if err := db.First(&schedulerCluster).Error; err != nil {
					t.Fatal(err)
				}

				var seedPeerCluster models.SeedPeerCluster
				if err := db.First(&seedPeerCluster).Error; err != nil {
					t.Fatal(err)
				}

				if err := db.Create([]*models.Job{
					{BaseModel: models.BaseModel{CreatedAt: now.Add(-2 * time.Hour)}, TaskID: "old-1", Type: "preheat", Args: models.JSONMap{}, SchedulerClusters: []models.SchedulerCluster{schedulerCluster}, SeedPeerClusters: []models.SeedPeerCluster{seedPeerCluster}},
					{BaseModel: models.BaseModel{CreatedAt: now.Add(-30 * time.Minute)}, TaskID: "new-1", Type: "preheat", Args: models.JSONMap{}, SchedulerClusters: []models.SchedulerCluster{schedulerCluster}, SeedPeerClusters: []models.SeedPeerCluster{seedPeerCluster}},
				}).Error; err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var jobs []models.Job
				assert.NoError(db.Where("type = ?", "preheat").Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal("new-1", jobs[0].TaskID)

				var schedulerClusterLinks, seedPeerClusterLinks int64
				assert.NoError(db.Table("job_scheduler_cluster").Count(&schedulerClusterLinks).Error)
				assert.NoError(db.Table("job_seed_peer_cluster").Count(&seedPeerClusterLinks).Error)
				assert.Equal(int64(1), schedulerClusterLinks)
				assert.Equal(int64(1), seedPeerClusterLinks)

				var gcJobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&gcJobs).Error)
				assert.Len(gcJobs, 1)
				assert.Equal(JobGCTaskID, gcJobs[0].TaskID)
				assert.Equal(GCStateSuccess, gcJobs[0].State)
				assert.EqualValues(1, gcJobs[0].Result["purged"])
			},
		},
		{
			name: "gc config without job section falls back to the default ttl",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {
				mockGCConfig(t, db, &models.GCConfig{Audit: &models.GCAuditConfig{TTL: time.Hour}})
				if err := db.Create([]*models.Job{
					{BaseModel: models.BaseModel{CreatedAt: now.Add(-models.DefaultGCJobTTL - time.Hour)}, TaskID: "old-1", Type: "preheat", Args: models.JSONMap{}},
					{BaseModel: models.BaseModel{CreatedAt: now.Add(-models.DefaultGCJobTTL + time.Hour)}, TaskID: "new-1", Type: "preheat", Args: models.JSONMap{}},
				}).Error; err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var jobs []models.Job
				assert.NoError(db.Where("type = ?", "preheat").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal("new-1", jobs[0].TaskID)

				var gcJobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&gcJobs).Error)
				assert.Len(gcJobs, 1)
				assert.EqualValues(1, gcJobs[0].Result["purged"])
			},
		},
		{
			name: "nothing to purge records a successful gc job",
			ctx:  context.WithValue(context.Background(), pkggc.ContextKeyTaskID, "task-job"),
			seed: func(t *testing.T, db *gorm.DB) {
				mockGCConfig(t, db, &models.GCConfig{Job: &models.GCJobConfig{TTL: time.Hour}})
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var gcJobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&gcJobs).Error)
				assert.Len(gcJobs, 1)
				assert.Equal("task-job", gcJobs[0].TaskID)
				assert.Equal(GCStateSuccess, gcJobs[0].State)
				assert.EqualValues(0, gcJobs[0].Result["purged"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := mockDB(t)
			tc.seed(t, db)

			err := NewJobGCTask(db).Runner.RunGC(tc.ctx)
			tc.expect(t, db, err)
		})
	}
}

func TestScheduler_RunGC(t *testing.T) {
	now := time.Now()
	keepAliveInterval := float64(time.Minute)

	tests := []struct {
		name   string
		ctx    context.Context
		seed   func(t *testing.T, db *gorm.DB)
		expect func(t *testing.T, db *gorm.DB, err error)
	}{
		{
			name: "stale inactive schedulers are swept and the rest kept",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {
				if err := db.Create([]*models.Scheduler{
					{BaseModel: models.BaseModel{UpdatedAt: now.Add(-time.Hour)}, Hostname: "stale-inactive", IP: "10.0.0.1", State: models.SchedulerStateInactive, SchedulerClusterID: 1},
					{BaseModel: models.BaseModel{UpdatedAt: now.Add(-time.Minute)}, Hostname: "fresh-inactive", IP: "10.0.0.2", State: models.SchedulerStateInactive, SchedulerClusterID: 1},
					{BaseModel: models.BaseModel{UpdatedAt: now.Add(-time.Hour)}, Hostname: "stale-active-without-config", IP: "10.0.0.3", State: models.SchedulerStateActive, LastKeepAliveAt: now.Add(-time.Hour), SchedulerClusterID: 1},
				}).Error; err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var schedulers []models.Scheduler
				assert.NoError(db.Order("id").Find(&schedulers).Error)
				assert.Len(schedulers, 2)
				assert.Equal("fresh-inactive", schedulers[0].Hostname)
				assert.Equal(models.SchedulerStateInactive, schedulers[0].State)
				assert.Equal("stale-active-without-config", schedulers[1].Hostname)
				assert.Equal(models.SchedulerStateActive, schedulers[1].State)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal(SchedulerGCTaskID, jobs[0].TaskID)
				assert.Equal(GCStateSuccess, jobs[0].State)
				assert.EqualValues(1, jobs[0].Result["purged"])
			},
		},
		{
			name: "active scheduler without keep alive for 3x interval is marked inactive",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {
				if err := db.Create([]*models.Scheduler{
					{Hostname: "dead", IP: "10.0.0.1", State: models.SchedulerStateActive, Config: models.JSONMap{"manager_keep_alive_interval": keepAliveInterval}, LastKeepAliveAt: now.Add(-4 * time.Minute), SchedulerClusterID: 1},
					{Hostname: "alive", IP: "10.0.0.2", State: models.SchedulerStateActive, Config: models.JSONMap{"manager_keep_alive_interval": keepAliveInterval}, LastKeepAliveAt: now.Add(-2 * time.Minute), SchedulerClusterID: 1},
					{Hostname: "no-interval", IP: "10.0.0.3", State: models.SchedulerStateActive, Config: models.JSONMap{"foo": "bar"}, LastKeepAliveAt: now.Add(-time.Hour), SchedulerClusterID: 1},
					{Hostname: "interval-not-a-number", IP: "10.0.0.4", State: models.SchedulerStateActive, Config: models.JSONMap{"manager_keep_alive_interval": "1m"}, LastKeepAliveAt: now.Add(-time.Hour), SchedulerClusterID: 1},
				}).Error; err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var schedulers []models.Scheduler
				assert.NoError(db.Order("id").Find(&schedulers).Error)
				assert.Len(schedulers, 4)
				assert.Equal(models.SchedulerStateInactive, schedulers[0].State)
				assert.Equal(models.SchedulerStateActive, schedulers[1].State)
				assert.Equal(models.SchedulerStateActive, schedulers[2].State)
				assert.Equal(models.SchedulerStateActive, schedulers[3].State)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.EqualValues(0, jobs[0].Result["purged"])
			},
		},
		{
			name: "task id and user id from context are recorded",
			ctx:  context.WithValue(context.WithValue(context.Background(), pkggc.ContextKeyUserID, uint(3)), pkggc.ContextKeyTaskID, "task-scheduler"),
			seed: func(t *testing.T, db *gorm.DB) {},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal("task-scheduler", jobs[0].TaskID)
				assert.Equal(uint(3), jobs[0].UserID)
				assert.Equal(SchedulerGCTaskID, jobs[0].Args["type"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := mockDB(t)
			tc.seed(t, db)

			err := NewSchedulerGCTask(db).Runner.RunGC(tc.ctx)
			tc.expect(t, db, err)
		})
	}
}

func TestSeedPeer_RunGC(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		ctx    context.Context
		seed   func(t *testing.T, db *gorm.DB)
		expect func(t *testing.T, db *gorm.DB, err error)
	}{
		{
			name: "stale inactive seed peers are deleted and the rest kept",
			ctx:  context.Background(),
			seed: func(t *testing.T, db *gorm.DB) {
				if err := db.Create([]*models.SeedPeer{
					{BaseModel: models.BaseModel{UpdatedAt: now.Add(-time.Hour)}, Hostname: "stale-inactive", IP: "10.0.0.1", State: models.SeedPeerStateInactive, SeedPeerClusterID: 1},
					{BaseModel: models.BaseModel{UpdatedAt: now.Add(-time.Minute)}, Hostname: "fresh-inactive", IP: "10.0.0.2", State: models.SeedPeerStateInactive, SeedPeerClusterID: 1},
					{BaseModel: models.BaseModel{UpdatedAt: now.Add(-time.Hour)}, Hostname: "stale-active", IP: "10.0.0.3", State: models.SeedPeerStateActive, SeedPeerClusterID: 1},
				}).Error; err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var seedPeers []models.SeedPeer
				assert.NoError(db.Order("id").Find(&seedPeers).Error)
				assert.Len(seedPeers, 2)
				assert.Equal("fresh-inactive", seedPeers[0].Hostname)
				assert.Equal("stale-active", seedPeers[1].Hostname)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal(SeedPeerGCTaskID, jobs[0].TaskID)
				assert.Equal(GCStateSuccess, jobs[0].State)
				assert.Equal(uint(0), jobs[0].UserID)
				assert.EqualValues(1, jobs[0].Result["purged"])
			},
		},
		{
			name: "task id and user id from context are recorded",
			ctx:  context.WithValue(context.WithValue(context.Background(), pkggc.ContextKeyUserID, uint(5)), pkggc.ContextKeyTaskID, "task-seed-peer"),
			seed: func(t *testing.T, db *gorm.DB) {},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal("task-seed-peer", jobs[0].TaskID)
				assert.Equal(uint(5), jobs[0].UserID)
				assert.Equal(SeedPeerGCTaskID, jobs[0].Args["type"])
				assert.EqualValues(0, jobs[0].Result["purged"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := mockDB(t)
			tc.seed(t, db)

			err := NewSeedPeerGCTask(db).Runner.RunGC(tc.ctx)
			tc.expect(t, db, err)
		})
	}
}

func TestJobRecorder_Record(t *testing.T) {
	tests := []struct {
		name   string
		init   bool
		result Result
		expect func(t *testing.T, db *gorm.DB, err error)
	}{
		{
			name:   "record before init",
			init:   false,
			result: Result{Purged: 1},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.Error(err)
				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Empty(jobs)
			},
		},
		{
			name:   "successful result",
			init:   true,
			result: Result{Purged: 3},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal(GCStateSuccess, jobs[0].State)
				assert.EqualValues(3, jobs[0].Result["purged"])
				assert.NotContains(jobs[0].Result, "error")
			},
		},
		{
			name:   "failed result",
			init:   true,
			result: Result{Purged: 2, Error: errors.New("boom")},
			expect: func(t *testing.T, db *gorm.DB, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var jobs []models.Job
				assert.NoError(db.Session(&gorm.Session{SkipHooks: true}).Where("type = ?", GCJobType).Order("id").Find(&jobs).Error)
				assert.Len(jobs, 1)
				assert.Equal(GCStateFailure, jobs[0].State)
				assert.EqualValues(2, jobs[0].Result["purged"])
				assert.Equal("boom", jobs[0].Result["error"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := mockDB(t)
			recorder := newJobRecorder(db)
			if tc.init {
				if err := recorder.Init(1, "task", models.JSONMap{"type": "test"}); err != nil {
					t.Fatal(err)
				}
			}

			err := recorder.Record(tc.result)
			tc.expect(t, db, err)
		})
	}
}
