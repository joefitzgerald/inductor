package tpl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// RootTemplate is a root template (.tpl) and all its associated partial templates
type RootTemplate struct {
	Template
	PartialTemplates []Template
}

// NewRootTemplate create a new RootTemplate instance complete with partial templates
func NewRootTemplate(path, osName string) *RootTemplate {
	rootTemplate := &RootTemplate{}
	rootTemplate.Path = path
	partials := make(map[string]*Template)

	// get all shared partial templates
	partialFiles := listPartialTemplatesFn(rootTemplate.Dir(), rootTemplate.BaseFilename())
	for _, f := range partialFiles {
		pt := &Template{Path: f}
		partials[pt.Filename()] = pt
	}

	// get all OS specific partial templates, overwriting any non-specific templates
	partialFilesOS := listPartialTemplatesOSSpecificFn(rootTemplate.Dir(), rootTemplate.BaseFilename(), osName)
	for _, f := range partialFilesOS {
		pt := &Template{Path: f}
		partials[pt.Filename()] = pt
	}

	// flatten the map of partials
	distinctPartials := []Template{}
	for _, p := range partials {
		distinctPartials = append(distinctPartials, *p)
	}

	// ensure a stable sort order so the output content is diffable
	sort.Sort(ByPath(distinctPartials))
	rootTemplate.PartialTemplates = distinctPartials

	return rootTemplate
}

// Content of this template and all of its partial templates
func (t *RootTemplate) Content() (string, error) {
	var buffer bytes.Buffer

	// write root template
	tpl, err := ioutil.ReadFile(t.Path)
	if err != nil {
		return "", err
	}
	buffer.Write(tpl)

	// write each partial template
	for _, pt := range t.PartialTemplates {
		// wrap the partial template define statement
		defineName := strings.TrimPrefix(pt.BaseFilename(), t.BaseFilename())
		defineName = strings.Replace(defineName, ".", "", -1)
		buffer.WriteString(fmt.Sprintf("\n{{define \"%s\"}}", defineName))

		err := pt.Content(&buffer)
		if err != nil {
			return "", err
		}

		buffer.WriteString("\n{{end}}")
	}

	return buffer.String(), nil
}

// FindPartialTemplate finds the root template by path if it exists
func (t *RootTemplate) FindPartialTemplate(path string) *Template {
	for _, pt := range t.PartialTemplates {
		if pt.Path == path {
			return &pt
		}
	}
	return nil
}

// for testing
var listPartialTemplatesFn = listPartialTemplates
var listPartialTemplatesOSSpecificFn = listPartialTemplatesOSSpecific

func listPartialTemplates(baseDir, baseFilename string) []string {
	return listFiles(fmt.Sprintf("%s/%s*.ptpl", baseDir, baseFilename))
}

func listPartialTemplatesOSSpecific(baseDir, baseFilename, osName string) []string {
	return listFiles(fmt.Sprintf("%s/%s/%s*.ptpl", baseDir, osName, baseFilename))
}

func listFiles(globPattern string) []string {
	templates, err := filepath.Glob(globPattern)
	// shouldn't happen unles we have a bad pattern
	if err != nil {
		panic(err)
	}
	sort.Strings(templates)
	return templates
}
