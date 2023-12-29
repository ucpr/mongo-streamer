package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ucpr/mongo-streamer/pkg/log"
)

var (
	// inject by ldflags
	BuildVersion   = ""
	BuildRevision  = ""
	BuildTimestamp = ""
)

const (
	gracefulShutdownTimeout = 5 * time.Second
)

func main() {
	log.Info("initializing mongo-streamer",
		slog.String("version", BuildVersion),
		slog.String("revision", BuildRevision),
		slog.String("timestamp", BuildTimestamp),
	)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	streamer, err := injectStreamer(ctx)
	if err != nil {
		log.Panic("failed to inject streamer", err)
	}
	srv, err := injectServer(ctx)
	if err != nil {
		log.Panic("failed to inject server", err)
	}

	go func() {
		if err := srv.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start http server", err)
		}
	}()

	go func() {
		streamer.Stream(ctx)
	}()

	<-ctx.Done()
	tctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(tctx); err != nil {
		log.Error("failed to shutdown http server", err)
	}
	if err := streamer.Close(tctx); err != nil {
		log.Error("failed to close change stream", err)
	}

	log.Info("successfully graceful shutdown")
}
