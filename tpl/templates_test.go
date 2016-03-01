package tpl

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListTemplatesScenario(t *testing.T) {
	// create temporary directory to store all our test template files
	tmpDir, err := ioutil.TempDir("", "inductor")
	if err != nil {
		t.Error("Couldn't create test temp dir: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// OS specific partial template overrides
	nanoautounattendxmldisksptpl := createTemplateFile(tmpDir, "nano/Autounattend.xml.disks.ptpl")
	nanopackerjsonprovisionersptpl := createTemplateFile(tmpDir, "nano/packer.json.provisioners.ptpl")

	// root template in a sub directory
	windowsxpautounattendxmltpl := createTemplateFile(tmpDir, "windowsxp/Autounattend.xml.tpl")

	// some non-template files to ignore
	createTemplateFile(tmpDir, "scripts/chef.bat")
	createTemplateFile(tmpDir, "scripts/winrm.ps1")
	createTemplateFile(tmpDir, "scripts/win-updates.ps1")
	createTemplateFile(tmpDir, "scripts/nano/cleanup.ps1")
	createTemplateFile(tmpDir, "scripts/nano/create.ps1")
	createTemplateFile(tmpDir, "scripts/nano/SetupComplete.cmd")
	createTemplateFile(tmpDir, "CHANGELOG.md")
	createTemplateFile(tmpDir, "inductor.json")
	createTemplateFile(tmpDir, "LICENSE.md")
	createTemplateFile(tmpDir, "README.md")

	// root Autounattend.xml template and shared root partial templates
	autounattendxmltpl := createTemplateFile(tmpDir, "Autounattend.xml.tpl")
	autounattendxmloobeptpl := createTemplateFile(tmpDir, "Autounattend.xml.oobe.ptpl")
	autounattendxmldisksptpl := createTemplateFile(tmpDir, "Autounattend.xml.disks.ptpl")

	// root packer.json template and shared root partial templates
	packerjsontpl := createTemplateFile(tmpDir, "packer.json.tpl")
	packerjsonbuildersptpl := createTemplateFile(tmpDir, "packer.json.builders.ptpl")
	createTemplateFile(tmpDir, "packer.json.provisioners.ptpl")

	// root Vagrantfile template without any sub templates or overrides
	vagrantfiletpl := createTemplateFile(tmpDir, "Vagrantfile.tpl")

	// reset template finder functions from unit tests
	listRootTemplatesFn = listRootTemplates
	listPartialTemplatesFn = listPartialTemplates
	listPartialTemplatesOSSpecificFn = listPartialTemplatesOSSpecific

	// act
	templates := New(tmpDir)

	// assert top level root templates
	autounattendXMLRootTemplate := findRootTemplate(templates.All, autounattendxmltpl)
	if autounattendXMLRootTemplate == nil {
		t.Errorf("Expected root templates to contain '%s'", autounattendxmltpl)
		return
	}
	windowsXpAutounattendXMLTemplate := findRootTemplate(templates.All, windowsxpautounattendxmltpl)
	if windowsXpAutounattendXMLTemplate == nil {
		t.Errorf("Expected root templates to contain '%s'", windowsXpAutounattendXMLTemplate)
		return
	}
	packerJSONTemplate := findRootTemplate(templates.All, packerjsontpl)
	if packerJSONTemplate == nil {
		t.Errorf("Expected root templates to contain '%s'", packerjsontpl)
		return
	}
	vagrantfileTemplate := findRootTemplate(templates.All, vagrantfiletpl)
	if vagrantfileTemplate == nil {
		t.Errorf("Expected root templates to contain '%s'", vagrantfiletpl)
		return
	}

	// assert autounattend Nano partial templates
	expectedAutounattendPartials := []string{
		autounattendxmloobeptpl,
		nanoautounattendxmldisksptpl,
	}
	autounattendNanoPartialTemplates := autounattendXMLRootTemplate.PartialTemplates("nano")
	if len(expectedAutounattendPartials) != len(autounattendNanoPartialTemplates) {
		t.Errorf("Expected %d partial templates, but got %d",
			len(expectedAutounattendPartials), len(autounattendNanoPartialTemplates))
	}
	for _, expectedFile := range expectedAutounattendPartials {
		if !containsPartialTemplate(autounattendNanoPartialTemplates, expectedFile) {
			t.Errorf("Expected partial templates to contain '%s'", expectedFile)
		}
	}

	// assert autounattend Windows2012R2 partial templates should have no overrides
	expectedWindows2012r2AutounattendPartials := []string{
		autounattendxmloobeptpl,
		autounattendxmldisksptpl,
	}
	autounattendWindows2012r2PartialTemplates := autounattendXMLRootTemplate.PartialTemplates("windows2012r2")
	if len(expectedWindows2012r2AutounattendPartials) != len(autounattendWindows2012r2PartialTemplates) {
		t.Errorf("Expected %d partial templates, but got %d",
			len(expectedWindows2012r2AutounattendPartials), len(autounattendWindows2012r2PartialTemplates))
	}
	for _, expectedFile := range expectedWindows2012r2AutounattendPartials {
		if !containsPartialTemplate(autounattendWindows2012r2PartialTemplates, expectedFile) {
			t.Errorf("Expected partial templates to contain '%s'", expectedFile)
		}
	}

	// assert packer.json Nano partial templates
	expectedPackerPartials := []string{
		nanopackerjsonprovisionersptpl,
		packerjsonbuildersptpl,
	}
	packerNanoPartialTemplates := packerJSONTemplate.PartialTemplates("nano")
	if len(expectedPackerPartials) != len(packerNanoPartialTemplates) {
		t.Errorf("Expected %d partial templates, but got %d",
			len(expectedPackerPartials), len(packerNanoPartialTemplates))
	}
	for _, expectedFile := range expectedPackerPartials {
		if !containsPartialTemplate(packerNanoPartialTemplates, expectedFile) {
			t.Errorf("Expected partial templates to contain '%s'", expectedFile)
		}
	}

	// assert the Vagrantfile template has no partials
	vagrantfileNanoPartialTemplates := vagrantfileTemplate.PartialTemplates("nano")
	if len(vagrantfileNanoPartialTemplates) != 0 {
		t.Errorf("Expected no Vagrantfile partial templates, but got %d", len(vagrantfileNanoPartialTemplates))
	}

	// the WindowsXP Autounattend in a subdir should have no partials
	autounattendWindowsXpPartialTemplates := windowsXpAutounattendXMLTemplate.PartialTemplates("windowsxp")
	if len(autounattendWindowsXpPartialTemplates) != 0 {
		t.Errorf("Expected no WindowsXP Autounattend partial templates, but got %d",
			len(autounattendWindowsXpPartialTemplates))
	}
}

func TestListTemplates(t *testing.T) {
	var expectedFiles = []string{
		"/Users/sneal/packer-windows/Autounattend.xml.tpl",
		"/Users/sneal/packer-windows/Vagrantfile.tpl",
	}
	listRootTemplatesFn = func(baseDir string) []string {
		return []string{
			"/Users/sneal/packer-windows/Autounattend.xml.tpl",
			"/Users/sneal/packer-windows/Vagrantfile.tpl",
		}
	}
	templates := New("/Users/sneal/packer-windows")
	if templates.BaseDir != "/Users/sneal/packer-windows" {
		t.Errorf("Unexpected a base dir of '%s'", templates.BaseDir)
	}
	for _, expectedFile := range expectedFiles {
		if !containsRootTemplate(templates.All, expectedFile) {
			t.Errorf("Expected partial templates to contain '%s'", expectedFile)
		}
	}
}

func findRootTemplate(pts []RootTemplate, path string) *RootTemplate {
	for _, pt := range pts {
		if pt.Path == path {
			return &pt
		}
	}
	return nil
}

func containsRootTemplate(pts []RootTemplate, e string) bool {
	for _, pt := range pts {
		if pt.Path == e {
			return true
		}
	}
	return false
}

func containsPartialTemplate(pts []PartialTemplate, e string) bool {
	for _, pt := range pts {
		if pt.Path == e {
			return true
		}
	}
	return false
}

func createTemplateFile(baseDir, path string) string {
	fullPath := filepath.Join(baseDir, path)
	fileName := filepath.Base(fullPath)
	fullDir := strings.TrimSuffix(fullPath, fileName)
	os.MkdirAll(fullDir, 0744)
	ioutil.WriteFile(fullPath, []byte(fileName), 0644)
	return fullPath
}
