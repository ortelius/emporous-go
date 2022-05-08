package parser

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

var _ Parser = &uorParser{}

type uorParser struct {
	filename      string
	templateFuncs []TemplatingFunc
}

// NewUORParser returns the UOR template implementation
// for the Parser.
func NewUORParser(filename string) Parser {
	return &uorParser{
		filename: filename,
	}
}

func (p *uorParser) GetLinkableData(data []byte) (template.Template, map[string]interface{}, error) {
	pattern := `\_\_uor\.(.*?)\_\_`
	templateSearch, _ := regexp.Compile(pattern)
	links := make(map[string]interface{})
	if templateSearch.Match(data) {
		found := templateSearch.FindAllSubmatch(data, -1)
		for _, t := range found {
			filename := strings.Trim(string(t[0]), "__uor.")
			filename = strings.Trim(filename, "__")
			formattedFilename := ConvertFilenameForGoTemplateValue(string(t[0]))
			// Set the template values to its original value
			// for now
			templateValue := unstructuredFileNameToGoTemplateValue(formattedFilename)
			subst := regexp.MustCompile(string(t[0]))
			data = subst.ReplaceAll(data, []byte(templateValue))
			links[formattedFilename] = filename
		}
	}

	t, err := template.New(p.filename).Parse(string(data))
	if err != nil {
		return template.Template{}, links, err
	}
	return *t, links, nil
}

func (p *uorParser) AddFuncs(tFuncs ...TemplatingFunc) {
	p.templateFuncs = append(p.templateFuncs, tFuncs...)
}

func unstructuredFileNameToGoTemplateValue(convertedFilename string) string {
	return fmt.Sprintf("{{.%s}}", convertedFilename)
}
