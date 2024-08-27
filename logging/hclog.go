package logging

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/hashicorp/go-hclog"
)

// check that HCLogger implements hclog.Logger interface
var _ hclog.Logger = &HCLogger{}

type HCLogger struct {
	name         string
	traceEnabled bool
	logger       *slog.Logger
}

func NewHCLogger(logger *slog.Logger, traceEnabled bool) *HCLogger {
	return &HCLogger{
		logger:       logger,
		traceEnabled: traceEnabled,
	}
}

func (l *HCLogger) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.Trace:
		l.Trace(msg, args...)
	case hclog.Debug:
		l.Debug(msg, args...)
	case hclog.Info:
		l.Info(msg, args...)
	case hclog.Warn:
		l.Warn(msg, args...)
	case hclog.Error:
		l.Error(msg, args...)
	}
}

func (l *HCLogger) Trace(msg string, args ...interface{}) {
	if !l.traceEnabled {
		return
	}

	l.logger.Debug(msg, args...)
}

func (l *HCLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *HCLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *HCLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *HCLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

func (l *HCLogger) GetLevel() hclog.Level {
	return hclog.Trace
}

func (l *HCLogger) IsTrace() bool {
	return l.traceEnabled && l.logger.Enabled(context.Background(), slog.LevelDebug)
}

func (l *HCLogger) IsDebug() bool {
	return l.logger.Enabled(context.Background(), slog.LevelDebug)
}

func (l *HCLogger) IsInfo() bool {
	return l.logger.Enabled(context.Background(), slog.LevelInfo)
}

func (l *HCLogger) IsWarn() bool {
	return l.logger.Enabled(context.Background(), slog.LevelWarn)
}

func (l *HCLogger) IsError() bool {
	return l.logger.Enabled(context.Background(), slog.LevelError)
}

func (l *HCLogger) With(args ...interface{}) hclog.Logger {
	return &HCLogger{
		name:         l.name,
		traceEnabled: l.traceEnabled,
		logger:       l.logger.With(args...),
	}
}

func (l *HCLogger) ImpliedArgs() []interface{} {
	return nil
}
func (l *HCLogger) ResetNamed(name string) hclog.Logger {
	return &HCLogger{
		name:         name,
		traceEnabled: l.traceEnabled,
		logger:       l.logger.With("name", name),
	}
}

func (l *HCLogger) Named(name string) hclog.Logger {
	return &HCLogger{
		name:         name,
		traceEnabled: l.traceEnabled,
		logger:       l.logger.With("name", name),
	}
}

func (l *HCLogger) Name() string {
	return l.name
}

func (l *HCLogger) SetLevel(level hclog.Level) {
}

func (l *HCLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return slog.NewLogLogger(l.logger.Handler(), slog.LevelDebug)
}

func (l *HCLogger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return os.Stderr
}
