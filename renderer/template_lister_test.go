package renderer

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFilterFilesToRenderWithoutCollision(t *testing.T) {
	var files = []string{
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.tpl",
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.windowsPE.tpl",
		"/Users/sneal/packer-windows/Autounattend.offlineServicing.tpl",
		"/Users/sneal/packer-windows/Autounattend.oobeSystem.tpl",
		"/Users/sneal/packer-windows/Autounattend.specialize.tpl",
	}
	filteredFiles := FilterFilesToRender(files, "windows2012r2")
	for _, expectedFile := range files {
		if !contains(filteredFiles, expectedFile) {
			t.Errorf("Expected the slice to contain: %s", expectedFile)
		}
	}
}

func TestFilterFilesToRenderWithCollision(t *testing.T) {
	var files = []string{
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.tpl",
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.windowsPE.tpl",
		"/Users/sneal/packer-windows/Autounattend.windowsPE.tpl",
		"/Users/sneal/packer-windows/Autounattend.offlineServicing.tpl",
	}
	var expected = []string{
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.tpl",
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.windowsPE.tpl",
		"/Users/sneal/packer-windows/Autounattend.offlineServicing.tpl",
	}
	filteredFiles := FilterFilesToRender(files, "windows2012r2")
	for _, expectedFile := range expected {
		if !contains(filteredFiles, expectedFile) {
			t.Errorf("Expected the slice to contain: %s", expectedFile)
		}
	}
}

func TestFilterFilesToRenderWithMultipleOS(t *testing.T) {
	var files = []string{
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.tpl",
		"/Users/sneal/packer-windows/Autounattend-windows2012.tpl",
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.windowsPE.tpl",
		"/Users/sneal/packer-windows/Autounattend-windows2012.windowsPE.tpl",
		"/Users/sneal/packer-windows/Autounattend.windowsPE.tpl",
		"/Users/sneal/packer-windows/Autounattend.offlineServicing.tpl",
	}
	var expected = []string{
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.tpl",
		"/Users/sneal/packer-windows/Autounattend-windows2012r2.windowsPE.tpl",
		"/Users/sneal/packer-windows/Autounattend.offlineServicing.tpl",
	}
	filteredFiles := FilterFilesToRender(files, "windows2012r2")
	t.Log(filteredFiles)
	if len(filteredFiles) != len(expected) {
		t.Errorf("Expected %d files to be returned, but got %d", len(expected), len(filteredFiles))
	}
	for _, expectedFile := range expected {
		if !contains(filteredFiles, expectedFile) {
			t.Errorf("Expected the slice to contain: %s", expectedFile)
		}
	}
}

func TestCanListFiles(t *testing.T) {
	// create temporary directory to store all our testNewPackerTemplate() files
	tmpDir, err := ioutil.TempDir("", "inductor")
	if err != nil {
		t.Error("Couldn't create test temp dir: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// stub out the autounattend files
	expected := make([]string, 0, 4)
	expected = append(expected, createTemplateFile(tmpDir, "Autounattend.tpl"))
	expected = append(expected, createTemplateFile(tmpDir, "Autounattend.windowsPE.tpl"))
	expected = append(expected, createTemplateFile(tmpDir, "Autounattend-windows2012r2.windowsPE.tpl"))
	createTemplateFile(tmpDir, "packer-windows2012r2.windowsPE.tpl")

	files, err := ListFiles(tmpDir, "Autounattend")
	if err != nil {
		t.Error("Couldn't list glob files:", err)
	}
	if len(files) != 3 {
		t.Errorf("Expected 3 files to be returned, but got %d", len(files))
	}
	for _, f := range expected {
		if !contains(files, f) {
			t.Errorf("Expected files to contain %s", f)
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
