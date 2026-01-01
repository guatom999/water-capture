package repositories

import (
	"context"
	"log"
	"time"

	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/jmoiron/sqlx"
)

type waterLevelRepository struct {
	db *sqlx.DB
}

// WaterLevelRepository interface
type WaterLevelRepositoryInterface interface {
	GetLatest(ctx context.Context) (*models.WaterLevel, error)
	GetAll(ctx context.Context, limit int) ([]models.LocationWithWaterLevel, error)
	GetByID(ctx context.Context, id string) (*models.WaterLevel, error)
	CreateWaterLevel(ctx context.Context, req *entities.WaterLevel) error
	DeleteWaterLevels(ctx context.Context, locationID int) error
}

func NewWaterLevelRepository(db *sqlx.DB) WaterLevelRepositoryInterface {
	return &waterLevelRepository{
		db: db,
	}
}

func (r *waterLevelRepository) GetLatest(ctx context.Context) (*models.WaterLevel, error) {
	return nil, nil
}
func (r *waterLevelRepository) GetAll(pctx context.Context, limit int) ([]models.LocationWithWaterLevel, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*20)
	defer cancel()

	query := `
        SELECT DISTINCT ON (l.id)
            l.id AS location_id,
            l.name AS location_name,
            l.description AS location_description,
            l.latitude,
            l.longitude,
            l.is_active,
            wl.id AS water_level_id,
            wl.level_cm,
			wl.image,
            wl.danger,
            wl.is_flooded,
            wl.measured_at,
            wl.note
        FROM locations l
        LEFT JOIN water_levels wl ON l.id = wl.location_id
        WHERE l.is_active = TRUE
        ORDER BY l.id, wl.measured_at DESC NULLS LAST
    `

	result := make([]models.LocationWithWaterLevel, 0)

	if err := r.db.SelectContext(ctx, &result, query); err != nil {
		return nil, err
	}

	return result, nil

}

func (r *waterLevelRepository) GetByID(ctx context.Context, id string) (*models.WaterLevel, error) {
	return nil, nil
}

func (r *waterLevelRepository) CreateWaterLevel(ctx context.Context, req *entities.WaterLevel) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `INSERT INTO water_levels(location_id, level_cm, image, danger, is_flooded, measured_at, note) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query, req.LocationID, req.LevelCm, req.Image, req.Danger, req.IsFlooded, req.MeasuredAt, req.Note)
	if err != nil {
		log.Printf("Error failed to insert into water_levels database %v", err.Error())
		return err
	}

	return nil
}

func (r *waterLevelRepository) DeleteWaterLevels(ctx context.Context, locationID int) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `DELETE FROM water_levels WHERE id IN (
		SELECT id FROM water_levels 
		WHERE location_id = $1 
		ORDER BY measured_at ASC 
		LIMIT 1
	)`

	_, err := r.db.ExecContext(ctx, query, locationID)
	if err != nil {
		log.Printf("Error failed to delete from water_levels database %v", err.Error())
		return err
	}

	return nil
}
