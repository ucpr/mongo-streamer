package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ucpr/mongo-streamer/internal/pubsub"
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
		name     string
		injector func(t *testing.T, mc *mc)
		event    []byte
		err      error
	}{
		{
			name: "success handle event",
			injector: func(t *testing.T, mc *mc) {
				t.Helper()

				mc.publishResult.EXPECT().Get(ctx).Return("id", nil).Times(1)
				mc.publisher.EXPECT().AsyncPublish(ctx, pubsub.Message{
					Data: []byte("data"),
				}).Return(mc.publishResult).Times(1)
			},
			event: []byte("data"),
			err:   nil,
		},
		{
			name: "failed to publish event",
			injector: func(t *testing.T, mc *mc) {
				t.Helper()

				mc.publishResult.EXPECT().Get(ctx).Return("", assert.AnError).Times(1)
				mc.publisher.EXPECT().AsyncPublish(ctx, pubsub.Message{
					Data: []byte("data"),
				}).Return(mc.publishResult).Times(1)
			},
			event: []byte("data"),
			err:   assert.AnError,
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

			h := NewHandler(mp)
			got := h.EventHandler(ctx, tt.event)
			assert.Equal(t, tt.err, got)
		})
	}
}
