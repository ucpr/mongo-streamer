package http

import (
	"context"
	"net/http"
	"time"

	"github.com/google/wire"

	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/metric"
)

//nolint:gochecknoglobals
var Set = wire.NewSet(
	NewServer,
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 1 * time.Second
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg *config.Metrics) *Server {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(health))
	metric.Register(mux)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &Server{
		srv: srv,
	}
}

// Serve starts the HTTP server.
func (s *Server) Serve() error {
	return s.srv.ListenAndServe()
}

// Shutdown gracefully shutdown the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
