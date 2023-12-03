// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/mongo"
)

// Injectors from wire.go:

func injectStreamer(ctx context.Context) (*Streamer, error) {
	mongoDB, err := config.NewMongoDB(ctx)
	if err != nil {
		return nil, err
	}
	client, err := mongo.NewClient(ctx, mongoDB)
	if err != nil {
		return nil, err
	}
	streamer, err := NewStreamer(ctx, client, mongoDB)
	if err != nil {
		return nil, err
	}
	return streamer, nil
}
