package tpl

import "testing"

func TestListPartialTemplates(t *testing.T) {
	var expectedFiles = []string{
		"/Users/sneal/packer-windows/autounattend.xml.disks.ptpl",
	}
	listPartialTemplatesFn = func(baseDir, baseFilename string) []string {
		return []string{
			"/Users/sneal/packer-windows/autounattend.xml.disks.ptpl",
		}
	}
	listPartialTemplatesOSSpecificFn = func(baseDir, baseFilename, osName string) []string {
		return []string{}
	}

	rootTemplate := &RootTemplate{}
	partialTemplates := rootTemplate.PartialTemplates("windows10")

	if len(expectedFiles) != len(partialTemplates) {
		t.Errorf("Expected %d files, but got %d", len(expectedFiles), len(partialTemplates))
	}
	for _, expectedFile := range expectedFiles {
		if !contains(partialTemplates, expectedFile) {
			t.Errorf("Expected partial templates to contain '%s'", expectedFile)
		}
	}
}

func TestListPartialTemplatesWithOSSpecificOverride(t *testing.T) {
	var expectedFiles = []string{
		"/Users/sneal/packer-windows/autounattend.xml.oobe.ptpl",
		"/Users/sneal/packer-windows/windows10/autounattend.xml.disks.ptpl",
	}
	listPartialTemplatesFn = func(baseDir, baseFilename string) []string {
		return []string{
			"/Users/sneal/packer-windows/autounattend.xml.oobe.ptpl",
			"/Users/sneal/packer-windows/autounattend.xml.disks.ptpl",
		}
	}
	listPartialTemplatesOSSpecificFn = func(baseDir, baseFilename, osName string) []string {
		return []string{
			"/Users/sneal/packer-windows/windows10/autounattend.xml.disks.ptpl",
		}
	}

	rootTemplate := &RootTemplate{}
	partialTemplates := rootTemplate.PartialTemplates("windows10")

	if len(expectedFiles) != len(partialTemplates) {
		t.Errorf("Expected %d files, but got %d", len(expectedFiles), len(partialTemplates))
	}
	for _, expectedFile := range expectedFiles {
		if !contains(partialTemplates, expectedFile) {
			t.Errorf("Expected partial templates to contain '%s'", expectedFile)
		}
	}
}

func contains(pts []PartialTemplate, e string) bool {
	for _, pt := range pts {
		if pt.Path == e {
			return true
		}
	}
	return false
}
