//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/mongo"
)

func injectStreamer(ctx context.Context) (*Streamer, error) {
	wire.Build(
		config.Set,
		mongo.Set,
		NewStreamer,
	)
	return nil, nil
}
