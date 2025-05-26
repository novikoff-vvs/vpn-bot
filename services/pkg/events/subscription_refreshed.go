package events

type SubscriptionRefreshed struct {
	UserUUID string `json:"user_uuid"`
	ChatId   int64  `json:"chat_id"`
}
