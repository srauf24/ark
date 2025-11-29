package main

import (
	"context"
	"errors"
	"fmt"
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
		if err := runMigrate(); err != nil {
			// Use a basic logger since we might not have loaded config yet
			l := zerolog.New(os.Stderr).With().Timestamp().Logger()
			l.Fatal().Err(err).Msg("migration command failed")
		}
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

func runMigrate() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: ark migrate <command> [args]\ncommands: up, status, validate")
	}

	cmd := os.Args[2]

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Initialize logger
	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()
	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	ctx := context.Background()

	switch cmd {
	case "up":
		log.Info().Msg("running database migrations...")
		if err := database.Migrate(ctx, &log, cfg); err != nil {
			return err
		}
		log.Info().Msg("migrations completed successfully")
	case "status":
		log.Info().Msg("checking migration status...")
		if err := database.Status(ctx, &log, cfg); err != nil {
			return err
		}
	case "validate":
		log.Info().Msg("validating database schema...")
		if err := database.Validate(ctx, &log, cfg); err != nil {
			return err
		}
		log.Info().Msg("schema validation passed")
	default:
		return fmt.Errorf("unknown migration command: %s", cmd)
	}

	return nil
}
