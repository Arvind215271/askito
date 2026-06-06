// ./internal/logger/logger.go
package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
	isDev bool
}

func New(env string) *Logger {
	isDev := env == "development"

	level := slog.LevelInfo
	if isDev {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handler slog.Handler

	if isDev {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
		isDev:  isDev,
	}
}

func (l *Logger) IsDev() bool {
	return l.isDev
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// fatal aka we should stop the server from continuing as it breaks the whole logic 
func (l *Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}