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
func Init(level Level, format Format, debug bool, attrs []slog.Attr) {
	var handler slog.Handler

	// If debug is enable, override the log level accordingly.
	if debug {
		level = LevelDebug
	}

	// Logs are written to STDOUT and will include source references if the
	// debug flag is set.
	opts := &slog.HandlerOptions{
		Level:     level.ToSlogLevel(),
		AddSource: debug,
	}

	switch format {
	case FormatText:
		handler = slog.NewTextHandler(os.Stdout, opts)
	case FormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	// TODO: Document which attributes are to be included in logs.
	handler = handler.WithAttrs(attrs)

	// Set the default structured logger with our opinionated one.
	slog.SetDefault(slog.New(handler))
}
