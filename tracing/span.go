package tracing

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

var tracer = opentracing.GlobalTracer()

// MaybeStartSpanFromContext creates a span from a given context, a noop tracer is returned if no tracer is found
func MaybeStartSpanFromContext(
	ctx context.Context,
	operationName string,
	opts ...opentracing.StartSpanOption,
) (opentracing.Span, context.Context) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span, ctx := opentracing.StartSpanFromContext(ctx, operationName, opts...)
		return span, ctx
	} else {
		span := tracer.StartSpan(operationName)
		return span, ctx
	}
}
