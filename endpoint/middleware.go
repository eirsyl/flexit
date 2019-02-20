package endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/metrics"
	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// Compose composes a endpoint from a enpoint and a slice of middlewares
func Compose(next Endpoint, name string, middlewares ...Middleware) (endpoint Endpoint) {
	endpoint = next
	for _, middleware := range middlewares {
		endpoint = middleware(name, endpoint)
	}
	return
}

// InstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func InstrumentingMiddleware(duration metrics.Histogram) Middleware {
	return func(name string, next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				duration.With("method", name, "success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(name string, next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.WithFields(&log.Fields{
					"transport_error": err,
					"took":            time.Since(begin),
				}).Info("request")
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// SentryMiddleware sends endpoint errors to sentry
func SentryMiddleware(client *raven.Client) Middleware {
	return func(name string, next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err != nil {
				fmt.Println("Sending error")
				client.CaptureError(err, nil)
			}
			return next(ctx, request)
		}
	}
}

// TraceServer returns a Middleware that wraps the `next` Endpoint in an
// OpenTracing Span called `operationName`.
//
// If `ctx` already has a Span, it is re-used and the operation name is
// overwritten. If `ctx` does not yet have a Span, one is created here.
func TraceServer(tracer opentracing.Tracer) Middleware {
	return func(name string, next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			serverSpan := opentracing.SpanFromContext(ctx)
			if serverSpan == nil {
				// All we can do is create a new root span.
				serverSpan = tracer.StartSpan(name)
			} else {
				serverSpan.SetOperationName(name)
			}
			defer func() {
				if err != nil {
					serverSpan.SetTag("error", true)
					serverSpan.SetTag("message", err.Error())
				}
				serverSpan.Finish()
			}()
			ext.SpanKindRPCServer.Set(serverSpan)
			ctx = opentracing.ContextWithSpan(ctx, serverSpan)
			return next(ctx, request)
		}
	}
}

// TraceClient returns a Middleware that wraps the `next` Endpoint in an
// OpenTracing Span called `operationName`.
func TraceClient(tracer opentracing.Tracer) Middleware {
	return func(name string, next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var clientSpan opentracing.Span
			if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
				clientSpan = tracer.StartSpan(
					name,
					opentracing.ChildOf(parentSpan.Context()),
				)
			} else {
				clientSpan = tracer.StartSpan(name)
			}
			defer func() {
				if err != nil {
					clientSpan.SetTag("error", true)
					clientSpan.SetTag("message", err.Error())
				}
				clientSpan.Finish()
			}()
			ext.SpanKindRPCClient.Set(clientSpan)
			ctx = opentracing.ContextWithSpan(ctx, clientSpan)
			return next(ctx, request)
		}
	}
}
