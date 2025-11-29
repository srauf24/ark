package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"ark/internal/config"
	"ark/internal/database"
	"ark/internal/handler"
	"ark/internal/logger"
	"ark/internal/repository"
	"ark/internal/router"
	"ark/internal/server"
	"ark/internal/service"

	"github.com/rs/zerolog"
)

const DefaultContextTimeout = 30

func main() {
	// Check for migrate subcommand before loading full config
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrate()
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize New Relic logger service
	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	if cfg.Primary.Env != "local" {
		if err := database.Migrate(context.Background(), &log, cfg); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database")
		}
	}

	// Initialize server
	srv, err := server.New(cfg, &log, loggerService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	// Initialize repositories, services, and handlers
	repos := repository.NewRepositories(srv)
	services, serviceErr := service.NewServices(srv, repos)
	if serviceErr != nil {
		log.Fatal().Err(serviceErr).Msg("could not create services")
	}
	handlers := handler.NewHandlers(srv, services)

	// Initialize router
	r := router.NewRouter(srv, handlers, services)

	// Setup HTTP server
	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Start server
	go func() {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}
	stop()
	cancel()

	log.Info().Msg("server exited properly")
}

// runMigrate handles the migrate subcommand
func runMigrate() {
	if len(os.Args) < 3 {
		printMigrateUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	// Load config for migration commands
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize logger
	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()
	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	ctx := context.Background()

	switch subcommand {
	case "up":
		log.Info().Msg("running database migrations")
		if err := database.Migrate(ctx, &log, cfg); err != nil {
			log.Fatal().Err(err).Msg("migration failed")
		}
		log.Info().Msg("migrations completed successfully")

	case "status":
		showMigrationStatus(ctx, &log, cfg)

	case "validate":
		validateMigrationSchema(ctx, &log, cfg)

	default:
		log.Error().Str("subcommand", subcommand).Msg("unknown migrate subcommand")
		printMigrateUsage()
		os.Exit(1)
	}
}

// showMigrationStatus displays the current migration version
func showMigrationStatus(ctx context.Context, log *zerolog.Logger, cfg *config.Config) {
	// Connect to database
	srv, err := server.New(cfg, log, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	// Query schema_version table
	var version int32
	err = srv.DB.Pool.QueryRow(ctx, "SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&version)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to query migration version")
	}

	log.Info().Int32("current_version", version).Msg("migration status")
}

// validateMigrationSchema validates that all expected tables exist
func validateMigrationSchema(ctx context.Context, log *zerolog.Logger, cfg *config.Config) {
	// Connect to database
	srv, err := server.New(cfg, log, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	// Check for expected tables
	expectedTables := []string{"assets", "asset_logs"}
	allExist := true

	for _, table := range expectedTables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)
		`
		err := srv.DB.Pool.QueryRow(ctx, query, table).Scan(&exists)
		if err != nil {
			log.Fatal().Err(err).Str("table", table).Msg("failed to check table existence")
		}

		if exists {
			log.Info().Str("table", table).Msg("table exists")
		} else {
			log.Error().Str("table", table).Msg("table does not exist")
			allExist = false
		}
	}

	if allExist {
		log.Info().Msg("schema validation passed")
	} else {
		log.Fatal().Msg("schema validation failed")
	}
}

// printMigrateUsage prints usage information for the migrate command
func printMigrateUsage() {
	usage := `Usage: ark migrate <subcommand>

Subcommands:
  up        Run pending database migrations
  status    Show current migration version
  validate  Validate that all expected tables exist

Examples:
  ark migrate up        # Run migrations
  ark migrate status    # Check current version
  ark migrate validate  # Verify schema
`
	println(usage)
}
