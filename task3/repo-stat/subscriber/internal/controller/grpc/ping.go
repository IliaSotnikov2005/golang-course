package grpc

import (
	"context"
	"log/slog"

	subscriberpb "github.com/IliaSotnikov2005/golang-course/task3/repo-stat/proto/subscriber"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/subscriber/internal/usecase"
)

type Server struct {
	subscriberpb.UnimplementedSubscriberServer
	log  *slog.Logger
	ping *usecase.Ping
}

func NewServer(log *slog.Logger, ping *usecase.Ping) *Server {
	return &Server{
		log:  log,
		ping: ping,
	}
}

func (s *Server) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	s.log.Debug("subscriberp ping request received")

	return &subscriberpb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}
