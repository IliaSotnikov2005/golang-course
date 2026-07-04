package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/domain"
)

type GetSubscriptionsInfoUseCase struct {
	processor RepositoryProvider
}

func NewGetSubscriptionsInfoUseCase(processor RepositoryProvider) *GetSubscriptionsInfoUseCase {
	return &GetSubscriptionsInfoUseCase{
		processor: processor,
	}
}

func (uc *GetSubscriptionsInfoUseCase) Execute(ctx context.Context) ([]domain.Repository, error) {
	return uc.processor.GetSubscriptionsInfo(ctx)
}
