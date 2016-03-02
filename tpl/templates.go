package tpl

import (
	"os"
	"path/filepath"
)

// Templates encompasses all templates within a base dir
type Templates struct {
	BaseDir string
	OSName  string
	All     []Template
}

// New creates a new Templates instance
func New(baseDir, osName string) *Templates {
	templates := &Templates{
		BaseDir: baseDir,
		OSName:  osName,
	}

	// find all root templates
	entries := listRootTemplatesFn(baseDir)
	for _, e := range entries {
		rootTemplate := NewRootTemplate(e, osName)
		templates.All = append(templates.All, *rootTemplate)
	}

	return templates
}

// FindTemplate finds the root template by path if it exists
// TODO: Should this be recursive?
func (t *Templates) FindTemplate(path string) *Template {
	for _, pt := range t.All {
		if pt.Path == path {
			return &pt
		}
	}
	return nil
}

// for testing
var listRootTemplatesFn = listRootTemplates

func listRootTemplates(baseDir string) []string {
	rootTemplates := []string{}
	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".tpl" {
			return nil
		}
		rootTemplates = append(rootTemplates, path)
		return nil
	})
	return rootTemplates
}
