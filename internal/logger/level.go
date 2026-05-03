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

package logger

import (
	"fmt"
	"log/slog"
)

// A Level is an abstraction over the [slog.Level] type that implements the
// [flag.Value] interface.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String satisfies the [flag.Value] interface, returning a string form of the
// underlying Level.
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
		return ""
	}
}

// Set satisfies the [flag.Value] interface and updates the Level. It accepts a
// string and may return an error if it cannot be parsed.
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
		return LevelStringError(str)
	}

	return nil
}

// ToSlogLevel returns the appropriate [slog.Level] for the Level.
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

// A LevelStringError records a [Level] string that was unable to be parsed.
type LevelStringError string

func (e LevelStringError) Error() string {
	return fmt.Sprintf("logger: invalid level string '%s'", string(e))
}
