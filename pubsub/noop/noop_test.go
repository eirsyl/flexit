package noop

import (
	"context"
	"testing"
)

func TestPublisher(t *testing.T) {
	publisher, err := NewNoopPublisher()
	if err != nil {
		t.Error(err)
	}

	if err = publisher.Publish(context.Background(), "chan", nil); err != nil {
		t.Error(err)
	}

	if err = publisher.PublishRaw(context.Background(), "chan", []byte{}); err != nil {
		t.Error(err)
	}

	if err = publisher.Close(); err != nil {
		t.Error(err)
	}
}

func TestSubscriber(t *testing.T) {
	subscriber, err := NewNoopSubscriber()
	if err != nil {
		t.Error(err)
	}

	if err = subscriber.Close(); err != nil {
		t.Error(err)
	}
}
