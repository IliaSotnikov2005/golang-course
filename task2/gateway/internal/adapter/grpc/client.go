package grpcclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	collectorpb "github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/api/proto/gen"
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
	log     *slog.Logger
}

func NewClient(address string, timeout time.Duration, log *slog.Logger) (*Client, error) {
	log.Info("Creating gRPC client", slog.String("address", address))
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
		log:     log,
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	log := c.log.With(
		slog.String("owner", owner),
		slog.String("repo", repo),
	)

	log.Debug("Calling collector.GetRepository")

	if _, ok := ctx.Deadline(); !ok && c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	resp, err := c.client.GetRepository(ctx, &collectorpb.GetRepositoryRequest{
		Owner: owner,
		Repo:  repo,
	})

	if err != nil {
		log.Error("gRPC call failed", slog.String("error", err.Error()))
		return nil, mapGRPCErrorToDomain(err)
	}

	log.Debug("gRPC call successful")

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
