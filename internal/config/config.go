package config

import (
	"context"

	"github.com/google/wire"
	envconfig "github.com/sethvargo/go-envconfig"
)

// Set is a Wire provider set that provides configuration.
var Set = wire.NewSet(
	NewMongoDB,
	NewPubSub,
	NewMetrics,
)

const (
	mongoDBPrefix = "MONGO_DB_"
	pubSubPrefix  = "PUBSUB_"
	mrtricsPrefix = "METRICS_"
)

type MongoDB struct {
	URI        string `env:"URI"`
	Password   string `env:"PASSWORD"`
	User       string `env:"USER"`
	Database   string `env:"DATABASE"`
	Collection string `env:"COLLECTION"`
}

type PubSub struct {
	// ProjectID is the id of the Google Cloud project to publish messages to.
	ProjectID string `env:"PROJECT_ID"`
	// TopicID is the id of the topic to publish messages to.
	TopicID string `env:"TOPIC_ID"`
	// PublishFormat is the format of the message to publish.
	// Supported format are: json, avro.
	PublishFormat string `env:"PUBLISH_FORMAT, default=json"`
}

type Metrics struct {
	Addr string `env:"ADDR"`
}

func NewMongoDB(ctx context.Context) (*MongoDB, error) {
	conf := &MongoDB{}
	pl := envconfig.PrefixLookuper(mongoDBPrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		return nil, err
	}

	return conf, nil
}

func NewPubSub(ctx context.Context) (*PubSub, error) {
	conf := &PubSub{}
	pl := envconfig.PrefixLookuper(pubSubPrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		return nil, err
	}

	return conf, nil
}

func NewMetrics(ctx context.Context) (*Metrics, error) {
	conf := &Metrics{}
	pl := envconfig.PrefixLookuper(mrtricsPrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		return nil, err
	}

	return conf, nil
}
