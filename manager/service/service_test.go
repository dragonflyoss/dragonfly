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
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/permission/rbac"
)

func mockDB(t *testing.T) *gorm.DB {
	assert := assert.New(t)
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.NewString())), &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormlogger.Discard,
	})
	assert.NoError(err)
	assert.NoError(db.AutoMigrate(
		&models.Audit{},
		&models.Job{},
		&models.SeedPeerCluster{},
		&models.SeedPeer{},
		&models.SchedulerCluster{},
		&models.Scheduler{},
		&models.User{},
		&models.Oauth{},
		&models.Config{},
		&models.Application{},
		&models.PersonalAccessToken{},
		&models.Peer{},
	))

	return db
}

func mockService(t *testing.T) *service {
	assert := assert.New(t)
	db := mockDB(t)
	enforcer, err := rbac.NewEnforcer(db)
	assert.NoError(err)

	return &service{db: db, enforcer: enforcer}
}

func mockSchedulerCluster(t *testing.T, db *gorm.DB, name string) models.SchedulerCluster {
	assert := assert.New(t)
	schedulerCluster := models.SchedulerCluster{
		Name:             name,
		Config:           models.JSONMap{"candidate_parent_limit": float64(4)},
		ClientConfig:     models.JSONMap{"load_limit": float64(50)},
		SeedClientConfig: models.JSONMap{},
		Scopes:           models.JSONMap{},
	}
	assert.NoError(db.Create(&schedulerCluster).Error)

	return schedulerCluster
}

func mockSeedPeerCluster(t *testing.T, db *gorm.DB, name string) models.SeedPeerCluster {
	assert := assert.New(t)
	seedPeerCluster := models.SeedPeerCluster{
		Name:   name,
		Config: models.JSONMap{"load_limit": float64(300)},
	}
	assert.NoError(db.Create(&seedPeerCluster).Error)

	return seedPeerCluster
}

func mockScheduler(t *testing.T, db *gorm.DB, schedulerClusterID uint, hostname, state string, features []string) models.Scheduler {
	assert := assert.New(t)
	scheduler := models.Scheduler{
		Hostname:           hostname,
		IDC:                "idc-" + hostname,
		IP:                 "127.0.0.1",
		Port:               8002,
		State:              state,
		Features:           features,
		SchedulerClusterID: schedulerClusterID,
	}
	assert.NoError(db.Create(&scheduler).Error)

	return scheduler
}

func mockSeedPeer(t *testing.T, db *gorm.DB, seedPeerClusterID uint, hostname string) models.SeedPeer {
	assert := assert.New(t)
	seedPeer := models.SeedPeer{
		Hostname:          hostname,
		Type:              "super",
		IP:                "127.0.0.1",
		Port:              4000,
		DownloadPort:      4001,
		SeedPeerClusterID: seedPeerClusterID,
	}
	assert.NoError(db.Create(&seedPeer).Error)

	return seedPeer
}

func mockUser(t *testing.T, db *gorm.DB, name string) models.User {
	assert := assert.New(t)
	user := models.User{
		Name:  name,
		Email: name + "@example.com",
		State: models.UserStateEnabled,
	}
	assert.NoError(db.Create(&user).Error)

	return user
}
