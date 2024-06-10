package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model

	TelegramChatID int64 `gorm:"unique"`

	ThreadID string `gorm:"unique"` // chatgpt thread id

	UserMessagesCount int32 // total messages from user in telegram chat
	PaidMessagesCount int32 // number of messages user paid for
}

func (c Chat) String() string {
	return fmt.Sprintf("Chat{TelegramChatID: %d, ThreadID: %s}", c.TelegramChatID, c.ThreadID)
}
