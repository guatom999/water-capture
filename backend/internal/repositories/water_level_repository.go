package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/jmoiron/sqlx"
)

type waterLevelRepository struct {
	db *sqlx.DB
}

// WaterLevelRepository interface
type WaterLevelRepositoryInterface interface {
	Create(ctx context.Context, waterLevel *models.WaterLevel) error
	GetLatest(ctx context.Context) (*models.WaterLevel, error)
	GetAll(ctx context.Context, limit int) (*models.MapMarkerResponse, error)
	GetByID(ctx context.Context, id string) (*models.WaterLevel, error)
}

func NewWaterLevelRepository(db *sqlx.DB) WaterLevelRepositoryInterface {
	return &waterLevelRepository{
		db: db,
	}
}

func (r *waterLevelRepository) Create(ctx context.Context, waterLevel *models.WaterLevel) error {
	return nil
}
func (r *waterLevelRepository) GetLatest(ctx context.Context) (*models.WaterLevel, error) {
	return nil, nil
}
func (r *waterLevelRepository) GetAll(pctx context.Context, limit int) (*models.MapMarkerResponse, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*20)
	defer cancel()

	query := "SELECT * FROM locations "

	result := make([]*models.Location, 0)
	locations := make([]*models.MapMarkers, 0)

	if err := r.db.SelectContext(ctx, &result, query); err != nil {
		return nil, err
	}

	for _, location := range result {

		locations = append(locations, &models.MapMarkers{
			ID:          fmt.Sprintf("%d", location.ID),
			Title:       location.Name,
			Latitude:    location.Latitude,
			Longitude:   location.Longitude,
			Description: location.Description,
		})

	}

	mapMarkers := &models.MapMarkerResponse{
		Markers: locations,
		Center:  models.Center{Latitude: 13.7563, Longitude: 100.5018},
		Zoom:    6,
	}

	return mapMarkers, nil
}
func (r *waterLevelRepository) GetByID(ctx context.Context, id string) (*models.WaterLevel, error) {
	return nil, nil
}
