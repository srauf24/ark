package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"ark/internal/model"
	"ark/internal/model/observation"
	"ark/internal/server"
)

type ObservationRepository struct {
	server *server.Server
}

func NewObservationRepository(server *server.Server) *ObservationRepository {
	return &ObservationRepository{server: server}
}

// note: typically the service layer translates the DTO to a model which is then used by the repository layer. In this case we are directly using the DTO in the repository layer for simplicity.
func (r *ObservationRepository) CreateObservation(ctx context.Context, userID string, payload *observation.CreateObservationPayload) (*observation.Observation, error) {
	stmt := `
        INSERT INTO
            observations (
                user_id,
                plant_id,
                date,
                height_cm,
                notes
            )
        VALUES
            (
                @user_id,
                @plant_id,
                COALESCE(@date, CURRENT_TIMESTAMP),
                @height_cm,
                @notes
            )
        RETURNING
        *
    `
	// use server.db to execute the query
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":   userID,
		"plant_id":  payload.PlantID,
		"date":      payload.Date,
		"height_cm": payload.HeightCM,
		"notes":     payload.Notes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create observation query for user_id=%s plant_id=%s: %w", userID, payload.PlantID, err)
	}
	// use pgx library to deserialize row into a struct from the data base (collect one row)
	observationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[observation.Observation])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:observations for user_id=%s plant_id=%s: %w", userID, payload.PlantID, err)
	}
	return &observationItem, nil
}

func (r *ObservationRepository) GetObservationByID(ctx context.Context, userID string, observationID uuid.UUID) (*observation.Observation, error) {
	stmt := `
        SELECT
            *
        FROM
            observations
        WHERE
            user_id = @user_id
            AND id = @observation_id
    `
	// use server.db to execute the query
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":        userID,
		"observation_id": observationID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get observation by id query for user_id=%s observation_id=%s: %w", userID, observationID, err)
	}
	// use pgx library to deserialize row into a struct from the data base (collect one row)
	observationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[observation.Observation])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:observations for user_id=%s observation_id=%s: %w", userID, observationID, err)
	}
	return &observationItem, nil
}

func (r *ObservationRepository) CheckObservationExists(ctx context.Context, userID string, observationID uuid.UUID) (*observation.Observation, error) {
	stmt := `
		SELECT
			*
		FROM
			observations
		WHERE
			id = @observation_id
			AND user_id = @user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"observation_id": observationID,
		"user_id":        userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check if observation exists for observation_id=%s user_id=%s: %w", observationID.String(), userID, err)
	}

	observationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[observation.Observation])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:observations for observation_id=%s user_id=%s: %w", observationID.String(), userID, err)
	}

	return &observationItem, nil
}

func (r *ObservationRepository) GetObservations(ctx context.Context, userID string, query *observation.GetObservationsQuery) (*model.PaginatedResponse[observation.Observation], error) {
	stmt := `
		SELECT
			o.*
		FROM observations o
	`

	args := pgx.NamedArgs{
		"user_id": userID,
	}
	conditions := []string{"o.user_id = @user_id"}

	// --- Dynamic filtering ---
	if query.Search != nil {
		conditions = append(conditions, "o.notes ILIKE @search")
		args["search"] = "%" + *query.Search + "%"
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	// --- Count query for pagination ---
	countStmt := "SELECT COUNT(*) FROM observations o"
	if len(conditions) > 0 {
		countStmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count observations for user_id=%s: %w", userID, err)
	}

	// --- Sorting ---
	if query.Sort != nil {
		stmt += " ORDER BY o." + *query.Sort
		if query.Order != nil && *query.Order == "desc" {
			stmt += " DESC"
		} else {
			stmt += " ASC"
		}
	} else {
		stmt += " ORDER BY o.created_at DESC"
	}

	// --- Pagination ---
	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	// --- Query execution ---
	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get observations query for user_id=%s: %w", userID, err)
	}
	defer rows.Close()

	// --- Collect into structs ---
	observations, err := pgx.CollectRows(rows, pgx.RowToStructByName[observation.Observation])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[observation.Observation]{
				Data:       []observation.Observation{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to map observation rows for user_id=%s: %w", userID, err)
	}

	// --- Build pagination response ---
	return &model.PaginatedResponse[observation.Observation]{
		Data:       observations,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}

func (r *ObservationRepository) UpdateObservation(ctx context.Context, userID string, observationID uuid.UUID, payload *observation.UpdateObservationPayload) (*observation.Observation, error) {
	// Build dynamic UPDATE statement based on which fields are provided
	updates := []string{}
	args := pgx.NamedArgs{
		"observation_id": observationID,
		"user_id":        userID,
	}

	if payload.HeightCM != nil {
		updates = append(updates, "height_cm = @height_cm")
		args["height_cm"] = payload.HeightCM
	}

	if payload.Notes != nil {
		updates = append(updates, "notes = @notes")
		args["notes"] = payload.Notes
	}

	// If no fields to update, return error
	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update for observation_id=%s", observationID.String())
	}

	// Always update the updated_at timestamp
	updates = append(updates, "updated_at = NOW()")

	stmt := `
		UPDATE observations
		SET ` + strings.Join(updates, ", ") + `
		WHERE id = @observation_id
			AND user_id = @user_id
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update observation query for observation_id=%s user_id=%s: %w", observationID.String(), userID, err)
	}

	observationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[observation.Observation])
	if err != nil {
		return nil, fmt.Errorf("failed to collect updated observation for observation_id=%s user_id=%s: %w", observationID.String(), userID, err)
	}

	return &observationItem, nil
}

func (r *ObservationRepository) DeleteObservation(ctx context.Context, userID string, observationID uuid.UUID) error {
	stmt := `
		DELETE FROM observations
		WHERE id = @observation_id
			AND user_id = @user_id
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"observation_id": observationID,
		"user_id":        userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete observation for observation_id=%s user_id=%s: %w", observationID.String(), userID, err)
	}

	// Check if any rows were affected
	if result.RowsAffected() == 0 {
		return fmt.Errorf("observation not found or not authorized for observation_id=%s user_id=%s", observationID.String(), userID)
	}

	return nil
}
