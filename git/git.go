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
	"log/slog"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
)

type Repo struct {
	repo *git.Repository
	err  error
}

func CloneRepo(store storage.Storer, url string) *Repo {
	slog.Info("cloning repo", "url", url)

	repo, err := git.Clone(store, nil, &git.CloneOptions{
		URL:          url,
		Progress:     io.Discard,
		Mirror:       true,
		NoCheckout:   true,
		SingleBranch: false,
	})

	if err != nil {
		panic(err)
	}

	return &Repo{
		repo: repo,
		err:  err,
	}
}

func OpenRepo(store storage.Storer, url string) *Repo {
	slog.Info("opening repo", "url", url)

	repo, err := git.Open(store, nil)

	return &Repo{
		repo: repo,
		err:  err,
	}
}

type Branch struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
}

func (r *Repo) Error() error {
	return r.err
}

func (r *Repo) GetBranches() ([]Branch, error) {
	if r.err != nil {
		return nil, r.err
	}

	iter, err := r.repo.Branches()

	if err != nil {
		r.err = err
		return nil, err
	}

	var branches []Branch
	iter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, Branch{
			Hash: ref.Hash().String(),
			Name: ref.Name().Short(),
		})
		return nil
	})
	return branches, nil
}

// Branches returns all the References that are branches.
func (r *Repo) Branches() ([]Branch, error) {
	iter, err := r.NewFilteredReferencesIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	})

	if err != nil {
		return nil, err
	}

	return ConsumeIter(iter, func(ref *plumbing.Reference) Branch {
		return Branch{
			Hash: ref.Hash().String(),
			Name: ref.Name().Short()[7:],
		}
	})
}

func (r *Repo) NewFilteredReferencesIter(fn func(ref *plumbing.Reference) bool) (storer.ReferenceIter, error) {
	iter, err := r.repo.Storer.IterReferences()

	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(r *plumbing.Reference) bool {
		return fn(r)
	}, iter), nil
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

func (r *Repo) GetHeadCommit() (*Commit, error) {
	if r.err != nil {
		return nil, r.err
	}

	head, err := r.repo.Head()

	if err != nil {
		r.err = err
		return nil, err
	}

	return r.GetCommit(head.Hash().String())
}

func (r *Repo) GetCommit(revision string) (*Commit, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.GetCommitByHash(r.ResolveRevOrHash(revision))
}

func (r *Repo) ResolveRevOrHash(revOrHash string) *plumbing.Hash {
	if hash, err := r.repo.ResolveRevision(plumbing.Revision("origin/" + revOrHash)); err == nil {
		return hash
	} else {
		hash := plumbing.NewHash(revOrHash)
		return &hash
	}
}

func (r *Repo) GetCommitByHash(hash *plumbing.Hash) (*Commit, error) {
	if r.err != nil {
		return nil, r.err
	}

	commit, err := r.repo.CommitObject(*hash)

	if err != nil {
		return nil, &RevisionNotFoundError{
			Revision: hash.String(),
		}
	}

	return &Commit{
		Hash:    commit.Hash.String(),
		Message: commit.Message,
		repo:    r,
		ptr:     commit,
	}, nil
}

func (r *Repo) GetLastCommit(filename string, from plumbing.Hash) (*Commit, error) {
	cIter := Must(r.repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
		PathFilter: func(path string) bool {
			return strings.HasPrefix(path, filename)
		},
	}))
	defer cIter.Close()

	var commit *object.Commit
	var err error

	for {
		commit, err = cIter.Next()

		if len(commit.ParentHashes) <= 1 {
			break // break out of the loop
		}
	}

	if err != nil {
		return nil, err
	}

	return &Commit{
		Hash:    commit.Hash.String(),
		Message: commit.Message,
		repo:    r,
		ptr:     commit,
	}, nil
}

type Commit struct {
	Hash    string
	Message string
	Date    string
	repo    *Repo
	ptr     *object.Commit
}

type Line = git.Line

// GetFileContents will attempt to read the contents of a file and return each
// line with the most recent blame information.
func (c *Commit) GetFileContents(filepath string, blame bool) ([]*Line, error) {
	var lines []*Line
	var err error

	if blame {
		file, _err := git.Blame(c.ptr, path.Clean(filepath))
		if _err != nil {
			err = _err
		} else {
			lines = file.Lines
		}
	} else {
		file, _err := c.ptr.File(path.Clean(filepath))
		if _err != nil {
			err = _err
		} else {
			_lines, _ := file.Lines()
			lines = Map(_lines, func(_line string) *Line {
				return &Line{
					Text: _line,
				}
			})
		}
	}

	if err != nil {
		return []*Line{}, &FileNotFoundError{
			Filepath: filepath,
		}
	}

	return lines, nil
}

func (c *Commit) NewNote(message string) (plumbing.Hash, error) {
	note := &object.Note{
		Message: message,
	}

	obj := c.repo.repo.Storer.NewEncodedObject()
	if err := note.Encode(obj); err != nil {
		return plumbing.ZeroHash, err
	}

	return c.repo.repo.Storer.SetEncodedObject(obj)
}

func Map[T, U any](s []T, fn func(T) U) []U {
	var out = make([]U, len(s))

	for i, v := range s {
		out[i] = fn(v)
	}

	return out
}

func (c *Commit) GetTree(dirpath string, includeCommits bool) (*TreeEntryMap, error) {
	tree, err := c.ptr.Tree()
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

		file := NewTreeEntry(path.Base(name), path.Join("blob", c.ptr.Hash.String(), name), entry.Hash.String(), entry.Mode.IsFile(), nil)

		if file.IsFile {
			treeEntryMap.AddFile(file)
		} else {
			treeEntryMap.AddDir(file)
		}

		paths[name] = file
	}

	if includeCommits {
		// Iterate through the commit history,
		iter := Must(c.repo.repo.Log(&git.LogOptions{
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
							ptr:     commit,
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

	err = treeEntryMap.InitHierarchy(path.Join("blob", c.ptr.Hash.String()))
	return treeEntryMap, err
}

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
