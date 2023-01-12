package errlist

import (
	"errors"
	"strings"
)

// Separator is used to separate error messages when calling Error on a list.
const Separator = ", "

// ErrList defines methods for lists of type error.
type ErrList interface {
	error
	Errors() []error
	Is(error) bool
}

// NewErrList returns a new ErrList based on a given list of
func NewErrList(list []error) ErrList {
	if len(list) == 0 {
		return nil
	}
	// In case of input error list contains nil
	var errs []error
	for _, e := range list {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errList(errs)
}

type errList []error

// visitErrFunc is a function that defines input and
// outputs when an error is visited.
type visitErrFunc func(error) bool

func (el errList) Error() string {
	if len(el) == 0 {
		return ""
	}
	if len(el) == 1 {
		return el[0].Error()
	}
	seenErrs := map[string]struct{}{}
	result := strings.Builder{}
	el.visitErr(func(err error) bool {
		msg := err.Error()
		if _, has := seenErrs[msg]; has {
			return false
		}
		seenErrs[msg] = struct{}{}
		if len(seenErrs) > 1 {
			result.WriteString(Separator)
		}
		result.WriteString(msg)
		return false
	})
	if len(seenErrs) == 1 {
		return result.String()
	}
	return "[" + result.String() + "]"
}

func (el errList) Is(target error) bool {
	return el.visitErr(func(err error) bool {
		return errors.Is(err, target)
	})
}

// visitsErr visits each error and stops if a match is found.
func (el errList) visitErr(f visitErrFunc) bool {
	for _, err := range el {
		switch err := err.(type) {
		case errList:
			if match := err.visitErr(f); match {
				return match
			}
		case ErrList:
			for _, nestedErr := range err.Errors() {
				if match := f(nestedErr); match {
					return match
				}
			}
		default:
			if match := f(err); match {
				return match
			}
		}
	}

	return false
}

func (el errList) Errors() []error {
	return el
}
