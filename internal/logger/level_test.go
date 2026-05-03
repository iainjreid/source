// Copyright 2026 Iain J. Reid
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
			if got, want := level.Set(tt.in), logger.LevelStringError(tt.in); got != want {
				t.Errorf("Level.Set(%v) = %v, want %s", tt.in, got, want)
			}
		})
	}
}
