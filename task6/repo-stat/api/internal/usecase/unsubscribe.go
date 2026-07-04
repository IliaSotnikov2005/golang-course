package usecase

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/domain"
)

type UnsubscribeUseCase struct {
	subClient Subscriber
}

func NewUnsubscribeUseCase(subClient Subscriber) *UnsubscribeUseCase {
	return &UnsubscribeUseCase{subClient: subClient}
}

func (uc *UnsubscribeUseCase) Execute(ctx context.Context, owner, repo string) error {
	if owner == "" || repo == "" {
		return fmt.Errorf("unsubscribe error: %w", domain.ErrInvalidInput)
	}

	err := uc.subClient.Unsubscribe(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("subscriber client error: %w", err)
	}

	return nil
}
