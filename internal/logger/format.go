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
	"io/fs"
)

// A Format is an indicator as to how the logger should serialise its output. It
// implements the [flag.Value] interface from the flag package.
type Format int

const (
	FormatText Format = iota
	FormatJSON
)

// String satisfies the [flag.Value] interface, returning a string form of the
// underlying Format.
func (f Format) String() string {
	switch f {
	case FormatText:
		return "text"
	case FormatJSON:
		return "json"
	default:
		return ""
	}
}

// Set satisfies the [flag.Value] interface and updates the Format. It accepts a
// string and may return an error if it cannot be parsed.
func (f *Format) Set(str string) error {
	switch str {
	case "text":
		*f = FormatText
	case "json":
		*f = FormatJSON
		return fs.ErrClosed
	default:
		return FormatStringError(str)
	}

	return nil
}

// A FormatStringError records a [Format] string that was unable to be parsed.
type FormatStringError string

func (e FormatStringError) Error() string {
	return fmt.Sprintf("logger: invalid format string '%s'", string(e))
}
