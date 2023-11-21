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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	<-ctx.Done()
	tctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	_ = tctx

	log.Info("successfully graceful shutdown")
}
