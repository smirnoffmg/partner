package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model

	Title          string
	Description    string
	Currency       string
	Amount         int64
	TelegramChatID int64

	Paid bool `gorm:"default:false"`
}

func (i *Invoice) String() string {
	return fmt.Sprintf("Invoice: %d, %v, %t", i.TelegramChatID, i.CreatedAt, i.Paid)
}
