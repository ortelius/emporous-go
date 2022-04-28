package parser

import "text/template"

// Parser will parse data to identify links referenced and
// provide a template for building.
type Parser interface {
	// GetLinkableData returns a template and a map of file
	// links and the associated variable in the template.
	GetLinkableData(data []byte) (template.Template, map[string]interface{})
}
