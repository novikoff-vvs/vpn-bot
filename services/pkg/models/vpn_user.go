package models

type VpnUser struct {
	ChatId         int64  `json:"chat_id"`
	Email          string `json:"email"`
	UUID           string `gorm:"id"`
	SubscriptionId string `json:"subId"`
}
