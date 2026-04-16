// Package logger is an internal package used to define opinionated project-wide
// logging defaults.
//
// Upstream consumers of code exported by this project should not be impacted by
// the stylistic logging choices made here, but by defining their own handlers
// for the default structured logging provider can still benefit from internal
// logging within this codebase.
package logger

import (
	"log/slog"
	"os"
)

// Init initialises the default log/slog logger.
func Init(level slog.Level, debug bool, attrs []slog.Attr) {
	var handler slog.Handler

	// If debug is enable, override the log level accordingly.
	if debug {
		level = slog.LevelDebug
	}

	// Logs are written to STDOUT and will include source references if the
	// debug flag is set.
	handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: debug,
	})

	// TODO: Document which attributes are to be included in logs.
	handler = handler.WithAttrs(attrs)

	// Set the default structured logger with our opinionated one.
	slog.SetDefault(slog.New(handler))
}
