package logger_test

import (
	"log/slog"
	"testing"

	"github.com/iainjreid/source/internal/logger"
)

var leveltests = []struct {
	in   string
	want slog.Level
}{
	{"debug", slog.LevelDebug},
	{"info", slog.LevelInfo},
	{"warn", slog.LevelWarn},
	{"error", slog.LevelError},
}

func TestLevelSetAndGet(t *testing.T) {
	var level logger.Level

	for _, tt := range leveltests {
		t.Run(tt.in, func(t *testing.T) {
			if err := level.Set(tt.in); err != nil {
				t.Errorf("Level.Set(%v) returned an unexpected error = %v", tt.in, err)
			}

			if got := level.ToSlogLevel(); got != tt.want {
				t.Errorf("Level(%s).ToSlogLevel() = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

var levelfailures = []struct {
	in string
}{
	{"invalid"},
	{""},
}

func TestInvalidLevel(t *testing.T) {
	var level logger.Level

	for _, tt := range levelfailures {
		t.Run(tt.in, func(t *testing.T) {
			if got, want := level.Set(tt.in), logger.InvalidLevelString(tt.in); got != want {
				t.Errorf("Level.Set(%v) = %v, want %s", tt.in, got, want)
			}
		})
	}
}
