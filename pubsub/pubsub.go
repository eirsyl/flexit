package pubsub

import (
	"context"
	"github.com/golang/protobuf/proto"
	"time"
)

type Publisher interface {
	// Publish will publish a message with context.
	Publish(context.Context, string, proto.Message) error
	// Publish will publish a raw byte array as a message with context.
	PublishRaw(context.Context, string, []byte) error
	// Close closes the underlying messaging provider
	Close() error
}

type MultiPublisher interface {
	Publisher
	// PublishMulti will publish multiple messages with a context.
	PublishMulti(context.Context, []string, []proto.Message) error
	// PublishMultiRaw will publish multiple raw byte array messages with a context.
	PublishMultiRaw(context.Context, []string, [][]byte) error
}

type Subscriber interface {
	// Start will return a channel of raw messages.
	Start() <-chan SubscriberMessage
	// Err will contain any errors returned from the consumer connection.
	Err() error
	// Close will initiate a graceful shutdown of the subscriber connection.
	Close() error
}

type SubscriberMessage interface {
	Message() []byte
	ExtendDoneDeadline(time.Duration) error
	Done() error
}
