package integration

import (
	"context"
	"testing"

	"ark/internal/config"
	"ark/internal/database"
	"ark/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMigration_CreatesAllTables verifies that running migrations creates
// all expected tables with correct structure.
func TestMigration_CreatesAllTables(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	cfg, log := setupTest(t)

	// Run migrations
	err := database.Migrate(ctx, &log, cfg)
	require.NoError(t, err, "migration should succeed")

	// Connect to database to verify
	conn := connectDB(t, cfg)
	defer conn.Close(ctx)

	// Verify assets table exists with correct columns
	var exists bool
	err = conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'assets'
		)
	`).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "assets table should exist")

	// Verify assets table has correct columns
	var columnCount int
	err = conn.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = 'public' 
		AND table_name = 'assets'
	`).Scan(&columnCount)
	require.NoError(t, err)
	assert.Equal(t, 8, columnCount, "assets table should have 8 columns")

	// Verify asset_logs table exists
	err = conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'asset_logs'
		)
	`).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "asset_logs table should exist")

	// Verify schema_version table shows version 1
	var version int32
	err = conn.QueryRow(ctx, "SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&version)
	require.NoError(t, err)
	assert.Equal(t, int32(1), version, "migration version should be 1")
}

// TestMigration_CreatesAllIndexes verifies that all expected indexes are created.
func TestMigration_CreatesAllIndexes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	cfg, log := setupTest(t)

	// Run migrations
	err := database.Migrate(ctx, &log, cfg)
	require.NoError(t, err, "migration should succeed")

	// Connect to database to verify
	conn := connectDB(t, cfg)
	defer conn.Close(ctx)

	// Expected indexes on assets table
	assetsIndexes := []string{
		"assets_pkey",
		"idx_assets_user_id",
		"idx_assets_name_trgm",
		"idx_assets_type",
	}

	for _, indexName := range assetsIndexes {
		var exists bool
		err = conn.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_indexes 
				WHERE schemaname = 'public' 
				AND tablename = 'assets'
				AND indexname = $1
			)
		`, indexName).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "index %s should exist on assets table", indexName)
	}

	// Expected indexes on asset_logs table
	assetLogsIndexes := []string{
		"asset_logs_pkey",
		"idx_asset_logs_user_id",
		"idx_asset_logs_asset_id",
		"idx_asset_logs_created_at",
		"idx_asset_logs_content_vector",
		"idx_asset_logs_tags",
	}

	for _, indexName := range assetLogsIndexes {
		var exists bool
		err = conn.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_indexes 
				WHERE schemaname = 'public' 
				AND tablename = 'asset_logs'
				AND indexname = $1
			)
		`, indexName).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "index %s should exist on asset_logs table", indexName)
	}
}

// TestMigration_Idempotent verifies that running migrations twice doesn't error.
func TestMigration_Idempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	cfg, log := setupTest(t)

	// Run migrations first time
	err := database.Migrate(ctx, &log, cfg)
	require.NoError(t, err, "first migration should succeed")

	// Run migrations second time (should be idempotent)
	err = database.Migrate(ctx, &log, cfg)
	require.NoError(t, err, "second migration should succeed (idempotent)")

	// Verify version is still 1
	conn := connectDB(t, cfg)
	defer conn.Close(ctx)

	var version int32
	err = conn.QueryRow(ctx, "SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&version)
	require.NoError(t, err)
	assert.Equal(t, int32(1), version, "migration version should still be 1")
}

// TestMigration_CreatesForeignKeys verifies that foreign key constraints are created.
func TestMigration_CreatesForeignKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	cfg, log := setupTest(t)

	// Run migrations
	err := database.Migrate(ctx, &log, cfg)
	require.NoError(t, err, "migration should succeed")

	// Connect to database to verify
	conn := connectDB(t, cfg)
	defer conn.Close(ctx)

	// Verify foreign key from asset_logs.asset_id to assets.id
	var exists bool
	err = conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.table_constraints 
			WHERE constraint_type = 'FOREIGN KEY'
			AND table_schema = 'public'
			AND table_name = 'asset_logs'
			AND constraint_name = 'asset_logs_asset_id_fkey'
		)
	`).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "foreign key constraint should exist on asset_logs.asset_id")

	// Verify CASCADE delete behavior
	var deleteRule string
	err = conn.QueryRow(ctx, `
		SELECT delete_rule 
		FROM information_schema.referential_constraints 
		WHERE constraint_schema = 'public'
		AND constraint_name = 'asset_logs_asset_id_fkey'
	`).Scan(&deleteRule)
	require.NoError(t, err)
	assert.Equal(t, "CASCADE", deleteRule, "foreign key should have CASCADE delete rule")
}

// setupTest creates a test configuration and logger.
// Tests will skip if database is not available.
func setupTest(t *testing.T) (*config.Config, logger.Logger) {
	t.Helper()

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Skipf("Skipping test: failed to load config: %v", err)
	}

	// Use test database if available
	if cfg.Database.Name == "ark" {
		t.Log("Warning: Using production database name. Consider setting ARK_DATABASE.NAME=ark_test")
	}

	loggerService := logger.NewLoggerService(cfg.Observability)
	t.Cleanup(func() {
		loggerService.Shutdown()
	})

	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	return cfg, log
}

// connectDB creates a database connection for verification.
func connectDB(t *testing.T, cfg *config.Config) *pgx.Conn {
	t.Helper()

	ctx := context.Background()
	dsn := cfg.Database.GetDSN()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	return conn
}
