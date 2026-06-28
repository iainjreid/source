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

type RefsResponse struct {
	Branches map[string]string `json:"branches"`
	Tags     map[string]string `json:"tags"`
	utils.Timeable
}

func NewRefsResponse() *RefsResponse {
	return &RefsResponse{
		Branches: make(map[string]string),
		Tags:     make(map[string]string),
	}
}

func (r *RefsResponse) AddBranch(name string, hash string) {
	r.Branches[name] = hash
}

func (r *RefsResponse) AddTag(name string, hash string) {
	r.Tags[name] = hash
}

func (h *Handlers) GetRefs(w http.ResponseWriter, r *http.Request) {
	repo := r.PathValue("repo")

	if repo == "" {
		h.SendErr(w, 400, "Must supply repo")
		return
	}

	if iter, err := h.Store.IterateRefs(r.Context(), repo); err != nil {
		h.SendErr(w, 500, err.Error())
	} else {
		defer iter.Close()

		refsResponse := NewRefsResponse()
		refsResponse.StartClock()

		err := iter.ForEach(func(ref *driver.Ref) error {
			switch {
			case ref.Name().IsBranch():
				refsResponse.AddBranch(ref.Name().String(), ref.Hash().String())
			case ref.Name().IsTag():
				refsResponse.AddTag(ref.Name().String(), ref.Hash().String())
			}
			return nil
		})

		if err != nil {
			h.SendErr(w, 500, err.Error())
		} else {
			refsResponse.StopClock()
			h.SendJSON(w, r, refsResponse)
		}
	}
}
