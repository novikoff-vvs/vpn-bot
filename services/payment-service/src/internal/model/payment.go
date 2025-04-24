package model

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	gorm.Model        `json:"metadata"`
	PaymentAmount     float64   `json:"payment_amount"`
	VpnUserID         uint      `json:"vpn_user_id"`
	IsPaid            bool      `json:"is_paid" gorm:"default:false"`
	InvoiceDate       time.Time `json:"invoice_date" time_format:"2020-07-17"`
	IsRemindedHasSent bool      `json:"is_reminded_has_sent" gorm:"default:false"`
	UUID              string    `json:"uuid" gorm:"unique"`
}
