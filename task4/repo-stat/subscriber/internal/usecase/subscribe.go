package usecase

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/subscriber/internal/domain"
)

type SubscribeUseCase struct {
	subscriptionRepository SubscriptionRepository
	github                 GithubClient
}

func NewSubscribeUseCase(subscriptionRepository SubscriptionRepository, github GithubClient) *SubscribeUseCase {
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
		return nil, fmt.Errorf("repository %s/%s does not exist", owner, repo)
	}

	sub := &domain.Subscription{
		Owner: owner,
		Repo:  repo,
	}

	return uc.subscriptionRepository.Save(ctx, sub)
}
