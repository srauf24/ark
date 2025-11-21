package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"ark/internal/errs"
	"ark/internal/model"
)

// LogRepository provides data access methods for asset_logs table.
// All methods enforce user isolation and most enforce asset scoping.
type LogRepository struct {
	db *pgxpool.Pool
}

// NewLogRepository creates a new LogRepository with the given database pool.
func NewLogRepository(db *pgxpool.Pool) *LogRepository {
	return &LogRepository{db: db}
}

// GetByID retrieves a single log by ID for the specified user.
// Returns NotFoundError if the log doesn't exist or belongs to another user.
// Note: Only user_id is checked (not asset_id) since log ID is globally unique.
func (r *LogRepository) GetByID(ctx context.Context, userID string, logID uuid.UUID) (*model.AssetLog, error) {
	query := `
		SELECT id, asset_id, user_id, content, tags, created_at, updated_at
		FROM asset_logs
		WHERE id = @logID AND user_id = @userID
	`

	args := pgx.NamedArgs{
		"logID":  logID,
		"userID": userID,
	}

	var log model.AssetLog
	err := r.db.QueryRow(ctx, query, args).Scan(
		&log.ID,
		&log.AssetID,
		&log.UserID,
		&log.Content,
		&log.Tags, // pgx handles []string ↔ text[] automatically
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewNotFoundError("log not found", false, nil)
		}
		return nil, fmt.Errorf("get log by id: %w", err)
	}

	return &log, nil
}

// buildLogWhereClause builds dynamic WHERE clause for ListByAsset/CountByAsset with filters
func buildLogWhereClause(params *model.LogQueryParams, args pgx.NamedArgs) string {
	clauses := []string{"user_id = @userID", "asset_id = @assetID"}

	// Tags filter: log must have ALL specified tags (AND logic)
	if len(params.Tags) > 0 {
		clauses = append(clauses, "tags @> @tags::text[]")
		args["tags"] = params.Tags
	}

	// Content search (case-insensitive)
	if params.Search != nil {
		searchPattern := "%" + *params.Search + "%"
		clauses = append(clauses, "content ILIKE @search")
		args["search"] = searchPattern
	}

	// Date range: created_at >= start_date
	if params.StartDate != nil {
		clauses = append(clauses, "created_at >= @startDate")
		args["startDate"] = *params.StartDate
	}

	// Date range: created_at <= end_date
	if params.EndDate != nil {
		clauses = append(clauses, "created_at <= @endDate")
		args["endDate"] = *params.EndDate
	}

	return "WHERE " + strings.Join(clauses, " AND ")
}

// validateLogSortBy prevents SQL injection by validating sort column
func validateLogSortBy(sortBy string) error {
	allowed := map[string]bool{
		"created_at": true,
		"updated_at": true,
	}
	if !allowed[sortBy] {
		return fmt.Errorf("invalid sort_by: %s", sortBy)
	}
	return nil
}

// ListByAsset retrieves logs for a specific asset with optional filtering and pagination
func (r *LogRepository) ListByAsset(ctx context.Context, userID string, assetID uuid.UUID, params *model.LogQueryParams) ([]*model.AssetLog, error) {
	// Validate sort parameters to prevent SQL injection
	if err := validateLogSortBy(params.SortBy); err != nil {
		return nil, err
	}
	if err := validateSortOrder(params.SortOrder); err != nil {
		return nil, err
	}

	// Build WHERE clause with named args
	args := pgx.NamedArgs{
		"userID":  userID,
		"assetID": assetID,
		"limit":   params.Limit,
		"offset":  params.Offset,
	}
	whereClause := buildLogWhereClause(params, args)

	// Build complete query with ORDER BY and LIMIT/OFFSET
	query := fmt.Sprintf(`
		SELECT id, asset_id, user_id, content, tags, created_at, updated_at
		FROM asset_logs
		%s
		ORDER BY %s %s
		LIMIT @limit OFFSET @offset
	`, whereClause, params.SortBy, params.SortOrder)

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("list logs by asset: %w", err)
	}
	defer rows.Close()

	logs := make([]*model.AssetLog, 0)
	for rows.Next() {
		var log model.AssetLog
		err := rows.Scan(
			&log.ID,
			&log.AssetID,
			&log.UserID,
			&log.Content,
			&log.Tags,
			&log.CreatedAt,
			&log.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan log: %w", err)
		}
		logs = append(logs, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate logs: %w", err)
	}

	return logs, nil
}

// CountByAsset returns the total number of logs matching the filters for a specific asset
func (r *LogRepository) CountByAsset(ctx context.Context, userID string, assetID uuid.UUID, params *model.LogQueryParams) (int64, error) {
	// Build WHERE clause (same logic as ListByAsset)
	args := pgx.NamedArgs{
		"userID":  userID,
		"assetID": assetID,
	}
	whereClause := buildLogWhereClause(params, args)

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM asset_logs
		%s
	`, whereClause)

	var count int64
	err := r.db.QueryRow(ctx, query, args).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count logs by asset: %w", err)
	}

	return count, nil
}

// Create inserts a new log for an asset
// Returns NotFoundError if the asset doesn't exist or doesn't belong to the user
func (r *LogRepository) Create(ctx context.Context, userID string, assetID uuid.UUID, req *model.CreateLogRequest) (*model.AssetLog, error) {
	query := `
		INSERT INTO asset_logs (asset_id, user_id, content, tags)
		VALUES (@assetID, @userID, @content, @tags)
		RETURNING id, asset_id, user_id, content, tags, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"assetID": assetID,
		"userID":  userID,
		"content": req.Content,
		"tags":    req.Tags, // nil becomes NULL, []string{} becomes empty array
	}

	var log model.AssetLog
	err := r.db.QueryRow(ctx, query, args).Scan(
		&log.ID,
		&log.AssetID,
		&log.UserID,
		&log.Content,
		&log.Tags,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err != nil {
		// Check for foreign key violation (asset doesn't exist or doesn't belong to user)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign_key_violation
			return nil, errs.NewNotFoundError("asset not found", false, nil)
		}
		return nil, fmt.Errorf("create log: %w", err)
	}

	return &log, nil
}

// buildLogUpdateSetClause builds dynamic SET clause for Update
func buildLogUpdateSetClause(req *model.UpdateLogRequest, args pgx.NamedArgs) string {
	setClauses := []string{"updated_at = now()"}

	if req.Content != nil {
		setClauses = append(setClauses, "content = @content")
		args["content"] = *req.Content
	}

	if req.Tags != nil {
		setClauses = append(setClauses, "tags = @tags")
		args["tags"] = *req.Tags
	}

	return strings.Join(setClauses, ", ")
}

// Update modifies an existing log (only non-nil fields are updated)
func (r *LogRepository) Update(ctx context.Context, userID string, logID uuid.UUID, req *model.UpdateLogRequest) (*model.AssetLog, error) {
	// Build SET clause dynamically based on non-nil fields
	args := pgx.NamedArgs{
		"logID":  logID,
		"userID": userID,
	}
	setClause := buildLogUpdateSetClause(req, args)

	query := fmt.Sprintf(`
		UPDATE asset_logs
		SET %s
		WHERE id = @logID AND user_id = @userID
		RETURNING id, asset_id, user_id, content, tags, created_at, updated_at
	`, setClause)

	var log model.AssetLog
	err := r.db.QueryRow(ctx, query, args).Scan(
		&log.ID,
		&log.AssetID,
		&log.UserID,
		&log.Content,
		&log.Tags,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewNotFoundError("log not found", false, nil)
		}
		return nil, fmt.Errorf("update log: %w", err)
	}

	return &log, nil
}

// Delete removes a log for a user
func (r *LogRepository) Delete(ctx context.Context, userID string, logID uuid.UUID) error {
	query := `
		DELETE FROM asset_logs
		WHERE id = @logID AND user_id = @userID
	`

	args := pgx.NamedArgs{
		"logID":  logID,
		"userID": userID,
	}

	result, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("delete log: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewNotFoundError("log not found", false, nil)
	}

	return nil
}
