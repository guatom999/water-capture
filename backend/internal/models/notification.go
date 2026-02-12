package models

type CreateSubscriptionRequest struct {
	UserId     int64  `json:"user_id" validate:"required,gt=0"`
	UserEmail  string `json:"user_email" validate:"required,email"`
	LocationID int64  `json:"location_id" validate:"required,gt=0"`
	Channel    string `json:"channel" validate:"required"`
	ProvinceID int64  `json:"province_id" validate:"required,gt=0"`
}

// type SubscriptionResponse struct {
//     ID        string    `json:"id"`
//     PlanID    string    `json:"plan_id"`
//     UserID    string    `json:"user_id"`
//     Status    string    `json:"status"`
//     Amount    float64   `json:"amount"`
//     CreatedAt time.Time `json:"created_at"`
// }
