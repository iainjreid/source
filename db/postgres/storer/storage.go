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
