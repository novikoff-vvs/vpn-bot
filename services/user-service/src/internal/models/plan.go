package models

import "time"

type Plan struct {
	Code         string    `gorm:"unique;primaryKey"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Price        float64   `gorm:"not null" json:"price"`
	DurationDays int       `gorm:"not null" json:"duration_days"`
	MaxDevices   int       `gorm:"not null" json:"max_devices"`
	Description  string    `gorm:"type:text" json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
