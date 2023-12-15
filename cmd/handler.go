package main

import (
	"context"
	"log/slog"

	"github.com/ucpr/mongo-streamer/internal/log"
	"github.com/ucpr/mongo-streamer/internal/pubsub"
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
	return nil
}
