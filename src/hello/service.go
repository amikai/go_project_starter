package hello

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type Service interface {
	SayHello(ctx context.Context, name string) (string, error)
}

func NewService(logger log.Logger, counter metrics.Counter) Service {
	var svc Service
	{
		svc = NewBasicService()
		svc = ServiceLoggingMW(logger)(svc)
		svc = ServiceInstrumentingMW(counter)(svc)
	}
	return svc
}

type basicService struct{}

func NewBasicService() Service {
	return &basicService{}
}

func (s *basicService) SayHello(ctc context.Context, name string) (string, error) {
	return "Hello " + name, nil
}
