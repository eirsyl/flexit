package transportgrpc

import (
	"context"

	flexitendpoint "github.com/eirsyl/flexit/endpoint"
	"github.com/eirsyl/flexit/examples/simple/pb"
	"github.com/eirsyl/flexit/examples/simple/pkg/endpoint"
	"github.com/eirsyl/flexit/examples/simple/pkg/service"
)

// NewGRPCServer creates a new server endpoint set
func NewGRPCServer(s service.Service, middleware []flexitendpoint.Middleware) service.Service {
	return &endpoint.Set{
		AddEndpoint:      flexitendpoint.Compose(MakeAddServerEndpoint(s), "Add", middleware...),
		SubtractEndpoint: flexitendpoint.Compose(MakeSubtractServerEndpoint(s), "Subtract", middleware...),
	}
}

// MakeAddServerEndpoint encodes the add endpoint
func MakeAddServerEndpoint(s service.Service) flexitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.AddRequest)
		return s.Add(ctx, req)
	}
}

// MakeSubtractServerEndpoint encodes the subtract endpoint
func MakeSubtractServerEndpoint(s service.Service) flexitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.SubtractRequest)
		return s.Subtract(ctx, req)
	}
}
