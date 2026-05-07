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
	"flag"
)

func String(flag *flag.FlagSet, value string, long, short string) *string {
	out := flag.String(long, value, "")
	if short != "" {
		flag.StringVar(out, short, value, "")
	}
	return out
}

func Int(flag *flag.FlagSet, value int, long, short string) *int {
	out := flag.Int(long, value, "")
	if short != "" {
		flag.IntVar(out, short, value, "")
	}
	return out
}

func Bool(flag *flag.FlagSet, value bool, long, short string) *bool {
	out := flag.Bool(long, value, "")
	if short != "" {
		flag.BoolVar(out, short, value, "")
	}
	return out
}
