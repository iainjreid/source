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
	"bytes"
	"log/slog"
	"regexp"
	"testing"

	"github.com/iainjreid/source/internal/logger"
)

func TestNewWithText(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf, slog.LevelInfo, false, false, nil)
	l.Info("message")

	if log, expr := buf.Bytes(), `^time=.* level=INFO msg=message\n$`; !satisfiesExpr(log, expr) {
		t.Errorf("expected %s to match regex %v", log, expr)
	}
}

func TestNewWithJSON(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf, slog.LevelInfo, true, false, nil)
	l.Info("message")

	if log, expr := buf.Bytes(), `^{"time":".*","level":"INFO","msg":"message"}\n$`; !satisfiesExpr(log, expr) {
		t.Errorf("expected %s to match regex %v", log, expr)
	}
}

func TestNewWithLogLevelWarn(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf, slog.LevelWarn, false, false, nil)
	l.Info("info message")
	l.Warn("warn message")

	if log, expr := buf.Bytes(), `^time=.* level=WARN msg="warn message"\n$`; !satisfiesExpr(log, expr) {
		t.Errorf("expected %s to match regex %v", log, expr)
	}
}

func TestNewWithDebugOverride(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf, slog.LevelError, false, true, nil)
	l.Debug("message")

	if log, expr := buf.Bytes(), `^time=.* level=DEBUG source=.* msg=message\n$`; !satisfiesExpr(log, expr) {
		t.Errorf("expected %s to match regex %v", log, expr)
	}
}

func TestInit(t *testing.T) {
	old := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(old)
	})

	logger.Init(slog.LevelInfo, false, false, nil)
	if slog.Default() == old {
		t.Error("expected default logger to have been updated")
	}
}

func satisfiesExpr(b []byte, expr string) bool {
	return regexp.MustCompile(expr).Match(b)
}
