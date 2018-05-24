package cmd

import (
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/eirsyl/flexit/app"
	"github.com/eirsyl/flexit/cmd"
	"github.com/eirsyl/flexit/debugserver"
	"github.com/eirsyl/flexit/examples/addsvc/pb"
	"github.com/eirsyl/flexit/examples/addsvc/pkg/endpoint"
	"github.com/eirsyl/flexit/examples/addsvc/pkg/service"
	"github.com/eirsyl/flexit/examples/addsvc/pkg/transport"
	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/metrics"
	"github.com/eirsyl/flexit/runtime"
	"github.com/eirsyl/flexit/tracing/jaeger"
	"github.com/getsentry/raven-go"
	"github.com/oklog/oklog/pkg/group"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var App = app.NewApp("eirsyl.flexit.addsvc", "Service for adding numbers")

func init() {
	cmd.StringConfig(RootCmd, "debugAddr", "", ":8080", "debug server listen addr")
	cmd.StringConfig(RootCmd, "grpcAddr", "", ":8090", "add api grpc server listen addr")
	cmd.StringConfig(RootCmd, "jaegerAddr", "", "127.0.0.1:6831", "jaeger agent addr")
	cmd.StringConfig(RootCmd, "sentryDsn", "", "", "sentry raven dsn")

	RootCmd.AddCommand(clientCmd)
}

var RootCmd = &cobra.Command{
	Use:   App.GetShortName(),
	Short: App.GetDescription(),
	PreRun: func(_ *cobra.Command, args []string) {
		// Validate flags
		cmd.CheckFlags(
			cmd.RequireString("debugAddr"),
			cmd.RequireString("grpcAddr"),
			cmd.RequireString("jaegerAddr"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		// Initialize logger
		var logger log.Logger
		{
			logger = log.NewLogrusLogger(true).WithField("app", App.GetName())
		}

		// Get flags
		var debugAddr, grpcAddr, jaegerAddr, sentryDsn string
		{
			debugAddr = viper.GetString("debugAddr")
			grpcAddr = viper.GetString("grpcAddr")
			jaegerAddr = viper.GetString("jaegerAddr")
			sentryDsn = viper.GetString("sentryDsn")
		}

		// Optimize runtime
		runtime.OptimizeRuntime(logger)

		// Instrumentation
		var additions metrics.Counter
		{
			// Business-level metrics.
			additions = metrics.NewCounterFrom(prometheus.CounterOpts{
				Namespace: App.GetShortName(),
				Subsystem: App.GetShortName(),
				Name:      "add_requests",
				Help:      "Total count of add requests.",
			}, []string{})
		}
		var duration metrics.Histogram
		{
			// Endpoint-level metrics.
			duration = metrics.NewSummaryFrom(prometheus.SummaryOpts{
				Namespace: App.GetShortName(),
				Subsystem: App.GetShortName(),
				Name:      "request_duration_seconds",
				Help:      "Request duration in seconds.",
			}, []string{"method", "success"})
		}

		// Initialize tracer
		var tracer opentracing.Tracer
		var closer io.Closer
		{
			var err error
			tracer, closer, err = jaeger.New(jaegerAddr, App.GetName(), true)
			if err != nil {
				logger.Error("err", err)
				os.Exit(1)
			}
		}
		defer closer.Close()

		var ravenClient *raven.Client
		{
			var err error
			ravenClient, err = raven.NewWithTags(sentryDsn, map[string]string{"service": App.GetName()})
			if err != nil {
				logger.Error("err", err)
				os.Exit(1)
			}
		}

		var g group.Group
		{
			debugLogger := logger.WithFields(&log.Fields{
				"transport": "debug/http",
				"addr":      debugAddr,
			})

			// Initialize debugserver
			var (
				debugServer = debugserver.New(debugLogger)
			)

			debugListener, err := net.Listen("tcp", debugAddr)
			if err != nil {
				debugLogger.Errorf("could not create listener: %v", err)
				os.Exit(1)
			}

			g.Add(func() error {
				debugLogger.Info("listening")
				return http.Serve(debugListener, debugServer)
			}, func(error) {
				debugListener.Close()
			})
		}
		{
			grpcLogger := logger.WithFields(&log.Fields{
				"transport": "addsvc/grpc",
				"addr":      grpcAddr,
			})

			// Initialize services
			var (
				service   = service.New(grpcLogger, additions)
				endpoint  = endpoint.New(service, grpcLogger, tracer, duration, ravenClient)
				transport = transport.NewGRPCServer(endpoint, tracer, grpcLogger, ravenClient)
			)

			grpcListener, err := net.Listen("tcp", grpcAddr)
			if err != nil {
				grpcLogger.Errorf("could not create listener: %v", err)
				os.Exit(1)
			}

			g.Add(func() error {
				grpcLogger.Info("listening")
				baseServer := grpc.NewServer()
				pb.RegisterAddServer(baseServer, transport)
				return baseServer.Serve(grpcListener)
			}, func(error) {
				grpcListener.Close()
			})
		}
		{
			cancelInterrupt := make(chan struct{})
			g.Add(func() error {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				select {
				case sig := <-c:
					logger.Errorf("received signal %s", sig)
					return nil
				case <-cancelInterrupt:
					return nil
				}
			}, func(error) {
				close(cancelInterrupt)
			})
		}
		return g.Run()

	},
}
