package subscription

import (
	"gorm.io/gorm"
	"time"
)

// Подписки (Subscriptions)
type Subscription struct {
	gorm.Model
	ID                uint               `gorm:"primaryKey"`
	Name              string             `gorm:"not null;size:255"`
	Description       string             `gorm:"type:text"`
	PriceMonthly      float64            `gorm:"type:decimal(10,2);not null"`
	PriceYearly       float64            `gorm:"type:decimal(10,2)"`
	Features          string             `gorm:"type:json"`
	IsActive          bool               `gorm:"default:true"`
	UserSubscriptions []UserSubscription `gorm:"foreignKey:SubscriptionID"`
}

// Подписки пользователей (UserSubscriptions)
type UserSubscription struct {
	gorm.Model
	ID              uint                  `gorm:"primaryKey"`
	UserID          uint                  `gorm:"not null;index"`
	SubscriptionID  uint                  `gorm:"not null;index"`
	StartDate       time.Time             `gorm:"not null"`
	EndDate         time.Time             `gorm:"not null;index"`
	NextPaymentDate time.Time             `gorm:"not null;index"`
	BillingCycle    string                `gorm:"size:20;not null"` // 'monthly' or 'yearly'
	Status          string                `gorm:"size:50;not null;default:'active'"`
	AutoRenewal     bool                  `gorm:"default:true"`
	LastPaymentID   *uint                 // Может быть nil
	User            User                  `gorm:"foreignKey:UserID"`
	Subscription    Subscription          `gorm:"foreignKey:SubscriptionID"`
	LastPayment     *Payment              `gorm:"foreignKey:LastPaymentID"`
	History         []SubscriptionHistory `gorm:"foreignKey:UserSubscriptionID"`
}
