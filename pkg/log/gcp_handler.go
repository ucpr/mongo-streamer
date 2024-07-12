package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"

	"cloud.google.com/go/logging"
	"go.opentelemetry.io/otel/trace"
)

const (
	cloudLoggingSourceKey       = "logging.googleapis.com/sourceLocation"
	cloudLoggingLabelsKey       = "logging.googleapis.com/labels"
	cloudLoggingTraceKey        = "logging.googleapis.com/trace"
	cloudLoggingTraceSpanKey    = "logging.googleapis.com/spanId"
	cloudLoggingTraceSampledKey = "logging.googleapis.com/trace_sampled"

	cloudLoggingMessageKey       = "message"
	cloudLoggingSeverityKey      = "severity"
	cloudLoggingEnvironmentLabel = "env"
	cloudLoggingServiceLabel     = "app"
	cloudLoggingTraceFormat      = "projects/%s/traces/%s"
)

type CloudLoggingHandler struct {
	handler     slog.Handler
	environment string
	service     string
	projectID   string
}

type CloudLoggingHandlerOptions struct {
	AddSource   bool
	Level       slog.Level
	Environment string
	Service     string
	ProjectID   string
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

// Ensure CloudLoggingHandler implements slog.Handler.
var _ slog.Handler = (*CloudLoggingHandler)(nil)

// NewCloudLoggingHandler creates a new CloudLoggingHandler.
func NewCloudLoggingHandler(out io.Writer, options *CloudLoggingHandlerOptions) *CloudLoggingHandler {
	replaceAttr := func(groups []string, attr slog.Attr) slog.Attr {
		cattr := attrReplacerForCloudLogging(groups, attr)
		if options.ReplaceAttr != nil {
			cattr = options.ReplaceAttr(groups, cattr)
		}
		return cattr
	}

	handler := slog.NewJSONHandler(out, &slog.HandlerOptions{
		AddSource:   options.AddSource,
		Level:       options.Level,
		ReplaceAttr: replaceAttr,
	})

	return &CloudLoggingHandler{
		handler:     handler,
		environment: options.Environment,
		service:     options.Service,
		projectID:   options.ProjectID,
	}
}

// Enabled implements the slog.Handler interface.
func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements the slog.Handler interface.
func (h *CloudLoggingHandler) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(
		slog.Group(cloudLoggingLabelsKey,
			slog.String(cloudLoggingServiceLabel, h.service),
			slog.String(cloudLoggingEnvironmentLabel, h.environment),
		),
	)
	record.Level = h.toCloudLoggingLevel(record.Level)

	// set trace
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		trace := fmt.Sprintf(cloudLoggingTraceFormat, h.projectID, sc.TraceID().String())
		record.AddAttrs(
			slog.String(cloudLoggingTraceKey, trace),
			slog.String(cloudLoggingTraceSpanKey, sc.SpanID().String()),
			slog.Bool(cloudLoggingTraceSampledKey, sc.TraceFlags().IsSampled()),
		)
	}

	return h.handler.Handle(ctx, record)
}

// WithAttrs implements the slog.Handler interface.
func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

// WithGroup implements the slog.Handler interface.
func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

func (h *CloudLoggingHandler) toCloudLoggingLevel(level slog.Level) slog.Level {
	switch level {
	case SeverityDebug:
		return slog.Level(logging.Debug)
	case SeverityInfo:
		return slog.Level(logging.Info)
	case SeverityNotice:
		return slog.Level(logging.Notice)
	case SeverityWarning:
		return slog.Level(logging.Warning)
	case SeverityError:
		return slog.Level(logging.Error)
	case SeverityCritical:
		return slog.Level(logging.Critical)
	}
	return slog.Level(logging.Info)
}

// attrReplacerForCloudLogging is default attribute replacer.
func attrReplacerForCloudLogging(_ []string, attr slog.Attr) slog.Attr {
	switch attr.Key {
	case slog.MessageKey:
		attr.Key = cloudLoggingMessageKey
	case slog.LevelKey:
		attr.Key = cloudLoggingSeverityKey
		attr.Value = slog.StringValue(logging.Severity(attr.Value.Any().(slog.Level)).String())
	case slog.SourceKey:
		attr.Key = cloudLoggingSourceKey
		// Replace the value of the "source" attribute with the value of the "sourceLocation" attribute.
		const skip = 9
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			return attr
		}
		v := fmt.Sprintf("%s:%d", file, line)
		attr.Value = slog.StringValue(v)
	}

	return attr
}
