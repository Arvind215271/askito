// ./internal/logger/logger.go
package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
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

// With returns a new Logger containing the provided structural attributes,
// while preserving the environment state flags.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		isDev:  l.isDev,
	}
}


// Fatal prints an error log message with structured fields and completely halts execution.
// It bypasses standard wrapper limitations to ensure that Go's 'AddSource: true' 
// option prints the exact line where Fatal was called, instead of masking it inside this package.
func (l *Logger) Fatal(msg string, args ...any) {
	// 1. Fast-path optimization: Check if Error level logging is enabled globally.
	// If it's disabled, skip the heavy stack tracing operations and exit immediately.
	if !l.Enabled(context.Background(), slog.LevelError) {
		os.Exit(1)
	}

	// 2. Allocate an array of size 1 to hold a single Program Counter (PC).
	// A Program Counter is a memory address pointing to a specific line of code.
	// Usually we learnt this in microprocessor course. Also in stack and heap memory of C and its working
	var pcs [1]uintptr

	// 3. Unwind the application's runtime execution call stack.
	// Skip frame 0 (runtime.Callers itself) and frame 1 (this Fatal helper function).
	// Frame 2 points precisely to the application code that triggered the fatal error.
	runtime.Callers(2, pcs[:])

	// 4. Manually construct a raw logging record from scratch.
	// We inject pcs[0] directly so that slog uses our custom-discovered source location 
	// rather than evaluating the current position inside internal/logger/logger.go.
	r := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])

	// 5. Append any variadic structured key-value arguments (e.g., "error", err) to the record.
	r.Add(args...)
	
	// 6. Push the manually assembled log record directly into the underlying handler 
	// (either TextHandler or JSONHandler) to format it and write it out to standard output.
	_ = l.Handler().Handle(context.Background(), r)

	// 7. Hard-crash the Go process with exit code 1.
	// This signals to Docker, Kubernetes, or systemd that the app failed unrecoverably.
	os.Exit(1)
}
