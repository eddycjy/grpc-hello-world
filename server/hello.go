package server

import (
	"golang.org/x/net/context"

	pb "github.com/EDDYCJY/grpc-hello-world/proto"
)

type helloService struct{}

func NewHelloService() *helloService {
	return &helloService{}
}

func (h helloService) SayHelloWorld(ctx context.Context, r *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{
		Message : r.Referer,
	}, nil
}