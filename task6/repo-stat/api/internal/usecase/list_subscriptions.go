package usecase

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/domain"
)

type ListSubscriptionsUseCase struct {
	subClient Subscriber
}

func NewListSubscriptionsUseCase(subClient Subscriber) *ListSubscriptionsUseCase {
	return &ListSubscriptionsUseCase{subClient: subClient}
}

func (uc *ListSubscriptionsUseCase) Execute(ctx context.Context) ([]domain.Subscription, error) {
	subs, err := uc.subClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("subscriber client error: %w", err)
	}

	return subs, nil
}
