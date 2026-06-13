package database

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSyncSequenceAfterExplicitInsertPostgres(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectExec(regexp.QuoteMeta("SELECT setval(pg_get_serial_sequence($1, 'id'), $2, true)")).
		WithArgs("scheduler_cluster", uint(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	require.NoError(t, err)

	require.NoError(t, syncSequenceAfterExplicitInsert(db, "scheduler_cluster", 1))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncSequenceAfterExplicitInsertNonPostgres(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	require.NoError(t, err)

	require.NoError(t, syncSequenceAfterExplicitInsert(db, "scheduler_cluster", 1))
}
