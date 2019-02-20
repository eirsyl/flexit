package cmd

import (
	"io"
	"os"

	flexitgrpc "github.com/eirsyl/flexit/transports/grpc"

	"github.com/eirsyl/flexit/examples/simple/pb"

	"fmt"
	"time"

	"context"

	"github.com/eirsyl/flexit/examples/simple/pkg/service"
	"github.com/eirsyl/flexit/examples/simple/pkg/transportgrpc"
	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/tracing/jaeger"
	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var clientCmd = &cobra.Command{
	Use: "client",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Initialize logger
		var logger log.Logger
		{
			logger = log.NewLogrusLogger(true).WithField("app", App.GetName())
		}

		// Initialize tracer
		var tracer opentracing.Tracer
		var closer io.Closer
		{
			var err error
			tracer, closer, err = jaeger.New("127.0.0.1:6831", "client", true)
			if err != nil {
				logger.Error("err", err)
				os.Exit(1)
			}
		}
		defer closer.Close()

		var ravenClient *raven.Client
		{
			var err error
			ravenClient, err = raven.NewWithTags("", nil)
			if err != nil {
				logger.Error("err", err)
				os.Exit(1)
			}
		}

		var (
			svc service.Service
			err error
		)

		conn, err := flexitgrpc.NewClient(
			"127.0.0.1:8090",
			[]flexitgrpc.ClientRequestFunc{
				flexitgrpc.ContextToGRPC(tracer, logger),
			},
			nil,
			[]flexitgrpc.ClientFinalizerFunc{flexitgrpc.SentryClientFinalizer(ravenClient)},
			grpc.WithInsecure(),
			grpc.WithTimeout(time.Second),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()

		svc = transportgrpc.NewGRPCClient(conn, tracer, log.NewLogrusLogger(true), ravenClient)

		sum, err := svc.Add(context.Background(), &pb.AddRequest{X: 100, Y: 200})
		if err != nil {
			logger.Error(err)
		}

		logger.Infof("Sum: %v", sum)

		div, err := svc.Subtract(context.Background(), &pb.SubtractRequest{X: 100, Y: 100})
		if err != nil {
			logger.Error(err)
		}

		logger.Infof("Sum: %v", div)

		return nil
	},
}
