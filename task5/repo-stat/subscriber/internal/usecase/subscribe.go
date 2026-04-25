package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/domain"
)

type SubscribeUseCase struct {
	subscriptionRepository SubscriptionRepository
	github                 GithubClient
	eventSender            SubscriptionEventSender
}

func NewSubscribeUseCase(subscriptionRepository SubscriptionRepository, github GithubClient, eventSender SubscriptionEventSender) *SubscribeUseCase {
	return &SubscribeUseCase{
		subscriptionRepository: subscriptionRepository,
		github:                 github,
	}
}

func (uc *SubscribeUseCase) Execute(ctx context.Context, owner, repo string) (*domain.Subscription, error) {
	exists, err := uc.github.Exists(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrRepositoryNotFound
	}

	sub := &domain.Subscription{
		Owner: owner,
		Repo:  repo,
	}

	sub, err = uc.subscriptionRepository.Save(ctx, sub)
	if err != nil {
		return nil, err
	}

	_ = uc.eventSender.NotifySubscribed(ctx, owner, repo)
	return sub, nil
}
