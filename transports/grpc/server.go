package grpc

import (
	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// NewBaseServer returns a new grpc server
func NewBaseServer(
	before []ServerRequestFunc,
	after []ServerResponseFunc,
	finalizer []ServerFinalizerFunc,
) *grpc.Server {

	server := grpc.NewServer(
		WithUnaryServerChain(
			ServerMiddlewareInterceptor(before, after, finalizer),
		),
		WithStreamServerChain(),
	)

	return server
}

// NewServer creates a new grpc serve, use this function in most cases.
func NewServer(
	tracer opentracing.Tracer,
	sentry *raven.Client,
) *grpc.Server {
	return NewBaseServer(
		[]ServerRequestFunc{
			GRPCToContext(tracer),
		},
		nil,
		[]ServerFinalizerFunc{
			SentryServerFinalizer(sentry),
		},
	)
}
