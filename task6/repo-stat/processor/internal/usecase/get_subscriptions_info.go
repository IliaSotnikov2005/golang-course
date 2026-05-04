package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/domain"
)

type GetSubscriptionsInfoUseCase struct {
	storage DataStorage
}

func NewGetSubscriptionsInfoUseCase(storage DataStorage) *GetSubscriptionsInfoUseCase {
	return &GetSubscriptionsInfoUseCase{storage: storage}
}

func (uc *GetSubscriptionsInfoUseCase) Execute(ctx context.Context) ([]domain.Repository, error) {
	return uc.storage.ListAll(ctx)
}
