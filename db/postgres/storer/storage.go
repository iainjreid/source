package storer

import (
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		ReferenceStorage: ReferenceStorage{pool: pool},
		ObjectStorage:    ObjectStorage{pool: pool},
	}
}
