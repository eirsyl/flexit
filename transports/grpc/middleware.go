package grpc

import (
	"context"
	"encoding/base64"
	"strings"

	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"
)

// SentryServerFinalizer sends an error to sentry if a grpc requests ended up with an error
func SentryServerFinalizer(client *raven.Client) ServerFinalizerFunc {
	return func(ctx context.Context, err error) {
		if err != nil {
			client.CaptureErrorAndWait(err, map[string]string{"transport": "grpc"})
		}
	}
}

// SentryClientFinalizer sends an error to sentry if a grpc requests ended up with an error
func SentryClientFinalizer(client *raven.Client) ClientFinalizerFunc {
	return func(ctx context.Context, err error) {
		if err != nil {
			client.CaptureErrorAndWait(err, map[string]string{"transport": "grpc"})
		}
	}
}

// ContextToGRPC returns a grpc RequestFunc that injects an OpenTracing Span
// found in `ctx` into the grpc Metadata. If no such Span can be found, the
// RequestFunc is a noop.
func ContextToGRPC(tracer opentracing.Tracer) ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			tracer.Inject(span.Context(), opentracing.HTTPHeaders, metadataReaderWriter{md}) // nolint: errcheck gas
		}
		return ctx
	}
}

// GRPCToContext returns a grpc RequestFunc that tries to join with an
// OpenTracing trace found in `req` and starts a new Span called
// `operationName` accordingly. If no trace could be found in `req`, the Span
// will be a trace root. The Span is incorporated in the returned Context and
// can be retrieved with opentracing.SpanFromContext(ctx).
func GRPCToContext(tracer opentracing.Tracer) ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		var span opentracing.Span

		wireContext, err := tracer.Extract(opentracing.HTTPHeaders, metadataReaderWriter{&md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			return ctx
		}

		operationName := ctx.Value(ContextKeyRequestMethod).(string)
		span = tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))

		return opentracing.ContextWithSpan(ctx, span)
	}
}

// A type that conforms to opentracing.TextMapReader and
// opentracing.TextMapWriter.
type metadataReaderWriter struct {
	*metadata.MD
}

func (w metadataReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	if strings.HasSuffix(key, "-bin") {
		val = string(base64.StdEncoding.EncodeToString([]byte(val)))
	}
	(*w.MD)[key] = append((*w.MD)[key], val)
}

func (w metadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range *w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}
