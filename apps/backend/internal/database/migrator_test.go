package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateSchema_AllTablesExist verifies that validateSchema returns nil
// when all expected tables exist in the database.
func TestValidateSchema_AllTablesExist(t *testing.T) {
	ctx := context.Background()

	// Setup: Create test database connection
	conn, err := pgx.Connect(ctx, getTestDSN(t))
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return
	}
	require.NoError(t, err, "failed to connect to test database")
	defer conn.Close(ctx)

	// Create the expected tables
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS assets (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id TEXT NOT NULL,
			name TEXT NOT NULL
		)
	`)
	require.NoError(t, err, "failed to create assets table")

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS asset_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			asset_id UUID NOT NULL,
			content TEXT NOT NULL
		)
	`)
	require.NoError(t, err, "failed to create asset_logs table")

	// Cleanup
	defer func() {
		conn.Exec(ctx, "DROP TABLE IF EXISTS asset_logs CASCADE")
		conn.Exec(ctx, "DROP TABLE IF EXISTS assets CASCADE")
	}()

	// Action: Call validateSchema
	err = validateSchema(ctx, conn)

	// Assert: Should return nil (no error)
	assert.NoError(t, err, "validateSchema should return nil when all tables exist")
}

// TestValidateSchema_MissingAssetsTable verifies that validateSchema returns
// an error when the assets table is missing.
func TestValidateSchema_MissingAssetsTable(t *testing.T) {
	ctx := context.Background()

	// Setup: Create test database connection
	conn, err := pgx.Connect(ctx, getTestDSN(t))
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return
	}
	require.NoError(t, err, "failed to connect to test database")
	defer conn.Close(ctx)

	// Create only asset_logs table (assets is missing)
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS asset_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			content TEXT NOT NULL
		)
	`)
	require.NoError(t, err, "failed to create asset_logs table")

	// Cleanup
	defer func() {
		conn.Exec(ctx, "DROP TABLE IF EXISTS asset_logs CASCADE")
		conn.Exec(ctx, "DROP TABLE IF EXISTS assets CASCADE")
	}()

	// Action: Call validateSchema
	err = validateSchema(ctx, conn)

	// Assert: Should return error mentioning "assets"
	assert.Error(t, err, "validateSchema should return error when assets table is missing")
	assert.Contains(t, err.Error(), "assets", "error message should mention 'assets' table")
}

// TestValidateSchema_MissingAssetLogsTable verifies that validateSchema returns
// an error when the asset_logs table is missing.
func TestValidateSchema_MissingAssetLogsTable(t *testing.T) {
	ctx := context.Background()

	// Setup: Create test database connection
	conn, err := pgx.Connect(ctx, getTestDSN(t))
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return
	}
	require.NoError(t, err, "failed to connect to test database")
	defer conn.Close(ctx)

	// Create only assets table (asset_logs is missing)
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS assets (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id TEXT NOT NULL,
			name TEXT NOT NULL
		)
	`)
	require.NoError(t, err, "failed to create assets table")

	// Cleanup
	defer func() {
		conn.Exec(ctx, "DROP TABLE IF EXISTS asset_logs CASCADE")
		conn.Exec(ctx, "DROP TABLE IF EXISTS assets CASCADE")
	}()

	// Action: Call validateSchema
	err = validateSchema(ctx, conn)

	// Assert: Should return error mentioning "asset_logs"
	assert.Error(t, err, "validateSchema should return error when asset_logs table is missing")
	assert.Contains(t, err.Error(), "asset_logs", "error message should mention 'asset_logs' table")
}

// getTestDSN returns the DSN for the test database.
// It uses environment variables or skips the test if not configured.
func getTestDSN(t *testing.T) string {
	t.Helper()

	// Check for test database configuration in environment
	// This matches the pattern used in other backend tests
	host := getEnvOrDefault("ARK_DATABASE.HOST", "localhost")
	port := getEnvOrDefault("ARK_DATABASE.PORT", "5432")
	user := getEnvOrDefault("ARK_DATABASE.USER", "postgres")
	password := getEnvOrDefault("ARK_DATABASE.PASSWORD", "")
	dbname := getEnvOrDefault("ARK_DATABASE.NAME", "ark")
	sslmode := getEnvOrDefault("ARK_DATABASE.SSL_MODE", "disable")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)

	return dsn
}

// getEnvOrDefault returns environment variable value or default if not set.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
