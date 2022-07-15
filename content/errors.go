package content

import "fmt"

type ErrNotStored struct {
	Reference string
}

func (e *ErrNotStored) Error() string {
	return fmt.Sprintf("descriptor for reference %s is not stored", e.Reference)
}
