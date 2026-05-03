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
	"testing"

	"github.com/iainjreid/source/internal/logger"
)

var formattests = []struct {
	in string
}{
	{"text"},
	{"json"},
}

func TestFormatSetAndGet(t *testing.T) {
	var format logger.Format

	for _, tt := range formattests {
		t.Run(tt.in, func(t *testing.T) {
			if err := format.Set(tt.in); err != nil {
				t.Errorf("Format.Set(%v) returned an unexpected error = %v", tt.in, err)
			}

			if got := format.String(); got != tt.in {
				t.Errorf("Format(%s).String() = %v, want %v", tt.in, got, tt.in)
			}
		})
	}
}

var formatfailures = []struct {
	in string
}{
	{"invalid"},
	{""},
}

func TestInvalidFormat(t *testing.T) {
	var format logger.Format

	for _, tt := range formatfailures {
		t.Run(tt.in, func(t *testing.T) {
			if got, want := format.Set(tt.in), logger.FormatStringError(tt.in); got != want {
				t.Errorf("Format.Set(%v) = %v, want %s", tt.in, got, want)
			}
		})
	}
}
