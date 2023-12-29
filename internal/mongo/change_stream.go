package mongo

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mmetric "github.com/ucpr/mongo-streamer/internal/metric/mongo"
	"github.com/ucpr/mongo-streamer/internal/persistent"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

type (
	// ChangeStream is a struct that represents a change stream.
	ChangeStream struct {
		cs           *mongo.ChangeStream
		handler      ChangeStreamHandler
		tokenManager persistent.StorageBuffer
		db           string
		col          string
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

// ChangeStreamParams is a struct that represents parameters for creating a ChangeStream.
type ChangeStreamParams struct {
	Client     *Client
	Handler    ChangeStreamHandler
	Storage    persistent.StorageBuffer
	Database   string
	Collection string
}

// NewChangeStream creates a new change stream instance.
func NewChangeStream(ctx context.Context, params ChangeStreamParams, opts ...ChangeStreamOption) (*ChangeStream, error) {
	chopts := &options.ChangeStreamOptions{}
	for _, opt := range opts {
		opt(&ChangeStreamOptions{chopts})
	}

	rt, err := params.Storage.Get()
	if err != nil {
		return nil, err
	}
	if rt != "" {
		chopts.SetResumeAfter(rt)
	}

	// TODO: refactor
	pipeline := mongo.Pipeline{}
	db, col := params.Database, params.Collection
	collection := params.Client.cli.Database(db).Collection(col)
	changeStream, err := collection.Watch(ctx, pipeline, chopts)
	if err != nil {
		// if resume token is not found, reset resume token and retry
		if errors.Is(err, mongo.ErrMissingResumeToken) {
			log.Warn("resume token is not found, reset resume token and retry", slog.String("db", db), slog.String("col", col))
			chopts.SetResumeAfter(nil)
			if err := params.Storage.Clear(); err != nil {
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
		handler:      params.Handler,
		tokenManager: params.Storage,
		db:           db,
		col:          col,
	}
	return cs, nil
}

// Run starts watching change stream.
func (c *ChangeStream) Run(ctx context.Context) {
	for c.cs.Next(ctx) {
		mmetric.ReceiveChangeStream(c.db, c.col)

		var streamObject bson.M
		if err := c.cs.Decode(&streamObject); err != nil {
			mmetric.HandleChangeEventFailed(c.db, c.col)
			log.Error("failed to decode steream object", slog.String("err", err.Error()))
			continue
		}

		// marshal stream object to json
		jb, err := bson.MarshalExtJSON(streamObject, false, false)
		if err != nil {
			mmetric.HandleChangeEventFailed(c.db, c.col)
			log.Error("failed to marshal stream object", slog.String("err", err.Error()))
			continue
		}
		mmetric.ReceiveBytes(c.db, c.col, len(jb))

		if err := c.handler(context.Background(), jb); err != nil {
			mmetric.HandleChangeEventFailed(c.db, c.col)
			log.Error("failed to handle change stream", slog.String("err", err.Error()))
			// TODO: If handle fails, the process is repeated again
			continue
		}
		mmetric.HandleChangeEventSuccess(c.db, c.col)

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
