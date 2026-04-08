package grpccontroller

import (
	"context"
	"log/slog"

	pb "github.com/IliaSotnikov2005/golang-course/task4/repo-stat/proto/subscriber"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/subscriber/internal/domain"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/subscriber/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedSubscriberServer
	log           *slog.Logger
	subscribeUC   *usecase.SubscribeUseCase
	unsubscribeUC *usecase.UnsubscribeUseCase
	listUC        *usecase.ListUseCase
	pingUC        *usecase.PingUseCase
}

func NewServer(log *slog.Logger, s *usecase.SubscribeUseCase, u *usecase.UnsubscribeUseCase, l *usecase.ListUseCase, p *usecase.PingUseCase) *Server {
	return &Server{
		log:           log,
		subscribeUC:   s,
		unsubscribeUC: u,
		listUC:        l,
		pingUC:        p,
	}
}

func (s *Server) Subscribe(ctx context.Context, req *pb.SubscribeRequest) (*pb.SubscribeResponse, error) {
	s.log.Info("grpc: Subscribe request", slog.String("owner", req.GetOwner()), slog.String("repo", req.GetRepo()))

	sub, err := s.subscribeUC.Execute(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		switch err {
		case domain.ErrRepositoryNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case domain.ErrSubscriptionAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		s.log.Error("subscribe execution failed", slog.String("err", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.SubscribeResponse{
		Id:    int64(sub.ID),
		Owner: sub.Owner,
		Repo:  sub.Repo,
	}, nil
}

func (s *Server) Unsubscribe(ctx context.Context, req *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error) {
	s.log.Info("grpc: Unsubscribe request", slog.String("owner", req.GetOwner()), slog.String("repo", req.GetRepo()))

	err := s.unsubscribeUC.Execute(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		s.log.Error("unsubscribe execution failed", slog.String("err", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.UnsubscribeResponse{
		Success: true,
	}, nil
}

func (s *Server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.log.Info("grpc: List request")

	subscriptions, err := s.listUC.Execute(ctx)
	if err != nil {
		s.log.Error("list execution failed", slog.String("err", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	pbSubs := make([]*pb.ListResponse_Subscription, 0, len(subscriptions))
	for _, sub := range subscriptions {
		pbSubs = append(pbSubs, &pb.ListResponse_Subscription{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return &pb.ListResponse{Subscriptions: pbSubs}, nil
}

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	s.log.Debug("grpc: Ping request")
	statusStr := s.pingUC.Execute(ctx)

	return &pb.PingResponse{Reply: string(statusStr)}, nil
}
