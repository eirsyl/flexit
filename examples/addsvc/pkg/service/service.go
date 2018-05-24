package service

import (
	"context"

	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/metrics"
)

type Service interface {
	Add(ctx context.Context, x, y int64) (int64, error)
}

func New(logger log.Logger, additions metrics.Counter) Service {
	var svc Service
	svc = NewBasicService()
	svc = LoggerMiddleware(logger)(svc)
	svc = InstrumentingMiddleware(additions)(svc)
	return svc
}

type basicService struct {
}

func NewBasicService() Service {
	return &basicService{}
}

func (s *basicService) Add(_ context.Context, x, y int64) (int64, error) {
	return x + y, nil
}
