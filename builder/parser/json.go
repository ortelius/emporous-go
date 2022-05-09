package parser

import (
	"fmt"
	"text/template"

	"github.com/buger/jsonparser"
)

var _ Parser = &jsonParser{}

type jsonParser struct {
	filename      string
	templateFuncs []TemplatingFunc
}

// NewJSONParser returns the JSON implementation
// for the Parser interface.
func NewJSONParser(filename string) Parser {
	return &jsonParser{
		filename: filename,
	}
}

// GetLinkableData returns a template and a map with template variable names mapped to the original content.
func (p *jsonParser) GetLinkableData(data []byte) (template.Template, map[string]interface{}, error) {
	links := make(map[string]interface{})
	handler := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) (err error) {
		// Generically determine when a value should be substituted.
		if p.evaluateTFuncs(string(value)) {
			filename := ConvertFilenameForGoTemplateValue(string(value))
			templateValue := fileNameToGoTemplateValue(filename)
			// The goal is to set the variable in the links
			// map to the original values to this information
			// can still be accessed.
			links[filename] = string(value)
			data, err = jsonparser.Set(data, []byte(templateValue), string(key))
			if err != nil {
				return err
			}
		}
		return nil
	}
	if err := jsonparser.ObjectEach(data, handler); err != nil {
		return template.Template{}, links, err
	}
	t, err := template.New(p.filename).Parse(string(data))
	if err != nil {
		return template.Template{}, links, err
	}
	return *t, links, nil
}

// AddFuncs adds functions used evaluate whether a value is an in-content link.
func (p *jsonParser) AddFuncs(tFuncs ...TemplatingFunc) {
	p.templateFuncs = append(p.templateFuncs, tFuncs...)
}

// If no functions are added no values will be considered
// links.
func (p *jsonParser) evaluateTFuncs(value string) bool {
	for _, f := range p.templateFuncs {
		if f(value) {
			return true
		}
	}
	return false
}

func fileNameToGoTemplateValue(convertedFilename string) string {
	return fmt.Sprintf("\"{{.%s}}\"", convertedFilename)
}
