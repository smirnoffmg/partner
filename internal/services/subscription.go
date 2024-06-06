package services

import (
	"github.com/rs/zerolog/log"
	"github.com/smirnoffmg/partner/internal/repositories"
)

type SubscriptionService struct {
	chatsRepo repositories.IChatRepo
}

func NewSubscriptionService(repo *repositories.ChatsRepo) (*SubscriptionService, error) {
	return &SubscriptionService{
		chatsRepo: repo,
	}, nil
}

func (s *SubscriptionService) IncreaseMessageCount(chatID int64) {
	if err := s.chatsRepo.IncreaseMessageCount(chatID); err != nil {
		log.Error().Err(err).Msgf("Cannot increase message count for chat with ID: %d", chatID)
	}
}
