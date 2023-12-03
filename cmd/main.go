package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/ucpr/mongo-streamer/internal/log"
)

const (
	gracefulShutdownTimeout = 5 * time.Second
)

func main() {
	log.Info("initializing mongo-streamer")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	streamer, err := injectStreamer(ctx)
	if err != nil {
		log.Panic("failed to inject streamer", err)
	}

	go func() {
		streamer.Stream(ctx)
	}()

	<-ctx.Done()
	tctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := streamer.Close(tctx); err != nil {
		log.Error("failed to close change stream", err)
	}

	log.Info("successfully graceful shutdown")
}
