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

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
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

	// Log migration details
	logger.Info().Int("count", len(m.Migrations)).Msg("loaded migration files")
	for _, mig := range m.Migrations {
		logger.Debug().Str("name", mig.Name).Msg("found migration")
	}

	from, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("retrieving current database migration version: %w", err)
	}

	logger.Info().Int32("current_version", from).Msg("starting migration")

	if err := m.Migrate(ctx); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	to, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("retrieving new database migration version: %w", err)
	}

	if from == to {
		logger.Info().Int32("version", to).Msg("database schema up to date")
	} else {
		logger.Info().Int32("from", from).Int32("to", to).Msg("migrated database schema")
	}

	// Validate schema after migration
	logger.Info().Msg("validating database schema")
	if err := validateSchema(ctx, conn); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}
	logger.Info().Msg("database schema validation passed")

	return nil
}

// Status prints the current migration status
func Status(ctx context.Context, logger *zerolog.Logger, cfg *config.Config) error {
	conn, err := connect(ctx, cfg)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	m, err := tern.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return err
	}

	version, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return err
	}

	logger.Info().Int32("version", version).Msg("current database version")
	return nil
}

// Validate checks if the schema is correct
func Validate(ctx context.Context, logger *zerolog.Logger, cfg *config.Config) error {
	conn, err := connect(ctx, cfg)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	return validateSchema(ctx, conn)
}

func connect(ctx context.Context, cfg *config.Config) (*pgx.Conn, error) {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))
	encodedPassword := url.QueryEscape(cfg.Database.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	return pgx.Connect(ctx, dsn)
}

func validateSchema(ctx context.Context, conn *pgx.Conn) error {
	// Check for critical tables
	tables := []string{"assets", "asset_logs", "schema_version"}

	for _, table := range tables {
		var exists bool
		err := conn.QueryRow(ctx,
			"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)",
			table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("checking table %s: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("missing required table: %s", table)
		}
	}
	return nil
}
