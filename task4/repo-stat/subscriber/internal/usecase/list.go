package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/subscriber/internal/domain"
)

type ListUseCase struct {
	manager SubscriptionRepository
}

func NewListUseCase(manager SubscriptionRepository) *ListUseCase {
	return &ListUseCase{
		manager: manager,
	}
}

func (uc *ListUseCase) Execute(ctx context.Context) ([]domain.Subscription, error) {
	return uc.manager.List(ctx)
}
