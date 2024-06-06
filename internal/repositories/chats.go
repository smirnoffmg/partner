package repositories

import (
	"github.com/rs/zerolog/log"
	entities "github.com/smirnoffmg/partner/internal/entities"
	"gorm.io/gorm"
)

type ChatsRepo struct {
	db *gorm.DB
}

func NewChatsRepo(db *gorm.DB) *ChatsRepo {
	return &ChatsRepo{
		db: db,
	}
}

func (r *ChatsRepo) GetOrCreate(chatID int64) (*entities.Chat, error) {
	var chat entities.Chat

	err := r.db.Where(entities.Chat{TelegramChatID: chatID}).FirstOrCreate(&chat).Error
	if err != nil {
		log.Error().Err(err).Msgf("Cannot find nor create chat with chatID: %d", chatID)

		return nil, err
	}

	return &chat, nil
}

func (r *ChatsRepo) Update(chatID int64, updates map[string]interface{}) error {
	chat, err := r.GetOrCreate(chatID)
	if err != nil {
		return err
	}

	err = r.db.Model(&chat).Omit("id", "telegram_chat_id").Updates(updates).Error
	if err != nil {
		log.Error().Err(err).Msgf("Cannot update chat with chatID: %d", chatID)

		return err
	}

	return nil
}

func (r *ChatsRepo) IncreaseMessageCount(chatID int64) error {
	err := r.db.Model(&entities.Chat{}).Where("telegram_chat_id = ?", chatID).UpdateColumn("user_messages_count", gorm.Expr("user_messages_count + ?", 1)).Error
	if err != nil {
		log.Error().Err(err).Msgf("Cannot increment messages count in chat with chatID: %d", chatID)
	}

	return err
}
