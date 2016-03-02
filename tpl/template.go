package tpl

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// Template represents a path to a template file
type Template struct {
	Path             string
	PartialTemplates []Templater
}

// NewRootTemplate create a new RootTemplate instance complete with partial templates
func NewRootTemplate(path, osName string) Templater {
	rootTemplate := &Template{}
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
	distinctPartials := []Templater{}
	for _, p := range partials {
		distinctPartials = append(distinctPartials, p)
	}

	// ensure a stable sort order so the output content is diffable
	sort.Sort(ByPath(distinctPartials))
	rootTemplate.PartialTemplates = distinctPartials

	return rootTemplate
}

// FullPath returns the full path to the template file
func (t *Template) FullPath() string {
	return t.Path
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

// Content of this template and all of its partial templates
func (t *Template) Content(buffer io.Writer) error {
	// write root template
	tpl, err := ioutil.ReadFile(t.Path)
	if err != nil {
		return err
	}
	buffer.Write(tpl)

	// write each partial template
	for _, pt := range t.PartialTemplates {
		// wrap the partial template define statement
		defineName := strings.TrimPrefix(pt.BaseFilename(), t.BaseFilename())
		defineName = strings.Replace(defineName, ".", "", -1)
		buffer.Write([]byte(fmt.Sprintf("\n{{define \"%s\"}}\n", defineName)))

		err := pt.Content(buffer)
		if err != nil {
			return err
		}

		buffer.Write([]byte("\n{{end}}"))
	}

	return nil
}

// FindTemplate finds the root template by path if it exists
func (t *Template) FindTemplate(path string) Templater {
	for _, pt := range t.PartialTemplates {
		if pt.FullPath() == path {
			return pt
		}
	}
	return nil
}

// ListTemplates lists all the associated partial templates
func (t *Template) ListTemplates() []Templater {
	return t.PartialTemplates
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

// ByPath is a type for sorting Templates
type ByPath []Templater

// Len is the count of elements
func (a ByPath) Len() int {
	return len(a)
}

// Swap the two values in the slice
func (a ByPath) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less compares two elements and returns true if i is less than j
func (a ByPath) Less(i, j int) bool {
	return a[i].FullPath() < a[j].FullPath()
}
