package mongo

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	// ChangeStream is a struct that represents a change stream.
	ChangeStream struct {
		cs      *mongo.ChangeStream
		handler ChangeStreamHandler
	}

	// ChangeStreamOptions is a struct that represents options for change stream.
	ChangeStreamOptions struct {
		*options.ChangeStreamOptions
	}

	// ChangeStreamOption is a type of option function for change stream.
	ChangeStreamOption func(opts *ChangeStreamOptions)

	// StreamObject is a struct that represents a change stream object.
	SteramObject struct{}

	// ChangeStreamHandler is a type of handler function that handles ChangeStream.
	ChangeStreamHandler func(ctx context.Context) error
)

// WithResumeToken sets the resume token for ChangeStream.
func WithResumeToken(resumeToken string) ChangeStreamOption {
	return func(o *ChangeStreamOptions) {
		o.ResumeAfter = resumeToken
	}
}

// WithBatchSize sets the batch size for ChangeStream.
func WithBatchSize(batchSize int32) ChangeStreamOption {
	return func(o *ChangeStreamOptions) {
		o.BatchSize = &batchSize
	}
}

// NewChangeStream creates a new change stream instance.
func NewChangeStream(ctx context.Context, cli *Client, db, col string, opts ...ChangeStreamOption) (*ChangeStream, error) {
	chopts := &options.ChangeStreamOptions{}
	for _, opt := range opts {
		opt(&ChangeStreamOptions{chopts})
	}

	pipeline := mongo.Pipeline{}
	changeStream, err := cli.cli.Database(db).Collection(col).Watch(ctx, pipeline, chopts)
	if err != nil {
		return nil, err
	}
	cs := &ChangeStream{cs: changeStream}
	return cs, nil
}

// Run starts watching change stream.
func (c *ChangeStream) Run() {
	for c.cs.Next(context.Background()) {
		if err := c.handler(context.Background()); err != nil {
			slog.Error("failed to handle change stream", slog.String("err", err.Error()))
		}
	}
}

// ResumeToken returns the resume token of the change stream.
func (c *ChangeStream) ResumeToken() string {
	return c.cs.ResumeToken().String()
}

// Close closes the change stream cursor.
func (c *ChangeStream) Close(ctx context.Context) error {
	return c.cs.Close(ctx)
}
