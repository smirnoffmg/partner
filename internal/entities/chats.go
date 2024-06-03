package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	PartnerName string

	TelegramChatID int64 `gorm:"unique"`

	ThreadID            string `gorm:"unique"` // chatgpt thread id
	LastThreadMessageID string // last message id in chatgpt thread

	UserMessagesCount int32 // total messages from user in telegram chat
}

func (c Chat) String() string {
	return fmt.Sprintf("Chat{PartnerName: %s, TelegramChatID: %d, ThreadID: %s}", c.PartnerName, c.TelegramChatID, c.ThreadID)
}
