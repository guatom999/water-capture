package services

import (
	"github.com/guatom999/self-boardcast/internal/repositories"
)

type (
	notificationService struct {
		notiRepo repositories.NotificationRepositoryInterface
	}

	NotificationServiceInterface interface {
		// SendNotification(ctx context.Context, notification *entities.Notification) error
	}
)

func NewNotificationService(notiRepo repositories.NotificationRepositoryInterface) NotificationServiceInterface {
	return &notificationService{notiRepo: notiRepo}
}
