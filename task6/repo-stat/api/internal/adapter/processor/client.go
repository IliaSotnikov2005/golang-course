package processor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/domain"
	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/utils"
	processorpb "github.com/IliaSotnikov2005/golang-course/task6/repo-stat/proto/processor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log    *slog.Logger
	conn   *grpc.ClientConn
	client processorpb.ProcessorServiceClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	log.Info("Creating collector gRPC client", slog.String("address", address))

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to collector: %w", err)
	}

	return &Client{
		log:    log,
		conn:   conn,
		client: processorpb.NewProcessorServiceClient(conn),
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	const operation = "adapter.processor.GetRepo"

	req := &processorpb.GetRepositoryRequest{
		Owner: owner,
		Repo:  repo,
	}

	resp, err := c.client.GetRepository(ctx, req)
	if err != nil {
		c.log.Error("processor grpc call failed",
			slog.String("op", operation),
			slog.Any("error", err),
		)
		return nil, utils.MapGRPCErrorToDomain(err)
	}

	return &domain.Repository{
		FullName:    resp.Info.GetFullName(),
		Description: resp.Info.GetDescription(),
		Stargazers:  int(resp.Info.GetStargazers()),
		Forks:       int(resp.Info.GetForks()),
		CreatedAt:   resp.Info.GetCreatedAt().AsTime(),
		HTMLURL:     resp.Info.GetHtmlUrl(),
	}, nil
}

func (c *Client) GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error) {
	resp, err := c.client.GetSubscriptionsInfo(ctx, &processorpb.GetSubscriptionsInfoRequest{})
	if err != nil {
		return nil, utils.MapGRPCErrorToDomain(err)
	}

	results := make([]domain.Repository, 0, len(resp.GetRepositories()))
	for _, r := range resp.GetRepositories() {
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
func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.client.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		c.log.Error("processor ping failed", slog.Any("error", err))
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) Close() error {
	return c.conn.Close()
}
