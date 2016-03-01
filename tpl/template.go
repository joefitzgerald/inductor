package tpl

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Template represents a path to a template file
type Template struct {
	Path string
}

// RootTemplate is a root template (.tpl) and all its associated partial templates
type RootTemplate struct {
	Template
	PartialTemplate []Template
}

// PartialTemplate is a partial template fragment (.ptpl)
type PartialTemplate struct {
	Template
}

// PartialTemplates returns all partial templates associated with this template
func (t *RootTemplate) PartialTemplates(osName string) []PartialTemplate {
	partials := make(map[string]*PartialTemplate)

	// get all shared partial templates
	partialFiles := listPartialTemplatesFn(t.Dir(), t.BaseFilename())
	for _, f := range partialFiles {
		pt := &PartialTemplate{Template{f}}
		partials[pt.Filename()] = pt
	}

	// get all OS specific partial templates, overwriting any non-specific templates
	partialFilesOS := listPartialTemplatesOSSpecificFn(t.Dir(), t.BaseFilename(), osName)
	for _, f := range partialFilesOS {
		pt := &PartialTemplate{Template{f}}
		partials[pt.Filename()] = pt
	}

	// flatten the map of partials
	distinctPartials := []PartialTemplate{}
	for _, p := range partials {
		distinctPartials = append(distinctPartials, *p)
	}
	return distinctPartials
}

// BaseFilename is the name of the file minus the file extension
func (t *Template) BaseFilename() string {
	ext := filepath.Ext(t.Path)
	return strings.TrimSuffix(t.Filename(), ext)
}

// Dir returns the template directory (no trailing slash)
func (t *Template) Dir() string {
	dir, _ := filepath.Split(t.Path)
	dir = strings.TrimSuffix(dir, "/")
	return dir
}

// Filename is the filename with extension of the file, no dir
func (t *Template) Filename() string {
	_, file := filepath.Split(t.Path)
	return file
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
	return templates
}
