package tpl

import "io"

// Templater can render a specific template and its partials
type Templater interface {
	FullPath() string
	BaseFilename() string
	Content(buffer io.Writer) error
	FindTemplate(path string) Templater
	ListTemplates() []Templater
}
