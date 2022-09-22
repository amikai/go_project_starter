package hello

import (
	"context"
	"log"

	"github.com/amikai/go_project_starter/src/hello/pb"
)

var _ pb.HelloServer = &Server{}

type Server struct {
	pb.UnimplementedHelloServer
}

func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", req.GetName())
	return &pb.HelloReply{Message: "Hello " + req.GetName()}, nil
}
