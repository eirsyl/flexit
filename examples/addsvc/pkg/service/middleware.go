package service

import (
	"context"

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

func (m *loggerMiddlewre) Add(ctx context.Context, x, y int64) (sum int64, err error) {
	defer func() {
		m.logger.WithFields(&log.Fields{
			"x":   x,
			"y":   y,
			"sum": sum,
			"err": err,
		}).Infof("Add")
	}()
	return m.next.Add(ctx, x, y)
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

func (m instrumentingMiddleware) Add(ctx context.Context, x, y int64) (int64, error) {
	v, err := m.next.Add(ctx, x, y)
	m.additions.Add(1)
	return v, err
}
