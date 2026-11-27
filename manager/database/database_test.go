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

package database

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

func mockDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "manager.db")), &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormlogger.Discard,
	})
	require.NoError(t, err)
	require.NoError(t, migrate(db))
	return db
}

func TestSeed_PersonalAccessToken(t *testing.T) {
	const pat = "dragonfly-pat-from-env"

	tests := []struct {
		name   string
		mock   func(t *testing.T, db *gorm.DB)
		expect func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "creates the token on first seed and keeps it on later seeds",
			mock: func(t *testing.T, db *gorm.DB) {},
			expect: func(t *testing.T, db *gorm.DB) {
				var tokens []models.PersonalAccessToken
				require.NoError(t, db.Unscoped().Find(&tokens).Error)
				require.Len(t, tokens, 1)
				assert.Equal(t, "default", tokens[0].Name)
				assert.Equal(t, pat, tokens[0].Token)
				assert.Equal(t, models.PersonalAccessTokenStateActive, tokens[0].State)
				assert.Equal(t, uint(1), tokens[0].UserID)
				assert.True(t, tokens[0].ExpiredAt.After(time.Now().AddDate(9, 0, 0)))
			},
		},
		{
			name: "leaves an inactive token with the same value untouched",
			mock: func(t *testing.T, db *gorm.DB) {
				require.NoError(t, db.Create(&models.PersonalAccessToken{
					Name:      "revoked",
					Token:     pat,
					Scopes:    types.DefaultPersonalAccessTokenScopes,
					State:     models.PersonalAccessTokenStateInactive,
					ExpiredAt: time.Now().Add(-time.Hour),
					UserID:    1,
				}).Error)
			},
			expect: func(t *testing.T, db *gorm.DB) {
				var tokens []models.PersonalAccessToken
				require.NoError(t, db.Unscoped().Find(&tokens).Error)
				require.Len(t, tokens, 1)
				assert.Equal(t, "revoked", tokens[0].Name)
				assert.Equal(t, models.PersonalAccessTokenStateInactive, tokens[0].State)
				assert.True(t, tokens[0].ExpiredAt.Before(time.Now()))
			},
		},
		{
			name: "does not recreate a soft-deleted token with the same value",
			mock: func(t *testing.T, db *gorm.DB) {
				token := models.PersonalAccessToken{
					Name:      "default",
					Token:     pat,
					Scopes:    types.DefaultPersonalAccessTokenScopes,
					State:     models.PersonalAccessTokenStateActive,
					ExpiredAt: time.Now().AddDate(10, 0, 0),
					UserID:    1,
				}
				require.NoError(t, db.Create(&token).Error)
				require.NoError(t, db.Delete(&token).Error)
			},
			expect: func(t *testing.T, db *gorm.DB) {
				var count int64
				require.NoError(t, db.Model(&models.PersonalAccessToken{}).Count(&count).Error)
				assert.Zero(t, count)
				require.NoError(t, db.Unscoped().Model(&models.PersonalAccessToken{}).Count(&count).Error)
				assert.Equal(t, int64(1), count)
			},
		},
		{
			name: "keeps an existing default token with a different value",
			mock: func(t *testing.T, db *gorm.DB) {
				require.NoError(t, db.Create(&models.PersonalAccessToken{
					Name:      "default",
					Token:     "previous-pat",
					Scopes:    types.DefaultPersonalAccessTokenScopes,
					State:     models.PersonalAccessTokenStateActive,
					ExpiredAt: time.Now().AddDate(10, 0, 0),
					UserID:    1,
				}).Error)
			},
			expect: func(t *testing.T, db *gorm.DB) {
				var tokens []models.PersonalAccessToken
				require.NoError(t, db.Unscoped().Find(&tokens).Error)
				require.Len(t, tokens, 1)
				assert.Equal(t, "default", tokens[0].Name)
				assert.Equal(t, "previous-pat", tokens[0].Token)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(DragonflyPATEnvName, pat)
			db := mockDB(t)
			tc.mock(t, db)

			// Seed runs on every manager start, so it must succeed on restart.
			require.NoError(t, seed(db))
			require.NoError(t, seed(db))
			tc.expect(t, db)
		})
	}
}
