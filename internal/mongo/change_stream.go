package mongo

import (
	"context"
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mmetric "github.com/ucpr/mongo-streamer/internal/metric/mongo"
	"github.com/ucpr/mongo-streamer/internal/model"
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
	ChangeStreamHandler func(ctx context.Context, event model.ChangeEvent) error
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
			log.Warn("Resume token is not found, reset resume token and retry", log.Fstring("db", db), log.Fstring("col", col))
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

		var streamObject model.ChangeEvent
		if err := c.cs.Decode(&streamObject); err != nil {
			mmetric.HandleChangeEventFailed(c.db, c.col)
			log.Error("Failed to decode steream object", log.Ferror(err))
			continue
		}

		// marshal stream object to json
		jb, err := json.Marshal(streamObject)
		if err != nil {
			// skip for metrics retention use
			log.Error("Failed to marshal stream object to json", log.Ferror(err))
		}
		mmetric.ReceiveBytes(c.db, c.col, len(jb))

		if err := c.handler(context.Background(), streamObject); err != nil {
			mmetric.HandleChangeEventFailed(c.db, c.col)
			log.Error("Failed to handle change stream", log.Ferror(err))
			// TODO: If handle fails, the process is repeated again
			continue
		}
		mmetric.HandleChangeEventSuccess(c.db, c.col)

		// save resume token
		if err := c.tokenManager.Set(c.resumeToken()); err != nil {
			log.Error("Failed to save resume token", log.Ferror(err))
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
