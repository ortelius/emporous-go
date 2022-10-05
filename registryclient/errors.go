package registryclient

import "fmt"

// ErrNoAvailableMirrors denotes that all registry mirrors are not accessible.
type ErrNoAvailableMirrors struct {
	Registry string
}

func (e *ErrNoAvailableMirrors) Error() string {
	return fmt.Sprintf("registry %q: no avaialble mirrors", e.Registry)
}
