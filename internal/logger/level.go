package logger

import (
	"fmt"
	"log/slog"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "debug"
	case LevelError:
		return "error"
	default:
		return fmt.Sprintf("%d", int(l))
	}
}

func (l *Level) Set(str string) error {
	switch str {
	case "debug":
		*l = LevelDebug
	case "info":
		*l = LevelInfo
	case "warn":
		*l = LevelWarn
	case "error":
		*l = LevelError
	default:
		return InvalidLevelString(str)
	}

	return nil
}

func (l Level) ToSlogLevel() slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type InvalidLevelString string

func (e InvalidLevelString) Error() string {
	return fmt.Sprintf("logger: invalid level string '%s'", string(e))
}
