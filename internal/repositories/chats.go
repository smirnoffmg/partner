package repositories

import (
	"gorm.io/gorm"

	entities "github.com/smirnoffmg/partner/internal/entities"
)

type ChatsRepo struct {
	DB *gorm.DB
}

func NewChatsRepo(db *gorm.DB) *ChatsRepo {
	return &ChatsRepo{
		DB: db,
	}
}

func (r *ChatsRepo) Create(chat *entities.Chat) error {
	return r.DB.Create(chat).Error
}

func (r *ChatsRepo) FindByTelegramChatID(chatID int64) (*entities.Chat, error) {
	var chat entities.Chat
	err := r.DB.Where("telegram_chat_id = ?", chatID).First(&chat).Error
	return &chat, err
}

func (r *ChatsRepo) FindByThreadID(threadID string) (*entities.Chat, error) {
	var chat entities.Chat
	err := r.DB.Where("thread_id = ?", threadID).First(&chat).Error
	return &chat, err
}

func (r *ChatsRepo) UpdateLastMessageID(threadID string, messageID string) error {
	return r.DB.Model(&entities.Chat{}).Where("thread_id = ?", threadID).Update("last_thread_message_id", messageID).Error
}

func (r *ChatsRepo) IncreaseMessageCount(chatID int64) {
	r.DB.Model(&entities.Chat{}).Where("telegram_chat_id = ?", chatID).Update("user_messages_count", gorm.Expr("user_messages_count + ?", 1))
}
