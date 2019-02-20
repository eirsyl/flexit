package grpc

import (
	"google.golang.org/grpc"
)

// NewClient returns a new grpc client
func NewClient(
	target string,
	before []ClientRequestFunc,
	after []ClientResponseFunc,
	finalizer []ClientFinalizerFunc,
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {

	options := []grpc.DialOption{
		WithUnaryClientChain(
			ClientMiddlewareInterceptor(before, after, finalizer),
		),
		WithStreamClientChain(),
	}

	conn, err := grpc.Dial(target, append(opts, options...)...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
