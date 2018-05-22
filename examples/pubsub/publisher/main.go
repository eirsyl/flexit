package main

import (
	"context"
	"fmt"
	"github.com/eirsyl/flexit/pubsub"
	"github.com/eirsyl/flexit/pubsub/nats"
)

func main() {

	var pub pubsub.Publisher
	var err error

	pub, err = nats.NewPublisher(&nats.Config{
		ClusterID: "test-cluster",
		ClientID:  "publisher-client",
		Urls:      "nats://127.0.0.1:4222",
	})

	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		pub.PublishRaw(context.Background(), "foo", []byte(fmt.Sprintf("Message %d", i)))
	}

	err = pub.Close()
	if err != nil {
		panic(err)
	}

}
