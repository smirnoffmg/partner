package services

import (
	"github.com/rs/zerolog/log"
	"github.com/smirnoffmg/partner/internal/repositories"
)

type SubscriptionService struct {
	repo              repositories.IChatRepo
	freeMessagesCount int32
}

func NewSubscriptionService(repo *repositories.ChatsRepo, freeMessagesCount int32) (*SubscriptionService, error) {
	return &SubscriptionService{
		repo:              repo,
		freeMessagesCount: freeMessagesCount,
	}, nil
}

func (s *SubscriptionService) IncreaseMessageCount(chatID int64) {
	if err := s.repo.IncreaseMessageCount(chatID); err != nil {
		log.Error().Err(err).Msgf("Cannot increase message count for chat with ID: %d", chatID)
	}
}

func (s *SubscriptionService) GetMessagesRemain(chatID int64) (int32, error) {
	chat, err := s.repo.GetOrCreate(chatID)
	if err != nil {
		log.Error().Err(err).Msgf("Cannot get chat with id: %d", chatID)
	}

	return s.freeMessagesCount + chat.PayedMessagesCount - chat.UserMessagesCount, nil
}
