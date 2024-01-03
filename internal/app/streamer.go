package app

import (
	"context"
	"encoding/json"

	"github.com/ucpr/mongo-streamer/internal/mongo"
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

func (e *Handler) EventHandler(ctx context.Context, event mongo.ChangeEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	res := e.pubsub.AsyncPublish(ctx, pubsub.Message{
		Data: data,
	})
	id, err := res.Get(ctx)
	if err != nil {
		return err
	}
	log.Info("successful publish event", log.Fstring("id", id))
	return nil
}
