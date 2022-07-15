package content

import "fmt"

// ErrNotStored denotes that a reference is not stored on the content store.
type ErrNotStored struct {
	Reference string
}

func (e *ErrNotStored) Error() string {
	return fmt.Sprintf("descriptor for reference %s is not stored", e.Reference)
}
