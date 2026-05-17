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

// Package storage provides behaviours associated with the storing of Git
// objects and references.
//
// For storage plugin implementations that use SQL, the provided embedded
// schemas and queries should hopefully be sufficiently portable to work in the
// vast majority of cases.
package storage

import (
	"context"
	_ "embed"

	"github.com/go-git/go-git/v5/storage"
)

// Schemas and queries that cover the base requirements when implementing a
// storage plugin. For their usage, refer to the documentation in each file.
var (
	//go:embed queries/count_refs.sql
	CountRefs string

	//go:embed queries/delete_ref.sql
	DeleteRef string

	//go:embed queries/get_ref.sql
	GetRefQuery string

	//go:embed queries/get_refs.sql
	GetRefsQuery string

	//go:embed queries/insert_ref.sql
	InsertRef string
)

// Three tables provide all of the core functionality currently offered by the
// project. In future more tables will be required to support user management
// and other behaviours.
var (
	//go:embed schemas/repos.sql
	ReposSchema string

	//go:embed schemas/objects.sql
	ObjectsSchema string

	//go:embed schemas/refs.sql
	RefsSchema string
)

// For those implementing their own storage plugin, the following schemas are
// optional. Be mindful that if you opt to not include these schemas, the join
// tables should also be omitted.
var (
	//go:embed schemas/topics.sql
	TopicsSchema string

	//go:embed schemas/languages.sql
	LanguagesSchema string

	//go:embed schemas/extensions.sql
	ExtensionsSchema string

	//go:embed schemas/licenses.sql
	LicensesSchema string
)

// Lastly, we have the join tables. Currently limited to providing support for
// repository topics at this time.
var (
	//go:embed schemas/repos_topics.sql
	ReposTopicsSchema string
)

// As the codebase is steared away from the hard dependency on internal code
// exposed by "go-git", some errors from the storage package need to be
// proxied (not sure of the Golang specific term for this?) to upstream imports
// clean.
var (
	ErrReferenceHasChanged = storage.ErrReferenceHasChanged
)

// A Repo represents a repository that may exist within the database, but we do
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

	Load(ctx context.Context) error

	IterObjects(ctx context.Context)

	// GetObject(ctx context.Context, objHash string)

	// IterRefs(ctx context.Context)

	// GetRef(ctx context.Context, refHash string)

	storage.Storer
}

// Storage is an interface representing the expected shape of a storage
// abstraction layer.
type Storage interface {
	// Protocol is expected to return the protocol that the Plugin supports. If
	// this protocol doesn't match the expected value (if the Plugin has been
	// installed under a different name) then the program will exit.
	Protocol() string

	// EnsureReady is expected to create or check for the existance of all
	// tables, indexes, and other functionality required for a given database to
	// accept writes.
	EnsureReady(context.Context) error

	// IterRepos(ctx context.Context, search any) iter.Seq[Repo]

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
