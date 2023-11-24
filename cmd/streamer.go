package main

import (
	"context"

	"github.com/ucpr/mongo-streamer/internal/log"
	"github.com/ucpr/mongo-streamer/internal/mongo"
)

type Streamer struct {
	cli *mongo.Client
	cs  *mongo.ChangeStream
}

func NewStreamer(ctx context.Context, cli *mongo.Client) (*Streamer, error) {
	cs, err := mongo.NewChangeStream(ctx, cli, "test", "test", eventHandler)
	if err != nil {
		return nil, err
	}

	return &Streamer{cli: cli, cs: cs}, nil
}

func (s *Streamer) Stream(ctx context.Context) {
	log.Info("starting change stream watcher")
	s.cs.Run()
}

func (s *Streamer) Close(ctx context.Context) error {
	return s.cs.Close(ctx)
}
