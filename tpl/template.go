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

// PartialTemplate is a partial template fragment (.ptpl)
type PartialTemplate struct {
	Template
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
