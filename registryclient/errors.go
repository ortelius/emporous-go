package registryclient

import "errors"

// ErrNoMatch is an error that is thrown when all files in a collection are
// filtered out by attribute filtering.
var ErrNoMatch = errors.New("no matching files")
