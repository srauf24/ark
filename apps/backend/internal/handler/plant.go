package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"ark/internal/middleware"
	"ark/internal/model/plant"
	"ark/internal/server"
	"ark/internal/service"
)

// PlantHandler handles plant-related HTTP requests
type PlantHandler struct {
	Handler
	plantService *service.PlantService
}

// NewPlantHandler creates a new PlantHandler
func NewPlantHandler(s *server.Server, plantService *service.PlantService) *PlantHandler {
	return &PlantHandler{
		Handler:      NewHandler(s),
		plantService: plantService,
	}
}

// CreatePlant handles POST /api/v1/plants
func (h *PlantHandler) CreatePlant(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *plant.CreatePlantPayload) (*plant.Plant, error) {
			userID := middleware.GetUserID(c)
			return h.plantService.CreatePlant(c, userID, payload)
		},
		http.StatusCreated,
		&plant.CreatePlantPayload{},
	)(c)
}

// GetPlantByID handles GET /api/v1/plants/:id
func (h *PlantHandler) GetPlantByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *plant.GetPlantByIDPayload) (*plant.PopulatedPlant, error) {
			userID := middleware.GetUserID(c)
			return h.plantService.GetPlantByID(c, userID, payload.ID)
		},
		http.StatusOK,
		&plant.GetPlantByIDPayload{},
	)(c)
}

// GetPlants handles GET /api/v1/plants
func (h *PlantHandler) GetPlants(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *plant.GetPlantsQuery) (interface{}, error) {
			userID := middleware.GetUserID(c)
			return h.plantService.GetPlants(c, userID, query)
		},
		http.StatusOK,
		&plant.GetPlantsQuery{},
	)(c)
}

// UpdatePlant handles PUT /api/v1/plants/:id
func (h *PlantHandler) UpdatePlant(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *plant.UpdatePlantPayload) (*plant.Plant, error) {
			userID := middleware.GetUserID(c)
			return h.plantService.UpdatePlant(c, userID, payload.ID, payload)
		},
		http.StatusOK,
		&plant.UpdatePlantPayload{},
	)(c)
}

// DeletePlant handles DELETE /api/v1/plants/:id
func (h *PlantHandler) DeletePlant(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *plant.DeletePlantPayload) error {
			userID := middleware.GetUserID(c)
			return h.plantService.DeletePlant(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&plant.DeletePlantPayload{},
	)(c)
}
