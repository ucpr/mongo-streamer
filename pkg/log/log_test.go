package log

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	gcpLogger := slog.New(newHandler(newHandlerOptions{
		format: formatGCP,
	}))
	gcpLogger.Info("test")

	jsonLogger := slog.New(newHandler(newHandlerOptions{
		format: formatJSON,
	}))
	jsonLogger.Info("test")
}

func Test_toLogLevel(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		in   slog.Level
		want slog.Value
	}{
		{
			name: "default",
			in:   slog.Level(100),
			want: slog.StringValue("DEFAULT"),
		},
		{
			name: "debug",
			in:   SeverityDebug,
			want: slog.StringValue("DEBUG"),
		},
		{
			name: "info",
			in:   SeverityInfo,
			want: slog.StringValue("INFO"),
		},
		{
			name: "notice",
			in:   SeverityNotice,
			want: slog.StringValue("NOTICE"),
		},
		{
			name: "warning",
			in:   SeverityWarning,
			want: slog.StringValue("WARNING"),
		},
		{
			name: "error",
			in:   SeverityError,
			want: slog.StringValue("ERROR"),
		},
		{
			name: "critical",
			in:   SeverityCritical,
			want: slog.StringValue("CRITICAL"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toLogLevel(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
