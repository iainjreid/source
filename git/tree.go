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
	"log"
	"path"
)

type TreeEntry struct {
	Name         string
	Path         string
	Hash         string
	IsFile       bool
	ChildEntries *TreeEntryMap
	Commit       *Commit

	// Used in the templates to control spacing
	Depth int
}

// NewTreeEntry creates an object that represents either a directory or a file
// in a Git repository.
func NewTreeEntry(name, path, hash string, isFile bool, commit *Commit) *TreeEntry {
	return &TreeEntry{
		Name:         name,
		Path:         path,
		Hash:         hash,
		ChildEntries: NewTreeEntryMap(),
		IsFile:       isFile,
		Commit:       commit,
		Depth:        -1,
	}
}

func (t *TreeEntry) AddChildDir(treeEntry *TreeEntry) {
	t.ChildEntries.AddDir(treeEntry)
}

func (t *TreeEntry) AddChildFile(treeEntry *TreeEntry) {
	t.ChildEntries.AddFile(treeEntry)
}

func (t *TreeEntry) ParentPath() string {
	return path.Dir(t.Path)
}

type TreeEntryMap struct {
	DirEntries  map[string]*TreeEntry
	FileEntries map[string]*TreeEntry

	Root *TreeEntry
}

// NewTreeEntryMap creates a container for directory and file entries.
func NewTreeEntryMap() *TreeEntryMap {
	return &TreeEntryMap{
		DirEntries:  make(map[string]*TreeEntry),
		FileEntries: make(map[string]*TreeEntry),
	}
}

func (t *TreeEntryMap) AddDir(treeEntry *TreeEntry) {
	t.DirEntries[path.Clean(treeEntry.Path)] = treeEntry
}

func (t *TreeEntryMap) RemoveDir(treeEntry *TreeEntry) {
	delete(t.DirEntries, path.Clean(treeEntry.Path))
}

func (t *TreeEntryMap) AddFile(treeEntry *TreeEntry) {
	t.FileEntries[path.Clean(treeEntry.Path)] = treeEntry
}

func (t *TreeEntryMap) GetDir(dirPath string) (*TreeEntry, error) {
	if dir, ok := t.DirEntries[path.Clean(dirPath)]; !ok {
		return nil, NewTreeEntryNotFoundError(dirPath)
	} else {
		return dir, nil
	}
}

func (t *TreeEntryMap) GetFile(filePath string) (*TreeEntry, error) {
	if dir, ok := t.FileEntries[path.Clean(filePath)]; !ok {
		return nil, NewTreeEntryNotFoundError(filePath)
	} else {
		return dir, nil
	}
}

func (t *TreeEntryMap) EnsureParents(treeEntry *TreeEntry, backoff int) {
	if backoff == 0 {
		panic("root too deep")
	}

	if treeEntry.Path == "." {
		return
	}

	if parent, _ := t.GetDir(treeEntry.ParentPath()); parent == nil {
		t.AddDir(&TreeEntry{
			Path:         path.Clean(treeEntry.ParentPath()),
			ChildEntries: NewTreeEntryMap(),
		})
	}

	t.EnsureParents(t.DirEntries[treeEntry.ParentPath()], backoff-1)
}

func (t *TreeEntryMap) SetDepth(treeEntry *TreeEntry, depth int) {
	for dirPath, dirEntry := range treeEntry.ChildEntries.DirEntries {
		if dirEntry.Depth > -1 {
			log.Fatalf("self-referencing tree: %s", dirPath)
		}

		dirEntry.Depth = depth + 1
		t.SetDepth(dirEntry, depth+1)
	}
}

func (t *TreeEntryMap) InitHierarchy(rootPath string) error {
	t.Root = &TreeEntry{
		Path:         path.Clean(rootPath),
		ChildEntries: NewTreeEntryMap(),
	}

	t.EnsureParents(t.Root, 99)
	t.AddDir(t.Root)

	// For every directory path we hold, determine its parent directory and assign
	// it as a child.
	for dirPath := range t.DirEntries {
		parentPath, _ := path.Split(dirPath)

		if dir, err := t.GetDir(parentPath); err == nil {
			dir.AddChildDir(t.DirEntries[dirPath])
		} else {
			return err
		}
	}

	// For every file path we hold, determine its parent directory and assign it
	// as a child.
	for filePath := range t.FileEntries {
		parentPath, _ := path.Split(filePath)

		if dir, err := t.GetDir(parentPath); err == nil {
			dir.AddChildFile(t.FileEntries[filePath])
		} else {
			return err
		}
	}

	// Whilst processing the hierarchy, the root directory will have been assigned
	// to itself as a known child.
	//
	// To avoid self-referencing, we will now remove the root directory from its
	// map of known children.
	//
	// TODO: If this function is improved to handle self-referencing hierarchies,
	//       this line may be removed, as it will no longer be required.
	t.Root.ChildEntries.RemoveDir(t.Root)

	// Assign depth values to every child directory.
	t.SetDepth(t.Root, 0)

	return nil
}
