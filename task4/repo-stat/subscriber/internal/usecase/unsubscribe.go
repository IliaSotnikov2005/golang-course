package usecase

import (
	"context"
)

type UnsubscribeUseCase struct {
	manager SubscriptionRepository
}

func NewUnsubscribeUseCase(manager SubscriptionRepository) *UnsubscribeUseCase {
	return &UnsubscribeUseCase{
		manager: manager,
	}
}

func (uc *UnsubscribeUseCase) Execute(ctx context.Context, owner, repo string) error {
	return uc.manager.Delete(ctx, owner, repo)
}
