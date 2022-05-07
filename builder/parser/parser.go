package parser

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

// Parser defines methods for data parsing to identify links referenced and
// provide a template for building.
type Parser interface {
	// GetLinkableData returns a template and a map with template
	// variable names mapped to the original content
	GetLinkableData([]byte) (template.Template, map[string]interface{}, error)
	// AddFuncs adds functions used evaluate
	// whether a value is an in-content link.
	// If no functions are added all values will be considered
	// links.
	AddFuncs(...TemplatingFunc)
}

// TemplatingFunc determine the condition
// that must be met for data to be templated
type TemplatingFunc func(interface{}) bool

// ErrInvalidFormat defines an error for unsupported format types
type ErrInvalidFormat struct {
	filename string
}

func (e *ErrInvalidFormat) Error() string {
	return fmt.Sprintf("format unsupported for filename: %s", e.filename)
}

// ByExtension returns a parser based on the extension of the filename.
func ByExtension(filename string) (Parser, error) {
	switch filepath.Ext(filename) {
	case ".json":
		return NewJSONParser(filename), nil
	}
	return nil, &ErrInvalidFormat{filename}
}

// ConvertFilenameForGoTemplateValue converts the current
// file string to a value that is an acceptable variable for Go templating
func ConvertFilenameForGoTemplateValue(filename string) string {
	filename = strings.Replace(filename, ".", "_", -1)
	filename = strings.Replace(filename, string(filepath.Separator), "_", -1)
	return filename
}
