package grpc

import (
	"google.golang.org/grpc"
)

// NewServer returns a new grpc server
func NewServer(
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
