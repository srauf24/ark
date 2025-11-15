package observation

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateObservationPayload struct {
	PlantID  uuid.UUID  `json:"plantId" validate:"required,uuid"`
	Date     *time.Time `json:"date,omitempty"`
	HeightCM *float64   `json:"heightCm" validate:"omitempty,gt=0"`
	Notes    *string    `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

func (p *CreateObservationPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type UpdateObservationPayload struct {
	ID       uuid.UUID `json:"id" validate:"required,uuid"`
	HeightCM *float64  `json:"heightCm,omitempty" validate:"omitempty,gt=0"`
	Notes    *string   `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

func (p *UpdateObservationPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type GetObservationsQuery struct {
	Page   *int    `query:"page" validate:"omitempty,min=1"`
	Limit  *int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort   *string `query:"sort" validate:"omitempty,oneof=created_at updated_at date height_cm sort_order"`
	Order  *string `query:"order" validate:"omitempty,oneof=asc desc"`
	Search *string `query:"search" validate:"omitempty,min=1"`
}

func (q *GetObservationsQuery) Validate() error {
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

type GetObservationByIDPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetObservationByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type DeleteObservationPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteObservationPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
