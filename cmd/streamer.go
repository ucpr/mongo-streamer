package main

import (
	"context"
	"time"

	"github.com/ucpr/mongo-streamer/internal/app"
	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/mongo"
	"github.com/ucpr/mongo-streamer/internal/persistent"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

type Streamer struct {
	cli *mongo.Client
	cs  *mongo.ChangeStream
	st  persistent.StorageBuffer
}

func NewStreamer(ctx context.Context, cli *mongo.Client, mcfg *config.MongoDB, h *app.Handler) (*Streamer, error) {
	stLog := persistent.NewLogWriter()
	st, err := persistent.NewBuffer(10, 5*time.Second, stLog)
	if err != nil {
		return nil, err
	}
	cs, err := mongo.NewChangeStream(ctx, mongo.ChangeStreamParams{
		Client:     cli,
		Handler:    h.EventHandler,
		Storage:    st,
		Database:   mcfg.Database,
		Collection: mcfg.Collection,
	})
	if err != nil {
		return nil, err
	}

	return &Streamer{
		cli: cli,
		cs:  cs,
		st:  st,
	}, nil
}

func (s *Streamer) Stream(ctx context.Context) {
	go func() {
		s.st.Watch(ctx)
	}()

	log.Info("Start change stream watcher")
	s.cs.Run(ctx)
}

func (s *Streamer) Close(ctx context.Context) error {
	if err := s.cs.Close(ctx); err != nil {
		return err
	}
	if err := s.st.Close(ctx); err != nil {
		return err
	}
	if err := s.cli.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
