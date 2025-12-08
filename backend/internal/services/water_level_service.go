package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/guatom999/self-boardcast/internal/repositories"
)

// WaterLevelService handles business logic
type waterLevelService struct {
	repo repositories.WaterLevelRepositoryInterface
}

type WaterLevelServiceInterface interface {
	ProcessImage(ctx context.Context, imageURL string) (*models.WaterLevel, error)
	GetByID(ctx context.Context, id string) (*models.WaterLevel, error)
	GetAllLocations(ctx context.Context, limit int) (*models.MapMarkerResponse, error)
}

// NewWaterLevelService creates new service
func NewWaterLevelService(repo repositories.WaterLevelRepositoryInterface) WaterLevelServiceInterface {
	return &waterLevelService{
		repo: repo,
	}
}

func (s *waterLevelService) GetAllLocations(ctx context.Context, limit int) (*models.MapMarkerResponse, error) {
	return s.repo.GetAll(ctx, limit)
}

func (s *waterLevelService) ProcessImage(ctx context.Context, imageURL string) (*models.WaterLevel, error) {
	// TODO: Call Python service to analyze image
	// Mock water level for now
	level := 1.2

	waterLevel := &models.WaterLevel{
		ID:       uuid.New().String(),
		ImageURL: imageURL,
		Level:    level,
		// Status:     models.DetermineStatus(level),
		DetectedAt: time.Now(),
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, waterLevel); err != nil {
		return nil, err
	}

	return waterLevel, nil
}

// GetByID returns water level by ID
func (s *waterLevelService) GetByID(ctx context.Context, id string) (*models.WaterLevel, error) {
	return s.repo.GetByID(ctx, id)
}
