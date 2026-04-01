package git

import "fmt"

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
