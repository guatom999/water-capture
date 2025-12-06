package models

import "time"

type WaterLevel struct {
	ID         string    `json:"id"`
	ImageURL   string    `json:"image_url"`
	Level      float64   `json:"level"`
	Status     string    `json:"status"`
	DetectedAt time.Time `json:"detected_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MapMarker struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Level      float64   `json:"level"`
	Icon       string    `json:"icon"`
	InfoWindow string    `json:"info_window"`
	Timestamp  time.Time `json:"timestamp"`
}

type MapMarkerResponse struct {
	Markers []MapMarker `json:"markers"`
	Center  Location    `json:"center"`
	Zoom    int         `json:"zoom"`
}

// // WaterLevelStatus constants
// const (
// 	StatusNormal  = "normal"
// 	StatusWarning = "warning"
// 	StatusDanger  = "danger"
// )

// // DetermineStatus calculates status based on water level
// func DetermineStatus(level float64) string {
// 	switch {
// 	case level >= 2.0:
// 		return StatusDanger
// 	case level >= 1.5:
// 		return StatusWarning
// 	default:
// 		return StatusNormal
// 	}
// }
