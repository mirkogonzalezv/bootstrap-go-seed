package pubsub

import (
	"context"
	"fmt"
	"sync"

	gcppubsub "cloud.google.com/go/pubsub/v2"
	"google.golang.org/api/option"
)

type gcpPubSub struct {
	client *gcppubsub.Client
	mu     sync.Mutex
	pubs   map[string]*gcppubsub.Publisher
}

func NewGCPPubSub(ctx context.Context, projectID string, opts ...option.ClientOption) (PubSub, error) {
	client, err := gcppubsub.NewClient(ctx, projectID, opts...)

	if err != nil {
		return nil, fmt.Errorf("pubsub gcp: %w", err)
	}

	return &gcpPubSub{client: client, pubs: make(map[string]*gcppubsub.Publisher)}, nil
}

func (p *gcpPubSub) publisher(topic string) *gcppubsub.Publisher {
	p.mu.Lock()

	defer p.mu.Unlock()

	pub, ok := p.pubs[topic]
	if !ok {
		pub = p.client.Publisher(topic)
		p.pubs[topic] = pub
	}

	return pub
}

func (p *gcpPubSub) Publish(ctx context.Context, topic string, data []byte, attrs map[string]string) error {
	pub := p.publisher(topic)

	result := pub.Publish(ctx, &gcppubsub.Message{
		Data:       data,
		Attributes: attrs,
	})

	_, err := result.Get(ctx)
	return err
}

func (p *gcpPubSub) Subscribe(ctx context.Context, subscription string, handler MessageHandler) error {
	sub := p.client.Subscriber(subscription)

	return sub.Receive(ctx, func(c context.Context, msg *gcppubsub.Message) {
		m := &Message{
			ID:         msg.ID,
			Data:       msg.Data,
			Attributes: msg.Attributes,
			AckFunc:    msg.Ack,
			NackFunc:   msg.Nack,
		}
		if err := handler(c, m); err != nil {
			msg.Nack()
			return
		}
		msg.Ack()
	})
}

func (p *gcpPubSub) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pub := range p.pubs {
		pub.Stop()
	}
	p.client.Close()
	return nil
}

func (p *gcpPubSub) Native() interface{} {
	return p.client
}
