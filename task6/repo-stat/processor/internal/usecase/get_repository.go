package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/processor/internal/domain"
)

type GetRepositoryUseCase struct {
	storage   DataStorage
	publisher EventPublisher
}

func NewGetRepositoryUseCase(storage DataStorage, publisher EventPublisher) *GetRepositoryUseCase {
	return &GetRepositoryUseCase{
		storage:   storage,
		publisher: publisher,
	}
}

func (u *GetRepositoryUseCase) Execute(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	fullName := owner + "/" + repo

	res, err := u.storage.GetByFullName(ctx, fullName)
	if err == nil {
		return res, nil
	}

	if errors.Is(err, domain.ErrNotFound) {
		err := u.publisher.PublishFetchRequest(ctx, owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to publish request: %w", err)
		}

		return nil, domain.ErrAccepted
	}

	return nil, err
}
