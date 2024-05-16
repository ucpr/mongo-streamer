package log

import (
	"fmt"
	"log/slog"
	"runtime"
)

func labelsAttr(service, env string) slog.Attr {
	return slog.Group("logging.googleapis.com/labels",
		slog.String("app", service),
		slog.String("env", env),
	)
}

// attrReplacerForCloudLogging replaces slog's default attributes for Google Cloud Logging.
func attrReplacerForCloudLogging(groups []string, attr slog.Attr) slog.Attr {
	switch attr.Key {
	case slog.MessageKey:
		attr.Key = "message"
	case slog.LevelKey:
		attr.Key = "severity"
		// Replace the value of the "severity" attribute with the value of the "level" attribute.
		level, ok := attr.Value.Any().(slog.Level)
		if ok {
			attr.Value = toLogLevel(level)
		}
	case slog.SourceKey:
		attr.Key = "logging.googleapis.com/sourceLocation"
		// Replace the value of the "source" attribute with the value of the "sourceLocation" attribute.
		const skip = 7
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			return attr
		}
		v := fmt.Sprintf("%s:%d", file, line)
		attr.Value = slog.StringValue(v)
	}

	return attr
}
