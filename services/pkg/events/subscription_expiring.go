package events

type SubscriptionExpiring struct {
	UserUUID      string `json:"user_uuid"`
	ChatId        int64  `json:"chat_id"`
	DaysRemaining int    `json:"days_remaining"`
}

