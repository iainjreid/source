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
	"fmt"
)

type TreeEntryNotFoundError struct {
	Name string
}

func NewTreeEntryNotFoundError(name string) *TreeEntryNotFoundError {
	return &TreeEntryNotFoundError{
		Name: name,
	}
}

func (e *TreeEntryNotFoundError) Error() string {
	return fmt.Sprintf("TreeEntry not found: %s", e.Name)
}

type RepoNotFoundError struct {
	Name string
}

func (e *RepoNotFoundError) Error() string {
	return fmt.Sprintf("repo not found: %s", e.Name)
}

type RevisionNotFoundError struct {
	Revision string
}

func (e *RevisionNotFoundError) Error() string {
	return fmt.Sprintf("revision not found: %s", e.Revision)
}

type FileNotFoundError struct {
	Filepath string
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("file not found: %s", e.Filepath)
}

type DirectoryNotFoundError struct {
	Dirpath string
}

func (e *DirectoryNotFoundError) Error() string {
	return fmt.Sprintf("directory not found: %s", e.Dirpath)
}
