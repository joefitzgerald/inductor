package filesys

import (
	"fmt"
	"path/filepath"
	"strings"
)

// TemplateKey returns the key used to group templates by
// path should end with .tpl
func TemplateKey(path string) string {
	if filepath.Ext(path) != ".tpl" {
		return ""
	}
	return strings.TrimSuffix(strings.ToLower(path), ".tpl")
}

// TemplateKeyForPartial returns the key associated with the partial template
// path should end with .ptpl
func TemplateKeyForPartial(templateKeys []string, path string) string {
	if filepath.Ext(path) != ".ptpl" {
		return ""
	}

	// Find the associated root template key, if any
	key := strings.TrimSuffix(strings.ToLower(path), ".ptpl")
	for !contains(templateKeys, key) {
		ext := filepath.Ext(key)
		if ext == "" {
			return ""
		}
		key = strings.TrimSuffix(key, ext)
	}
	return key
}

// ListTemplates returns all templates and partial templates grouped by
// the associated full template full path for the specified OS
func ListTemplates(baseDir string, osName string) map[string]map[string]string {
	templates := make(map[string]map[string]string)

	// find all root templates
	entries := listTemplatesFn(baseDir)
	rootTemplateKeys := []string{}
	for _, e := range entries {
		key := TemplateKey(e)
		rootTemplateKeys = append(rootTemplateKeys, key)
		templates[key] = make(map[string]string)
		templates[key][e] = e
	}

	// find all partial templates and associate them with a root template
	entries = listPartialTemplatesFn(baseDir)
	for _, e := range entries {
		pathWithoutOS := removeOsDir(e, osName)
		key := TemplateKeyForPartial(rootTemplateKeys, pathWithoutOS)
		if key == "" {
			continue
		}

		// OS specific entries are longer, ensure we always keep those
		if existingEntry, ok := templates[key][pathWithoutOS]; ok {
			if len(e) > len(existingEntry) {
				templates[key][pathWithoutOS] = e
			}
		} else {
			templates[key][pathWithoutOS] = e
		}
	}
	return templates
}

// for testing
var listTemplatesFn = listTemplates
var listPartialTemplatesFn = listPartialTemplates

func listTemplates(baseDir string) []string {
	return listFiles(fmt.Sprintf("%s/**/*.tpl", baseDir))
}

func listPartialTemplates(baseDir string) []string {
	return listFiles(fmt.Sprintf("%s/**/*.ptpl", baseDir))
}

func listFiles(globPattern string) []string {
	templates, err := filepath.Glob(globPattern)
	// shouldn't happen unles we have a bad pattern
	if err != nil {
		panic(err)
	}
	return templates
}

func removeOsDir(path string, osName string) string {
	dir, file := filepath.Split(path)
	dir = strings.TrimSuffix(strings.TrimSuffix(dir, "/"), osName)
	return filepath.Join(dir, file)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
