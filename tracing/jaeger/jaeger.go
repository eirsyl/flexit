package jaeger

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

// New creates a new jaeger tracer, zipkin header propagation is optional (istio environments)
func New(host string, service string, zipkinHeaders bool) (opentracing.Tracer, io.Closer, error) {
	var opts []jaeger.TracerOption

	if zipkinHeaders {
		zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
		injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
		extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)
		zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)
		opts = []jaeger.TracerOption{
			injector,
			extractor,
			zipkinSharedRPCSpan,
		}
	}

	// Metrics
	metricsFactory := prometheus.New()
	metrics := jaeger.NewMetrics(metricsFactory, nil)
	opts = append(opts, jaeger.TracerOptions.Metrics(metrics))

	sender, err := jaeger.NewUDPTransport(host, 0)
	if err != nil {
		return nil, nil, err
	}

	tracer, closer := jaeger.NewTracer(
		service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
		),
		opts...,
	)

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}
