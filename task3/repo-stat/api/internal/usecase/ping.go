package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/domain"
)

type PingUseCase struct {
	processor  Pinger
	subscriber Pinger
}

func NewPingUseCase(processor Pinger, subscriber Pinger) *PingUseCase {
	return &PingUseCase{
		processor:  processor,
		subscriber: subscriber,
	}
}

func (u *PingUseCase) Execute(ctx context.Context) (domain.PingResponse, bool) {
	pStatus := u.processor.Ping(ctx)
	sStatus := u.subscriber.Ping(ctx)

	isOk := pStatus == domain.PingStatusUp && sStatus == domain.PingStatusUp

	statusText := "ok"
	if !isOk {
		statusText = "degraded"
	}

	return domain.PingResponse{
		Status: statusText,
		Services: []domain.ServiceStatus{
			{Name: "processor", Status: string(pStatus)},
			{Name: "subscriber", Status: string(sStatus)},
		},
	}, isOk
}
