package model

type Plan struct {
	Code  string  `gorm:"unique;primaryKey"`
	Price float64 `gorm:"not null" json:"price"`
}
