package service

import (
	"context"
	"errors"

	"github.com/eirsyl/flexit/examples/simple/pb"

	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/metrics"
)

// Service defines the exposed methods
type Service pb.SimpleServer

// New returns a new inatnce of this service with middleware
func New(logger log.Logger, additions metrics.Counter) Service {
	var svc Service
	svc = NewBasicService()
	svc = LoggerMiddleware(logger)(svc)
	svc = InstrumentingMiddleware(additions)(svc)
	return svc
}

type basicService struct{}

// NewBasicService creates the service without middleware
func NewBasicService() Service {
	return &basicService{}
}

func (s *basicService) Add(_ context.Context, r *pb.AddRequest) (*pb.AddResponse, error) {
	return &pb.AddResponse{Sum: r.X + r.Y}, errors.New("OVERFLOW ERROR")
}

func (s *basicService) Subtract(_ context.Context, r *pb.SubtractRequest) (*pb.SubtractResponse, error) {
	return &pb.SubtractResponse{Sum: r.X - r.Y, Err: "OVERFLOW ERROR"}, nil
}
