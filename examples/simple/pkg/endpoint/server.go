package endpoint

import (
	"context"

	"github.com/eirsyl/flexit/endpoint"
	flexitendpoint "github.com/eirsyl/flexit/endpoint"
	"github.com/eirsyl/flexit/examples/simple/pb"
	"github.com/eirsyl/flexit/examples/simple/pkg/service"
	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/metrics"
	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
)

// New creates a new server endpoint set
func New(s service.Service, logger log.Logger, tracer opentracing.Tracer, duration metrics.Histogram, ravenClient *raven.Client) service.Service {
	middlewares := []flexitendpoint.Middleware{
		endpoint.TraceServer(tracer),
		endpoint.LoggingMiddleware(logger),
		endpoint.InstrumentingMiddleware(duration),
		endpoint.SentryMiddleware(ravenClient),
	}

	return &Set{
		AddEndpoint:      flexitendpoint.Compose(MakeAddEndpoint(s), "Add", middlewares...),
		SubtractEndpoint: flexitendpoint.Compose(MakeSubtractEndpoint(s), "Subtract", middlewares...),
	}
}

func MakeAddEndpoint(s service.Service) flexitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.AddRequest)
		return s.Add(ctx, req)
	}
}

func MakeSubtractEndpoint(s service.Service) flexitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.SubtractRequest)
		return s.Subtract(ctx, req)
	}
}
