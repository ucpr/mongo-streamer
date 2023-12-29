package app

import (
	"context"
	"log/slog"

	"github.com/ucpr/mongo-streamer/internal/pubsub"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

type Handler struct {
	pubsub pubsub.Publisher
}

func NewHandler(ps pubsub.Publisher) *Handler {
	return &Handler{
		pubsub: ps,
	}
}

func (e *Handler) EventHandler(ctx context.Context, event []byte) error {
	res := e.pubsub.AsyncPublish(ctx, pubsub.Message{
		Data: event,
	})
	id, err := res.Get(ctx)
	if err != nil {
		return err
	}
	log.Info("successful publish event", slog.String("id", id))
	return nil
}
