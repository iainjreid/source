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

package driver

import (
	"context"

	"github.com/go-git/go-git/v5/storage"
)

type Driver interface {
	Open(context.Context, string) (Store, error)
}

// Repo represents a repository that may exist within the database, but we do
// not need to know this to still perform optimistic queries on its behalf.
//
// For databases that support joins, or equivalent behaviour, knowing just the
// name of the repository will allow us to defer the existance check to the
// database itself when the time comes.
type Repo interface {
	Name() string

	Description() string

	Create(ctx context.Context) error

	Exists() (bool, error)

	storage.Storer
}

// Store is an interface representing the expected shape of a storage
// abstraction layer.
type Store interface {
	// Protocol is expected to return the protocol that the Plugin supports. If
	// this protocol doesn't match the expected value (if the Plugin has been
	// installed under a different name) then the program will exit.
	Protocol() string

	// EnsureReady is expected to create or check for the existance of all
	// tables, indexes, and other functionality required for a given database to
	// accept writes.
	EnsureReady(context.Context) error

	// Repo is expected to return a struct that implements the [Repo] interface.
	// The repository may or may not exist, but it is not for this method to
	// retermine that.
	//
	// If the repository data is required at any time, it will be loaded using
	// the appropriate method on the [Repo] interface.
	Repo(name string) Repo

	// ListRepos should return a slice of the available repositories stored with
	// the database.
	ListRepos(context.Context) ([]Repo, error)
}
