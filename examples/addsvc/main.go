package main

import (
	"github.com/eirsyl/flexit/app"
	"fmt"
	"os"
	"github.com/eirsyl/flexit/examples/addsvc/pkg/service"
)

func main() {
	a := app.New("io.whale.webhook")
	srv := service.NewWebhookService()
	if err := a.Run(srv); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
