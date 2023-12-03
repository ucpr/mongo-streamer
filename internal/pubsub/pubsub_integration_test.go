//go:build integration

package pubsub

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
)

const (
	testProjectID      = "dummy-project"
	testTopicID        = "dyummy-topic"
	testSubscriptionID = "dummy-subscription"
)

func TestMain(m *testing.M) {
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")

	cli, err := pubsub.NewClient(context.Background(), testProjectID)
	if err != nil {
		panic(err)
	}
	// check topic exists
	if exists, _ := cli.Topic(testTopicID).Exists(context.Background()); exists {
		return
	}
	// create topic
	topic, err := cli.CreateTopic(context.Background(), testTopicID)
	if err != nil {
		panic(err)
	}
	// check subscription exists
	if exists, _ := cli.Subscription(testSubscriptionID).Exists(context.Background()); exists {
		return
	}
	// create subscription
	_, err = cli.CreateSubscription(context.Background(), testSubscriptionID, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 10 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	m.Run()
}
