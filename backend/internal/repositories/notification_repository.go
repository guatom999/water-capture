package repositories

import (
	"context"

	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/jmoiron/sqlx"
)

type (
	notificationRepository struct {
		db *sqlx.DB
	}

	NotificationRepositoryInterface interface {
		GetSubscriptionsByLocationID(ctx context.Context, locationID int) ([]entities.NotificationSubscription, error)
		GetActiveSubscriptions(ctx context.Context) ([]entities.NotificationSubscription, error)
		CreateSubscription(ctx context.Context, sub *entities.NotificationSubscription) error
		LogNotification(ctx context.Context, log *entities.NotificationLog) error
	}
)

func NewNotificationRepository(db *sqlx.DB) NotificationRepositoryInterface {
	return &notificationRepository{
		db: db,
	}
}

// GetSubscriptionsByLocationID ดึง subscriptions ตาม location_id
func (r *notificationRepository) GetSubscriptionsByLocationID(ctx context.Context, locationID int) ([]entities.NotificationSubscription, error) {
	var subscriptions []entities.NotificationSubscription

	query := `
		SELECT id, user_id, location_id, channel, target, threshold_level, is_active, created_at, updated_at
		FROM notification_subscriptions
		WHERE (location_id = $1 OR location_id IS NULL)
		AND is_active = true
	`

	if err := r.db.SelectContext(ctx, &subscriptions, query, locationID); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// GetActiveSubscriptions ดึง subscriptions ที่ active ทั้งหมด
func (r *notificationRepository) GetActiveSubscriptions(ctx context.Context) ([]entities.NotificationSubscription, error) {
	var subscriptions []entities.NotificationSubscription

	query := `
		SELECT id, user_id, location_id, channel, target, threshold_level, is_active, created_at, updated_at
		FROM notification_subscriptions
		WHERE is_active = true
	`

	if err := r.db.SelectContext(ctx, &subscriptions, query); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// CreateSubscription สร้าง subscription ใหม่
func (r *notificationRepository) CreateSubscription(ctx context.Context, sub *entities.NotificationSubscription) error {
	query := `
		INSERT INTO notification_subscriptions (user_id, location_id, channel, target, threshold_level, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	return r.db.QueryRowContext(ctx, query,
		sub.UserID,
		sub.LocationID,
		sub.Channel,
		sub.Target,
		sub.ThresholdLevel,
		sub.IsActive,
	).Scan(&sub.ID)
}

// LogNotification บันทึก log การส่ง notification
func (r *notificationRepository) LogNotification(ctx context.Context, log *entities.NotificationLog) error {
	query := `
		INSERT INTO notification_logs (subscription_id, location_id, water_level, message, channel, status, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	return r.db.QueryRowContext(ctx, query,
		log.SubscriptionID,
		log.LocationID,
		log.WaterLevel,
		log.Message,
		log.Channel,
		log.Status,
		log.ErrorMessage,
	).Scan(&log.ID)
}
