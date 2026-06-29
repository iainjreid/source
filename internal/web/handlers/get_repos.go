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

	"github.com/iainjreid/source/internal/utils"
	"github.com/iainjreid/source/storage/driver"
)

type ReposResponse struct {
	Repos []driver.Repo `json:"repos"`
	utils.Timeable
}

func NewReposResponse() *ReposResponse {
	return &ReposResponse{}
}

func (r *ReposResponse) AddRepo(repo driver.Repo) {
	r.Repos = append(r.Repos, repo)
}

func (h *Handlers) GetRepos(w http.ResponseWriter, r *http.Request) {
	refsResponse := NewReposResponse()
	refsResponse.StartClock()

	if repos, err := h.Store.ListRepos(r.Context()); err != nil {
		h.SendErr(w, 500, err.Error())
	} else {
		for _, repo := range repos {
			refsResponse.AddRepo(repo)
		}

		if err != nil {
			h.SendErr(w, 500, err.Error())
		} else {
			refsResponse.StopClock()
			h.SendJSON(w, r, refsResponse)
		}
	}
}
