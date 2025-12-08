package models

import "time"

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

type WaterLevel struct {
	ID         string    `json:"id"`
	ImageURL   string    `json:"image_url"`
	Level      float64   `json:"level"`
	Status     string    `json:"status"`
	DetectedAt time.Time `json:"detected_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// type C struct {
// 	Name      string  `json:"name"`
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// }

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
