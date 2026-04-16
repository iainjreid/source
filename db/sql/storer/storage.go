package storer

import (
	"database/sql"

	"github.com/go-git/go-git/v5/storage/memory"
)

// Storage is a struct implementation of the [storage.Storer] interface.
type Storage struct {
	ObjectStorage
	ReferenceStorage

	// The following structs are duplicated from the go-git memory storage
	// implementation, however they are omitted from the declaration below.
	memory.IndexStorage
	memory.ShallowStorage
	memory.ModuleStorage
	memory.ConfigStorage
}

// NewStorage creates a generic SQL-backed storage driver that uses the supplied
// database client for all operations.
func NewStorage(client *sql.DB) *Storage {
	return &Storage{
		ReferenceStorage: ReferenceStorage{client: client},
		ObjectStorage:    ObjectStorage{client: client},
	}
}
