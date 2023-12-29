package main

import (
	"context"
	"log/slog"

	"github.com/ucpr/mongo-streamer/internal/pubsub"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

type EventHandler struct {
	pubsub pubsub.Publisher
}

func NewEventHandler(ps pubsub.Publisher) *EventHandler {
	return &EventHandler{
		pubsub: ps,
	}
}

func (e *EventHandler) EventHandler(ctx context.Context, event []byte) error {
	log.Info("event", slog.String("event", string(event)))
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
