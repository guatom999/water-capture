package models

import (
	"database/sql"
	"time"
)

type Location struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"` // nullable
	Latitude    float64   `json:"latitude" db:"latitude"`
	Longitude   float64   `json:"longitude" db:"longitude"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateWaterLevelReq struct {
	LocationID int64     `json:"location_id"`
	LevelCm    float64   `json:"level_cm"`
	Image      string    `json:"image"`
	Danger     string    `json:"danger"`
	IsFlooded  bool      `json:"is_flooded"`
	MeasuredAt time.Time `json:"measured_at"`
	Note       string    `json:"note"`
}

type WaterLevel struct {
	ID         string    `json:"id"`
	ImageURL   string    `json:"image_url"`
	Level      float64   `json:"level"`
	Status     string    `json:"status"`
	DetectedAt time.Time `json:"detected_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type MapMarker struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Level       float64   `json:"level"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type MapMarkerResponse struct {
	Center  Center        `json:"center"`
	Zoom    int           `json:"zoom"`
	Markers []*MapMarkers `json:"markers"`
}

type Center struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MapMarkers struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

type PredictWater struct {
	FileName   string  `json:"file_name"`
	WaterLevel float64 `json:"water_level"`
}

type LocationWithWaterLevel struct {
	StationID           int64   `db:"station_id"`
	LocationName        string  `db:"location_name"`
	LocationDescription string  `db:"location_description"`
	ProvinceID          int     `db:"province_id"`
	Latitude            float64 `db:"latitude"`
	Longitude           float64 `db:"longitude"`
	IsActive            bool    `db:"is_active"`
	BankLevel           float64 `db:"bank_level"`

	WaterLevelID *int64     `db:"water_level_id"`
	LevelCm      *float64   `db:"level_cm"`
	Image        *string    `db:"image"`
	Danger       *string    `db:"danger"`
	IsFlooded    *bool      `db:"is_flooded"`
	MeasuredAt   *time.Time `db:"measured_at"`
	Note         *string    `db:"note"`
}

type LocationWithWaterLevelRes struct {
	StationID           int64    `json:"station_id"`
	LocationName        string   `json:"location_name"`
	LocationDescription string   `json:"location_description"`
	ProvinceID          int      `json:"province_id"`
	Latitude            float64  `json:"latitude"`
	Longitude           float64  `json:"longitude"`
	IsActive            bool     `json:"is_active"`
	BankLevel           float64  `json:"bank_level"`
	WaterLevelID        *int64   `json:"water_level_id"`
	LevelCm             *float64 `json:"level_cm"`
	Image               *string  `json:"image"`
	Danger              *string  `json:"danger"`
	IsFlooded           *bool    `json:"is_flooded"`
	MeasuredAt          string   `json:"measured_at"`
	Note                *string  `json:"note"`
	// Station             []int    `json:"station"`
	// StationName         string   `json:"station_name"`
	// StationLatitude     string   `json:"station_lat"`
	// StationLongtitude   string   `json:"station_long"`
}

type WaterLocationDetailRes struct {
	LocationID int64          `json:"location_id"`
	LevelCm    float64        `json:"level_cm"`
	Image      string         `json:"image"`
	Danger     string         `json:"danger"`
	IsFlooded  bool           `json:"is_flooded"`
	Source     sql.NullString `json:"source"`
	MeasuredAt time.Time      `json:"measured_at"`
	Note       string         `json:"note"`
}
