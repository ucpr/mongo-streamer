package mongo

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ucpr/mongo-streamer/internal/config"
)

// Set is a Wire provider set that provides a MongoDB client.
var Set = wire.NewSet(
	NewClient,
)

// Client is a MongoDB client.
type Client struct {
	cli *mongo.Client
	db  string
}

// NewClient creates a new MongoDB client.
func NewClient(ctx context.Context, cfg *config.MongoDB) (*Client, error) {
	api := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetServerAPIOptions(api)

	// Set authentication options if user and password are provided.
	if cfg.User != "" && cfg.Password != "" {
		opts.SetAuth(options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		})
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB instance: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB instance: %w", err)
	}

	return &Client{
		cli: client,
		db:  cfg.Database,
	}, nil
}

// Disconnect disconnects the client from the MongoDB instance.
func (c *Client) Disconnect(ctx context.Context) error {
	return c.cli.Disconnect(ctx)
}

// Collection returns a collection from the MongoDB client.
func (c *Client) Collection(name string) *mongo.Collection {
	return c.cli.Database(c.db).Collection(name)
}
