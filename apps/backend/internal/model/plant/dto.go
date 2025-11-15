package plant

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreatePlantRequest needs to match the Plant model's custom types
// optional fields have omit empty

// Represents the incoming request body from the API (what the user sends when creating a plant) - defines what the user is allowed to send in a POST request.
// Used in  handler or service layer, right after you parse the JSON request:
type CreatePlantPayload struct {
	Name        string     `json:"name" validate:"required,max=100"`
	Species     string     `json:"species" validate:"required,max=100"`
	Location    *string    `json:"location,omitempty" validate:"omitempty,max=255"`
	PlantedDate *time.Time `json:"plantedDate,omitempty"`
	Notes       *string    `json:"notes,omitempty" validate:"omitempty,max=1000"`
	Metadata    *Metadata  `json:"metadata,omitempty"`
}

func (p *CreatePlantPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// UpdatePlantRequest - all fields are optional for partial updates
type UpdatePlantPayload struct {
	ID          uuid.UUID  `json:"id" validate:"required,uuid"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,max=100"`
	Species     *string    `json:"species,omitempty" validate:"omitempty,max=100"`
	Location    *string    `json:"location,omitempty" validate:"omitempty,max=255"`
	PlantedDate *time.Time `json:"plantedDate,omitempty"`
	Notes       *string    `json:"notes,omitempty" validate:"omitempty,max=1000"`
	Metadata    *Metadata  `json:"metadata,omitempty"`
}

func (p *UpdatePlantPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type GetPlantsQuery struct {
	Page        *int       `query:"page" validate:"omitempty,min=1"`
	Limit       *int       `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort        *string    `query:"sort" validate:"omitempty,oneof=created_at updated_at name species location planted_date sort_order"`
	Order       *string    `query:"order" validate:"omitempty,oneof=asc desc"`
	Search      *string    `query:"search" validate:"omitempty,min=1"`
	Species     *string    `query:"species" validate:"omitempty,min=1"`
	Location    *string    `query:"location" validate:"omitempty,min=1"`
	PlantedFrom *time.Time `query:"plantedFrom"` // filter p.planted_date >= @planted_from
	PlantedTo   *time.Time `query:"plantedTo"`   // filter p.planted_date <= @planted_to
}

func (q *GetPlantsQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(q); err != nil {
		return err
	}

	// Set defaults for pagination
	if q.Page == nil {
		defaultPage := 1
		q.Page = &defaultPage
	}
	if q.Limit == nil {
		defaultLimit := 20
		q.Limit = &defaultLimit
	}
	if q.Sort == nil {
		defaultSort := "created_at"
		q.Sort = &defaultSort
	}
	if q.Order == nil {
		defaultOrder := "desc"
		q.Order = &defaultOrder
	}

	return nil
}

type GetPlantByIDPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetPlantByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type DeletePlantPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeletePlantPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// high level idea: stats may include total plants, plants by species, etc.
type GetPlantsStatsPayload struct{}

func (p *GetPlantsStatsPayload) Validate() error {
	return nil
}
