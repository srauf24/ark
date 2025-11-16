package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"ark/internal/middleware"
	"ark/internal/model/observation"
	"ark/internal/server"
	"ark/internal/service"
)

// ObservationHandler handles observation-related HTTP requests
type ObservationHandler struct {
	Handler
	observationService *service.ObservationService
}

// NewObservationHandler creates a new ObservationHandler
func NewObservationHandler(s *server.Server, observationService *service.ObservationService) *ObservationHandler {
	return &ObservationHandler{
		Handler:            NewHandler(s),
		observationService: observationService,
	}
}

// CreateObservation handles POST /api/v1/observations
func (h *ObservationHandler) CreateObservation(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *observation.CreateObservationPayload) (*observation.Observation, error) {
			userID := middleware.GetUserID(c)
			return h.observationService.CreateObservation(c, userID, payload)
		},
		http.StatusCreated,
		&observation.CreateObservationPayload{},
	)(c)
}

// GetObservationByID handles GET /api/v1/observations/:id
func (h *ObservationHandler) GetObservationByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *observation.GetObservationByIDPayload) (*observation.Observation, error) {
			userID := middleware.GetUserID(c)
			return h.observationService.GetObservationByID(c, userID, payload.ID)
		},
		http.StatusOK,
		&observation.GetObservationByIDPayload{},
	)(c)
}

// GetObservations handles GET /api/v1/observations
func (h *ObservationHandler) GetObservations(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *observation.GetObservationsQuery) (interface{}, error) {
			userID := middleware.GetUserID(c)
			return h.observationService.GetObservations(c, userID, query)
		},
		http.StatusOK,
		&observation.GetObservationsQuery{},
	)(c)
}

// UpdateObservation handles PUT /api/v1/observations/:id
func (h *ObservationHandler) UpdateObservation(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *observation.UpdateObservationPayload) (*observation.Observation, error) {
			userID := middleware.GetUserID(c)
			return h.observationService.UpdateObservation(c, userID, payload.ID, payload)
		},
		http.StatusOK,
		&observation.UpdateObservationPayload{},
	)(c)
}

// DeleteObservation handles DELETE /api/v1/observations/:id
func (h *ObservationHandler) DeleteObservation(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *observation.DeleteObservationPayload) error {
			userID := middleware.GetUserID(c)
			return h.observationService.DeleteObservation(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&observation.DeleteObservationPayload{},
	)(c)
}
