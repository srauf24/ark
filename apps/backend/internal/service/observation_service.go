package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"ark/internal/lib/weather"
	"ark/internal/middleware"
	"ark/internal/model"
	"ark/internal/model/observation"
	"ark/internal/repository"
)

// ObservationService handles business logic for observation operations
type ObservationService struct {
	observationRepo *repository.ObservationRepository
	plantRepo       *repository.PlantRepository
	weatherClient   *weather.Client // Reserved for future weather enrichment
}

// NewObservationService creates a new ObservationService
func NewObservationService(
	observationRepo *repository.ObservationRepository,
	plantRepo *repository.PlantRepository,
	weatherClient *weather.Client,
) *ObservationService {
	return &ObservationService{
		observationRepo: observationRepo,
		plantRepo:       plantRepo,
		weatherClient:   weatherClient,
	}
}

// CreateObservation creates a new observation after validation
// Future: Add weather enrichment when plant metadata contains lat/lon coordinates
func (s *ObservationService) CreateObservation(ctx echo.Context, userID string, payload *observation.CreateObservationPayload) (*observation.Observation, error) {
	logger := middleware.GetLogger(ctx)

	// Validate payload
	if err := payload.Validate(); err != nil {
		logger.Error().Err(err).Msg("observation validation failed")
		return nil, fmt.Errorf("observation service: validation failed: %w", err)
	}

	logger.Debug().Msg("observation validation passed")

	// Check plant exists and belongs to user
	_, err := s.plantRepo.CheckPlantExists(ctx.Request().Context(), userID, payload.PlantID)
	if err != nil {
		logger.Error().Err(err).Msg("plant validation failed")
		return nil, fmt.Errorf("observation service: plant not found or unauthorized: %w", err)
	}

	logger.Debug().Msg("plant validation passed")

	// TODO: Future weather enrichment hook
	// When implementing weather enrichment:
	// 1. Extract lat/lon from plant.Metadata (if present)
	// 2. Call s.enrichPlantWithWeather(ctx, userID, plantItem)
	// 3. Update plant metadata with weather snapshot (best-effort, don't block on failure)

	// Create observation
	createdObservation, err := s.observationRepo.CreateObservation(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create observation")
		return nil, fmt.Errorf("observation service: failed to create observation: %w", err)
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "observation_created").
		Str("observation_id", createdObservation.ID.String()).
		Str("plant_id", createdObservation.PlantID.String()).
		Time("date", createdObservation.Date).
		Float64("height_cm", func() float64 {
			if createdObservation.HeightCM != nil {
				return *createdObservation.HeightCM
			}
			return 0
		}()).
		Msg("Observation created successfully")

	return createdObservation, nil
}

// GetObservationByID retrieves an observation by ID
func (s *ObservationService) GetObservationByID(ctx echo.Context, userID string, observationID uuid.UUID) (*observation.Observation, error) {
	logger := middleware.GetLogger(ctx)

	// Validate inputs
	if userID == "" {
		logger.Error().Msg("user_id is required")
		return nil, fmt.Errorf("observation service: user_id is required")
	}
	if observationID == uuid.Nil {
		logger.Error().Msg("observation_id is required")
		return nil, fmt.Errorf("observation service: observation_id is required")
	}

	// Call repository to get observation
	observationItem, err := s.observationRepo.GetObservationByID(ctx.Request().Context(), userID, observationID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch observation by ID")
		return nil, fmt.Errorf("observation service: failed to get observation by id: %w", err)
	}

	return observationItem, nil
}

// GetObservations retrieves a paginated list of observations
func (s *ObservationService) GetObservations(ctx echo.Context, userID string, query *observation.GetObservationsQuery) (*model.PaginatedResponse[observation.Observation], error) {
	logger := middleware.GetLogger(ctx)

	// Validate user ID
	if userID == "" {
		logger.Error().Msg("user_id is required")
		return nil, fmt.Errorf("observation service: user_id is required")
	}

	// Validate and set defaults for query
	if err := query.Validate(); err != nil {
		logger.Error().Err(err).Msg("query validation failed")
		return nil, fmt.Errorf("observation service: query validation failed: %w", err)
	}

	// Call repository to get observations
	observations, err := s.observationRepo.GetObservations(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch observations")
		return nil, fmt.Errorf("observation service: failed to get observations: %w", err)
	}

	return observations, nil
}

// UpdateObservation updates an existing observation
func (s *ObservationService) UpdateObservation(ctx echo.Context, userID string, observationID uuid.UUID, payload *observation.UpdateObservationPayload) (*observation.Observation, error) {
	logger := middleware.GetLogger(ctx)

	// Validate payload
	if err := payload.Validate(); err != nil {
		logger.Error().Err(err).Msg("observation validation failed")
		return nil, fmt.Errorf("observation service: validation failed: %w", err)
	}

	// Validate user ID and observation ID
	if userID == "" {
		logger.Error().Msg("user_id is required")
		return nil, fmt.Errorf("observation service: user_id is required")
	}
	if observationID == uuid.Nil {
		logger.Error().Msg("observation_id is required")
		return nil, fmt.Errorf("observation service: observation_id is required")
	}

	logger.Debug().Msg("observation validation passed")

	// Check if observation exists and belongs to user
	observationItem, err := s.observationRepo.CheckObservationExists(ctx.Request().Context(), userID, observationID)
	if err != nil {
		logger.Error().Err(err).Msg("observation validation failed")
		return nil, fmt.Errorf("observation service: observation not found or unauthorized: %w", err)
	}

	logger.Debug().Msg("observation existence check passed")

	// Verify plant ownership
	_, err = s.plantRepo.CheckPlantExists(ctx.Request().Context(), userID, observationItem.PlantID)
	if err != nil {
		logger.Error().Err(err).Msg("plant validation failed")
		return nil, fmt.Errorf("observation service: plant not found or unauthorized: %w", err)
	}

	logger.Debug().Msg("plant validation passed")

	// Call repository to update observation
	updatedObservation, err := s.observationRepo.UpdateObservation(ctx.Request().Context(), userID, observationID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update observation")
		return nil, fmt.Errorf("observation service: failed to update observation: %w", err)
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "observation_updated").
		Str("observation_id", updatedObservation.ID.String()).
		Str("plant_id", updatedObservation.PlantID.String()).
		Time("date", updatedObservation.Date).
		Float64("height_cm", func() float64 {
			if updatedObservation.HeightCM != nil {
				return *updatedObservation.HeightCM
			}
			return 0
		}()).
		Msg("Observation updated successfully")

	return updatedObservation, nil
}

// DeleteObservation deletes an observation
func (s *ObservationService) DeleteObservation(ctx echo.Context, userID string, observationID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	// Validate user ID and observation ID
	if userID == "" {
		logger.Error().Msg("user_id is required")
		return fmt.Errorf("observation service: user_id is required")
	}
	if observationID == uuid.Nil {
		logger.Error().Msg("observation_id is required")
		return fmt.Errorf("observation service: observation_id is required")
	}

	// Check if observation exists and belongs to user
	observationItem, err := s.observationRepo.CheckObservationExists(ctx.Request().Context(), userID, observationID)
	if err != nil {
		logger.Error().Err(err).Msg("observation validation failed")
		return fmt.Errorf("observation service: observation not found or unauthorized: %w", err)
	}

	// Verify plant ownership
	_, err = s.plantRepo.CheckPlantExists(ctx.Request().Context(), userID, observationItem.PlantID)
	if err != nil {
		logger.Error().Err(err).Msg("plant validation failed")
		return fmt.Errorf("observation service: plant not found or unauthorized: %w", err)
	}

	// Call repository to delete observation
	if err := s.observationRepo.DeleteObservation(ctx.Request().Context(), userID, observationID); err != nil {
		logger.Error().Err(err).Msg("failed to delete observation")
		return fmt.Errorf("observation service: failed to delete observation: %w", err)
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "observation_deleted").
		Str("observation_id", observationID.String()).
		Msg("Observation deleted successfully")

	return nil
}

// ============================================================
// FUTURE: Weather Enrichment Methods (commented out for MVP)
// ============================================================
// Uncomment and implement these methods when adding weather enrichment:
//
// // enrichPlantWithWeather attempts to fetch weather data and update plant metadata
// // This is a best-effort operation - failures are logged but don't block observation creation
// func (s *ObservationService) enrichPlantWithWeather(ctx context.Context, userID string, plantItem *plant.Plant) error {
// 	// 1. Parse plant metadata to extract latitude and longitude
// 	// 2. If coordinates present, call s.weatherClient.FetchWeatherSafe(ctx, lat, lon)
// 	// 3. If weather data retrieved, update plant.Metadata with lastWeatherSnapshot
// 	// 4. Call s.plantRepo.UpdatePlant() to persist updated metadata
// 	// 5. Log errors but don't fail - this is best-effort enrichment
// 	return nil
// }
//
// // extractCoordinatesFromMetadata extracts lat/lon from plant metadata JSONB
// func extractCoordinatesFromMetadata(metadata json.RawMessage) (lat, lon float64, ok bool) {
// 	// Parse metadata and look for latitude/longitude fields
// 	return 0, 0, false
// }
