package grpcclient

import (
	"context"
	"fmt"
	"time"

	collectorpb "github.com/IliaSotnikov2005/golang-course/task2/api/proto/gen"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn    *grpc.ClientConn
	client  collectorpb.CollectorServiceClient
	timeout time.Duration
}

func NewClient(address string, timeout time.Duration) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to collector: %w", err)
	}

	client := collectorpb.NewCollectorServiceClient(conn)

	return &Client{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	resp, err := c.client.GetRepository(ctx, &collectorpb.GetRepositoryRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return nil, mapGRPCErrorToDomain(err)
	}

	return &domain.Repository{
		Name:        resp.Name,
		Description: resp.Description,
		Stargazers:  int(resp.Stargazers),
		Forks:       int(resp.Forks),
		CreatedAt:   resp.CreatedAt.AsTime(),
		HTMLURL:     resp.HtmlUrl,
	}, nil
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
		return fmt.Errorf("%w: %v", domain.ErrNotFound, st.Message())
	case codes.PermissionDenied:
		return fmt.Errorf("%w: %v", domain.ErrForbidden, st.Message())
	case codes.Unauthenticated:
		return fmt.Errorf("%w: %v", domain.ErrUnauthorized, st.Message())
	case codes.ResourceExhausted:
		return fmt.Errorf("%w: %v", domain.ErrRateLimit, st.Message())
	case codes.InvalidArgument:
		return fmt.Errorf("%w: %v", domain.ErrInvalidInput, st.Message())
	case codes.DeadlineExceeded:
		return fmt.Errorf("%w: %v", domain.ErrTimeout, st.Message())
	case codes.Unavailable:
		return fmt.Errorf("collector service unavailable: %w", err)
	default:
		return fmt.Errorf("%w: %v", domain.ErrInternal, st.Message())
	}
}
