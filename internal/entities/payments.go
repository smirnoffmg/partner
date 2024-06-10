package entities

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model

	InvoicePayload          string `json:"invoice_payload" gorm:"not null"`
	TotalAmount             int    `json:"total_amount" gorm:"not null"`
	Currency                string `json:"currency" gorm:"not null"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id" gorm:"not null"`
	ProviderPaymentChargeID string `json:"provider_payment_charge_id" gorm:"not null"`

	Success bool
}
