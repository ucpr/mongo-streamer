package mongo

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ucpr/mongo-streamer/internal/log"
	"github.com/ucpr/mongo-streamer/internal/persistent"
)

type (
	// ChangeStream is a struct that represents a change stream.
	ChangeStream struct {
		cs           *mongo.ChangeStream
		handler      ChangeStreamHandler
		tokenManager persistent.StorageBuffer
	}

	// ChangeStreamOptions is a struct that represents options for change stream.
	ChangeStreamOptions struct {
		*options.ChangeStreamOptions
	}

	// ChangeStreamOption is a type of option function for change stream.
	ChangeStreamOption func(opts *ChangeStreamOptions)

	// ChangeStreamHandler is a type of handler function that handles ChangeStream.
	ChangeStreamHandler func(ctx context.Context, event []byte) error
)

// WithBatchSize sets the batch size for ChangeStream.
func WithBatchSize(batchSize int32) ChangeStreamOption {
	return func(o *ChangeStreamOptions) {
		o.BatchSize = &batchSize
	}
}

// NewChangeStream creates a new change stream instance.
func NewChangeStream(ctx context.Context, cli *Client, db, col string, handler ChangeStreamHandler, st persistent.StorageBuffer, opts ...ChangeStreamOption) (*ChangeStream, error) {
	chopts := &options.ChangeStreamOptions{}
	for _, opt := range opts {
		opt(&ChangeStreamOptions{chopts})
	}

	rt, err := st.Get()
	if err != nil {
		return nil, err
	}
	if rt != "" {
		chopts.SetResumeAfter(rt)
	}

	// TODO: refactor
	pipeline := mongo.Pipeline{}
	collection := cli.cli.Database(db).Collection(col)
	changeStream, err := collection.Watch(ctx, pipeline, chopts)
	if err != nil {
		// if resume token is not found, reset resume token and retry
		if errors.Is(err, mongo.ErrMissingResumeToken) {
			log.Warn("resume token is not found, reset resume token and retry", slog.String("db", db), slog.String("col", col))
			chopts.SetResumeAfter(nil)
			if err := st.Clear(); err != nil {
				return nil, err
			}

			changeStream, err = collection.Watch(ctx, pipeline, chopts)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	cs := &ChangeStream{
		cs:           changeStream,
		handler:      handler,
		tokenManager: st,
	}
	return cs, nil
}

// Run starts watching change stream.
func (c *ChangeStream) Run(ctx context.Context) {
	for c.cs.Next(ctx) {
		var streamObject bson.M
		if err := c.cs.Decode(&streamObject); err != nil {
			log.Error("failed to decode steream object", slog.String("err", err.Error()))
			continue
		}

		// marshal stream object to json
		jb, err := bson.MarshalExtJSON(streamObject, false, false)
		if err != nil {
			log.Error("failed to marshal stream object", slog.String("err", err.Error()))
			continue
		}

		if err := c.handler(context.Background(), jb); err != nil {
			log.Error("failed to handle change stream", slog.String("err", err.Error()))
			// TODO: If handle fails, the process is repeated again
			continue
		}

		// save resume token
		if err := c.tokenManager.Set(c.resumeToken()); err != nil {
			log.Error("failed to save resume token", slog.String("err", err.Error()))
			continue
		}
	}
}

// resumeToken returns the resume token of the change stream.
func (c *ChangeStream) resumeToken() string {
	return c.cs.ResumeToken().String()
}

// Close closes the change stream cursor.
func (c *ChangeStream) Close(ctx context.Context) error {
	return c.cs.Close(ctx)
}
