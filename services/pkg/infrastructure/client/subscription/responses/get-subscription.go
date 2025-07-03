package responses

import "time"

type GetSubscriptionResponse struct {
	Result GetSubscriptionResult `json:"result"`
}

type GetSubscriptionResult struct {
	ID        int        `json:"ID"`
	UserUUID  string     `json:"UserUUID"`
	User      User       `json:"User"`
	PlanCode  string     `json:"PlanCode"`
	Plan      Plan       `json:"Plan"`
	StartedAt time.Time  `json:"StartedAt"`
	ExpiresAt time.Time  `json:"ExpiresAt"`
	IsActive  bool       `json:"IsActive"`
	AutoRenew bool       `json:"AutoRenew"`
	CreatedAt time.Time  `json:"CreatedAt"`
	UpdatedAt time.Time  `json:"UpdatedAt"`
	DeletedAt *time.Time `json:"DeletedAt"` // Nullable
}

// User represents the user data
type User struct {
	UUID         string      `json:"UUID"`
	Email        string      `json:"Email"`
	Subscription interface{} `json:"Subscription"` // null, change type if needed
	ChatID       int64       `json:"ChatId"`
	CreatedAt    time.Time   `json:"CreatedAt"`
	IsActive     bool        `json:"IsActive"`
	UpdatedAt    time.Time   `json:"UpdatedAt"`
	DeletedAt    *time.Time  `json:"DeletedAt"` // Nullable
}

// Plan represents the subscription plan
type Plan struct {
	Code         string    `json:"Code"`
	Name         string    `json:"name"`
	Price        int       `json:"price"`
	DurationDays int       `json:"duration_days"`
	MaxDevices   int       `json:"max_devices"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}
