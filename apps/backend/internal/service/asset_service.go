package service

import (
	"context"

	"ark/internal/model"
	"ark/internal/repository"
	"ark/internal/validation"

	"github.com/google/uuid"
)

type AssetService struct {
	repo *repository.AssetRepository
}

func NewAssetService(repo *repository.AssetRepository) *AssetService {
	return &AssetService{
		repo: repo,
	}
}

func (s *AssetService) List(ctx context.Context, userID string, params *model.AssetQueryParams) (*model.AssetListResponse, error) {
	params.SetDefaults()

	assets, err := s.repo.List(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	return model.NewAssetListResponse(assets, total, params.Limit, params.Offset), nil
}

func (s *AssetService) GetByID(ctx context.Context, userID string, assetID uuid.UUID) (*model.AssetResponse, error) {
	asset, err := s.repo.GetByID(ctx, userID, assetID)
	if err != nil {
		return nil, err
	}

	return model.NewAssetResponse(asset), nil
}

func (s *AssetService) Create(ctx context.Context, userID string, req *model.CreateAssetRequest) (*model.AssetResponse, error) {
	// Business Validation
	if err := validation.ValidateAssetType(req.Type); err != nil {
		return nil, err
	}

	if err := validation.ValidateMetadataJSON(req.Metadata); err != nil {
		return nil, err
	}

	asset, err := s.repo.Create(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return model.NewAssetResponse(asset), nil
}

func (s *AssetService) Update(ctx context.Context, userID string, assetID uuid.UUID, req *model.UpdateAssetRequest) (*model.AssetResponse, error) {
	// Business Validation
	if err := validation.ValidateAssetType(req.Type); err != nil {
		return nil, err
	}

	if err := validation.ValidateMetadataJSON(req.Metadata); err != nil {
		return nil, err
	}

	asset, err := s.repo.Update(ctx, userID, assetID, req)
	if err != nil {
		return nil, err
	}

	return model.NewAssetResponse(asset), nil
}

func (s *AssetService) Delete(ctx context.Context, userID string, assetID uuid.UUID) error {
	return s.repo.Delete(ctx, userID, assetID)
}
