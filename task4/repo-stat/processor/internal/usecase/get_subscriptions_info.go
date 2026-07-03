package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/domain"
)

type GetSubscriptionsInfoUseCase struct {
	collector RepositoryProvider
}

func NewGetSubscriptionsInfoUseCase(collector RepositoryProvider) *GetSubscriptionsInfoUseCase {
	return &GetSubscriptionsInfoUseCase{collector: collector}
}

func (guc *GetSubscriptionsInfoUseCase) Execute(ctx context.Context) ([]domain.Repository, error) {
	return guc.collector.GetSubscriptionsInfo(ctx)
}
