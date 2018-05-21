package nats

import (
	"context"
	"errors"
	"github.com/nats-io/go-nats-streaming"
	"github.com/eirsyl/flexit/pubsub"
	"github.com/golang/protobuf/proto"
	"time"
)

var (
	ClusterIDRequired = errors.New("clusterid required")
	ClientIDRequired  = errors.New("clientid required")
	UrlsRequired      = errors.New("urls required")
	SubjectRequired = errors.New("subject required")
)

// PUBLISHER

type natsPublisher struct {
	sc stan.Conn
}

func NewPublisher(cfg *Config) (pubsub.Publisher, error) {
	p := &natsPublisher{}

	if cfg.ClusterID == "" {
		return p, ClusterIDRequired
	}
	if cfg.ClientID == "" {
		return p, ClientIDRequired
	}
	if cfg.Urls == "" {
		return p, UrlsRequired
	}

	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(cfg.Urls))
	if err != nil {
		return p, err
	}

	p.sc = sc

	return p, nil
}

func (p *natsPublisher) Publish(ctx context.Context, key string, m proto.Message) error {
	mb, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return p.PublishRaw(ctx, key, mb)
}

func (p *natsPublisher) PublishRaw(_ context.Context, key string, m []byte) error {
	return p.sc.Publish(key, m)
}

func (p *natsPublisher) Close() error {
	return p.sc.Close()
}

// SUBSCRIBER

type natsSubscriberMessage struct {
	message []byte
	sec uint64
	ack func() error
}

func(sm *natsSubscriberMessage) Message() []byte {
	return sm.message
}

func(sm *natsSubscriberMessage) ExtendDoneDeadline(duration time.Duration) error {
	return nil
}

func(sm *natsSubscriberMessage) Done() error {
	return sm.ack()
}

type natsSubscriber struct {
	sc stan.Conn
	so stan.SubscriptionOption

	subject string
	group string
	kerr error
	stop chan chan error
}

func NewSubscriber(cfg *Config, subject, group string, so stan.SubscriptionOption) (pubsub.Subscriber, error) {
	s := &natsSubscriber{
		stop: make(chan chan error, 1),
	}

	if cfg.ClusterID == "" {
		return s, ClusterIDRequired
	}
	if cfg.ClientID == "" {
		return s, ClientIDRequired
	}
	if cfg.Urls == "" {
		return s, UrlsRequired
	}
	if subject == "" {
		return s, SubjectRequired
	}

	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(cfg.Urls))
	if err != nil {
		return s, err
	}

	s.sc = sc
	s.so = so
	s.subject = subject
	s.group = group

	stan.StartWithLastReceived()
	return s, nil
}

func (s *natsSubscriber) Start() <- chan pubsub.SubscriberMessage {
	output := make(chan pubsub.SubscriberMessage)

	go func(s *natsSubscriber, sc stan.Conn, so stan.SubscriptionOption) {
		defer close(output)

		opts := []stan.SubscriptionOption{
			stan.SetManualAckMode(),
			stan.AckWait(60*time.Second),
			so,
		}

		handler := func(m *stan.Msg) {
			output <- &natsSubscriberMessage{
				message: m.Data,
				sec: m.Sequence,
				ack: m.Ack,
			}
		}

		var sub stan.Subscription
		var err error
		if s.group == "" {
			sub, err = sc.Subscribe(s.subject, handler, opts...)

		} else {
			sub, err = sc.QueueSubscribe(
				s.subject,
				s.group,
				handler,
				opts...
			)
		}

		if err != nil {
			s.kerr = err
			return
		}

		select {
			case exit := <- s.stop:
				exit <- sub.Close()
		}

	}(s, s.sc, s.so)

	return output
}

func (s *natsSubscriber) Err() error {
	return s.kerr
}

func (s *natsSubscriber) Close() error {
	exit := make(chan error)
	s.stop <- exit
	err := <-exit
	if err != nil {
		return err
	}
	return s.sc.Close()
}
