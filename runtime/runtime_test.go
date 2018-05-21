package runtime

import (
	"github.com/eirsyl/flexit/log"
	"testing"
)

func TestOptimizeRuntime(t *testing.T) {
	logger := log.NewLogrusLogger(true)
	OptimizeRuntime(logger)
}
