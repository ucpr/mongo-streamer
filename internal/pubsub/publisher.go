package pubsub

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
)

const (
	publisherByteThreshold  = 5000
	publisherCountThreshold = 10
	publisherDelayThreshold = 100 * time.Millisecond
)

type Message struct {
	Data        []byte
	Attributes  map[string]string
	OrderingKey string
}

// Pulisher is an interface for PubSub Publisher.
type Publisher interface {
	AsyncPublish(ctx context.Context, msg Message) *pubsub.PublishResult
}

// PubSubPublisher is a publisher for Google Cloud Pub/Sub.
type PubSubPublisher struct {
	cli   *pubsub.Client
	topic *pubsub.Topic
}

// Ensure that PubSubPublisher implements Publisher.
//
//nolint:gochecknoglobals
var _ Publisher = (*PubSubPublisher)(nil)

// NewPublisher creates a new publisher.
func NewPublisher(ctx context.Context, projectID, topicID string) (*PubSubPublisher, error) {
	cli, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	topic := cli.Topic(topicID)
	// Set the default values for PublishSettings.
	topic.PublishSettings.ByteThreshold = publisherByteThreshold
	topic.PublishSettings.CountThreshold = publisherCountThreshold
	topic.PublishSettings.DelayThreshold = publisherDelayThreshold

	return &PubSubPublisher{
		cli:   cli,
		topic: topic,
	}, nil
}

// Publish publishes a message to the topic.
func (p *PubSubPublisher) AsyncPublish(ctx context.Context, msg Message) *pubsub.PublishResult {
	result := p.topic.Publish(ctx, &pubsub.Message{
		Data:        msg.Data,
		Attributes:  msg.Attributes,
		OrderingKey: msg.OrderingKey,
	})
	return result
}

// Close closes the publisher.
func (p *PubSubPublisher) Close() error {
	return p.cli.Close()
}
