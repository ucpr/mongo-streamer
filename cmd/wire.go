//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/http"
	"github.com/ucpr/mongo-streamer/internal/mongo"
	"github.com/ucpr/mongo-streamer/internal/pubsub"
)

func injectStreamer(ctx context.Context) (*Streamer, error) {
	wire.Build(
		config.Set,
		mongo.Set,
		pubsub.Set,
		NewStreamer,
		NewEventHandler,
	)
	return nil, nil
}

func injectServer(ctx context.Context) (*http.Server, error) {
	wire.Build(
		config.Set,
		http.Set,
	)
	return nil, nil
}
