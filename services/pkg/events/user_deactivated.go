package events

type UserDeactivated struct {
	UserUUID string `json:"user_uuid"`
	ChatId   int64  `json:"chat_id"`
}
