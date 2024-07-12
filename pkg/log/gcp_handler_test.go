package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"testing"

	"cloud.google.com/go/logging"
	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestCloudLoggingHandler_Enabled(t *testing.T) {
	t.Parallel()

	handler := NewCloudLoggingHandler(io.Discard, &CloudLoggingHandlerOptions{
		Level: SeverityInfo,
	})
	assert.Equal(t, false, handler.Enabled(context.Background(), SeverityDebug))
	assert.Equal(t, true, handler.Enabled(context.Background(), SeverityInfo))
}

func TestCloudLoggingHandler_WithAttr(t *testing.T) {
	t.Parallel()

	handler := NewCloudLoggingHandler(io.Discard, &CloudLoggingHandlerOptions{
		Level: SeverityInfo,
	})
	assert.NotNil(t, handler.WithAttrs([]slog.Attr{}))
}

func TestCloudLoggingHandler_WithGroup(t *testing.T) {
	t.Parallel()

	handler := NewCloudLoggingHandler(io.Discard, &CloudLoggingHandlerOptions{
		Level: SeverityInfo,
	})
	assert.NotNil(t, handler.WithGroup("group"))
}

func TestCloudLoggingHandler_toCloudLoggingLevel(t *testing.T) {
	t.Parallel()
	handler := NewCloudLoggingHandler(io.Discard, &CloudLoggingHandlerOptions{
		AddSource:   true,
		Level:       SeverityInfo,
		Environment: "test",
		Service:     "test",
		ProjectID:   "test",
	})

	patterns := []struct {
		name string
		in   slog.Level
		want slog.Level
	}{
		{
			name: "undefined log level",
			in:   slog.Level(100),
			want: slog.Level(logging.Info),
		},
		{
			name: "info",
			in:   SeverityInfo,
			want: slog.Level(logging.Info),
		},
		{
			name: "debug",
			in:   SeverityDebug,
			want: slog.Level(logging.Debug),
		},
		{
			name: "notice",
			in:   SeverityNotice,
			want: slog.Level(logging.Notice),
		},
		{
			name: "warning",
			in:   SeverityWarning,
			want: slog.Level(logging.Warning),
		},
		{
			name: "error",
			in:   SeverityError,
			want: slog.Level(logging.Error),
		},
		{
			name: "critical",
			in:   SeverityCritical,
			want: slog.Level(logging.Critical),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := handler.toCloudLoggingLevel(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCloudLoggingHandler(t *testing.T) {
	t.Parallel()

	const env, service, projectID = "env", "service", "projectID"
	patterns := []struct {
		name        string
		setupLogger func(buf *bytes.Buffer) *slog.Logger
		ctx         func() context.Context
		wantValues  map[string]any
		wantHasKeys []string
	}{
		{
			name: "success without trace",
			setupLogger: func(buf *bytes.Buffer) *slog.Logger {
				handler := NewCloudLoggingHandler(buf, &CloudLoggingHandlerOptions{
					AddSource:   true,
					Level:       SeverityInfo,
					Environment: env,
					Service:     service,
					ProjectID:   projectID,
				})
				return slog.New(handler)
			},
			ctx: func() context.Context {
				return context.Background()
			},
			wantValues: map[string]any{
				cloudLoggingMessageKey: "message",
				cloudLoggingLabelsKey: map[string]any{
					cloudLoggingServiceLabel:     service,
					cloudLoggingEnvironmentLabel: env,
				},
				cloudLoggingSeverityKey: logging.Info.String(),
			},
			wantHasKeys: []string{
				cloudLoggingSourceKey,
			},
		},
		{
			name: "success with replaceAttr",
			setupLogger: func(buf *bytes.Buffer) *slog.Logger {
				handler := NewCloudLoggingHandler(buf, &CloudLoggingHandlerOptions{
					AddSource:   true,
					Level:       SeverityInfo,
					Environment: env,
					Service:     service,
					ProjectID:   projectID,
					ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
						if cloudLoggingMessageKey == attr.Key {
							return slog.String("custom_attr", "message")
						}
						return attr
					},
				})
				return slog.New(handler)
			},
			ctx: func() context.Context {
				return context.Background()
			},
			wantValues: map[string]any{
				"custom_attr": "message",
				cloudLoggingLabelsKey: map[string]any{
					cloudLoggingServiceLabel:     service,
					cloudLoggingEnvironmentLabel: env,
				},
				cloudLoggingSeverityKey: logging.Info.String(),
			},
			wantHasKeys: []string{},
		},
		{
			name: "success with tracing",
			setupLogger: func(buf *bytes.Buffer) *slog.Logger {
				handler := NewCloudLoggingHandler(buf, &CloudLoggingHandlerOptions{
					AddSource:   true,
					Level:       SeverityInfo,
					Environment: env,
					Service:     service,
					ProjectID:   projectID,
				})
				return slog.New(handler)
			},
			ctx: func() context.Context {
				spanRecorder := tracetest.NewSpanRecorder()
				traceProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))

				tracer := traceProvider.Tracer("test-tracer")
				ctx, span := tracer.Start(context.Background(), "test-span")
				span.End()

				return ctx
			},
			wantValues: map[string]any{
				cloudLoggingMessageKey: "message",
				cloudLoggingLabelsKey: map[string]any{
					cloudLoggingServiceLabel:     service,
					cloudLoggingEnvironmentLabel: env,
				},
				cloudLoggingSeverityKey: logging.Info.String(),
			},
			wantHasKeys: []string{
				cloudLoggingSourceKey,
				cloudLoggingTraceKey,
				cloudLoggingTraceSpanKey,
				cloudLoggingTraceSampledKey,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			logger := tt.setupLogger(buf)
			logger.InfoContext(tt.ctx(), "message")

			got := map[string]any{}
			err := json.Unmarshal(buf.Bytes(), &got)
			assert.NoError(t, err)

			for key, want := range tt.wantValues {
				assert.Equal(t, want, got[key])
			}
			for _, key := range tt.wantHasKeys {
				assert.Contains(t, got, key)
			}
		})
	}
}
