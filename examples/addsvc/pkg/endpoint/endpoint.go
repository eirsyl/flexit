package endpoint

import (
	"context"

	"github.com/eirsyl/flexit/endpoint"
	"github.com/eirsyl/flexit/examples/addsvc/pkg/service"
	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/metrics"
	"github.com/eirsyl/flexit/tracing"
	"github.com/getsentry/raven-go"
	"github.com/opentracing/opentracing-go"
)

type Set struct {
	AddEndpoint endpoint.Endpoint
}

func composeEndpoints(ep endpoint.Endpoint, name string, tracer opentracing.Tracer, logger log.Logger, duration metrics.Histogram, ravenClient *raven.Client) endpoint.Endpoint {
	ep = tracing.TraceServer(tracer, name)(ep)
	ep = endpoint.LoggingMiddleware(logger.WithField("method", name))(ep)
	ep = endpoint.InstrumentingMiddleware(duration.With("method", name))(ep)
	ep = endpoint.SentryMiddleware(ravenClient)(ep)
	return ep
}

func New(svc service.Service, logger log.Logger, tracer opentracing.Tracer, duration metrics.Histogram, ravenClient *raven.Client) Set {
	return Set{
		AddEndpoint: composeEndpoints(MakeAddEndpoint(svc), "Add", tracer, logger, duration, ravenClient),
	}
}

// Add

func MakeAddEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddRequest)
		sum, err := s.Add(ctx, req.X, req.Y)
		return AddResponse{Sum: sum, Err: err}, nil
	}
}

func (s *Set) Add(ctx context.Context, x, y int64) (int64, error) {
	resp, err := s.AddEndpoint(ctx, AddRequest{X: x, Y: y})
	if err != nil {
		return 0, err
	}
	response := resp.(AddResponse)
	return response.Sum, response.Err
}

type AddRequest struct {
	X int64
	Y int64
}

type AddResponse struct {
	Sum int64
	Err error
}

func (r AddResponse) Failed() error { return r.Err }
