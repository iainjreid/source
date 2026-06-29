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
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/google/uuid"
	"github.com/iainjreid/source/storage/driver"
)

func CloneRepo(ctx context.Context, store driver.Store, url string) (*Repo, error) {
	slog.Info("cloning repo", "url", url)

	name, err := GetRepoName(url)
	if err != nil {
		return nil, err
	}

	exists, err := store.RepoExists(ctx, name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("Repo already exists")
	}

	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	if err := store.CreateRepo(ctx, driver.Repo{
		ID:   id.String(),
		Name: name,
	}); err != nil {
		return nil, fmt.Errorf("error whilst creating repo: %v", err)
	}

	storer, err := store.ToStorer(ctx, name)
	if err != nil {
		return nil, err
	}

	repo, err := git.Clone(storer, nil, &git.CloneOptions{
		URL:          url,
		Progress:     io.Discard,
		Mirror:       true,
		NoCheckout:   true,
		SingleBranch: false,
	})

	if err != nil {
		slog.Error("error whilst cloning the repository", "err", err)
	}

	return &Repo{
		Name: name,
		Repo: repo,
	}, err
}

func ConsumeIter[T interface{}](iter storer.ReferenceIter, formatter func(ref *plumbing.Reference) T) ([]T, error) {
	defer iter.Close()

	var out []T

	err := iter.ForEach(func(ref *plumbing.Reference) error {
		out = append(out, formatter(ref))
		return nil
	})

	return out, err
}

type Line = git.Line

func NextNonMergeCommit(iter object.CommitIter) (*object.Commit, error) {
	var commit *object.Commit
	var err error

	for {
		commit, err = iter.Next()

		if len(commit.ParentHashes) <= 1 {
			break // break out of the loop
		}
	}

	return commit, err
}

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val

}

func GetRepoName(repoUrl string) (string, error) {
	u, err := url.Parse(repoUrl)
	if err != nil {
		return "", err
	}

	name := path.Base(u.Path)
	name = strings.TrimSuffix(name, ".git")

	return name, nil
}
