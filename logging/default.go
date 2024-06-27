package logging

import (
	"log/slog"
	"os"
)

var DefaultLogger = NewDefaultLogger()

func NewDefaultLogger() *slog.Logger {
	return NewLogger("json", "info")
}

func NewLogger(format, level string) *slog.Logger {
	o := &slog.HandlerOptions{
		Level: getLogLevel(level),
	}

	var h slog.Handler

	if format == "text" {
		h = slog.NewTextHandler(os.Stdout, o)
	} else {
		h = slog.NewJSONHandler(os.Stdout, o)
	}

	return slog.New(h)
}

func getLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
