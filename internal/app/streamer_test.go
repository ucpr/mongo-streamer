package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ucpr/mongo-streamer/internal/config"
	"github.com/ucpr/mongo-streamer/internal/model"
	"github.com/ucpr/mongo-streamer/internal/pubsub/mock"
)

func TestHandler_EventHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	type mc struct {
		publisher     *mock.MockPublisher
		publishResult *mock.MockPublishResult
	}
	patterns := []struct {
		name          string
		injector      func(t *testing.T, mc *mc)
		event         model.ChangeEvent
		publishFormat string
		err           error
	}{
		{
			name: "success handle event with json format",
			injector: func(t *testing.T, mc *mc) {
				t.Helper()

				mc.publishResult.EXPECT().Get(ctx).Return("id", nil).Times(1)
				mc.publisher.EXPECT().AsyncPublish(ctx, gomock.Any()).Return(mc.publishResult).Times(1)
			},
			event: model.ChangeEvent{
				ID: "id",
			},
			publishFormat: config.PubSubPublishFormatJSON,
			err:           nil,
		},
		{
			name: "success handle event with avro format",
			injector: func(t *testing.T, mc *mc) {
				t.Helper()

				mc.publishResult.EXPECT().Get(ctx).Return("id", nil).Times(1)
				mc.publisher.EXPECT().AsyncPublish(ctx, gomock.Any()).Return(mc.publishResult).Times(1)
			},
			event: model.ChangeEvent{
				ID: "id",
			},
			publishFormat: config.PubSubPublishFormatAvro,
			err:           nil,
		},
		{
			name: "invalid publish format",
			injector: func(t *testing.T, mc *mc) {
				t.Helper()
			},
			event: model.ChangeEvent{
				ID: "id",
			},
			publishFormat: "invalid_format",
			err:           ErrInvalidPublishFormat,
		},
		{
			name: "failed to publish event",
			injector: func(t *testing.T, mc *mc) {
				t.Helper()

				mc.publishResult.EXPECT().Get(ctx).Return("", assert.AnError).Times(1)
				mc.publisher.EXPECT().AsyncPublish(ctx, gomock.Any()).Return(mc.publishResult).Times(1)
			},
			event: model.ChangeEvent{
				ID: "id",
			},
			publishFormat: config.PubSubPublishFormatJSON,
			err:           assert.AnError,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mp := mock.NewMockPublisher(ctrl)
			mpr := mock.NewMockPublishResult(ctrl)
			tt.injector(t, &mc{
				publisher:     mp,
				publishResult: mpr,
			})

			h := NewHandler(mp, &config.PubSub{
				PublishFormat: tt.publishFormat,
			})
			err := h.EventHandler(ctx, tt.event)
			assert.Equal(t, tt.err, err)
		})
	}
}
