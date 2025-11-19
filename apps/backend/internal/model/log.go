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
