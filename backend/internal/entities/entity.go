package entities

import "time"

type WaterLevel struct {
	ID         int64     `db:"id" json:"id"`
	LocationID int64     `db:"location_id" json:"location_id"`
	LevelCm    float64   `db:"level_cm" json:"level_cm"`
	Image      string    `db:"image" json:"image"`
	Danger     string    `db:"danger" json:"danger"`
	IsFlooded  bool      `db:"is_flooded" json:"is_flooded"`
	Source     string    `db:"source" json:"source"`
	MeasuredAt time.Time `db:"measured_at" json:"measured_at"`
	Note       string    `db:"note" json:"note"`
}
