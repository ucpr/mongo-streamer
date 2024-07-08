//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/mongo"
)

func inject(ctx context.Context) (*App, error) {
	wire.Build(
		config.Set,
		mongo.Set,
		NewApp,
	)
	return nil, nil
}
