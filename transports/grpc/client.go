package grpc

import (
	"context"

	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// NewBasicClient returns a new basic grpc client
func NewBasicClient(
	ctx context.Context,
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

	conn, err := grpc.DialContext(ctx, target, append(opts, options...)...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// NewClient returns a new grpc client with tracing and sentry built in. Use this function in most cases
func NewClient(
	ctx context.Context,
	target string,
	tracer opentracing.Tracer,
	sentry *raven.Client,
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	return NewBasicClient(
		ctx,
		target,
		[]ClientRequestFunc{
			ContextToGRPC(tracer),
		},
		nil,
		[]ClientFinalizerFunc{
			SentryClientFinalizer(sentry),
		},
		opts...,
	)
}
