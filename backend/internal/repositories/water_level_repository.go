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
	// GetLatest(ctx context.Context) (*models.WaterLevel, error)
	GetAll(ctx context.Context, limit int) ([]models.LocationWithWaterLevel, error)
	GetByLocationID(ctx context.Context, locationID int) ([]*entities.WaterLevel, error)
	CreateWaterLevel(ctx context.Context, req *entities.WaterLevel) error
	// DeleteOldestWaterLevels(ctx context.Context, locationID int, keepLatest int) error

	MarkForDeletion(ctx context.Context, id int64, scheduledAt time.Time) error
	GetPendingDeletions(ctx context.Context) ([]*entities.WaterLevel, error)

	// // Phase 2: Hard delete methods
	HardDelete(ctx context.Context, id int64) error
	// HardDeleteBatch(ctx context.Context, ids []int64) error

	// // Recovery methods
	// CancelDeletion(ctx context.Context, id int64) error
	// GetFailedDeletions(ctx context.Context) ([]*entities.WaterLevel, error)
}

func NewWaterLevelRepository(db *sqlx.DB) WaterLevelRepositoryInterface {
	return &waterLevelRepository{
		db: db,
	}
}

//	func (r *waterLevelRepository) GetLatest(ctx context.Context) (*models.WaterLevel, error) {
//		return nil, nil
//	}
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
			l.bank_level,
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
		log.Printf("Error failed to select from locations database %v", err.Error())
		return nil, err
	}

	return result, nil

}

func (r *waterLevelRepository) GetByLocationID(ctx context.Context, locationID int) ([]*entities.WaterLevel, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	query := `SELECT * FROM water_levels WHERE location_id = $1`

	result := make([]*entities.WaterLevel, 0)
	if err := r.db.SelectContext(ctx, &result, query, locationID); err != nil {
		log.Printf("Error failed to select from water_levels database %v", err.Error())
		return nil, err
	}

	return result, nil
}

func (r *waterLevelRepository) GetLastWaterLevelWithLimit(ctx context.Context, locationID int, limit int) (*entities.WaterLevel, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	query := `SELECT * FROM water_levels WHERE location_id = $1 ORDER BY measured_at DESC LIMIT $2`

	result := &entities.WaterLevel{}
	if err := r.db.GetContext(ctx, result, query, locationID); err != nil {
		log.Printf("Error failed to select from water_levels database %v", err.Error())
		return nil, err
	}

	return result, nil
}

func (r *waterLevelRepository) CreateWaterLevel(ctx context.Context, req *entities.WaterLevel) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `INSERT INTO water_levels(location_id, level_cm, image, danger, is_flooded, source, measured_at, note, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query, req.LocationID, req.LevelCm, req.Image, req.Danger, req.IsFlooded, req.Source, req.MeasuredAt, req.Note, "ACTIVE")
	if err != nil {
		log.Printf("Error failed to insert into water_levels database %v", err.Error())
		return err
	}

	return nil
}

// func (r *waterLevelRepository) DeleteOldestWaterLevels(ctx context.Context, locationID int, keepLatest int) error {

// 	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
// 	defer cancel()

// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		log.Printf("Error failed to begin transaction %v", err.Error())
// 		return fmt.Errorf("error failed to begin transaction %v", err.Error())
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	query := `
//         DELETE FROM water_levels
//         WHERE id IN (
//             SELECT id
//             FROM water_levels
//             WHERE location_id = $1
//             ORDER BY measured_at DESC
//             OFFSET $2
//         )
//     `

// 	_, err = tx.ExecContext(ctx, query, locationID, keepLatest)
// 	if err != nil {
// 		log.Printf("Error failed to delete from water_levels database %v", err.Error())
// 		return fmt.Errorf("error failed to delete from water_levels database %v", err.Error())
// 	}

// 	if err = tx.Commit(); err != nil {
// 		log.Printf("Error failed to commit transaction %v", err.Error())
// 		return fmt.Errorf("error failed to commit transaction %v", err.Error())
// 	}

// 	return nil
// }

func (r *waterLevelRepository) MarkForDeletion(ctx context.Context, id int64, scheduledAt time.Time) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `
	    UPDATE water_levels 
	    SET status = 'PENDING_DELETION',
	        scheduled_delete_at = $1
        WHERE location_id = $2 
          	AND id IN (
              SELECT id FROM (
                  SELECT id 
                  FROM water_levels 
                  WHERE location_id = $2
                  ORDER BY measured_at DESC
				  OFFSET 5
              ) AS keep_records
          )
	`

	_, err := r.db.ExecContext(ctx, query, scheduledAt, id)
	if err != nil {
		log.Printf("Error failed to update water_levels database %v", err.Error())
		return err
	}

	return nil
}

func (r *waterLevelRepository) GetPendingDeletions(ctx context.Context) ([]*entities.WaterLevel, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `SELECT * FROM water_levels WHERE status = 'PENDING_DELETION'`

	result := make([]*entities.WaterLevel, 0)
	if err := r.db.SelectContext(ctx, &result, query); err != nil {
		log.Printf("Error failed to select from water_levels database %v", err.Error())
		return nil, err
	}

	return result, nil
}

func (r *waterLevelRepository) HardDelete(ctx context.Context, id int64) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `DELETE FROM water_levels WHERE id = $1 AND status = 'PENDING_DELETION'`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Error failed to delete from water_levels database %v", err.Error())
		return err
	}

	return nil
}

func (r *waterLevelRepository) CancelDeletion(ctx context.Context, id int64) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `UPDATE water_levels SET status = 'ACTIVE', scheduled_delete_at = NULL, updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Error failed to update water_levels database %v", err.Error())
		return err
	}

	return nil
}

// func (r *waterLevelRepository) DeleteOldestWaterLevels(ctx context.Context, locationID int) error {

// 	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
// 	defer cancel()

// 	query := `DELETE FROM water_levels WHERE id IN (
// 		SELECT id FROM water_levels
// 		WHERE location_id = $1
// 		ORDER BY measured_at ASC
// 		LIMIT 1
// 	)`

// 	_, err := r.db.ExecContext(ctx, query, locationID)
// 	if err != nil {
// 		log.Printf("Error failed to delete from water_levels database %v", err.Error())
// 		return err
// 	}

// 	return nil
// }
