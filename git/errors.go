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
