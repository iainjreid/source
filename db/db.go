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
