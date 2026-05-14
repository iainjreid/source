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
	"testing"

	"github.com/iainjreid/source/cmd/src/internal/cli"
)

func TestStringSlice(t *testing.T) {
	testTable := []struct {
		in, out string
	}{
		{
			"the",
			"the",
		},
		{
			"cat",
			"the,cat",
		},
		{
			"in the",
			"the,cat,in the",
		},
		{
			"hat",
			"the,cat,in the,hat",
		},
	}

	var stringSlice cli.StringSlice

	for _, tt := range testTable {
		t.Run(tt.in, func(t *testing.T) {
			if err := stringSlice.Set(tt.in); err != nil {
				t.Errorf("StringSlice.Set(%v) returned an unexpected error = %v", tt.in, err)
			}

			if got := stringSlice.String(); got != tt.out {
				t.Errorf("StringSlice().String() = %v, want %v", got, tt.out)
			}
		})
	}
}
