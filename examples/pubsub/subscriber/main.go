package main

import (
	"fmt"
	"github.com/eirsyl/flexit/pubsub"
	"github.com/eirsyl/flexit/pubsub/nats"
	"github.com/nats-io/go-nats-streaming"
	"time"
)

func main() {

	var sub pubsub.Subscriber
	var err error

	sub, err = nats.NewSubscriber(&nats.Config{
		ClusterID: "test-cluster",
		ClientID:  "worker",
		Urls:      "nats://127.0.0.1:4222",
	}, "foo", "workers", stan.DurableName("workers"))

	if err != nil {
		panic(err)
	}

	messages := sub.Start()

	go func() {
		for {
			msg := <-messages
			if msg != nil {
				fmt.Printf("Received message: %s \n", msg.Message())
				msg.Done()
			}
		}
	}()

	// Wait 30 secs and close
	time.Sleep(30 * time.Second)
	err = sub.Close()
	if err != nil {
		panic(err)
	}

}
