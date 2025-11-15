package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/srauf24/gardenjournal/internal/middleware"
	"github.com/srauf24/gardenjournal/internal/model"
	"github.com/srauf24/gardenjournal/internal/model/plant"
	"github.com/srauf24/gardenjournal/internal/repository"
)

// PlantService handles business logic for plant operations
type PlantService struct {
	plantRepo *repository.PlantRepository
}

// NewPlantService creates a new PlantService
func NewPlantService(plantRepo *repository.PlantRepository) *PlantService {
	return &PlantService{
		plantRepo: plantRepo,
	}
}

// CreatePlant creates a new plant after validation
func (s *PlantService) CreatePlant(ctx echo.Context, userID string, payload *plant.CreatePlantPayload) (*plant.Plant, error) {
	logger := middleware.GetLogger(ctx)

	// Validate payload
	if err := payload.Validate(); err != nil {
		logger.Error().Err(err).Msg("plant validation failed")
		return nil, fmt.Errorf("plant service: validation failed: %w", err)
	}

	logger.Debug().Msg("plant validation passed")

	// Call repository to create plant (ctx.Request().Context() gives context.Context for repository layer)
	createdPlant, err := s.plantRepo.CreatePlant(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create plant")
		return nil, fmt.Errorf("plant service: failed to create plant: %w", err)
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "plant_created").
		Str("plant_id", createdPlant.ID.String()).
		Str("name", createdPlant.Name).
		Str("species", createdPlant.Species).
		Str("location", func() string {
			if createdPlant.Location != nil {
				return *createdPlant.Location
			}
			return ""
		}()).
		Msg("Plant created successfully")

	return createdPlant, nil
}

// GetPlantByID retrieves a plant by ID
func (s *PlantService) GetPlantByID(ctx echo.Context, userID string, plantID uuid.UUID) (*plant.PopulatedPlant, error) {
	logger := middleware.GetLogger(ctx)

	// Validate inputs
	if userID == "" {
		logger.Error().Msg("user_id is required")
		return nil, fmt.Errorf("plant service: user_id is required")
	}
	if plantID == uuid.Nil {
		logger.Error().Msg("plant_id is required")
		return nil, fmt.Errorf("plant service: plant_id is required")
	}

	// Call repository to get plant
	plantItem, err := s.plantRepo.GetPlantByID(ctx.Request().Context(), userID, plantID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch plant by ID")
		return nil, fmt.Errorf("plant service: failed to get plant by id: %w", err)
	}

	return plantItem, nil
}

// GetPlants retrieves a paginated list of plants
func (s *PlantService) GetPlants(ctx echo.Context, userID string, query *plant.GetPlantsQuery) (*model.PaginatedResponse[plant.PopulatedPlant], error) {
	logger := middleware.GetLogger(ctx)

	// Validate user ID
	if userID == "" {
		logger.Error().Msg("user_id is required")
		return nil, fmt.Errorf("plant service: user_id is required")
	}

	// Validate and set defaults for query
	if err := query.Validate(); err != nil {
		logger.Error().Err(err).Msg("query validation failed")
		return nil, fmt.Errorf("plant service: query validation failed: %w", err)
	}

	// Call repository to get plants
	plants, err := s.plantRepo.GetPlants(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch plants")
		return nil, fmt.Errorf("plant service: failed to get plants: %w", err)
	}

	return plants, nil
}

// UpdatePlant updates an existing plant
func (s *PlantService) UpdatePlant(ctx echo.Context, userID string, plantID uuid.UUID, payload *plant.UpdatePlantPayload) (*plant.Plant, error) {
	logger := middleware.GetLogger(ctx)

	// Validate payload
	if err := payload.Validate(); err != nil {
		logger.Error().Err(err).Msg("plant validation failed")
		return nil, fmt.Errorf("plant service: validation failed: %w", err)
	}

	// Validate user ID and plant ID
	if userID == "" {
		logger.Error().Msg("user_id is required")
		return nil, fmt.Errorf("plant service: user_id is required")
	}
	if plantID == uuid.Nil {
		logger.Error().Msg("plant_id is required")
		return nil, fmt.Errorf("plant service: plant_id is required")
	}

	logger.Debug().Msg("plant validation passed")

	// Check if plant exists and belongs to user
	_, err := s.plantRepo.CheckPlantExists(ctx.Request().Context(), userID, plantID)
	if err != nil {
		logger.Error().Err(err).Msg("plant validation failed")
		return nil, fmt.Errorf("plant service: plant not found or unauthorized: %w", err)
	}

	logger.Debug().Msg("plant existence check passed")

	// Call repository to update plant
	updatedPlant, err := s.plantRepo.UpdatePlant(ctx.Request().Context(), userID, plantID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update plant")
		return nil, fmt.Errorf("plant service: failed to update plant: %w", err)
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "plant_updated").
		Str("plant_id", updatedPlant.ID.String()).
		Str("name", updatedPlant.Name).
		Str("species", updatedPlant.Species).
		Str("location", func() string {
			if updatedPlant.Location != nil {
				return *updatedPlant.Location
			}
			return ""
		}()).
		Msg("Plant updated successfully")

	return updatedPlant, nil
}

// DeletePlant deletes a plant
func (s *PlantService) DeletePlant(ctx echo.Context, userID string, plantID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	// Validate user ID and plant ID
	if userID == "" {
		logger.Error().Msg("user_id is required")
		return fmt.Errorf("plant service: user_id is required")
	}
	if plantID == uuid.Nil {
		logger.Error().Msg("plant_id is required")
		return fmt.Errorf("plant service: plant_id is required")
	}

	// Check if plant exists and belongs to user
	_, err := s.plantRepo.CheckPlantExists(ctx.Request().Context(), userID, plantID)
	if err != nil {
		logger.Error().Err(err).Msg("plant validation failed")
		return fmt.Errorf("plant service: plant not found or unauthorized: %w", err)
	}

	// Call repository to delete plant
	if err := s.plantRepo.DeletePlant(ctx.Request().Context(), userID, plantID); err != nil {
		logger.Error().Err(err).Msg("failed to delete plant")
		return fmt.Errorf("plant service: failed to delete plant: %w", err)
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "plant_deleted").
		Str("plant_id", plantID.String()).
		Msg("Plant deleted successfully")

	return nil
}
