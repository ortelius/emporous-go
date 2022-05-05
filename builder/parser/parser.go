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
	GetLinkableData(data []byte, fileIndex map[string]struct{}) (template.Template, map[string]interface{}, error)
}

// ErrInvalidFormat defines an error for unsupported format types
type ErrInvalidFormat struct {
	filename string
}

func (e *ErrInvalidFormat) Error() string {
	return fmt.Sprintf("format unrecognized by filename: %s", e.filename)
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
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	ext = strings.TrimPrefix(ext, ".")
	return fmt.Sprintf("%s_%s", base, ext)
}
