package services

import (
	"context"
	"errors"
	"log"

	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/guatom999/self-boardcast/internal/repositories"
)

type (
	notificationService struct {
		authRepo repositories.AuthRepositoryInterface
		notiRepo repositories.NotificationRepositoryInterface
	}

	NotificationServiceInterface interface {
		// SendNotification(ctx context.Context, notification *entities.Notification) error
		CreateSubscription(ctx context.Context, subscription *models.CreateSubscriptionRequest) error
	}
)

func NewNotificationService(authRepo repositories.AuthRepositoryInterface, notiRepo repositories.NotificationRepositoryInterface) NotificationServiceInterface {
	return &notificationService{authRepo: authRepo, notiRepo: notiRepo}
}

func (s *notificationService) CreateSubscription(ctx context.Context, subscription *models.CreateSubscriptionRequest) error {

	if s.notiRepo.IsUserAlreaySubscript(ctx, subscription.UserId) {
		return errors.New("user already has an active subscription")
	}

	user := &entities.User{}

	user, err := s.authRepo.GetUserByEmail(ctx, subscription.UserEmail)
	if err != nil {
		log.Printf("Error: user not found %s", err.Error())
		return err
	}

	if err := s.notiRepo.CreateSubscription(ctx, &entities.NotificationSubscription{
		UserID:     user.ID,
		LocationID: subscription.LocationID,
		Channel:    subscription.Channel,
		ProvinceID: subscription.ProvinceID,
		IsActive:   true,
	}); err != nil {
		log.Printf("Error: create subscription failed%s", err.Error())
		return err
	}

	return nil
}
