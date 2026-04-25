package kafka

import (
	"context"
	"encoding/json"

	dtokafka "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/domain/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Dispatcher struct {
	client *kgo.Client
	topic  string
}

func NewDispatcher(client *kgo.Client, topic string) *Dispatcher {
	return &Dispatcher{client: client, topic: topic}
}

func (d *Dispatcher) Dispatch(ctx context.Context, owner, repo string) error {
	req := dtokafka.RepoRequest{Owner: owner, Repo: repo}
	val, _ := json.Marshal(req)
	return d.client.ProduceSync(ctx, &kgo.Record{
		Topic: d.topic,
		Value: val,
	}).FirstErr()
}
