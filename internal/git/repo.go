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

package git

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/iainjreid/source/storage/driver"
)

type Repo struct {
	Name        string
	Description string

	Repo *git.Repository
	// Store driver.Store
}

// This struct is subject to removal, and will be merged with the Ref type from
// the driver package in future.
type Branch struct {
	RepoName string `json:"-"`
	Hash     string `json:"hash"`
	Name     string `json:"name"`
}

func OpenRepoById(ctx context.Context, store driver.Store, id string) (*Repo, error) {
	storer, err := store.ToStorer(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Repo{
		Repo: &git.Repository{
			Storer: storer,
		},
	}, nil
}

func OpenRepo(ctx context.Context, store driver.Store, name string) (*Repo, error) {
	repo, err := store.GetRepo(ctx, name)
	if err != nil {
		return nil, err
	}

	return OpenRepoById(ctx, store, repo.ID)
}

// Branches returns all the References that are branches.
func (r *Repo) Branches() ([]Branch, error) {
	iter, err := r.NewFilteredReferencesIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsBranch()
	})

	if err != nil {
		return nil, err
	}

	return ConsumeIter(iter, func(ref *plumbing.Reference) Branch {
		return Branch{
			RepoName: r.Name,
			Hash:     ref.Hash().String(),
			Name:     ref.Name().Short(),
		}
	})
}

func (r *Repo) NewFilteredReferencesIter(fn func(ref *plumbing.Reference) bool) (storer.ReferenceIter, error) {
	iter, err := r.Repo.Storer.IterReferences()

	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(r *plumbing.Reference) bool {
		return fn(r)
	}, iter), nil
}

func (r *Repo) GetCommit(revision string) (*Commit, error) {
	return r.GetCommitByHash(r.ResolveRevOrHash(revision))
}

func (r *Repo) ResolveRevOrHash(revOrHash string) *plumbing.Hash {
	if hash, err := r.Repo.ResolveRevision(plumbing.Revision(revOrHash)); err == nil {
		return hash
	} else {
		hash := plumbing.NewHash(revOrHash)
		return &hash
	}
}

func (r *Repo) GetCommitByHash(hash *plumbing.Hash) (*Commit, error) {
	commit, err := r.Repo.CommitObject(*hash)

	if err != nil {
		return nil, &RevisionNotFoundError{
			Revision: hash.String(),
		}
	}

	return &Commit{
		Hash:    commit.Hash.String(),
		Message: commit.Message,
		repo:    r,
		Commit:  commit,
	}, nil
}
