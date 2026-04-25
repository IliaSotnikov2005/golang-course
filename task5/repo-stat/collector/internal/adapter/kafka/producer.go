package kafka

import (
	"context"
	"encoding/json"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/domain"
	dtokafka "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/domain/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Adapter struct {
	client *kgo.Client
	topic  string
}

func NewAdapter(client *kgo.Client, topic string) *Adapter {
	return &Adapter{client: client, topic: topic}
}

func (a *Adapter) Send(ctx context.Context, repo *domain.Repository, opErr error) error {
	var resp dtokafka.RepoResponse
	if opErr != nil {
		resp.Error = opErr.Error()
	} else {
		resp = dtokafka.RepoResponse{
			FullName:    repo.FullName,
			Description: repo.Description,
			Stargazers:  repo.Stargazers,
			Forks:       repo.Forks,
			CreatedAt:   repo.CreatedAt,
			HTMLURL:     repo.HTMLURL,
			Error:       "",
		}
	}

	val, _ := json.Marshal(resp)
	return a.client.ProduceSync(ctx, &kgo.Record{
		Topic: a.topic,
		Value: val,
	}).FirstErr()
}
func (a *Adapter) Close() {
	a.client.Close()
}
