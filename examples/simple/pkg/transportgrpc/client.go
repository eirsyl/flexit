package transportgrpc

import (
	"context"

	flexitendpoint "github.com/eirsyl/flexit/endpoint"
	"github.com/eirsyl/flexit/examples/simple/pb"
	"github.com/eirsyl/flexit/examples/simple/pkg/endpoint"
	"github.com/eirsyl/flexit/examples/simple/pkg/service"
	"github.com/eirsyl/flexit/log"
	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, tracer opentracing.Tracer, logger log.Logger, raven *raven.Client) service.Service {
	middlewares := []flexitendpoint.Middleware{
		flexitendpoint.TraceClient(tracer),
		flexitendpoint.SentryMiddleware(raven),
	}

	client := pb.NewSimpleClient(conn)

	return &endpoint.Set{
		AddEndpoint:      flexitendpoint.Compose(MakeAddEndpoint(client), "Add", middlewares...),
		SubtractEndpoint: flexitendpoint.Compose(MakeSubtractEndpoint(client), "Subtract", middlewares...),
	}
}

func MakeAddEndpoint(client pb.SimpleClient) flexitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.AddRequest)
		return client.Add(ctx, req)
	}
}

func MakeSubtractEndpoint(client pb.SimpleClient) flexitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.SubtractRequest)
		return client.Subtract(ctx, req)
	}
}
