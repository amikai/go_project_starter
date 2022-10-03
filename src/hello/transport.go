package hello

import (
	"context"
	"errors"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"

	hellov1 "github.com/amikai/go_project_starter/src/internal/gen/hello/v1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	sayHello grpctransport.Handler
	hellov1.UnimplementedHelloServiceServer
}

// NewGRPCServer makes a set of endpoints available as a gRPC hello server.
func NewGRPCServer(endpoints Set, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) hellov1.HelloServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	if zipkinTracer != nil {
		// Zipkin GRPC Server Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit gRPC server as a
		// ServerOption.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path if used in combination with the Go kit gRPC Interceptor.
		//
		// In this example, we demonstrate a global Zipkin tracing service with
		// Go kit gRPC Interceptor.
		options = append(options, zipkin.GRPCServerTrace(zipkinTracer))
	}

	return &grpcServer{
		sayHello: grpctransport.NewServer(
			endpoints.HelloEndpoint,
			decodeGRPCSayHelloRequest,
			encodeGRPCSayHelloResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "SayHello", logger)))...,
		),
	}
}

func (s *grpcServer) SayHello(ctx context.Context, req *hellov1.SayHelloRequest) (*hellov1.SayHelloResponse, error) {
	_, rep, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*hellov1.SayHelloResponse), nil
}

// decodeGRPCSumRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sayHello request to a user-domain sayHello request. Primarily useful in a server.
func decodeGRPCSayHelloRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*hellov1.SayHelloRequest)
	return SayHelloRequest{Name: req.Name}, nil
}

func decodeGRPCSayHelloResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*hellov1.SayHelloResponse)
	return SayHelloResponse{Message: reply.Message, Err: str2err(reply.Err)}, nil
}

// encodeGRPCSayHelloRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain sum request to a gRPC sayHello request. Primarily useful in a client.
func encodeGRPCSayHelloRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(SayHelloRequest)
	return &hellov1.SayHelloRequest{Name: req.Name}, nil
}

// encodeGRPCSayHelloResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain sayHello response to a gRPC concat reply. Primarily useful in a
// server.
func encodeGRPCSayHelloResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SayHelloResponse)
	return &hellov1.SayHelloResponse{Message: resp.Message, Err: err2str(resp.Err)}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
