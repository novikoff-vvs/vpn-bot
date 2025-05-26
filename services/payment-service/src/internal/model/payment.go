package model

import (
	"time"
)

type Payment struct {
	PaymentAmount  float64 `json:"payment_amount"`
	OperationId    string  `json:"operation_id"`
	WithdrawAmount float64 `json:"withdraw_amount"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Token          string
	UUID           string `gorm:"unique;primaryKey"`
}
