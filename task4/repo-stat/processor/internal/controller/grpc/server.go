package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/domain"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/usecase"
	processorpb "github.com/IliaSotnikov2005/golang-course/task4/repo-stat/proto/processor"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	processorpb.UnimplementedProcessorServiceServer
	log                         *slog.Logger
	getRepositoryUseCase        usecase.GetRepositoryUseCase
	getSubscriptionsInfoUseCase usecase.GetSubscriptionsInfoUseCase
	pingUseCase                 usecase.PingUseCase
}

func NewHandler(
	log *slog.Logger,
	getRepositoryUseCase *usecase.GetRepositoryUseCase,
	getSubscribtionsInfoUseCase *usecase.GetSubscriptionsInfoUseCase,
	pingUseCase *usecase.PingUseCase,
) *Handler {
	return &Handler{
		log:                         log,
		getRepositoryUseCase:        *getRepositoryUseCase,
		getSubscriptionsInfoUseCase: *getSubscribtionsInfoUseCase,
		pingUseCase:                 *pingUseCase,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *processorpb.GetRepositoryRequest) (*processorpb.GetRepositoryResponse, error) {
	const operation = "grpccontroller.Handler.GetRepository"
	log := h.log.With(slog.String("operation", operation))
	log.Debug("grpc: GetRepository request", "owner", req.GetOwner(), "repo", req.GetRepo())

	repo, err := h.getRepositoryUseCase.Execute(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		h.log.Error("usecase error", slog.String("error", err.Error()))
	}

	return &processorpb.GetRepositoryResponse{
		Info: &processorpb.RepositoryInfo{
			FullName:    repo.FullName,
			Description: repo.Description,
			Stargazers:  int32(repo.Stargazers),
			Forks:       int32(repo.Forks),
			CreatedAt:   timestamppb.New(repo.CreatedAt),
			HtmlUrl:     repo.HTMLURL,
		}}, nil
}

func (h *Handler) GetSubscriptionsInfo(ctx context.Context, req *processorpb.GetSubscriptionsInfoRequest) (*processorpb.GetSubscriptionsInfoResponse, error) {
	const operation = "grpccontroller.Handler.GetSubscriptionsInfo"
	log := h.log.With(slog.String("operation", operation))
	log.Debug("grpc: GetSubscriptionsInfo request")

	repositories, err := h.getSubscriptionsInfoUseCase.Execute(ctx)
	if err != nil {
		log.Error("usecase error", slog.String("error", err.Error()))
		return nil, err
	}

	pbRepos := make([]*processorpb.RepositoryInfo, 0, len(repositories))

	for _, repo := range repositories {
		pbRepos = append(pbRepos, &processorpb.RepositoryInfo{
			FullName:    repo.FullName,
			Description: repo.Description,
			Stargazers:  int32(repo.Stargazers),
			Forks:       int32(repo.Forks),
			CreatedAt:   timestamppb.New(repo.CreatedAt),
			HtmlUrl:     repo.HTMLURL,
		})
	}

	return &processorpb.GetSubscriptionsInfoResponse{
		Repositories: pbRepos,
	}, nil
}
func (h *Handler) Ping(ctx context.Context, req *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	const operation = "grpccontroller.Handler.Ping"
	log := h.log.With(slog.String("operation", operation))
	log.Debug("grpc: Ping request")

	status := h.pingUseCase.Execute(ctx)

	if status != domain.PingStatusUp {
		return &processorpb.PingResponse{Reply: string(status)}, nil
	}

	return &processorpb.PingResponse{Reply: string(status)}, nil
}
