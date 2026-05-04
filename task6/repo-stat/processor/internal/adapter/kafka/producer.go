package kafka

import (
	"context"
	"encoding/json"

	dtokafka "github.com/IliaSotnikov2005/golang-course/task6/repo-stat/platform/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Publisher struct {
	client *kgo.Client
	topic  string
}

func NewPublisher(client *kgo.Client, topic string) *Publisher {
	return &Publisher{client: client, topic: topic}
}

func (p *Publisher) PublishFetchRequest(ctx context.Context, owner, repo string) error {
	payload := dtokafka.RepoRequest{Owner: owner, Repo: repo}
	val, _ := json.Marshal(payload)

	return p.client.ProduceSync(ctx, &kgo.Record{
		Topic: p.topic,
		Value: val,
	}).FirstErr()
}
