package http

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ucpr/mongo-streamer/internal/config"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	cfg := &config.Metrics{}

	got := NewServer(cfg)
	assert.NotNil(t, got)
}
