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

package cli

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"
)

func PrepareBinary(t *testing.T) func(int, ...string) (*bytes.Buffer, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	bin := filepath.Join(t.TempDir(), "src")

	cmd := exec.Command("go", "build", "-o", bin, ".")

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Run()

	return func(exitCode int, arg ...string) (*bytes.Buffer, *bytes.Buffer) {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		cmd = exec.Command(bin, arg...)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()

		switch err := err.(type) {
		case *exec.ExitError:
			if got, want := err.ExitCode(), exitCode; got != want {
				t.Errorf("Unexpected exit code: %d", got)
			}
		default:
			t.Errorf("Unexpected error type: %v", err)
		}

		return stdout, stderr
	}
}
