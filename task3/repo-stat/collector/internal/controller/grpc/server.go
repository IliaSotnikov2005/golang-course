package grpccontroller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/collector/internal/domain"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/collector/internal/usecase"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/proto/collector"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	repo, err := h.getRepositoryUseCase.Execute(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		h.log.Error("usecase error", slog.String("error", err.Error()))
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, domain.ErrInvalidInput):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &collector.GetRepositoryResponse{
		FullName:    repo.FullName,
		Description: repo.Description,
		Stargazers:  int32(repo.Stargazers),
		Forks:       int32(repo.Forks),
		CreatedAt:   timestamppb.New(repo.CreatedAt),
		HtmlUrl:     repo.HTMLURL,
	}, nil
}

func (h *Handler) Ping(ctx context.Context, req *collector.PingRequest) (*collector.PingResponse, error) {
	status := h.pingUseCase.Execute(ctx)

	if status != domain.PingStatusUp {
		return &collector.PingResponse{Reply: string(status)}, nil
	}

	return &collector.PingResponse{Reply: string(status)}, nil
}
