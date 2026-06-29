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
	"io"
	"log"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Given a commit is a type of Git object, the nomenclature here does not work.
// In future, this struct will be merged with other object types.
type Commit struct {
	Hash    string
	Message string
	Date    string
	repo    *Repo
	Commit  *object.Commit
}

// GetFileContents will attempt to read the contents of a file and return each
// line with the most recent blame information.
func (c *Commit) GetFileContents(filepath string, blame bool) ([]*Line, int64, error) {
	var lines []*Line
	var size int64
	var err error

	if blame {
		file, _err := git.Blame(c.Commit, path.Clean(filepath))
		if _err != nil {
			err = _err
		} else {
			lines = file.Lines
		}
	} else {
		file, _err := c.Commit.File(path.Clean(filepath))
		if _err != nil {
			err = _err
		} else {
			_lines, _ := file.Lines()
			size = file.Size
			lines = Map(_lines, func(_line string) *Line {
				return &Line{
					Text: _line,
				}
			})
		}
	}

	if err != nil {
		return []*Line{}, size, &FileNotFoundError{
			Filepath: filepath,
		}
	}

	return lines, size, nil
}

func (c *Commit) NewNote(message string) (plumbing.Hash, error) {
	note := &object.Note{
		Message: message,
	}

	obj := c.repo.Repo.Storer.NewEncodedObject()
	if err := note.Encode(obj); err != nil {
		return plumbing.ZeroHash, err
	}

	return c.repo.Repo.Storer.SetEncodedObject(obj)
}

func Map[T, U any](s []T, fn func(T) U) []U {
	var out = make([]U, len(s))

	for i, v := range s {
		out[i] = fn(v)
	}

	return out
}

func (c *Commit) GetTree(dirpath string, includeCommits bool) (*TreeEntryMap, error) {
	tree, err := c.Commit.Tree()
	if err != nil {
		return nil, &DirectoryNotFoundError{
			Dirpath: dirpath,
		}
	}

	walker := object.NewTreeWalker(tree, true, make(map[plumbing.Hash]bool))

	defer walker.Close()

	treeEntryMap := NewTreeEntryMap()

	var paths = make(map[string]*TreeEntry)

	// Iterate through the tree entries and assign them to a map, keyed by their
	// full paths.
	for {
		name, entry, err := walker.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		file := NewTreeEntry(path.Base(name), name, entry.Hash.String(), entry.Mode.IsFile(), nil)

		if file.IsFile {
			treeEntryMap.AddFile(file)
		} else {
			treeEntryMap.AddDir(file)
		}

		paths[name] = file
	}

	if includeCommits {
		// Iterate through the commit history,
		iter := Must(c.repo.Repo.Log(&git.LogOptions{
			Order: git.LogOrderCommitterTime,
			PathFilter: func(filepath string) bool {
				for path := range paths {
					if strings.HasPrefix(filepath, path) {
						return true
					}
				}

				return false
			},
		}))

		for {
			log.Println("loading next commit...")

			commit := Must(NextNonMergeCommit(iter))
			fileStats := Must(commit.Stats())

			for _, fileStat := range fileStats {
				for path := range paths {
					if strings.HasPrefix(fileStat.Name, path) {
						paths[path].Commit = &Commit{
							Hash:    commit.Hash.String(),
							Message: commit.Message,
							Date:    commit.Author.When.Format(object.DateFormat),
							repo:    c.repo,
							Commit:  commit,
						}

						delete(paths, path)
					}
				}
			}

			if len(paths) == 0 {
				break
			}
		}

		iter.Close()
	}

	// sort.Slice(dirs, func(i, j int) bool {
	// 	return dirs[i].Name < dirs[j].Name
	// })

	// sort.Slice(files, func(i, j int) bool {
	// 	return files[i].Name < files[j].Name
	// })

	err = treeEntryMap.InitHierarchy("")
	return treeEntryMap, err
}
