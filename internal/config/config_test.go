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
			},
			want: &MongoDB{
				URI:      "mongodb://localhost:27017",
				Password: "pass",
				User:     "root",
				Database: "database",
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
