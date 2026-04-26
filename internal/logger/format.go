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
)

type Format int

const (
	FormatText Format = iota
	FormatJSON
)

func (f Format) String() string {
	switch f {
	case FormatText:
		return "text"
	case FormatJSON:
		return "json"
	default:
		return fmt.Sprintf("%d", int(f))
	}
}

func (f *Format) Set(str string) error {
	switch str {
	case "text":
		*f = FormatText
	case "json":
		*f = FormatJSON
	default:
		return InvalidFormatString(str)
	}

	return nil
}

type InvalidFormatString string

func (e InvalidFormatString) Error() string {
	return fmt.Sprintf("logger: invalid format string '%s'", string(e))
}
