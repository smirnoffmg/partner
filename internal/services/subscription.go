package services

import "github.com/smirnoffmg/partner/internal/repositories"

type SubscriptionService struct {
	chatsRepo *repositories.ChatsRepo
}

func NewSubscriptionService(repo *repositories.ChatsRepo) (svc *SubscriptionService, err error) {
	svc = &SubscriptionService{
		chatsRepo: repo,
	}
	return
}
