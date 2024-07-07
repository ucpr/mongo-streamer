package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ucpr/mongo-streamer/pkg/log"
	"github.com/ucpr/mongo-streamer/pkg/stamp"
)

const (
	gracefulShutdownTimeout = 5 * time.Second
)

func main() {
	log.Info("Initialize mongo-streamer",
		log.Fstring("version", stamp.BuildVersion),
		log.Fstring("revision", stamp.BuildRevision),
		log.Fstring("timestamp", stamp.BuildTimestamp),
	)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	streamer, err := injectStreamer(ctx)
	if err != nil {
		log.Panic("Failed to inject streamer", log.Ferror(err))
	}
	srv, err := injectServer(ctx)
	if err != nil {
		log.Panic("Failed to inject server", log.Ferror(err))
	}

	go func() {
		if err := srv.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed to start http server", log.Ferror(err))
		}
	}()

	go func() {
		streamer.Stream(ctx)
	}()

	<-ctx.Done()
	tctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(tctx); err != nil {
		log.Error("Failed to shutdown http server", log.Ferror(err))
	}
	if err := streamer.Close(tctx); err != nil {
		log.Error("Failed to close change stream", log.Ferror(err))
	}

	log.Info("Successfully graceful shutdown")
}
