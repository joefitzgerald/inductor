package tpl

import (
	"os"
	"path/filepath"
)

// Templates encompasses all templates within a base dir
type templates struct {
	baseDir string
	osName  string
	all     []Templater
}

// New creates a new Templates instance
func New(baseDir, osName string) TemplateContainer {
	templates := &templates{
		baseDir: baseDir,
		osName:  osName,
	}

	// find all root templates
	entries := listRootTemplatesFn(baseDir)
	for _, e := range entries {
		rootTemplate := NewRootTemplate(e, osName)
		templates.all = append(templates.all, rootTemplate)
	}

	return templates
}

// ListTemplates returns all root templates
func (t *templates) ListTemplates() []Templater {
	return t.all
}

// FindTemplate finds the root template by path if it exists
func (t *templates) FindTemplate(path string) Templater {
	for _, pt := range t.all {
		if pt.FullPath() == path {
			return pt
		}
	}
	return nil
}

// for testing
var listRootTemplatesFn = listRootTemplates

func listRootTemplates(baseDir string) []string {
	rootTemplates := []string{}
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".tpl" {
			return nil
		}
		rootTemplates = append(rootTemplates, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return rootTemplates
}
