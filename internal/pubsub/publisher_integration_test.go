//go:build integration

package pubsub

import (
	"context"
	"strconv"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ucpr/mongo-streamer/internal/config"
)

//nolint:paralleltest
func TestPublisher_AsyncPublish(t *testing.T) {
	ctx := context.Background()

	cli, err := pubsub.NewClient(ctx, testProjectID)
	require.NoError(t, err)
	publisher, err := NewPublisher(ctx, &config.PubSub{
		ProjectID: testProjectID,
		TopicID:   testTopicID,
	})
	require.NoError(t, err)

	type (
		input struct {
			data       []byte
			attributes map[string]string
		}
		want struct {
			data       []byte
			attributes map[string]string
		}
	)
	patterns := []struct {
		name string
		in   input
		want want
	}{
		{
			name: "success",
			in: input{
				data: []byte("hoge"),
				attributes: map[string]string{
					"foo": "bar",
				},
			},
			want: want{
				data: []byte("hoge"),
				attributes: map[string]string{
					"foo": "bar",
				},
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ret := publisher.AsyncPublish(ctx, Message{
				Data:       tt.in.data,
				Attributes: tt.in.attributes,
			})

			select {
			case <-ret.Ready():
				serverID, err := ret.Get(ctx)
				require.NoError(t, err)
				num, err := strconv.Atoi(serverID)
				require.NoError(t, err)
				assert.GreaterOrEqual(t, num, 1)
			case <-ctx.Done():
				t.Fatal("timeout")
			}

			err = cli.Subscription(testSubscriptionID).Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
				assert.Equal(t, tt.want.data, msg.Data)
				assert.Equal(t, tt.want.attributes, msg.Attributes)
				msg.Ack()
			})
			require.NoError(t, err)
		})
	}
}
