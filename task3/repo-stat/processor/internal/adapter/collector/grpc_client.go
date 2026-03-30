package collector

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/processor/internal/domain"
	collectorpb "github.com/IliaSotnikov2005/golang-course/task3/repo-stat/proto/collector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type collectorAdapter struct {
	client collectorpb.CollectorServiceClient
	conn   *grpc.ClientConn
}

func NewCollectorAdapter(address string) (*collectorAdapter, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to collector: %w", err)
	}

	return &collectorAdapter{
		client: collectorpb.NewCollectorServiceClient(conn),
		conn:   conn,
	}, nil
}

func (ca *collectorAdapter) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	req := &collectorpb.GetRepositoryRequest{
		Owner: owner,
		Repo:  repo,
	}

	response, err := ca.client.GetRepository(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domain.Repository{
		Name:        response.Name,
		Description: response.Description,
		Stargazers:  int(response.Stargazers),
		Forks:       int(response.Forks),
		CreatedAt:   response.CreatedAt.AsTime(),
	}, nil
}

func (ca *collectorAdapter) Close() error {
	return ca.conn.Close()
}
