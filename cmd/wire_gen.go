// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/ucpr/mongo-streamer/internal/app"
	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/http"
	"github.com/ucpr/mongo-streamer/internal/mongo"
	"github.com/ucpr/mongo-streamer/internal/pubsub"
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
	pubSub, err := config.NewPubSub(ctx)
	if err != nil {
		return nil, err
	}
	pubSubPublisher, err := pubsub.NewPublisher(ctx, pubSub)
	if err != nil {
		return nil, err
	}
	handler := app.NewHandler(pubSubPublisher, pubSub)
	streamer, err := NewStreamer(ctx, client, mongoDB, handler)
	if err != nil {
		return nil, err
	}
	return streamer, nil
}

func injectServer(ctx context.Context) (*http.Server, error) {
	metrics, err := config.NewMetrics(ctx)
	if err != nil {
		return nil, err
	}
	server := http.NewServer(metrics)
	return server, nil
}
