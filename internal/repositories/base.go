package repositories

import (
	"github.com/smirnoffmg/partner/internal/entities"
)

type IChatRepo interface {
	GetOrCreate(chatID int64) (*entities.Chat, error)
	Update(chatID int64, updates map[string]interface{}) error
	IncreaseMessageCount(chatID int64) error
}
