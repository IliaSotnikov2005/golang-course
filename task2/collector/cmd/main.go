package cmd

import (
	"log"
	"net"

	collectorpb "github.com/IliaSotnikov2005/golang-course/task2/api/proto/gen"
	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/adapter/github"
	grpcadapter "github.com/IliaSotnikov2005/golang-course/task2/collector/internal/adapter/grpc"
	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	githubClient := github.NewClient()
	getRepoUseCase := usecase.NewGetRepositoryUseCase(githubClient)
	grpcHandler := grpcadapter.NewHandler(getRepoUseCase)

	lis, err := net.Listen("tcp", ":50505")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	collectorpb.RegisterCollectorServiceServer(s, grpcHandler)

	log.Println("Collector service starting on :50505")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
