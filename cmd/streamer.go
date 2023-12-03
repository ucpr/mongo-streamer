package main

import (
	"context"
	"log/slog"

	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/log"
	"github.com/ucpr/mongo-streamer/internal/mongo"
)

type Streamer struct {
	cli *mongo.Client
	cs  *mongo.ChangeStream
}

func NewStreamer(ctx context.Context, cli *mongo.Client, mcfg *config.MongoDB) (*Streamer, error) {
	cs, err := mongo.NewChangeStream(ctx, cli, mcfg.Database, mcfg.Collection, eventHandler)
	if err != nil {
		return nil, err
	}

	return &Streamer{cli: cli, cs: cs}, nil
}

func eventHandler(ctx context.Context, event []byte) error {
	log.Info("event", slog.String("event", string(event)))
	return nil
}

func (s *Streamer) Stream(ctx context.Context) {
	log.Info("starting change stream watcher")
	s.cs.Run(ctx)
}

func (s *Streamer) Close(ctx context.Context) error {
	return s.cs.Close(ctx)
}
