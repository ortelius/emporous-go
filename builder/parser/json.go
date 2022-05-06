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
// for the Parser.
func NewJSONParser(filename string) Parser {
	return &jsonParser{
		filename: filename,
	}
}

func (p *jsonParser) GetLinkableData(data []byte) (template.Template, map[string]interface{}, error) {
	links := make(map[string]interface{})
	handler := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) (err error) {
		// Generically determine when a value should be substituted
		if p.evaluateTFuncs(string(value)) {
			filename := ConvertFilenameForGoTemplateValue(string(value))
			templateValue := fileNameToGoTemplateValue(filename)
			// Set the template values to its original value
			// for now
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

func (p *jsonParser) AddFuncs(tFuncs ...TemplatingFunc) {
	p.templateFuncs = append(p.templateFuncs, tFuncs...)
}

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
