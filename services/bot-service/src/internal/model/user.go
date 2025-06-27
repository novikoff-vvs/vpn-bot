package model

import "gorm.io/gorm"

type NotifyUser struct {
	gorm.Model
	ChatId int64 `gorm:"unique"`
}
