package main

import (
	"github.com/eirsyl/flexit/examples/addsvc/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
