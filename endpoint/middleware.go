package endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/metrics"
	raven "github.com/getsentry/raven-go"
)

// Compose composes a endpoint from a enpoint and a slice of middlewares
func Compose(next Endpoint, middlewares ...Middleware) (endpoint Endpoint) {
	endpoint = next
	for _, middleware := range middlewares {
		endpoint = middleware(endpoint)
	}
	return
}

// InstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func InstrumentingMiddleware(duration metrics.Histogram) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Endpoint) Endpoint {
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
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err != nil {
				fmt.Println("Sending error")
				client.CaptureError(err, nil)
			}
			return next(ctx, request)
		}
	}
}
