package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	logger slog.Logger
}

func NewLogger(cfg *Config) *Logger {
	opts := &slog.HandlerOptions{
		Level: getLoggerLevel(cfg.Level),
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return &Logger{logger: *logger}
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func getLoggerLevel(LogLevel string) slog.Level {
	switch strings.ToLower(LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "fatal":
		return slog.LevelError
	case "panic":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
