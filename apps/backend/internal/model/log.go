package model

import (
	"time"

	"github.com/google/uuid"
)

// AssetLog represents a configuration change or troubleshooting log for an asset
type AssetLog struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AssetID   uuid.UUID `json:"asset_id" db:"asset_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	Tags      []string  `json:"tags,omitempty" db:"tags"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateLogRequest is the DTO for creating a new log entry
type CreateLogRequest struct {
	Content string   `json:"content" validate:"required,min=2,max=10000"`
	Tags    []string `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
}

// UpdateLogRequest is the DTO for updating an existing log entry
type UpdateLogRequest struct {
	Content *string   `json:"content,omitempty" validate:"omitempty,min=2,max=10000"`
	Tags    *[]string `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
}

// LogResponse is the DTO for single log responses
type LogResponse struct {
	ID        uuid.UUID `json:"id"`
	AssetID   uuid.UUID `json:"asset_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewLogResponse converts an AssetLog domain model to LogResponse DTO
func NewLogResponse(log *AssetLog) *LogResponse {
	if log == nil {
		return nil
	}

	return &LogResponse{
		ID:        log.ID,
		AssetID:   log.AssetID,
		UserID:    log.UserID,
		Content:   log.Content,
		Tags:      log.Tags,
		CreatedAt: log.CreatedAt,
		UpdatedAt: log.UpdatedAt,
	}
}

// LogListResponse is the DTO for paginated log list responses
type LogListResponse struct {
	Logs   []LogResponse `json:"logs"`
	Total  int64         `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

// NewLogListResponse converts a slice of AssetLogs to LogListResponse with pagination metadata
func NewLogListResponse(logs []*AssetLog, total int64, limit, offset int) *LogListResponse {
	responses := make([]LogResponse, 0, len(logs))
	for _, log := range logs {
		if resp := NewLogResponse(log); resp != nil {
			responses = append(responses, *resp)
		}
	}

	return &LogListResponse{
		Logs:   responses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}

// LogQueryParams represents query parameters for listing logs
type LogQueryParams struct {
	Limit     int        `query:"limit" validate:"omitempty,min=1,max=200"`
	Offset    int        `query:"offset" validate:"omitempty,min=0"`
	Tags      []string   `query:"tags" validate:"omitempty,dive,max=50"`
	Search    *string    `query:"search" validate:"omitempty,max=100"`
	StartDate *time.Time `query:"start_date"`
	EndDate   *time.Time `query:"end_date"`
	SortBy    string     `query:"sort_by" validate:"omitempty,oneof=created_at updated_at"`
	SortOrder string     `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// SetDefaults sets default values for LogQueryParams
func (q *LogQueryParams) SetDefaults() {
	if q.Limit == 0 {
		q.Limit = DefaultLogLimit
	}
	if q.Limit > MaxLogLimit {
		q.Limit = MaxLogLimit
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
}
