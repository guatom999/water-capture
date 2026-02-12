package repositories

import (
	"context"
	"time"

	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/jmoiron/sqlx"
)

type (
	notificationRepository struct {
		db *sqlx.DB
	}

	NotificationRepositoryInterface interface {
		GetUserSubscriptions(ctx context.Context, userID int64) (*entities.NotificationSubscription, error)
		GetSubscriptionsByLocationID(ctx context.Context, locationID int) ([]entities.NotificationSubscription, error)
		GetActiveSubscriptions(ctx context.Context) ([]entities.NotificationSubscription, error)
		CreateSubscription(ctx context.Context, sub *entities.NotificationSubscription) error
		LogNotification(ctx context.Context, log *entities.NotificationLog) error
		// IsSubscriptionExist(ctx context.Context, userID int64, locationID int) (bool, error)
		IsUserAlreaySubscript(ctx context.Context, userID int64) bool
	}
)

func NewNotificationRepository(db *sqlx.DB) NotificationRepositoryInterface {
	return &notificationRepository{
		db: db,
	}
}

func (r *notificationRepository) GetUserSubscriptions(ctx context.Context, userID int64) (*entities.NotificationSubscription, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	result := entities.NotificationSubscription{}

	query := `
		SELECT id, user_id, location_id, channel, province_id, is_active, created_at, updated_at
		FROM notification_subscriptions
		WHERE user_id = $1
	`

	if err := r.db.GetContext(ctx, &result, query, userID); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSubscriptionsByLocationID ดึง subscriptions ตาม location_id
func (r *notificationRepository) GetSubscriptionsByLocationID(ctx context.Context, locationID int) ([]entities.NotificationSubscription, error) {
	var subscriptions []entities.NotificationSubscription

	query := `
		SELECT id, user_id, location_id, channel, province_id, is_active, created_at, updated_at
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
		SELECT id, user_id, location_id, channel, province_id, is_active, created_at, updated_at
		FROM notification_subscriptions
		WHERE is_active = true
	`

	if err := r.db.SelectContext(ctx, &subscriptions, query); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *notificationRepository) CreateSubscription(ctx context.Context, sub *entities.NotificationSubscription) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO notification_subscriptions (user_id, location_id, channel, province_id, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	return r.db.QueryRowContext(ctx, query,
		sub.UserID,
		sub.LocationID,
		sub.Channel,
		sub.ProvinceID,
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

func (r *notificationRepository) IsUserAlreaySubscript(ctx context.Context, userID int64) bool {
	var count int
	query := `
		SELECT COUNT(*)
		FROM notification_subscriptions
		WHERE user_id = $1
	`
	if err := r.db.GetContext(ctx, &count, query, userID); err != nil {
		return false
	}
	return count > 0
}
