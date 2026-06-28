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

package utils_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/iainjreid/source/internal/utils"
)

func TestTimeableStartStopClock(t *testing.T) {
	var tm utils.Timeable

	tm.StartClock()

	time.Sleep(20 * time.Millisecond)

	tm.StopClock()

	if tm.TimeElapsed < 20 {
		t.Fatalf("TimeElapsed = %.1fms, want >= 20ms", tm.TimeElapsed)
	}
}

func TestTimeableStartStopClockImmediately(t *testing.T) {
	var tm utils.Timeable

	tm.StartClock()
	tm.StopClock()

	if tm.TimeElapsed < 0 {
		t.Fatalf("TimeElapsed = %.1fms, want >= 0", tm.TimeElapsed)
	}
}

func TestTimeableToJSON(t *testing.T) {
	tm := utils.Timeable{
		TimeElapsed: 12.3,
	}

	data, err := json.Marshal(tm)
	if err != nil {
		t.Fatal(err)
	}

	const want = `{"timeElapsed":12.3}`

	if string(data) != want {
		t.Fatalf("got %q, want %q", data, want)
	}
}
