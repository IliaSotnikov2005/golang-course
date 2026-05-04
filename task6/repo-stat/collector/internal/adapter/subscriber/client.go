package subscriber

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/collector/internal/domain"
	pb "github.com/IliaSotnikov2005/golang-course/task6/repo-stat/proto/subscriber"
)

type Client struct {
	client pb.SubscriberClient
}

func NewClient(grpcClient pb.SubscriberClient) *Client {
	return &Client{
		client: grpcClient,
	}
}

func (c *Client) GetSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	resp, err := c.client.List(ctx, &pb.ListRequest{})
	if err != nil {
		return nil, fmt.Errorf("grpc list error: %w", err)
	}

	res := make([]domain.Subscription, 0, len(resp.Subscriptions))
	for _, s := range resp.Subscriptions {
		res = append(res, domain.Subscription{
			Owner: s.Owner,
			Repo:  s.Repo,
		})
	}

	return res, nil
}
