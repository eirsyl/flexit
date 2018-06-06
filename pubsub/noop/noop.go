package noop

import (
	"context"
	"github.com/eirsyl/flexit/pubsub"
	"github.com/golang/protobuf/proto"
	"time"
)

// Publisher

type noopPublisher struct{}

func NewNoopPublisher() (pubsub.Publisher, error) {
	return &noopPublisher{}, nil
}

func (p *noopPublisher) Publish(context.Context, string, proto.Message) error {
	return nil
}

func (p *noopPublisher) PublishRaw(context.Context, string, []byte) error {
	return nil
}

func (p *noopPublisher) Close() error {
	return nil
}

// Subscriber

type noopSubscriber struct{}

type NoopSubscriberMessage struct{}

func (sm *NoopSubscriberMessage) Message() []byte {
	return []byte{}
}

func ExtendDoneDeadline(time.Duration) error {
	return nil
}

func Done() error {
	return nil
}

func NewNoopSubscriber() (pubsub.Subscriber, error) {
	return &noopSubscriber{}, nil
}

func (s *noopSubscriber) Start() <-chan pubsub.SubscriberMessage {
	return make(chan pubsub.SubscriberMessage)
}

func (s *noopSubscriber) Err() error {
	return nil
}

func (s *noopSubscriber) Close() error {
	return nil
}
