package collector

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/domain"
	collectorpb "github.com/IliaSotnikov2005/golang-course/task4/repo-stat/proto/collector"
)

type collectorAdapter struct {
	client collectorpb.CollectorServiceClient
}

func NewCollectorAdapter(grpcClient collectorpb.CollectorServiceClient) *collectorAdapter {
	return &collectorAdapter{
		client: grpcClient,
	}
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
		FullName:    response.Info.FullName,
		Description: response.Info.Description,
		Stargazers:  int(response.Info.Stargazers),
		Forks:       int(response.Info.Forks),
		CreatedAt:   response.Info.CreatedAt.AsTime(),
	}, nil
}

func (ca *collectorAdapter) GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error) {
	req := &collectorpb.GetSubscriptionsInfoRequest{}

	response, err := ca.client.GetSubscriptionsInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	results := make([]domain.Repository, 0, len(response.GetRepositories()))
	for _, r := range response.GetRepositories() {
		results = append(results, domain.Repository{
			FullName:    r.GetFullName(),
			Description: r.GetDescription(),
			Stargazers:  int(r.GetStargazers()),
			Forks:       int(r.GetForks()),
			CreatedAt:   r.GetCreatedAt().AsTime(),
			HTMLURL:     r.GetHtmlUrl(),
		})
	}

	return results, nil
}
