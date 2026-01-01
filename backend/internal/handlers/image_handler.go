package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/guatom999/self-boardcast/internal/config"
	"github.com/guatom999/self-boardcast/internal/utils"
	"github.com/labstack/echo/v4"
)

type ImageHandler struct {
	cfg *config.Config
}

func NewImageHandler(cfg *config.Config) *ImageHandler {
	return &ImageHandler{
		cfg: cfg,
	}
}

// ServeImage serves an image file with security validations
func (h *ImageHandler) ServeImage(c echo.Context) error {
	filename := c.Param("filename")

	// Validate filename to prevent path traversal
	if err := utils.ValidateImagePath(filename); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid filename",
		})
	}

	// Get safe file path
	filepath, err := utils.GetSafeFilePath(h.cfg.App.UploadDir, filename)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid file path",
		})
	}

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Image not found",
		})
	}

	// Validate file type
	if err := utils.ValidateImageType(filepath); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Invalid file type",
		})
	}

	// Set cache headers for better performance
	c.Response().Header().Set("Cache-Control", "public, max-age=86400") // 24 hours

	// Serve the file
	return c.File(filepath)
}

func (h *ImageHandler) HealthCheck(c echo.Context) error {

	fmt.Println("h.cfg.App.UploadDir", h.cfg.App.UploadDir)

	if _, err := os.Stat(h.cfg.App.UploadDir); os.IsNotExist(err) {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status":  "unhealthy",
			"message": fmt.Sprint("Upload directory not found"),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
		// "upload_dir": h.cfg.App.UploadDir,
	})
}
