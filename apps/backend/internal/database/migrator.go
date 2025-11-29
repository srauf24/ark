package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"strconv"

	"ark/internal/config"

	"github.com/jackc/pgx/v5"
	tern "github.com/jackc/tern/v2/migrate"
	"github.com/rs/zerolog"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(ctx context.Context, logger *zerolog.Logger, cfg *config.Config) error {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	// URL-encode the password
	encodedPassword := url.QueryEscape(cfg.Database.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	logger.Info().Msg("connecting to database for migrations")
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer conn.Close(ctx)

	m, err := tern.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return fmt.Errorf("constructing database migrator: %w", err)
	}

	subtree, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("retrieving database migrations subtree: %w", err)
	}

	if err := m.LoadMigrations(subtree); err != nil {
		return fmt.Errorf("loading database migrations: %w", err)
	}

	// Log loaded migrations
	logger.Info().
		Int("migration_count", len(m.Migrations)).
		Msg("loaded migration files")

	// Log each migration file name
	for i, migration := range m.Migrations {
		logger.Debug().
			Int("sequence", i+1).
			Str("name", migration.Name).
			Msg("migration file loaded")
	}

	from, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("retrieving current database migration version: %w", err)
	}

	to := int32(len(m.Migrations))

	if from == to {
		logger.Info().
			Int32("version", from).
			Msg("database schema up to date")
	} else {
		logger.Info().
			Int32("from_version", from).
			Int32("to_version", to).
			Msg("starting database migration")

		if err := m.Migrate(ctx); err != nil {
			return fmt.Errorf("executing migrations: %w", err)
		}

		logger.Info().
			Int32("from_version", from).
			Int32("to_version", to).
			Msg("database migration completed successfully")
	}

	// Validate schema after migration
	logger.Info().Msg("validating database schema")
	if err := validateSchema(ctx, conn); err != nil {
		logger.Error().Err(err).Msg("schema validation failed")
		return fmt.Errorf("schema validation failed: %w", err)
	}
	logger.Info().Msg("schema validation passed")

	return nil
}

// validateSchema checks that all expected tables exist after migration.
// This ensures migrations executed successfully and created the required schema.
func validateSchema(ctx context.Context, conn *pgx.Conn) error {
	expectedTables := []string{"assets", "asset_logs"}

	for _, table := range expectedTables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)
		`
		err := conn.QueryRow(ctx, query, table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("checking if table %s exists: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("table %s does not exist after migration", table)
		}
	}

	return nil
}
