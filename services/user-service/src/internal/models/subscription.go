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
	UserUUID  string         `gorm:"not null;index"` // связь с User
	User      *User          `gorm:"foreignKey:UserUUID;references:UUID"`
	PlanCode  string         `gorm:"index"` // связь с Plan
	Plan      Plan           `gorm:"foreignKey:PlanCode;references:Code"`
	Type      string         `gorm:"not null;size:50"` // например: "free", "pro", "premium"
	StartedAt time.Time      `gorm:"not null"`
	ExpiresAt time.Time      `gorm:"not null"`
	IsActive  bool           `gorm:"default:true"`  // удобен для отключения по крону
	AutoRenew bool           `gorm:"default:false"` // если используешь автопродление
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
