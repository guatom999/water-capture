package repositories

import (
	"context"

	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/stretchr/testify/mock"
)

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) GetSubscriptionsByLocationID(ctx context.Context, locationID int) ([]entities.NotificationSubscription, error) {
	args := m.Called(ctx, locationID)
	if v := args.Get(0); v != nil {
		if res, ok := v.([]entities.NotificationSubscription); ok {
			return res, args.Error(1)
		}
	}
	return nil, args.Error(1)
}

func (m *MockNotificationRepository) GetActiveSubscriptions(ctx context.Context) ([]entities.NotificationSubscription, error) {
	args := m.Called(ctx)
	if v := args.Get(0); v != nil {
		if res, ok := v.([]entities.NotificationSubscription); ok {
			return res, args.Error(1)
		}
	}
	return nil, args.Error(1)
}

func (m *MockNotificationRepository) CreateSubscription(ctx context.Context, sub *entities.NotificationSubscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockNotificationRepository) LogNotification(ctx context.Context, log *entities.NotificationLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}
