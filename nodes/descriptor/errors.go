package descriptor

import "fmt"

// ParseError defines an error when parsing attributes into Properties.
type ParseError struct {
	Key string
	Err error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parse property key %q: %v", e.Key, e.Err)
}
