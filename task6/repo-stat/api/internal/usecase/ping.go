package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/domain"
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
	resChan := make(chan domain.ServiceStatus, 2)

	go func() {
		resChan <- domain.ServiceStatus{Name: "processor", Status: string(u.processor.Ping(ctx))}
	}()

	go func() {
		resChan <- domain.ServiceStatus{Name: "subscriber", Status: string(u.subscriber.Ping(ctx))}
	}()

	var services []domain.ServiceStatus
	allUp := true

	for range 2 {
		res := <-resChan
		services = append(services, res)
		if domain.PingStatus(res.Status) != domain.PingStatusUp {
			allUp = false
		}
	}

	statusText := "ok"
	if !allUp {
		statusText = "degraded"
	}

	return domain.PingResponse{
		Status:   statusText,
		Services: services,
	}, allUp
}
