package entities

import (
	"time"
)

type NotificationSubscription struct {
	ID             int       `json:"id" db:"id"`
	UserID         *int      `json:"user_id" db:"user_id"`
	LocationID     *int      `json:"location_id" db:"location_id"`
	Channel        string    `json:"channel" db:"channel"`                 // 'email', 'line', 'sms'
	Target         string    `json:"target" db:"target"`                   // email address, LINE user ID, phone
	ThresholdLevel float64   `json:"threshold_level" db:"threshold_level"` // alert when level exceeds this
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type NotificationLog struct {
	ID             int        `json:"id" db:"id"`
	SubscriptionID *int       `json:"subscription_id" db:"subscription_id"`
	LocationID     int        `json:"location_id" db:"location_id"`
	WaterLevel     float64    `json:"water_level" db:"water_level"`
	Message        string     `json:"message" db:"message"`
	Channel        string     `json:"channel" db:"channel"`
	Status         string     `json:"status" db:"status"` // 'pending', 'sent', 'failed'
	SentAt         *time.Time `json:"sent_at" db:"sent_at"`
	ErrorMessage   *string    `json:"error_message" db:"error_message"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}
