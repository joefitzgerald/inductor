package tpl

// TemplateContainer contains all templates and partial for a specific OS
type TemplateContainer interface {
	FindTemplate(path string) Templater
	ListTemplates() []Templater
}
