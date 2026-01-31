package handlers

import "github.com/guatom999/self-boardcast/internal/services"

type (
	notificationHandler struct {
		notificationService services.NotificationServiceInterface
	}

	NotificationHandlerInterface interface {
	}
)

func NewNotificationHandler(notificationService services.NotificationServiceInterface) NotificationHandlerInterface {
	return &notificationHandler{notificationService: notificationService}
}
