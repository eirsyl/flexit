package runtime

import (
	"github.com/eirsyl/flexit/log"
	"go.uber.org/automaxprocs/maxprocs"
)

func OptimizeRuntime(logger log.Logger) {
	maxprocs.Set(maxprocs.Logger(logger.Infof))
}
