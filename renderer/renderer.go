package renderer

import "github.com/joefitzgerald/inductor/tpl"

// Renderer will render the given set of templates to disk
type Renderer interface {
	Render(templates tpl.TemplateContainer) error
}
