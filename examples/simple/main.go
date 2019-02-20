package main

import (
	"os"

	"github.com/eirsyl/flexit/examples/simple/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
