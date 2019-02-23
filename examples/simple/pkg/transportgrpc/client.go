package transportgrpc

import (
	"context"

	flexitendpoint "github.com/eirsyl/flexit/endpoint"
	"github.com/eirsyl/flexit/examples/simple/pb"
	"github.com/eirsyl/flexit/examples/simple/pkg/endpoint"
	"github.com/eirsyl/flexit/examples/simple/pkg/service"
	"google.golang.org/grpc"
)

// NewGRPCClient creates a new client endpoint set
func NewGRPCClient(conn *grpc.ClientConn, middleware []flexitendpoint.Middleware) service.Service {
	client := pb.NewSimpleClient(conn)

	return &endpoint.Set{
		AddEndpoint:      flexitendpoint.Compose(MakeAddClientEndpoint(client), "Add", middleware...),
		SubtractEndpoint: flexitendpoint.Compose(MakeSubtractClientEndpoint(client), "Subtract", middleware...),
	}
}

// MakeAddClientEndpoint creates the add endpoint
func MakeAddClientEndpoint(client pb.SimpleClient) flexitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.AddRequest)
		return client.Add(ctx, req)
	}
}

// MakeSubtractClientEndpoint creates the subtract endpoint
func MakeSubtractClientEndpoint(client pb.SimpleClient) flexitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.SubtractRequest)
		return client.Subtract(ctx, req)
	}
}
