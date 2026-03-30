package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/collector/internal/domain"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/collector/internal/usecase"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/proto/collector"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	collector.UnimplementedCollectorServiceServer
	log                  *slog.Logger
	getRepositoryUseCase *usecase.GetRepositoryUseCase
	pingUseCase          *usecase.PingUseCase
}

func NewHandler(
	log *slog.Logger,
	getRepositoryUseCase *usecase.GetRepositoryUseCase,
	pingUseCase *usecase.PingUseCase,
) *Handler {
	return &Handler{
		log:                  log,
		getRepositoryUseCase: getRepositoryUseCase,
		pingUseCase:          pingUseCase,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *collector.GetRepositoryRequest) (*collector.GetRepositoryResponse, error) {
	const operation = "grpccontroller.Handler.GetRepository"
	log := h.log.With(slog.String("operation", operation))
	log.Debug("grpc: GetRepository request", "owner", req.GetOwner(), "repo", req.GetRepo())

	repo, err := h.getRepositoryUseCase.Execute(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		h.log.Error("usecase error", slog.String("error", err.Error()))
	}

	return &collector.GetRepositoryResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stargazers:  int32(repo.Stargazers),
		Forks:       int32(repo.Forks),
		CreatedAt:   timestamppb.New(repo.CreatedAt),
		HtmlUrl:     repo.HTMLURL,
	}, nil
}

func (h *Handler) Ping(ctx context.Context, req *collector.PingRequest) (*collector.PingResponse, error) {
	const operation = "grpccontroller.Handler.Ping"
	log := h.log.With(slog.String("operation", operation))
	log.Debug("grpc: Ping request")

	status := h.pingUseCase.Execute(ctx)

	if status != domain.PingStatusUp {
		return &collector.PingResponse{Reply: string(status)}, nil
	}

	return &collector.PingResponse{Reply: string(status)}, nil
}
