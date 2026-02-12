package handlers

import (
	"net/http"

	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/guatom999/self-boardcast/internal/services"
	"github.com/labstack/echo/v4"
)

type (
	notificationHandler struct {
		notificationService services.NotificationServiceInterface
	}

	NotificationHandlerInterface interface {
		CreateSubscription(c echo.Context) error
	}
)

func NewNotificationHandler(notificationService services.NotificationServiceInterface) NotificationHandlerInterface {
	return &notificationHandler{notificationService: notificationService}
}

func (h *notificationHandler) CreateSubscription(c echo.Context) error {

	ctx := c.Request().Context()

	req := new(models.CreateSubscriptionRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := h.notificationService.CreateSubscription(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Subscription created successfully",
	})
}
