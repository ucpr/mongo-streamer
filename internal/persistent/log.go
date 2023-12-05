package persistent

import (
	"context"
	"log/slog"

	"github.com/ucpr/mongo-streamer/internal/log"
)

// Log is a writer that writes to the log.
//
// This is for testing purposes and may be removed and should
// not be used in a production environment.
type Log struct{}

var _ Storage = (*Log)(nil)

func NewLogWriter() *Log {
	return &Log{}
}

func (l *Log) Write(s string) error {
	log.Info("persistent: write data", slog.String("data", s))
	return nil
}

func (l *Log) Clear() error {
	log.Info("persistent: clear data")
	return nil
}

func (l *Log) Read() (string, error) {
	log.Info("persistent: read data")
	return "", nil
}

func (l *Log) Close(ctx context.Context) error {
	return nil
}
