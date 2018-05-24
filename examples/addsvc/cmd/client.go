package cmd

import (
	"io"
	"os"

	"fmt"
	"time"

	"context"

	"github.com/eirsyl/flexit/examples/addsvc/pkg/service"
	"github.com/eirsyl/flexit/examples/addsvc/pkg/transport"
	"github.com/eirsyl/flexit/log"
	"github.com/eirsyl/flexit/tracing/jaeger"
	"github.com/getsentry/raven-go"
	"github.com/opentracing/opentracing-go"
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

		conn, err := grpc.Dial("127.0.0.1:8090", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		svc = transport.NewGRPCClient(conn, tracer, log.NewLogrusLogger(true), ravenClient)

		sum, err := svc.Add(context.Background(), 10, 20)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		logger.Infof("Sum: %v", sum)

		return nil
	},
}
