package filesys

import (
	"log"
	"testing"
)

func TestCreateTemplateKey(t *testing.T) {
	path := "/Users/sneal/packer-windows/Autounattend.xml.tpl"
	key := TemplateKey(path)
	if key != "/users/sneal/packer-windows/autounattend.xml" {
		t.Errorf("unexpected template key '%s'", key)
	}
}

func TestCreateTemplateKeyWithWrongFileExt(t *testing.T) {
	path := "/Users/sneal/packer-windows/Autounattend.xml.foo"
	key := TemplateKey(path)
	if key != "" {
		t.Errorf("Expected no template key, but got: '%s'", key)
	}
}

func TestCreateTemplateKeyFromPartialTemplate(t *testing.T) {
	templates := []string{
		"/users/sneal/packer-windows/autounattend.xml",
		"/users/sneal/packer-windows/foo/autounattend.xml",
		"/users/sneal/packer-windows/packer.json",
	}
	path := "/Users/sneal/packer-windows/Autounattend.xml.foo.baz.ptpl"
	key := TemplateKeyForPartial(templates, path)
	if key != "/users/sneal/packer-windows/autounattend.xml" {
		t.Errorf("unexpected template key '%s'", key)
	}
}

func TestCreateTemplateKeyFromPartialTemplateWithWrongFileExt(t *testing.T) {
	templates := []string{}
	path := "/Users/sneal/packer-windows/Autounattend.xml.foo.baz"
	key := TemplateKeyForPartial(templates, path)
	if key != "" {
		t.Errorf("Expected no template key, but got: '%s'", key)
	}
}

func TestCreateTemplateKeyFromPartialTemplateWithNoAssociatedTemplate(t *testing.T) {
	templates := []string{}
	path := "/Users/sneal/packer-windows/Autounattend.xml.foo.baz.ptpl"
	key := TemplateKeyForPartial(templates, path)
	if key != "" {
		t.Errorf("Expected no template key, but got: '%s'", key)
	}
}

func TestRemoveOsDirFromPath(t *testing.T) {
	path := "/Users/sneal/packer-windows/windows10/Autounattend.xml.foo.baz.ptpl"
	expectedPath := "/Users/sneal/packer-windows/Autounattend.xml.foo.baz.ptpl"
	pathWithoutOs := removeOsDir(path, "windows10")
	if pathWithoutOs != expectedPath {
		t.Errorf("Expected '%s', but got '%s'", expectedPath, pathWithoutOs)
	}
}

func TestListTemplates(t *testing.T) {
	var expectedFiles = []string{
		"/Users/sneal/packer-windows/Autounattend.xml.tpl",
		"/Users/sneal/packer-windows/autounattend.xml.disks.ptpl",
	}
	listTemplatesFn = func(baseDir string) []string {
		return []string{
			"/Users/sneal/packer-windows/Autounattend.xml.tpl",
			"/Users/sneal/packer-windows/packer.json.tpl",
		}
	}
	listPartialTemplatesFn = func(baseDir string) []string {
		return []string{"/Users/sneal/packer-windows/autounattend.xml.disks.ptpl"}
	}
	templates := ListTemplates("/Users/sneal/packer-windows", "windows10")
	t.Log(templates)
	files, ok := templates["/users/sneal/packer-windows/autounattend.xml"]
	if !ok {
		t.Errorf("Templates didn't contain the autounattend.xml key")
	}
	if len(expectedFiles) != len(files) {
		t.Errorf("Expected %d files, but got %d", len(expectedFiles), len(files))
	}
	for _, expectedFile := range expectedFiles {
		if files[removeOsDir(expectedFile, "ignored")] != expectedFile {
			t.Errorf("Expected it to contain: %s", expectedFile)
		}
	}
}

func TestListTemplatesWithoutRootTemplate(t *testing.T) {
	listTemplatesFn = func(baseDir string) []string {
		return []string{}
	}
	listPartialTemplatesFn = func(baseDir string) []string {
		return []string{"/Users/sneal/packer-windows/autounattend.xml.disks.ptpl"}
	}
	templates := ListTemplates("/Users/sneal/packer-windows", "windows10")
	if len(templates) != 0 {
		t.Error("Expected no root templates")
	}
}

func TestListTemplatesWithOsOverride(t *testing.T) {
	var expectedFiles = []string{
		"/Users/sneal/packer-windows/Autounattend.xml.tpl",
		"/Users/sneal/packer-windows/autounattend.xml.oobe.ptpl",
		"/Users/sneal/packer-windows/windows10/autounattend.xml.disks.ptpl",
	}
	listTemplatesFn = func(baseDir string) []string {
		return []string{
			"/Users/sneal/packer-windows/Autounattend.xml.tpl",
		}
	}
	listPartialTemplatesFn = func(baseDir string) []string {
		return []string{
			"/Users/sneal/packer-windows/autounattend.xml.oobe.ptpl",
			"/Users/sneal/packer-windows/windows10/autounattend.xml.disks.ptpl",
			"/Users/sneal/packer-windows/autounattend.xml.disks.ptpl",
		}
	}
	templates := ListTemplates("/Users/sneal/packer-windows", "windows10")
	files, ok := templates["/users/sneal/packer-windows/autounattend.xml"]
	if !ok {
		t.Errorf("Templates didn't contain the autounattend.xml key")
	}
	log.Println(templates)
	if len(expectedFiles) != len(files) {
		t.Errorf("Expected %d files, but got %d", len(expectedFiles), len(files))
	}
	for _, expectedFile := range expectedFiles {
		if files[removeOsDir(expectedFile, "windows10")] != expectedFile {
			t.Errorf("Expected it to contain: %s", expectedFile)
		}
	}
}
