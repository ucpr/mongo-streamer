package log

import (
	"context"
	"log/slog"
	"os"

	"cloud.google.com/go/logging"
)

//nolint:gochecknoglobals
var (
	SeverityDefault  = slog.Level(logging.Default)
	SeverityDebug    = slog.Level(logging.Debug)
	SeverityInfo     = slog.Level(logging.Info)
	SeverityNotice   = slog.Level(logging.Notice)
	SeverityWarning  = slog.Level(logging.Warning)
	SeverityError    = slog.Level(logging.Error)
	SeverityCritical = slog.Level(logging.Critical)
)

// logger is the global logger.
// it is initialized by init() and should not be modified.
var logger *slog.Logger

// init initializes the logger.
func init() {
	logFormat := os.Getenv("LOG_FORMAT")
	service := os.Getenv("SERVICE")
	env := os.Getenv("ENV")

	handler := newHandler(logFormat, service, env)
	logger = slog.New(handler)
}

// newHandler returns a slog.Handler based on the given format.
func newHandler(format, service, env string) slog.Handler {
	switch format {
	case "gcp":
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,
			Level:       SeverityDefault,
			ReplaceAttr: attrReplacerForCloudLogging,
		})
		handler.WithAttrs([]slog.Attr{
			slog.Group("logging.googleapis.com/labels",
				slog.String("app", service),
				slog.String("env", env),
			),
		})
		return handler
	case "json":
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       SeverityDefault,
			ReplaceAttr: attrReplacerForDefault,
		})
	}

	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       SeverityDefault,
		ReplaceAttr: attrReplacerForDefault,
	})
}

// attrReplacerForDefault is default attribute replacer.
func attrReplacerForDefault(groups []string, attr slog.Attr) slog.Attr {
	// Replace the value of the "severity" attribute with the value of the "level" attribute.
	level, ok := attr.Value.Any().(slog.Level)
	if ok {
		attr.Value = toLogLevel(level)
	}
	return attr
}

// attrReplacerForCloudLogging replaces slog's default attributes for Google Cloud Logging.
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
	// Replace the value of the "severity" attribute with the value of the "level" attribute.
	level, ok := attr.Value.Any().(slog.Level)
	if ok {
		attr.Value = toLogLevel(level)
	}

	return attr
}

// toLogLevel converts a slog.Level to a slog.Value.
func toLogLevel(level slog.Level) slog.Value {
	ls := "DEFAULT"

	switch level {
	case SeverityDebug:
		ls = "DEBUG"
	case SeverityInfo:
		ls = "INFO"
	case SeverityNotice:
		ls = "NOTICE"
	case SeverityWarning:
		ls = "WARNING"
	case SeverityError:
		ls = "ERROR"
	case SeverityCritical:
		ls = "CRITICAL"
	}

	return slog.StringValue(ls)
}

// Debug logs a debug message.
func Debug(msg string, attrs ...any) {
	logger.Log(context.Background(), SeverityDebug, msg, attrs...)
}

// Info logs an info message.
func Info(msg string, attrs ...any) {
	logger.Log(context.Background(), SeverityInfo, msg, attrs...)
}

// Notice logs a notice message.
func Notice(msg string, attrs ...any) {
	logger.Log(context.Background(), SeverityNotice, msg, attrs...)
}

// Warn logs a warning message.
func Warn(msg string, attrs ...any) {
	logger.Log(context.Background(), SeverityWarning, msg, attrs...)
}

// Error logs an error message.
func Error(msg string, attrs ...any) {
	logger.Log(context.Background(), SeverityError, msg, attrs...)
}

// Critical logs a critical message.
func Critical(msg string, attrs ...any) {
	logger.Log(context.Background(), SeverityCritical, msg, attrs...)
}

// Panic logs a critical message and panics.
func Panic(msg string, attrs ...any) {
	logger.Log(context.Background(), SeverityCritical, msg, attrs...)
	panic(msg)
}
