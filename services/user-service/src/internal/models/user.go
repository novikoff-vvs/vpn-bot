package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UUID         string         `gorm:"unique;primaryKey"`
	Email        string         `gorm:"unique;not null;size:255"`
	Subscription *Subscription  `gorm:"foreignKey:UserUUID;references:UUID;constraint:OnDelete:CASCADE;"`
	ChatId       int64          `gorm:"unique"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoCreateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
