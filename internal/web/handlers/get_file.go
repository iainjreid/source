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

package handlers

import (
	"net/http"
	"path"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/iainjreid/source/internal/git"
	"github.com/iainjreid/source/internal/utils"
)

type BlobResponse struct {
	FileName         string   `json:"fileName"`
	FileBytes        int64    `json:"fileBytes"`
	FileLines        []string `json:"fileLines"`
	LineCount        int      `json:"lineCount"`
	LastCommitHash   string   `json:"lastCommitHash"`
	LastCommitMsg    string   `json:"lastCommitMsg"`
	LastCommitAuthor string   `json:"lastCommitAuthor"`
	utils.Timeable

	repo   *git.Repo
	commit *git.Commit
}

func (b *BlobResponse) LoadLastCommit(rev string, filepath string) error {
	iter, err := b.repo.Repo.Log(&gogit.LogOptions{
		From:     plumbing.NewHash(rev),
		FileName: &filepath,
	})
	if err != nil {
		return err
	}
	defer iter.Close()

	var result *object.Commit

	err = iter.ForEach(func(commit *object.Commit) error {
		tree, err := commit.Tree()
		if err != nil {
			return err
		}

		hash, err := fileHash(tree, filepath)
		if err != nil {
			return nil
		}

		if commit.NumParents() == 0 {
			result = commit
			return storer.ErrStop
		}

		parent, err := commit.Parent(0)
		if err != nil {
			return err
		}

		parentTree, err := parent.Tree()
		if err != nil {
			return err
		}

		parentHash, err := fileHash(parentTree, filepath)
		if err != nil {
			result = commit
			return storer.ErrStop
		}

		if hash != parentHash {
			result = commit
			return storer.ErrStop
		}

		return nil
	})

	if err != nil && err != storer.ErrStop {
		return err
	}

	b.commit = &git.Commit{
		Commit: result,
	}

	b.LastCommitAuthor = result.Author.Email
	b.LastCommitHash = result.Hash.String()
	b.LastCommitMsg = result.Message

	return err
}

func fileHash(tree *object.Tree, path string) (plumbing.Hash, error) {
	file, err := tree.File(path)
	if err != nil {
		return plumbing.ZeroHash, err
	}
	return file.Hash, nil
}

func (b *BlobResponse) LoadFile(filepath string) error {
	lines, size, err := b.commit.GetFileContents(filepath, false)

	if err != nil {
		return err
	}

	for _, line := range lines {
		b.FileLines = append(b.FileLines, line.Text)
	}

	_, b.FileName = path.Split(filepath)
	b.FileBytes = size
	b.LineCount = len(lines)

	return err
}

func (h *Handlers) GetBlob(w http.ResponseWriter, r *http.Request) {
	repo := r.PathValue("repo")
	hash := r.PathValue("hash")
	filepath := r.PathValue("filepath")

	if repo == "" || hash == "" || filepath == "" {
		h.SendErr(w, 400, "Must supply repo, hash, and filepath")
		return
	}

	if repo, err := git.OpenRepoById(r.Context(), h.Store, repo); err != nil {
		h.SendErr(w, 500, err.Error())
	} else {
		blobResponse := BlobResponse{
			repo: repo,
		}
		blobResponse.StartClock()

		if err := blobResponse.LoadLastCommit(hash, filepath); err != nil {
			h.SendErr(w, 500, err.Error())
			return
		}

		if err := blobResponse.LoadFile(filepath); err != nil {
			h.SendErr(w, 500, err.Error())
			return
		}

		blobResponse.StopClock()
		h.SendJSON(w, r, blobResponse)
	}
}

func (h *Handlers) GetRawBlob(w http.ResponseWriter, r *http.Request) {
	repo := r.PathValue("repo")
	hash := r.PathValue("hash")
	filepath := r.PathValue("filepath")

	if repo == "" || hash == "" || filepath == "" {
		h.SendErr(w, 400, "Must supply repo, hash, and filepath")
		return
	}

	if repo, err := git.OpenRepoById(r.Context(), h.Store, repo); err != nil {
		h.SendErr(w, 500, err.Error())
	} else {
		_hash := plumbing.NewHash(hash)
		commit, err := repo.GetCommitByHash(&_hash)
		if err != nil {
			h.SendErr(w, 500, err.Error())
			return
		}

		file, err := commit.Commit.File(path.Clean(filepath))
		if err != nil {
			h.SendErr(w, 500, err.Error())
			return
		}

		contents, err := file.Reader()
		if err != nil {
			h.SendErr(w, 500, err.Error())
			return
		}

		h.SendPlain(w, r, contents)
	}
}
