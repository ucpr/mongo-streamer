package mongo

import (
	"context"
	"log/slog"
	"time"

	"github.com/ucpr/mongo-streamer/internal/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	StreamObject struct {
		ID struct {
			Data string `bson:"_data"`
		} `bson:"_id"`
		OperationType  string           `bson:"operationType"`
		ClusterTime    time.Time        `bson:"clusterTime"`
		CollectionUUID primitive.Binary `bson:"collectionUUID"`
		WallTime       time.Time        `bson:"wallTime"`
		FullDocument   []byte           `bson:"fullDocument"`
		Namespace      struct {
			DB   string `bson:"db"`
			Coll string `bson:"coll"`
		} `bson:"ns"`
		DocumentKey struct {
			ID primitive.ObjectID `bson:"_id"`
		} `bson:"documentKey"`
	}

	// ChangeStreamHandler is a type of handler function that handles ChangeStream.
	ChangeStreamHandler func(ctx context.Context, event *StreamObject) error
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
func NewChangeStream(ctx context.Context, cli *Client, db, col string, handler ChangeStreamHandler, opts ...ChangeStreamOption) (*ChangeStream, error) {
	chopts := &options.ChangeStreamOptions{}
	for _, opt := range opts {
		opt(&ChangeStreamOptions{chopts})
	}

	pipeline := mongo.Pipeline{}
	changeStream, err := cli.cli.Database(db).Collection(col).Watch(ctx, pipeline, chopts)
	if err != nil {
		return nil, err
	}
	cs := &ChangeStream{
		cs:      changeStream,
		handler: handler,
	}
	return cs, nil
}

// Run starts watching change stream.
func (c *ChangeStream) Run() {
	for c.cs.Next(context.Background()) {
		var streamObject StreamObject
		if err := c.cs.Decode(&streamObject); err != nil {
			log.Error("failed to decode steream object", slog.String("err", err.Error()))
			continue
		}

		if err := c.handler(context.Background(), &streamObject); err != nil {
			log.Error("failed to handle change stream", slog.String("err", err.Error()))
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
