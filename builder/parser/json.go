package parser

import (
	"fmt"
	"text/template"

	"github.com/buger/jsonparser"
)

var _ Parser = &jsonParser{}

type jsonParser struct {
	filename string
}

// NewJSONParser returns the JSON implementation
// for the Parser.
func NewJSONParser(filename string) Parser {
	return &jsonParser{
		filename: filename,
	}
}

func (p *jsonParser) GetLinkableData(data []byte, fileIndex map[string]struct{}) (template.Template, map[string]interface{}, error) {
	links := make(map[string]interface{})
	handler := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) (err error) {
		// Currently determining whether to replace the value based on
		// if it is found in the current workspace.
		// QUESTION(jpower432): Would it be better to have in content
		// annotation to determine what should be substituted.
		if _, found := fileIndex[string(value)]; found {
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
	return *t, links, err
}

func fileNameToGoTemplateValue(convertedFilename string) string {
	return fmt.Sprintf("\"{{.%s}}\"", convertedFilename)
}
