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

// Package db provides behaviours associated with the storing of Git objects and
// references.
package db

import (
	"context"
)

// Type is an interface representing the expected shape of a database
// abstraction layer.
type DB interface {
	// TODO: Remove this method when repository importing is supported.
	HardReset(context.Context) error

	// EnsureReady is expected to create or check for the existance of all
	// tables, indexes, and other functionality required for a given database to
	// accept writes.
	EnsureReady(context.Context) error
}
