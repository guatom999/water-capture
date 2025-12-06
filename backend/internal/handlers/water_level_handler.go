package handlers

import (
	"net/http"

	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/guatom999/self-boardcast/internal/services"
	"github.com/labstack/echo/v4"
)

// WaterLevelHandler handles HTTP requests
type WaterLevelHandler struct {
	service services.WaterLevelServiceInterface
}

func NewMapHandler(service services.WaterLevelServiceInterface) *WaterLevelHandler {
	return &WaterLevelHandler{
		service: service,
	}
}

func (h *WaterLevelHandler) GetMapMarkers(c echo.Context) error {

	markers := []models.MapMarker{
		{
			ID:         "marker-1",
			Title:      "ปตท. สยาม",
			Latitude:   13.7465,
			Longitude:  100.5320,
			Level:      1.85,
			Icon:       "http://maps.google.com/mapfiles/ms/icons/yellow-dot.png",
			InfoWindow: "ปตท. สยาม\nLevel: 1.85 m\nStatus: warning",
		},
		{
			ID:         "marker-2",
			Title:      "สะพานเพชรบุรี",
			Latitude:   13.7563,
			Longitude:  100.5018,
			Level:      2.15,
			Icon:       "http://maps.google.com/mapfiles/ms/icons/red-dot.png",
			InfoWindow: "สะพานเพชรบุรี\nLevel: 2.15 m\nStatus: danger",
		},
		{
			ID:         "marker-3",
			Title:      "สะพานพระราม 4",
			Latitude:   13.7311,
			Longitude:  100.5642,
			Level:      1.20,
			Icon:       "http://maps.google.com/mapfiles/ms/icons/green-dot.png",
			InfoWindow: "สะพานพระราม 4\nLevel: 1.20 m\nStatus: normal",
		},
		{
			ID:         "marker-4",
			Title:      "ข่วง ชั้นนอก",
			Latitude:   13.7880,
			Longitude:  100.4948,
			Level:      1.65,
			Icon:       "http://maps.google.com/mapfiles/ms/icons/yellow-dot.png",
			InfoWindow: "ข่วง ชั้นนอก\nLevel: 1.65 m\nStatus: warning",
		},
		{
			ID:         "marker-5",
			Title:      "โรงเรียนสกลนคร",
			Latitude:   13.7219,
			Longitude:  100.5534,
			Level:      0.95,
			Icon:       "http://maps.google.com/mapfiles/ms/icons/green-dot.png",
			InfoWindow: "โรงเรียนสกลนคร\nLevel: 0.95 m\nStatus: normal",
		},
	}

	return c.JSON(http.StatusOK, map[string]any{
		"markers": markers,
		"center": map[string]any{
			"latitude":  13.7563,
			"longitude": 100.5018,
			"name":      "Bangkok",
		},
		"zoom": 12,
	})

}
