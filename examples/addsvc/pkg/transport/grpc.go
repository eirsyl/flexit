package transport

import (
	"context"
	"errors"

	flexitendpoint "github.com/eirsyl/flexit/endpoint"
	"github.com/eirsyl/flexit/examples/addsvc/pb"
	"github.com/eirsyl/flexit/examples/addsvc/pkg/endpoint"
	"github.com/eirsyl/flexit/examples/addsvc/pkg/service"
	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/tracing"
	grpctransport "github.com/eirsyl/flexit/transports/grpc"
	"github.com/getsentry/raven-go"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type grpcServer struct {
	add grpctransport.Handler
}

func NewGRPCServer(endpoints endpoint.Set, tracer opentracing.Tracer, logger log.Logger, raven *raven.Client) pb.AddServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
		grpctransport.ServerFinalizer(grpctransport.SentryServerFinalizer(raven)),
	}

	return &grpcServer{
		add: grpctransport.NewServer(
			endpoints.AddEndpoint,
			decodeGRPCAddRequest,
			encodeGRPCAddResponse,
			append(options, grpctransport.ServerBefore(tracing.GRPCToContext(tracer, "Add", logger)))...,
		),
	}
}

func (s *grpcServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	_, rep, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.AddResponse), nil
}

func NewGRPCClient(conn *grpc.ClientConn, tracer opentracing.Tracer, logger log.Logger, raven *raven.Client) service.Service {
	options := []grpctransport.ClientOption{
		grpctransport.ClientFinalizer(grpctransport.SentryClientFinalizer(raven)),
	}

	var addEndpoint flexitendpoint.Endpoint
	{
		addEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"Add",
			encodeGRPCAddRequest,
			decodeGRPCAddResponse,
			pb.AddResponse{},
			append(options, grpctransport.ClientBefore(tracing.ContextToGRPC(tracer, logger)))...,
		).Endpoint()
		addEndpoint = tracing.TraceClient(tracer, "Add")(addEndpoint)
	}

	return &endpoint.Set{
		AddEndpoint: addEndpoint,
	}
}

func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AddRequest)
	return endpoint.AddRequest{X: int64(req.X), Y: int64(req.Y)}, nil
}

func encodeGRPCAddResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.AddResponse)
	return &pb.AddResponse{Sum: int64(resp.Sum), Err: err2str(resp.Err)}, nil
}

func encodeGRPCAddRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoint.AddRequest)
	return &pb.AddRequest{X: int64(req.X), Y: int64(req.Y)}, nil
}

func decodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.AddResponse)
	return endpoint.AddResponse{Sum: int64(reply.Sum), Err: str2err(reply.Err)}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
