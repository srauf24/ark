package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ark/internal/errs"
	"ark/internal/model"
)

// AssetRepository provides data access methods for assets table.
// All methods enforce user isolation - assets are scoped to the requesting user.
type AssetRepository struct {
	db *pgxpool.Pool
}

// NewAssetRepository creates a new AssetRepository with the given database pool.
func NewAssetRepository(db *pgxpool.Pool) *AssetRepository {
	return &AssetRepository{db: db}
}

// GetByID retrieves a single asset by ID for the specified user.
// Returns NotFoundError if the asset doesn't exist or belongs to another user.
// This dual-key lookup (id AND user_id) prevents unauthorized access.
func (r *AssetRepository) GetByID(ctx context.Context, userID string, assetID uuid.UUID) (*model.Asset, error) {
	query := `
		SELECT id, user_id, name, type, hostname, metadata, created_at, updated_at
		FROM assets
		WHERE id = $1 AND user_id = $2
	`

	var asset model.Asset
	err := r.db.QueryRow(ctx, query, assetID, userID).Scan(
		&asset.ID,
		&asset.UserID,
		&asset.Name,
		&asset.Type,     // pointer - handles NULL
		&asset.Hostname, // pointer - handles NULL
		&asset.Metadata, // json.RawMessage - handles NULL
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewNotFoundError("asset not found", false, nil)
		}
		return nil, fmt.Errorf("get asset by id: %w", err)
	}

	return &asset, nil
}

// buildAssetWhereClause builds dynamic WHERE clause for List/Count with filters
func buildAssetWhereClause(params *model.AssetQueryParams, args *[]any, paramNum *int) string {
	clauses := []string{"user_id = $1"}

	if params.Type != nil {
		clauses = append(clauses, fmt.Sprintf("type = $%d", *paramNum))
		*args = append(*args, *params.Type)
		*paramNum++
	}

	if params.Search != nil {
		searchPattern := "%" + *params.Search + "%"
		clauses = append(clauses, fmt.Sprintf("(name ILIKE $%d OR hostname ILIKE $%d)", *paramNum, *paramNum+1))
		*args = append(*args, searchPattern, searchPattern)
		*paramNum += 2
	}

	return "WHERE " + strings.Join(clauses, " AND ")
}

// validateAssetSortBy prevents SQL injection by validating sort column
func validateAssetSortBy(sortBy string) error {
	allowed := map[string]bool{
		"name":       true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowed[sortBy] {
		return fmt.Errorf("invalid sort_by: %s", sortBy)
	}
	return nil
}

// validateSortOrder prevents SQL injection by validating sort direction
func validateSortOrder(sortOrder string) error {
	if sortOrder != "asc" && sortOrder != "desc" {
		return fmt.Errorf("invalid sort_order: %s", sortOrder)
	}
	return nil
}

// List retrieves assets for a user with optional filtering and pagination
func (r *AssetRepository) List(ctx context.Context, userID string, params *model.AssetQueryParams) ([]*model.Asset, error) {
	// Validate sort parameters to prevent SQL injection
	if err := validateAssetSortBy(params.SortBy); err != nil {
		return nil, err
	}
	if err := validateSortOrder(params.SortOrder); err != nil {
		return nil, err
	}

	// Build WHERE clause
	args := []any{userID}
	paramNum := 2
	whereClause := buildAssetWhereClause(params, &args, &paramNum)

	// Build complete query with ORDER BY and LIMIT/OFFSET
	query := fmt.Sprintf(`
		SELECT id, user_id, name, type, hostname, metadata, created_at, updated_at
		FROM assets
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, params.SortBy, params.SortOrder, paramNum, paramNum+1)

	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list assets: %w", err)
	}
	defer rows.Close()

	assets := make([]*model.Asset, 0)
	for rows.Next() {
		var asset model.Asset
		err := rows.Scan(
			&asset.ID,
			&asset.UserID,
			&asset.Name,
			&asset.Type,
			&asset.Hostname,
			&asset.Metadata,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan asset: %w", err)
		}
		assets = append(assets, &asset)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate assets: %w", err)
	}

	return assets, nil
}

// Count returns the total number of assets matching the filters
func (r *AssetRepository) Count(ctx context.Context, userID string, params *model.AssetQueryParams) (int64, error) {
	// Build WHERE clause (same logic as List)
	args := []any{userID}
	paramNum := 2
	whereClause := buildAssetWhereClause(params, &args, &paramNum)

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM assets
		%s
	`, whereClause)

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count assets: %w", err)
	}

	return count, nil
}

// Create inserts a new asset for a user
func (r *AssetRepository) Create(ctx context.Context, userID string, req *model.CreateAssetRequest) (*model.Asset, error) {
	query := `
		INSERT INTO assets (user_id, name, type, hostname, metadata)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, type, hostname, metadata, created_at, updated_at
	`

	var asset model.Asset
	err := r.db.QueryRow(ctx, query,
		userID,
		req.Name,
		req.Type,     // nil becomes NULL
		req.Hostname, // nil becomes NULL
		req.Metadata, // nil becomes NULL
	).Scan(
		&asset.ID,
		&asset.UserID,
		&asset.Name,
		&asset.Type,
		&asset.Hostname,
		&asset.Metadata,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create asset: %w", err)
	}

	return &asset, nil
}

// buildAssetUpdateSetClause builds dynamic SET clause for Update
func buildAssetUpdateSetClause(req *model.UpdateAssetRequest, args *[]any, paramNum *int) string {
	setClauses := []string{"updated_at = now()"}

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", *paramNum))
		*args = append(*args, *req.Name)
		*paramNum++
	}

	if req.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", *paramNum))
		*args = append(*args, *req.Type)
		*paramNum++
	}

	if req.Hostname != nil {
		setClauses = append(setClauses, fmt.Sprintf("hostname = $%d", *paramNum))
		*args = append(*args, *req.Hostname)
		*paramNum++
	}

	if req.Metadata != nil {
		setClauses = append(setClauses, fmt.Sprintf("metadata = $%d", *paramNum))
		*args = append(*args, *req.Metadata)
		*paramNum++
	}

	return strings.Join(setClauses, ", ")
}

// Update modifies an existing asset (only non-nil fields are updated)
func (r *AssetRepository) Update(ctx context.Context, userID string, assetID uuid.UUID, req *model.UpdateAssetRequest) (*model.Asset, error) {
	// Build SET clause dynamically based on non-nil fields
	args := []any{assetID, userID}
	paramNum := 3
	setClause := buildAssetUpdateSetClause(req, &args, &paramNum)

	query := fmt.Sprintf(`
		UPDATE assets
		SET %s
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, name, type, hostname, metadata, created_at, updated_at
	`, setClause)

	var asset model.Asset
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&asset.ID,
		&asset.UserID,
		&asset.Name,
		&asset.Type,
		&asset.Hostname,
		&asset.Metadata,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewNotFoundError("asset not found", false, nil)
		}
		return nil, fmt.Errorf("update asset: %w", err)
	}

	return &asset, nil
}

// Delete removes an asset for a user
func (r *AssetRepository) Delete(ctx context.Context, userID string, assetID uuid.UUID) error {
	query := `
		DELETE FROM assets
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(ctx, query, assetID, userID)
	if err != nil {
		return fmt.Errorf("delete asset: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewNotFoundError("asset not found", false, nil)
	}

	return nil
}
