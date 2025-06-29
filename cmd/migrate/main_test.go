package main

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"interview/internal/config"
	"interview/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectDB(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.DatabaseConfig
		verbose  bool
		wantErr  bool
	}{
		{
			name: "Valid config verbose true",
			cfg: config.DatabaseConfig{
				Host:     "localhost",
				Port:     3306,
				User:     "root",
				Password: "password",
				Name:     "test_db",
			},
			verbose: true,
			wantErr: true, // Will fail because we don't have MySQL running in test
		},
		{
			name: "Valid config verbose false",
			cfg: config.DatabaseConfig{
				Host:     "localhost",
				Port:     3306,
				User:     "root",
				Password: "password",
				Name:     "test_db",
			},
			verbose: false,
			wantErr: true, // Will fail because we don't have MySQL running in test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := connectDB(tt.cfg, tt.verbose)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	return db
}

func TestMigrateUp(t *testing.T) {
	db := setupTestDB(t)

	err := migrateUp(db)
	assert.NoError(t, err)

	// Check if table was created
	assert.True(t, db.Migrator().HasTable(&models.Transaction{}))
}

func TestMigrateDown(t *testing.T) {
	db := setupTestDB(t)

	// First create the table
	err := migrateUp(db)
	require.NoError(t, err)
	assert.True(t, db.Migrator().HasTable(&models.Transaction{}))

	// Then drop it
	err = migrateDown(db)
	assert.NoError(t, err)
	assert.False(t, db.Migrator().HasTable(&models.Transaction{}))
}

func TestMigrateReset(t *testing.T) {
	db := setupTestDB(t)

	// First create the table
	err := migrateUp(db)
	require.NoError(t, err)
	assert.True(t, db.Migrator().HasTable(&models.Transaction{}))

	// Reset should drop and recreate
	err = migrateReset(db)
	assert.NoError(t, err)
	assert.True(t, db.Migrator().HasTable(&models.Transaction{}))
}

func TestCreateIndexes(t *testing.T) {
	db := setupTestDB(t)

	// First create the table
	err := db.AutoMigrate(&models.Transaction{})
	require.NoError(t, err)

	err = createIndexes(db)
	assert.NoError(t, err)
}

func TestMigrateStatus(t *testing.T) {
	db := setupTestDB(t)

	// Test with no tables
	err := migrateStatus(db)
	assert.NoError(t, err)

	// Test with tables
	err = migrateUp(db)
	require.NoError(t, err)

	err = migrateStatus(db)
	assert.NoError(t, err)
}

func TestMigrateDown_NoTable(t *testing.T) {
	db := setupTestDB(t)

	// Try to drop table that doesn't exist
	err := migrateDown(db)
	// SQLite doesn't error on DROP TABLE IF NOT EXISTS behavior
	// This test ensures the function handles missing tables gracefully
	assert.NoError(t, err)
}

func TestMigrateUp_MultipleRuns(t *testing.T) {
	db := setupTestDB(t)

	// Run migration multiple times
	err := migrateUp(db)
	assert.NoError(t, err)

	err = migrateUp(db)
	assert.NoError(t, err) // Should not error on second run

	assert.True(t, db.Migrator().HasTable(&models.Transaction{}))
}

func TestCreateIndexes_AlreadyExists(t *testing.T) {
	db := setupTestDB(t)

	// Create table first
	err := db.AutoMigrate(&models.Transaction{})
	require.NoError(t, err)

	// Run createIndexes multiple times
	err = createIndexes(db)
	assert.NoError(t, err)

	err = createIndexes(db)
	assert.NoError(t, err) // Should not error if index already exists
}
