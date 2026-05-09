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

package main

import (
	"testing"

	"github.com/iainjreid/source/cmd/src/internal/cli"
)

func TestBinary(t *testing.T) {
	test := cli.PrepareBinary(t)

	stdout, stderr := test(cli.BadUsage)
	if got, want := stdout.String(), ""; got != want {
		t.Errorf("stdout want '%s' got '%s'", want, got)
	}

	if got, want := stderr.String(), "Run 'src --help' for usage.\n"; got != want {
		t.Errorf("stderr want '%s' got '%s'", want, got)
	}

	stdout, stderr = test(cli.BadUsage, "start")
	if got, want := stdout.String(), ""; got != want {
		t.Errorf("stdout want '%s' got '%s'", want, got)
	}

	if got, want := stderr.String(), "Run 'src start --help' for usage.\n"; got != want {
		t.Errorf("stderr want '%s' got '%s'", want, got)
	}

	stdout, stderr = test(cli.BadUsage, "manage")
	if got, want := stdout.String(), ""; got != want {
		t.Errorf("stdout want '%s' got '%s'", want, got)
	}

	if got, want := stderr.String(), "Run 'src manage --help' for usage.\n"; got != want {
		t.Errorf("stderr want '%s' got '%s'", want, got)
	}

	stdout, stderr = test(cli.BadUsage, "verify")
	if got, want := stdout.String(), ""; got != want {
		t.Errorf("stdout want '%s' got '%s'", want, got)
	}

	if got, want := stderr.String(), "Run 'src verify --help' for usage.\n"; got != want {
		t.Errorf("stderr want '%s' got '%s'", want, got)
	}
}
