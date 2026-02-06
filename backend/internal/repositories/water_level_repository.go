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
	CreateWaterLevel(ctx context.Context, req []*entities.WaterLevel) error
	CreateProvince(ctx context.Context, req []*entities.Province) error
	CreateStationLocation(ctx context.Context, req []*entities.Location) error
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
			l.station_id as station_id,
            l.name AS location_name,
            l.description AS location_description,
			l.province_id,
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
        LEFT JOIN water_levels wl ON l.station_id = wl.station_id
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

func (r *waterLevelRepository) GetByLocationID(ctx context.Context, stationID int) ([]*entities.WaterLevel, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	query := `SELECT * FROM water_levels WHERE station_id = $1`

	result := make([]*entities.WaterLevel, 0)
	if err := r.db.SelectContext(ctx, &result, query, stationID); err != nil {
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

func (r *waterLevelRepository) CreateProvince(ctx context.Context, req []*entities.Province) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if len(req) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("Error failed to begin transaction %v", err.Error())
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO provinces(name, code) VALUES ($1, $2)`

	for _, province := range req {
		_, err = tx.ExecContext(ctx, query, province.Name, province.Code)
		if err != nil {
			log.Printf("Error failed to insert into provinces database for province_id %d: %v", province.ID, err.Error())
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error failed to commit transaction %v", err.Error())
		return err
	}

	return nil
}

func (r *waterLevelRepository) CreateStationLocation(ctx context.Context, req []*entities.Location) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if len(req) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("Error failed to begin transaction %v", err.Error())
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO locations(station_id, name, description, latitude, longitude, is_active, bank_level, province_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, location := range req {
		_, err = tx.ExecContext(ctx, query, location.StationID, location.Name, location.Description, location.Latitude, location.Longitude, location.IsActive, location.BankLevel, location.ProvinceID)
		if err != nil {
			log.Printf("Error failed to insert into locations %v", err.Error())
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error failed to commit transaction %v", err.Error())
		return err
	}

	return nil
}

func (r *waterLevelRepository) CreateWaterLevel(ctx context.Context, req []*entities.WaterLevel) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if len(req) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("Error failed to begin transaction %v", err.Error())
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO water_levels(station_id, level_cm, image, danger, is_flooded, source, measured_at, note, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	for _, waterLevel := range req {
		_, err = tx.ExecContext(ctx, query, waterLevel.StationID, waterLevel.LevelCm, waterLevel.Image, waterLevel.Danger, waterLevel.IsFlooded, waterLevel.Source, waterLevel.MeasuredAt, waterLevel.Note, "ACTIVE")
		if err != nil {
			log.Printf("Error failed to insert into water_levels database %v", err.Error())
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error failed to commit transaction %v", err.Error())
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
