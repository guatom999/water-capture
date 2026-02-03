package entities

import (
	"database/sql"
	"time"
)

type Province struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Code      string    `db:"code" json:"code"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Location struct {
	ID          int64     `db:"id" json:"id"`
	StationID   int64     `db:"station_id" json:"station_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	ProvinceID  int       `db:"province_id" json:"province_id"`
	BankLevel   float64   `db:"bank_level" json:"bank_level"`
	Latitude    float64   `db:"latitude" json:"latitude"`
	Longitude   float64   `db:"longitude" json:"longitude"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type WaterLevel struct {
	ID                int64          `db:"id" json:"id"`
	StationID         int64          `db:"station_id" json:"station_id"`
	LevelCm           float64        `db:"level_cm" json:"level_cm"`
	Image             string         `db:"image" json:"image"`
	Danger            string         `db:"danger" json:"danger"`
	IsFlooded         bool           `db:"is_flooded" json:"is_flooded"`
	Source            sql.NullString `db:"source" json:"source"`
	MeasuredAt        time.Time      `db:"measured_at" json:"measured_at"`
	Note              string         `db:"note" json:"note"`
	Status            string         `db:"status"` // "active", "pending_deletion", "deleted"
	DeletedAt         sql.NullTime   `db:"deleted_at"`
	ScheduledDeleteAt sql.NullTime   `db:"scheduled_delete_at"`
}
