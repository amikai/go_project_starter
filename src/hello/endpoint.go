package hello

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

var _ endpoint.Failer = SayHelloResponse{}

type SayHelloRequest struct {
	Name string `json:"name"`
}

type SayHelloResponse struct {
	Message string `json:"message"`
	Err     error  `json:"error,omitempty"`
}

// Failed implements endpoint.Failer.
func (r SayHelloResponse) Failed() error { return r.Err }

type Set struct {
	HelloEndpoint endpoint.Endpoint
}

func NewEndpoints(svc Service, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) Set {
	var helloEndpoint endpoint.Endpoint
	{
		helloEndpoint = makeSayHelloEndpoint(svc)
		// SayHello is limited to 1 request per second with burst of 1 request.
		// Note, rate is defined as a time interval between requests.
		helloEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(helloEndpoint)
		helloEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(helloEndpoint)
		helloEndpoint = opentracing.TraceServer(otTracer, "hello")(helloEndpoint)
		if zipkinTracer != nil {
			helloEndpoint = zipkin.TraceEndpoint(zipkinTracer, "hello")(helloEndpoint)
		}
		helloEndpoint = EndpointLoggingMW(log.With(logger, "method", "hello"))(helloEndpoint)
		helloEndpoint = EndpointInstrumentingMW(duration.With("method", "hello"))(helloEndpoint)
	}
	return Set{
		HelloEndpoint: helloEndpoint,
	}
}

func makeSayHelloEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SayHelloRequest)
		msg, err := s.SayHello(ctx, req.Name)
		return SayHelloResponse{Message: msg, Err: err}, nil
	}
}
