package attributes

import (
	"errors"
	"fmt"
)

// ParseError defines an error when parsing attributes into Properties.
type ParseError struct {
	Key string
	Err error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parse property key %q: %v", e.Key, e.Err)
}

// ErrInvalidAttribute defines the error thrown when an attribute has an invalid
// type.
var ErrInvalidAttribute = errors.New("invalid attribute type")
