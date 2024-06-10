package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model

	TelegramChatID int64 `gorm:"unique"`

	ThreadID string `gorm:"unique"` // chatgpt thread id

	UserMessagesCount int32 `gorm:"default:0"`
	PaidMessagesCount int32 `gorm:"default:0"`
}

func (c Chat) String() string {
	return fmt.Sprintf("Chat{TelegramChatID: %d, ThreadID: %s}", c.TelegramChatID, c.ThreadID)
}
