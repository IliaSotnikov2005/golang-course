package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/domain"
)

type PingUseCase struct{}

func NewPingUseCase() *PingUseCase {
	return &PingUseCase{}
}

func (u *PingUseCase) Execute(ctx context.Context) domain.PingStatus {
	return domain.PingStatusUp
}
