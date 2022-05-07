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
	pattern := `\[\[\.uor\..*?\]\]`
	templateSearch, _ := regexp.Compile(pattern)
	links := make(map[string]interface{})
	if templateSearch.Match(data) {
		found := templateSearch.FindAllSubmatch(data, -1)
		fmt.Printf("found  = %q\n", found)
		fmt.Println("\nMatched")
		for i, t := range found {
			fmt.Printf("%d, %q\n", i, t)
			filename := strings.Trim(string(t[0]), "[[.uor$")
			filename = strings.Trim(filename, "]]")
			formattedFilename := ConvertFilenameForGoTemplateValue(filename)
			//templateValue := fileNameToGoTemplateValue(filename)
			// Set the template values to its original value
			// for now
			links[formattedFilename] = filename
			fmt.Printf("links: %s\n", links)
		}
	}

	t, err := template.New(p.filename).Parse(string(data))
	fmt.Printf("t = %q\n", t)
	if err != nil {
		return template.Template{}, links, err
	}
	return *t, links, nil
}

func (p *uorParser) AddFuncs(tFuncs ...TemplatingFunc) {
	p.templateFuncs = append(p.templateFuncs, tFuncs...)
}

func (p *uorParser) evaluateTFuncs(value string) bool {
	for _, f := range p.templateFuncs {
		if f(value) {
			return true
		}
	}
	return false
}
