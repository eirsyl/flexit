package grpc

import (
	"context"
	"github.com/getsentry/raven-go"
)

func SentryServerFinalizer(client *raven.Client) ServerFinalizerFunc {
	return func(ctx context.Context, err error) {
		if err != nil {
			client.CaptureError(err, map[string]string{"transport": "grpc"})
		}
	}
}

func SentryClientFinalizer(client *raven.Client) ClientFinalizerFunc {
	return func(ctx context.Context, err error) {
		if err != nil {
			client.CaptureErrorAndWait(err, map[string]string{"transport": "grpc"})
		}
	}
}
