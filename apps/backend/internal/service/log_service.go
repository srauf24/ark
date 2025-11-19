package service

import (
	"context"
	"strings"

	"ark/internal/model"
	"ark/internal/repository"

	"github.com/google/uuid"
)

type LogService struct {
	logRepo   *repository.LogRepository
	assetRepo *repository.AssetRepository
}

func NewLogService(logRepo *repository.LogRepository, assetRepo *repository.AssetRepository) *LogService {
	return &LogService{
		logRepo:   logRepo,
		assetRepo: assetRepo,
	}
}

// processTags cleans up tags: trim, lowercase, deduplicate, and limit to 20
func processTags(tags []string) []string {
	if tags == nil {
		return nil
	}

	uniqueTags := make(map[string]bool)
	var processed []string

	for _, tag := range tags {
		cleanTag := strings.ToLower(strings.TrimSpace(tag))
		if cleanTag != "" && !uniqueTags[cleanTag] {
			uniqueTags[cleanTag] = true
			processed = append(processed, cleanTag)
		}
	}

	// Limit to 20 tags
	if len(processed) > 20 {
		return processed[:20]
	}

	return processed
}

func (s *LogService) ListByAsset(ctx context.Context, userID string, assetID uuid.UUID, params *model.LogQueryParams) (*model.LogListResponse, error) {
	// Verify asset ownership
	_, err := s.assetRepo.GetByID(ctx, userID, assetID)
	if err != nil {
		return nil, err
	}

	params.SetDefaults()

	logs, err := s.logRepo.ListByAsset(ctx, userID, assetID, params)
	if err != nil {
		return nil, err
	}

	total, err := s.logRepo.CountByAsset(ctx, userID, assetID, params)
	if err != nil {
		return nil, err
	}

	return model.NewLogListResponse(logs, total, params.Limit, params.Offset), nil
}

func (s *LogService) GetByID(ctx context.Context, userID string, logID uuid.UUID) (*model.LogResponse, error) {
	log, err := s.logRepo.GetByID(ctx, userID, logID)
	if err != nil {
		return nil, err
	}

	return model.NewLogResponse(log), nil
}

func (s *LogService) Create(ctx context.Context, userID string, assetID uuid.UUID, req *model.CreateLogRequest) (*model.LogResponse, error) {
	// Verify asset ownership
	_, err := s.assetRepo.GetByID(ctx, userID, assetID)
	if err != nil {
		return nil, err
	}

	// Process tags
	if req.Tags != nil {
		req.Tags = processTags(req.Tags)
	}

	log, err := s.logRepo.Create(ctx, userID, assetID, req)
	if err != nil {
		return nil, err
	}

	return model.NewLogResponse(log), nil
}

func (s *LogService) Update(ctx context.Context, userID string, logID uuid.UUID, req *model.UpdateLogRequest) (*model.LogResponse, error) {
	// Process tags if present
	if req.Tags != nil {
		processed := processTags(*req.Tags)
		req.Tags = &processed
	}

	log, err := s.logRepo.Update(ctx, userID, logID, req)
	if err != nil {
		return nil, err
	}

	return model.NewLogResponse(log), nil
}

func (s *LogService) Delete(ctx context.Context, userID string, logID uuid.UUID) error {
	return s.logRepo.Delete(ctx, userID, logID)
}
