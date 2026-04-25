package kafka

import (
	"context"
	"encoding/json"

	dtokafka "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/platform/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

type EventSender struct {
	client *kgo.Client
	topic  string
}

func NewEventSender(client *kgo.Client, topic string) *EventSender {
	return &EventSender{
		client: client,
		topic:  topic,
	}
}

func (s *EventSender) NotifySubscribed(ctx context.Context, owner, repo string) error {
	payload := dtokafka.RepoRequest{Owner: owner, Repo: repo}
	val, _ := json.Marshal(payload)

	return s.client.ProduceSync(ctx, &kgo.Record{
		Topic: s.topic,
		Value: val,
	}).FirstErr()
}
