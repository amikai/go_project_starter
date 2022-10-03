package hello

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency
// and returns a service Middleware.
func ServiceLoggingMW(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) SayHello(ctx context.Context, name string) (msg string, err error) {
	defer func() {
		mw.logger.Log("method", "SayHello", "msg", msg, "err", err)
	}()
	return mw.next.SayHello(ctx, name)
}

// InstrumentingMiddleware returns a service middleware that instruments
// the number of integers summed and characters concatenated over the lifetime of
// the service.
func ServiceInstrumentingMW(counter metrics.Counter) Middleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			counter: counter,
			next:    next,
		}
	}
}

type instrumentingMiddleware struct {
	counter metrics.Counter
	next    Service
}

func (mw instrumentingMiddleware) SayHello(ctx context.Context, name string) (string, error) {
	msg, err := mw.next.SayHello(ctx, name)
	mw.counter.Add(1)
	return msg, err
}
