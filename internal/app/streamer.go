package app

import (
	"context"
	"errors"

	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/model"
	"github.com/ucpr/mongo-streamer/internal/pubsub"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

var ErrInvalidPublishFormat = errors.New("handler: invalid publish format")

type Handler struct {
	pubsub pubsub.Publisher
	pcfg   *config.PubSub
}

func NewHandler(ps pubsub.Publisher, pcfg *config.PubSub) *Handler {
	return &Handler{
		pubsub: ps,
		pcfg:   pcfg,
	}
}

func (e *Handler) EventHandler(ctx context.Context, event model.ChangeEvent) error {
	data, err := e.marshalEventData(event)
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

func (e *Handler) marshalEventData(event model.ChangeEvent) ([]byte, error) {
	switch e.pcfg.PublishFormat {
	case config.PubSubPublishFormatJSON:
		return event.JSON()
	case config.PubSubPublishFormatAvro:
		return event.Avro()
	default:
		return nil, ErrInvalidPublishFormat
	}
}
