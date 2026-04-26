package usecase

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/api/internal/domain"
)

type SubscribeUseCase struct {
	subscriberClient Subscriber
}

func NewSubscribeUseCase(subscriber Subscriber) *SubscribeUseCase {
	return &SubscribeUseCase{subscriberClient: subscriber}
}

func (suc *SubscribeUseCase) Execute(ctx context.Context, owner, repo string) error {
	if owner == "" || repo == "" {
		return fmt.Errorf("subscribe error: %w", domain.ErrInvalidInput)
	}

	if err := suc.subscriberClient.Subscribe(ctx, owner, repo); err != nil {
		return fmt.Errorf("subscriber client error: %w", err)
	}

	return nil
}
