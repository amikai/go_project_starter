package hello

import (
	"context"
	"log"

	hellov1 "github.com/amikai/go_project_starter/src/internal/gen/hello/v1"
)

type Server struct {
	hellov1.UnimplementedHelloServiceServer
}

func (s *Server) SayHello(ctx context.Context, req *hellov1.SayHelloRequest) (*hellov1.SayHelloResponse, error) {
	log.Printf("Received: %v", req.GetName())
	return &hellov1.SayHelloResponse{Message: "Hello " + req.GetName()}, nil
}
