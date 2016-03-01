package tpl

import (
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Template represents a path to a template file
type Template struct {
	Path string
}

// Content writes the partial template to the specified writer
func (t *Template) Content(buffer io.Writer) error {
	buffer.Write([]byte("\n"))
	tpl, err := ioutil.ReadFile(t.Path)
	if err != nil {
		return err
	}
	buffer.Write(tpl)
	return nil
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

// ByPath is a type for sorting Templates
type ByPath []Template

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
	return a[i].Path < a[j].Path
}
