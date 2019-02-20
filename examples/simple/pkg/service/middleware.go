package service

import (
	"context"

	"github.com/eirsyl/flexit/examples/simple/pb"

	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/metrics"
)

type Middleware func(Service) Service

func LoggerMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggerMiddlewre{logger, next}
	}
}

type loggerMiddlewre struct {
	logger log.Logger
	next   Service
}

func (m *loggerMiddlewre) Add(ctx context.Context, r *pb.AddRequest) (resp *pb.AddResponse, err error) {
	defer func() {
		m.logger.WithFields(&log.Fields{
			"x":   r.X,
			"y":   r.Y,
			"sum": resp.Sum,
			"err": resp.Err,
		}).Infof("Add")
	}()
	return m.next.Add(ctx, r)
}

func (m *loggerMiddlewre) Subtract(ctx context.Context, r *pb.SubtractRequest) (resp *pb.SubtractResponse, err error) {
	defer func() {
		m.logger.WithFields(&log.Fields{
			"x":   r.X,
			"y":   r.Y,
			"sum": resp.Sum,
			"err": resp.Err,
		}).Infof("Subtract")
	}()
	return m.next.Subtract(ctx, r)
}

func InstrumentingMiddleware(additions metrics.Counter) Middleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			additions: additions,
			next:      next,
		}
	}
}

type instrumentingMiddleware struct {
	additions metrics.Counter
	next      Service
}

func (m instrumentingMiddleware) Add(ctx context.Context, r *pb.AddRequest) (*pb.AddResponse, error) {
	m.additions.Add(1)
	return m.next.Add(ctx, r)
}

func (m instrumentingMiddleware) Subtract(ctx context.Context, r *pb.SubtractRequest) (*pb.SubtractResponse, error) {
	m.additions.Add(1)
	return m.next.Subtract(ctx, r)
}
