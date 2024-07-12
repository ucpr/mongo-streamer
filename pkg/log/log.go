package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

//nolint:gochecknoglobals
var (
	SeverityDebug    = slog.LevelDebug // -4
	SeverityInfo     = slog.LevelInfo  // 0
	SeverityNotice   = slog.Level(2)   // 2
	SeverityWarning  = slog.LevelWarn  // 4
	SeverityError    = slog.LevelError // 8
	SeverityCritical = slog.Level(10)  // 10
)

// format is the log format type
type format string

//nolint:unused
const (
	// formatGCP is the Google Cloud Platform format.
	formatGCP format = "gcp"
	// formatJSON is the JSON format.
	formatJSON format = "json"
	// formatText is the text format.
	formatText format = "text"
)

// logger is the global logger.
// it is initialized by init() and should not be modified.
//
//nolint:gochecknoglobals
var logger *slog.Logger

// init initializes the logger.
func init() {
	logFormat := strings.ToLower(os.Getenv("LOG_FORMAT"))
	service := os.Getenv("SERVICE")
	env := os.Getenv("ENV")
	googleProjectID := os.Getenv("GOOGLE_PROJECT_ID")

	handler := newHandler(newHandlerOptions{
		format:          format(logFormat),
		service:         service,
		env:             env,
		googleProjectID: googleProjectID,
	})
	logger = slog.New(handler)
}

type newHandlerOptions struct {
	format          format
	service         string
	env             string
	googleProjectID string
}

// newHandler returns a slog.Handler based on the given format.
func newHandler(opts newHandlerOptions) slog.Handler {
	switch opts.format {
	case formatGCP:
		return NewCloudLoggingHandler(os.Stdout, &CloudLoggingHandlerOptions{
			AddSource:   true,
			Level:       SeverityInfo,
			Environment: opts.env,
			Service:     opts.service,
			ProjectID:   opts.googleProjectID,
		})
	case formatJSON:
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       SeverityInfo,
			ReplaceAttr: attrReplacerForDefault,
		})
	}

	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       SeverityInfo,
		ReplaceAttr: attrReplacerForDefault,
	})
}

// attrReplacerForDefault is default attribute replacer.
func attrReplacerForDefault(groups []string, attr slog.Attr) slog.Attr {
	switch attr.Key {
	case slog.LevelKey:
		attr.Value = toLogLevel(attr.Value.Any().(slog.Level))
	case slog.SourceKey:
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
	DebugContext(context.Background(), msg, attrs...)
}

// DebugContext logs a debug message with a context.
func DebugContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityDebug, msg, attrs...)
}

// Info logs an info message.
func Info(msg string, attrs ...any) {
	InfoContext(context.Background(), msg, attrs...)
}

// InfoContext logs an info message with a context.
func InfoContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityInfo, msg, attrs...)
}

// Notice logs a notice message.
func Notice(msg string, attrs ...any) {
	NoticeContext(context.Background(), msg, attrs...)
}

// NoticeContext logs a notice message with a context.
func NoticeContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityNotice, msg, attrs...)
}

// Warn logs a warning message.
func Warn(msg string, attrs ...any) {
	WarnContext(context.Background(), msg, attrs...)
}

// WarnContext logs a warning message with a context.
func WarnContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityWarning, msg, attrs...)
}

// Error logs an error message.
func Error(msg string, attrs ...any) {
	ErrorContext(context.Background(), msg, attrs...)
}

// ErrorContext logs an error message with a context.
func ErrorContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityError, msg, attrs...)
}

// Critical logs a critical message.
func Critical(msg string, attrs ...any) {
	CriticalContext(context.Background(), msg, attrs...)
}

// CriticalContext logs a critical message with a context.
func CriticalContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityCritical, msg, attrs...)
}

// Panic logs a critical message and panics.
func Panic(msg string, attrs ...any) {
	PanicContext(context.Background(), msg, attrs...)
	panic(msg)
}

// PanicContext logs a critical message with a context and panics.
func PanicContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityCritical, msg, attrs...)
	panic(msg)
}

// Fatal logs a critical message and exits.
func Fatal(msg string, attrs ...any) {
	FatalContext(context.Background(), msg, attrs...)
	os.Exit(1)
}

// FatalContext logs a critical message with a context and exits.
func FatalContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityCritical, msg, attrs...)
	os.Exit(1)
}
