package models

import (
	"gorm.io/gorm"
	"time"
)

const (
	SubscriptionTestCode = "test"
	SubscriptionBaseCode = "base"
)

type Subscription struct {
	ID        uint           `gorm:"primaryKey"`
	UserUUID  string         `gorm:"not null;index;"`
	User      *User          `gorm:"foreignKey:UserUUID;references:UUID;constraint:OnDelete:CASCADE;"` // Ключевой момент
	PlanCode  string         `gorm:"index"`
	Plan      Plan           `gorm:"foreignKey:PlanCode;references:Code;constraint:OnDelete:SET NULL;"`
	StartedAt time.Time      `gorm:"not null"`
	ExpiresAt time.Time      `gorm:"not null"`
	IsActive  bool           `gorm:"default:true"`
	AutoRenew bool           `gorm:"default:false"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
