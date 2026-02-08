package handlers

import (
	"context"
	"net/http"

	"github.com/guatom999/self-boardcast/internal/services"
	"github.com/labstack/echo/v4"
)

// WaterLevelHandler handles HTTP requests
type waterLevelHandler struct {
	service services.WaterLevelServiceInterface
}

type WaterLevelHandlerInterface interface {
	GetMapMarkers(c echo.Context) error
	GetSectionDetail(c echo.Context) error
}

func NewMapHandler(service services.WaterLevelServiceInterface) WaterLevelHandlerInterface {
	return &waterLevelHandler{
		service: service,
	}
}

func (h *waterLevelHandler) GetMapMarkers(c echo.Context) error {

	ctx := context.Background()

	markers, err := h.service.GetAllLocations(ctx, 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"markers": markers,
	})

}

func (h *waterLevelHandler) GetSectionDetail(c echo.Context) error {

	ctx := context.Background()

	stationID := c.QueryParam("station_id")

	markers, err := h.service.GetByLocationID(ctx, stationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"markers": markers,
	})

}
