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
	"log/slog"
	"net/http"

	"github.com/iainjreid/source/internal/git"
	"github.com/iainjreid/source/internal/utils"
)

type TreeResponse struct {
	git.TreeEntryMap
	utils.Timeable

	repo   *git.Repo
	commit *git.Commit
}

func (t *TreeResponse) LoadCommit(rev string) error {
	var err error

	if t.commit, err = t.repo.GetCommit(rev); err != nil {
		return err
	}

	return err
}

func (t *TreeResponse) LoadTree() error {
	treeEntryMap, err := t.commit.GetTree("", false)

	if err == nil {
		t.FileEntries = treeEntryMap.Root.ChildEntries.FileEntries
		t.DirEntries = treeEntryMap.Root.ChildEntries.DirEntries
	}

	return err
}

func (h *Handlers) GetTree(w http.ResponseWriter, r *http.Request) {
	slog.Info("get tree")
	repo := r.PathValue("repo")
	hash := r.PathValue("hash")

	if repo == "" || hash == "" {
		h.SendErr(w, 400, "Must supply repo and hash")
		return
	}
	slog.Info("get tree")

	if repo, err := git.OpenRepo(r.Context(), h.Store, repo); err != nil {
		h.SendErr(w, 500, err.Error())
	} else {
		treeResponse := TreeResponse{
			repo: repo,
		}
		treeResponse.StartClock()

		if err := treeResponse.LoadCommit(hash); err != nil {
			h.SendErr(w, 500, err.Error())
			return
		}

		if err := treeResponse.LoadTree(); err != nil {
			h.SendErr(w, 500, err.Error())
			return
		}

		treeResponse.StopClock()
		h.SendJSON(w, r, treeResponse)
	}
}
