package entities

import (
	"time"
)

type NotificationSubscription struct {
	ID         int       `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	LocationID int64     `json:"location_id" db:"location_id"`
	Channel    string    `json:"channel" db:"channel"`
	ProvinceID int64     `json:"province_id" db:"province_id"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type NotificationLog struct {
	ID             int        `json:"id" db:"id"`
	SubscriptionID int64      `json:"subscription_id" db:"subscription_id"`
	LocationID     int64      `json:"location_id" db:"location_id"`
	WaterLevel     float64    `json:"water_level" db:"water_level"`
	Message        string     `json:"message" db:"message"`
	Channel        string     `json:"channel" db:"channel"`
	Status         string     `json:"status" db:"status"` // 'pending', 'sent', 'failed'
	SentAt         *time.Time `json:"sent_at" db:"sent_at"`
	ErrorMessage   *string    `json:"error_message" db:"error_message"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}
