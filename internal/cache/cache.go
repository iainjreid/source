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

package cache

import (
	"sync"

	"github.com/go-git/go-git/v5/plumbing"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/iainjreid/source/internal/utils"
)

type Namespace struct {
	Refs *lru.Cache[string, *plumbing.Reference]
	Objs *lru.Cache[string, plumbing.EncodedObject]
}

var (
	mu         sync.RWMutex
	namespaces = map[string]Namespace{}
)

func Register(name string) Namespace {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := namespaces[name]; !exists {
		namespaces[name] = Namespace{
			Refs: utils.Must(lru.New[string, *plumbing.Reference](1028)),
			Objs: utils.Must(lru.New[string, plumbing.EncodedObject](1028)),
		}
	}

	return namespaces[name]
}
