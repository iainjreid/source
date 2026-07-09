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

package config

import (
	"log/slog"
	"strconv"
)

// String is a string-backed configuration value.
type String string

// StringValue returns a new String initialised with the provided value.
func StringValue(str string) *String {
	val := String(str)
	return &val
}

// String satisfies fmt.Stringer and flag.Value.
func (s *String) String() string {
	return s.Value()
}

// Value returns the underlying string value.
func (s *String) Value() string {
	return string(*s)
}

// Set satisfies flag.Value.
func (s *String) Set(str string) error {
	*s = String(str)
	return nil
}

// Int is an integer-backed configuration value.
type Int int

// IntValue returns a new Int initialised with the provided value.
func IntValue(i int) *Int {
	val := Int(i)
	return &val
}

// String satisfies fmt.Stringer and flag.Value.
func (i *Int) String() string {
	return strconv.Itoa(i.Value())
}

// Value returns the underlying integer value.
func (i *Int) Value() int {
	return int(*i)
}

// Set satisfies flag.Value.
func (i *Int) Set(str string) error {
	v, err := strconv.ParseInt(str, 0, strconv.IntSize)
	if err == nil {
		*i = Int(v)
	}

	return err
}

// SlogLevel is a slog.Level-backed configuration value.
type SlogLevel slog.Level

// SlogLevelValue returns a new SlogLevel initialised with the provided value.
func SlogLevelValue(l slog.Level) *SlogLevel {
	val := SlogLevel(l)
	return &val
}

// String satisfies fmt.Stringer and flag.Value.
func (s *SlogLevel) String() string {
	return s.Value().String()
}

// Value returns the underlying slog.Level.
func (s *SlogLevel) Value() slog.Level {
	return slog.Level(*s)
}

// Set satisfies flag.Value.
func (s *SlogLevel) Set(str string) error {
	var level slog.Level

	if err := level.UnmarshalText([]byte(str)); err != nil {
		return err
	}

	*s = SlogLevel(level)

	return nil
}
