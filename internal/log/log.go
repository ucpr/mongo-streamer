package log

import (
	"log/slog"
	"os"

	"cloud.google.com/go/logging"
)

var (
	CloudLoggingSeverityDefault  = slog.Level(logging.Default)
	CloudLoggingSeverityDebug    = slog.Level(logging.Debug)
	CloudLoggingSeverityInfo     = slog.Level(logging.Info)
	CloudLoggingSeverityNotice   = slog.Level(logging.Notice)
	CloudLoggingSeverityWarning  = slog.Level(logging.Warning)
	CloudLoggingSeverityError    = slog.Level(logging.Error)
	CloudLoggingSeverityCritical = slog.Level(logging.Critical)
)

func init() {
	logFormat := os.Getenv("LOG_FORMAT")
	handler := newHandler(logFormat)
	slog.SetDefault(slog.New(handler))
}

func newHandler(format string) slog.Handler {
	switch format {
	case "gcp":
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,
			Level:       CloudLoggingSeverityDefault,
			ReplaceAttr: attrReplacerForCloudLogging,
		})
		handler.WithAttrs([]slog.Attr{
			slog.Group("logging.googleapis.com/labels",
				slog.String("app", os.Getenv("SERVICE")),
				slog.String("env", os.Getenv("ENV")),
			),
		})
		return handler
	case "json":
		return slog.NewJSONHandler(os.Stdout, nil)
	}

	return slog.NewTextHandler(os.Stdout, nil)
}

func attrReplacerForCloudLogging(groups []string, attr slog.Attr) slog.Attr {
	switch attr.Key {
	case slog.MessageKey:
		attr.Key = "message"
	case slog.LevelKey:
		attr.Key = "severity"
		attr.Value = slog.StringValue(logging.Severity(attr.Value.Any().(slog.Level)).String())
	case slog.SourceKey:
		attr.Key = "logging.googleapis.com/sourceLocation"
	}
	return attr
}
