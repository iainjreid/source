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
	"fmt"
	"net/url"
	"sync"

	"github.com/go-git/go-git/v5/storage"
	"github.com/iainjreid/source/storage/driver"
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

	//go:embed queries/update_ref.sql
	UpdateRef string
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
// proxied (not sure of the Golang specific term for this?) to keep upstream
// imports clean.
var (
	ErrReferenceHasChanged = storage.ErrReferenceHasChanged
)

var (
	mu      sync.RWMutex
	drivers = map[string]driver.Driver{}
)

func Register(scheme string, driver driver.Driver) {
	mu.Lock()
	defer mu.Unlock()

	if driver == nil {
		panic("driver cannot be nil")
	}

	if _, exists := drivers[scheme]; exists {
		panic("driver already registered: " + scheme)
	}

	drivers[scheme] = driver
}

func Lookup(scheme string) (driver.Driver, error) {
	mu.RLock()
	defer mu.RUnlock()

	d, ok := drivers[scheme]
	if !ok {
		return nil, fmt.Errorf("unsupported database scheme %q", scheme)
	}

	return d, nil
}

// Open retrieves the driver from the registry that matches the scheme in the
// provided URL.
func Open(ctx context.Context, rawURL string) (driver.Store, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	d, err := Lookup(u.Scheme)
	if err != nil {
		return nil, err
	}

	return d.Open(ctx, rawURL)
}
