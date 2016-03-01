package tpl

import (
	"os"
	"path/filepath"
)

// Templates encompasses all templates within a base dir
type Templates struct {
	BaseDir string
	All     []RootTemplate
}

// New creates a new Templates instance
func New(baseDir string) *Templates {
	templates := &Templates{
		BaseDir: baseDir,
	}

	// find all root templates
	entries := listRootTemplatesFn(baseDir)
	for _, e := range entries {
		rootTemplate := RootTemplate{}
		rootTemplate.Path = e
		templates.All = append(templates.All, rootTemplate)
	}

	return templates
}

// FindTemplate finds the root template by path if it exists
func (t *Templates) FindTemplate(path string) *RootTemplate {
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
