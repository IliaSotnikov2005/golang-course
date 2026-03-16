package grpcadapter

import (
	"context"

	collectorpb "github.com/IliaSotnikov2005/golang-course/task2/api/proto/gen"
	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	collectorpb.UnimplementedCollectorServiceServer
	getRepoUseCase *usecase.GetRepositoryUseCase
}

func NewHandler(getRepoUseCase *usecase.GetRepositoryUseCase) *Handler {
	return &Handler{
		getRepoUseCase: getRepoUseCase,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *collectorpb.GetRepositoryRequest) (*collectorpb.GetRepositoryResponse, error) {
	repo, err := h.getRepoUseCase.Execute(ctx, req.Owner, req.Repo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &collectorpb.GetRepositoryResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stargazers:  int32(repo.Stargazers),
		Forks:       int32(repo.Forks),
		CreatedAt:   timestamppb.New(repo.CreatedAt),
		HtmlUrl:     repo.HTMLURL,
	}, nil
}
