package processor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/api/internal/domain"
	processorpb "github.com/IliaSotnikov2005/golang-course/task4/repo-stat/proto/processor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
		return nil, mapGRPCErrorToDomain(err)
	}

	return &domain.Repository{
		Name:        resp.GetName(),
		Description: resp.GetDescription(),
		Stargazers:  int(resp.GetStargazers()),
		Forks:       int(resp.GetForks()),
		CreatedAt:   resp.GetCreatedAt().AsTime(),
		HTMLURL:     resp.GetHtmlUrl(),
	}, nil
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

func mapGRPCErrorToDomain(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	switch st.Code() {
	case codes.NotFound:
		return fmt.Errorf("%w: %s", domain.ErrNotFound, st.Message())

	case codes.PermissionDenied:
		return fmt.Errorf("%w: %s", domain.ErrForbidden, st.Message())

	case codes.Unauthenticated:
		return fmt.Errorf("%w: %s", domain.ErrUnauthorized, st.Message())

	case codes.ResourceExhausted:
		return fmt.Errorf("%w: %s", domain.ErrRateLimit, st.Message())

	case codes.InvalidArgument:
		return fmt.Errorf("%w: %s", domain.ErrInvalidInput, st.Message())

	case codes.DeadlineExceeded:
		return fmt.Errorf("%w: request to collector service timed out", domain.ErrTimeout)

	case codes.Unavailable:
		return fmt.Errorf("%w: collector service is unavailable", domain.ErrInternal)

	case codes.Internal:
		return fmt.Errorf("%w: collector service internal error: %s", domain.ErrInternal, st.Message())

	default:
		return fmt.Errorf("%w: unexpected error from collector: %s", domain.ErrInternal, st.Message())
	}
}
