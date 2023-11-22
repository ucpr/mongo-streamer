package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/log"
	"github.com/ucpr/mongo-streamer/internal/mongo"
)

const (
	gracefulShutdownTimeout = 5 * time.Second
)

func eventHandler(ctx context.Context, event []byte) error {
	log.Info("event", slog.String("event", string(event)))
	return nil
}

func main() {
	log.Info("initializing mongo-streamer")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	cfg, err := config.NewMongoDB(ctx)
	if err != nil {
		log.Panic("failed to get environment variables", err)
	}

	mcli, err := mongo.NewClient(ctx, cfg)
	if err != nil {
		log.Panic("failed to create mongo client", err)
	}

	cs, err := mongo.NewChangeStream(ctx, mcli, "test", "test", eventHandler)
	if err != nil {
		log.Panic("failed to create change stream", err)
	}
	go func() {
		log.Info("starting change stream watcher")
		cs.Run()
	}()
	defer func() {
		if err := cs.Close(ctx); err != nil {
			log.Error("failed to close change stream", err)
		}
	}()

	<-ctx.Done()
	tctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	_ = tctx

	log.Info("successfully graceful shutdown")
}
