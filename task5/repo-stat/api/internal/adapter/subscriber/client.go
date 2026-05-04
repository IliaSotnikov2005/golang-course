package subscriber

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/api/internal/domain"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/api/internal/utils"

	subscirberpb "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/proto/subscriber"

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

func (c *Client) Subscribe(ctx context.Context, owner, repo string) error {
	_, err := c.client.Subscribe(ctx, &subscirberpb.SubscribeRequest{Owner: owner, Repo: repo})
	return utils.MapGRPCErrorToDomain(err)
}

func (c *Client) Unsubscribe(ctx context.Context, owner, repo string) error {
	_, err := c.client.Unsubscribe(ctx, &subscirberpb.UnsubscribeRequest{Owner: owner, Repo: repo})
	return utils.MapGRPCErrorToDomain(err)
}

func (c *Client) List(ctx context.Context) ([]domain.Subscription, error) {
	resp, err := c.client.List(ctx, &subscirberpb.ListRequest{})
	if err != nil {
		return nil, utils.MapGRPCErrorToDomain(err)
	}

	subs := make([]domain.Subscription, 0, len(resp.GetSubscriptions()))
	for _, sub := range resp.GetSubscriptions() {
		subs = append(subs, domain.Subscription{
			Owner: sub.GetOwner(),
			Repo:  sub.GetRepo(),
		})
	}

	return subs, nil
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
