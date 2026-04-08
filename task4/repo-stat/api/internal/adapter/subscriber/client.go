package subscriber

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/api/internal/domain"

	subscirberpb "github.com/IliaSotnikov2005/golang-course/task4/repo-stat/proto/subscriber"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log    *slog.Logger
	conn   *grpc.ClientConn
	client subscirberpb.SubscriberClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	log.Info("Creating subscriber gRPC client", slog.String("address", address))

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to subscriber: %w", err)
	}

	return &Client{
		log:    log,
		conn:   conn,
		client: subscirberpb.NewSubscriberClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.client.Ping(ctx, &subscirberpb.PingRequest{})
	if err != nil {
		c.log.Error("subscriber ping failed", slog.Any("error", err))
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) Close() error {
	return c.conn.Close()
}
