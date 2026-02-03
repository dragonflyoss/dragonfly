/*
 *     Copyright 2020 The Dragonfly Authors
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
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/config"

	_ "github.com/polardb/polardbx-connector-go"
)

func newPolardb(cfg *config.Config) (*gorm.DB, error) {

	dsn, _ := formatPolardbDSN(&cfg.Database.Polardb)

	// Initialize gorm logger.
	logLevel := gormlogger.Info
	if !cfg.Verbose {
		logLevel = gormlogger.Warn
	}
	gormLogger := zapgorm2.New(logger.CoreLogger.Desugar()).LogMode(logLevel)

	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DriverName: "polardbx",
			DSN:        dsn}),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   gormLogger,
		})

	if err != nil {
		return nil, err
	}

	// Run migration.
	if err := migrate(db); err != nil {
		fmt.Println("migrate err")
		return nil, err
	}

	// Run migration.
	if cfg.Database.Polardb.Migrate {
		if err := migrate(db); err != nil {
			return nil, err
		}
	}

	// Run seed.
	if err := seed(db); err != nil {
		return nil, err
	}

	return db, nil
}

// formatPolardbDSN formats PolarDB DSN.
func formatPolardbDSN(cfg *config.PolardbConfig) (string, error) {
	// dsn := "root:123456@tcp(local:8527)/manager?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.AddrList + ")/" + cfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	return dsn, nil
}
