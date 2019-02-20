package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ServerMethodNameInterceptor is a grpc UnaryInterceptor that injects the method name into
// context so it can be consumed by middlewares.
func ServerMethodNameInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx = context.WithValue(ctx, ContextKeyRequestMethod, info.FullMethod)
	return handler(ctx, req)
}

// ServerMiddlewareInterceptor intercepts server requests and executes before, after and finalizer functions
func ServerMiddlewareInterceptor(
	before []ServerRequestFunc,
	after []ServerResponseFunc,
	finalizer []ServerFinalizerFunc,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}

		if len(finalizer) > 0 {
			defer func() {
				for _, f := range finalizer {
					f(ctx, err)
				}
			}()
		}

		for _, f := range before {
			ctx = f(ctx, md)
		}

		resp, err = handler(ctx, req)

		var mdHeader, mdTrailer metadata.MD
		for _, f := range after {
			ctx = f(ctx, &mdHeader, &mdTrailer)
		}

		if len(mdHeader) > 0 {
			if err = grpc.SendHeader(ctx, mdHeader); err != nil {
				return nil, err
			}
		}

		if len(mdTrailer) > 0 {
			if err = grpc.SetTrailer(ctx, mdTrailer); err != nil {
				return nil, err
			}
		}

		return resp, err
	}
}

// ClientMiddlewareInterceptor intercepts client requests and executes before, after and finalizer functions
func ClientMiddlewareInterceptor(
	before []ClientRequestFunc,
	after []ClientResponseFunc,
	finalizer []ClientFinalizerFunc,
) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		if len(finalizer) > 0 {
			defer func() {
				for _, f := range finalizer {
					f(ctx, err)
				}
			}()
		}

		ctx = context.WithValue(ctx, ContextKeyRequestMethod, method)

		md := &metadata.MD{}
		for _, f := range before {
			ctx = f(ctx, md)
		}
		ctx = metadata.NewOutgoingContext(ctx, *md)

		var header, trailer metadata.MD
		callOptions := []grpc.CallOption{
			grpc.Header(&header),
			grpc.Trailer(&trailer),
		}
		err = invoker(
			ctx, method, req, reply, cc, append(callOptions, opts...)...,
		)

		for _, f := range after {
			ctx = f(ctx, header, trailer)
		}

		return err
	}
}
