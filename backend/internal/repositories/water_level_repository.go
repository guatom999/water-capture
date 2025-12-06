package repositories

import (
	"context"

	"github.com/guatom999/self-boardcast/internal/models"
)

type waterLevelRepository struct {
}

// WaterLevelRepository interface
type WaterLevelRepositoryInterface interface {
	Create(ctx context.Context, waterLevel *models.WaterLevel) error
	GetLatest(ctx context.Context) (*models.WaterLevel, error)
	GetAll(ctx context.Context, limit int) ([]*models.WaterLevel, error)
	GetByID(ctx context.Context, id string) (*models.WaterLevel, error)
}

func NewWaterLevelRepository() WaterLevelRepositoryInterface {
	return &waterLevelRepository{}
}

func (r *waterLevelRepository) Create(ctx context.Context, waterLevel *models.WaterLevel) error {
	return nil
}
func (r *waterLevelRepository) GetLatest(ctx context.Context) (*models.WaterLevel, error) {
	return nil, nil
}
func (r *waterLevelRepository) GetAll(ctx context.Context, limit int) ([]*models.WaterLevel, error) {
	return nil, nil
}
func (r *waterLevelRepository) GetByID(ctx context.Context, id string) (*models.WaterLevel, error) {
	return nil, nil
}
