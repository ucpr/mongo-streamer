package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ucpr/mongo-streamer/internal/config"
)

type Client struct {
	cli *mongo.Client
	db  string
}

func NewClient(ctx context.Context, cfg config.MongoDB) (*Client, error) {
	api := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetAuth(options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		}).
		SetServerAPIOptions(api)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB instance: %w", err)
	}

	return &Client{
		cli: client,
		db:  cfg.Database,
	}, nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.cli.Disconnect(ctx)
}
