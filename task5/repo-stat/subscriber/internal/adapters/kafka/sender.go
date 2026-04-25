package kafka

import (
	"context"
	"encoding/json"
	"fmt"

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
	val, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal kafka payload: %w", err)
	}

	results := s.client.ProduceSync(ctx, &kgo.Record{
		Topic: s.topic,
		Value: val,
	})

	if err := results.FirstErr(); err != nil {
		return fmt.Errorf("kafka produce failed: %w", err)
	}

	return nil
}
