package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongoDB(t *testing.T) {
	ctx := context.Background()

	patterns := []struct {
		name  string
		setup func(t *testing.T)
		want  *MongoDB
	}{
		{
			name: "default",
			setup: func(t *testing.T) {
				t.Helper()
			},
			want: &MongoDB{},
		},
		{
			name: "set envs",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv("MONGO_DB_URI", "mongodb://localhost:27017")
				t.Setenv("MONGO_DB_USER", "root")
				t.Setenv("MONGO_DB_PASSWORD", "pass")
				t.Setenv("MONGO_DB_DATABASE", "database")
				t.Setenv("MONGO_DB_COLLECTION", "col")
			},
			want: &MongoDB{
				URI:        "mongodb://localhost:27017",
				Password:   "pass",
				User:       "root",
				Database:   "database",
				Collection: "col",
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			got, err := NewMongoDB(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPubSub(t *testing.T) {
	ctx := context.Background()

	patterns := []struct {
		name  string
		setup func(t *testing.T)
		want  *PubSub
	}{
		{
			name: "default",
			setup: func(t *testing.T) {
				t.Helper()
			},
			want: &PubSub{
				PublishFormat: "json",
			},
		},
		{
			name: "set envs",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv("PUBSUB_PROJECT_ID", "project")
				t.Setenv("PUBSUB_TOPIC_ID", "topic")
				t.Setenv("PUBSUB_PUBLISH_FORMAT", "avro")
			},
			want: &PubSub{
				ProjectID:     "project",
				TopicID:       "topic",
				PublishFormat: "avro",
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			got, err := NewPubSub(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMetrics(t *testing.T) {
	ctx := context.Background()

	patterns := []struct {
		name  string
		setup func(t *testing.T)
		want  *Metrics
	}{
		{
			name: "default",
			setup: func(t *testing.T) {
				t.Helper()
			},
			want: &Metrics{},
		},
		{
			name: "set envs",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv("METRICS_ADDR", "localhost:8080")
			},
			want: &Metrics{
				Addr: "localhost:8080",
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			got, err := NewMetrics(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
