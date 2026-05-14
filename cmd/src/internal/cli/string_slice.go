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
	"strings"
)

// A StringSlice
type StringSlice []string

// String satisfies the [flag.Value] interface, returning a comma delimited
// string made up of the underlying slice.
func (s *StringSlice) String() string {
	return strings.Join(*s, ",")
}

// Set satisfies the [flag.Value] interface and appends the provided string to
// the underlying slice.
func (s *StringSlice) Set(str string) error {
	*s = append(*s, str)
	return nil
}
