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

package cli_test

import (
	"log/slog"
	"testing"

	"github.com/iainjreid/source/cmd/src/internal/cli"
)

func TestResolveLoggerFlags(t *testing.T) {
	testTable := []struct {
		explain        string
		quiet, verbose bool
		level          slog.Level
		err            error
	}{
		{
			"choosing quiet should result in error level logs",
			true, false, slog.LevelError, nil,
		},
		{
			"choosing verbose should result in debug level logs",
			false, true, slog.LevelDebug, nil,
		},
		{
			"no configuration should result in the expected default",
			false, false, slog.LevelInfo, nil,
		},
		{
			"incompatible values should result in an error",
			true, true, slog.LevelInfo, cli.LoggerFlagsError,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.explain, func(t *testing.T) {
			if level, err := cli.ResolveLoggerFlags(&tt.quiet, &tt.verbose); level != tt.level || err != tt.err {
				t.Errorf("ResolveLoggerFlags(%v, %v) = (%v, %v), want (%v, %v)", tt.quiet, tt.verbose, level, err, tt.level, tt.err)
			}
		})
	}
}
