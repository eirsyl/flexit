package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eirsyl/flexit/log"
)

// Interrupt returns a exeute and interrupt function used by oklog/run to manage
// goroutine lifecycle of an app.
func Interrupt(logger log.Logger) (func() error, func(error)) {
	cancelInterrupt := make(chan struct{})

	execute := func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			logger.Errorf("received signal %s", sig)
			return nil
		case <-cancelInterrupt:
			return nil
		}
	}

	interrupt := func(error) {
		close(cancelInterrupt)
	}

	return execute, interrupt
}
