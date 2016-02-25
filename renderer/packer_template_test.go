package renderer

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestCanLoadTemplates(t *testing.T) {
	// create temporary directory to store all our testNewPackerTemplate() files
	tmpDir, err := ioutil.TempDir("", "inductor")
	if err != nil {
		t.Error("Couldn't create test temp dir: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// stub out the vagrant, packer, and autounattend files
	createTemplateFile(tmpDir, "packer.tpl")
	createTemplateFile(tmpDir, "Autounattend.tpl")
	createTemplateFile(tmpDir, "Vagrantfile.tpl")

	template, err := NewPackerTemplateWithOverrides(tmpDir, "windows2012r2")
	if err != nil {
		t.Error("Error loading Packer templates with overrides", err)
	}
	if template.PackerTpl != "packer.tpl" {
		t.Errorf("The packer.tpl wasn't properly read in, got: %s", template.PackerTpl)
	}
	if template.AutounattendTpl != "Autounattend.tpl" {
		t.Errorf("The Autounattend.tpl wasn't properly read in, got: %s", template.AutounattendTpl)
	}
	if template.VagrantfileTpl != "Vagrantfile.tpl" {
		t.Errorf("The Vagrantfile.tpl wasn't properly read in, got: %s", template.VagrantfileTpl)
	}
}

func TestCanLoadAutounattendMultiTemplates(t *testing.T) {
	// create temporary directory to store all our testNewPackerTemplate() files
	tmpDir, err := ioutil.TempDir("", "inductor")
	if err != nil {
		t.Error("Couldn't create test temp dir: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// stub out the vagrant, packer, and autounattend files
	createTemplateFile(tmpDir, "packer.tpl")
	createTemplateFile(tmpDir, "Vagrantfile.tpl")
	createTemplateFile(tmpDir, "Autounattend.tpl")
	createTemplateFile(tmpDir, "Autounattend-windows2012r2.windowsPE.tpl")
	createTemplateFile(tmpDir, "Autounattend-windows2012r2.oobe.tpl")

	template, err := NewPackerTemplateWithOverrides(tmpDir, "windows2012r2")
	if err != nil {
		t.Error("Error loading Packer templates with overrides", err)
	}
	if !strings.Contains(template.AutounattendTpl, "Autounattend.tpl") {
		t.Errorf("The main Autounattend.xml.tpl wasn't properly read in, got: %s", template.AutounattendTpl)
	}
	if !strings.Contains(template.AutounattendTpl, "Autounattend-windows2012r2.windowsPE.tpl") {
		t.Errorf("The Autounattend-windows2012r2.windowsPE.tpl wasn't properly read in, got: %s", template.AutounattendTpl)
	}
	if !strings.Contains(template.AutounattendTpl, "Autounattend-windows2012r2.oobe.tpl") {
		t.Errorf("The Autounattend-windows2012r2.oobe.tpl wasn't properly read in, got: %s", template.AutounattendTpl)
	}
}

func TestCanLoadPackerJSONMultiTemplates(t *testing.T) {
	// create temporary directory to store all our testNewPackerTemplate() files
	tmpDir, err := ioutil.TempDir("", "inductor")
	if err != nil {
		t.Error("Couldn't create test temp dir: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// stub out the vagrant, packer, and autounattend files
	createTemplateFile(tmpDir, "Vagrantfile.tpl")
	createTemplateFile(tmpDir, "Autounattend.tpl")
	createTemplateFile(tmpDir, "packer.tpl")
	createTemplateFile(tmpDir, "packer.vbox.tpl")
	createTemplateFile(tmpDir, "packer.vmware.tpl")

	template, err := NewPackerTemplateWithOverrides(tmpDir, "windows2012r2")
	if err != nil {
		t.Error("Error loading Packer templates with overrides", err)
	}
	if !strings.Contains(template.PackerTpl, "packer.tpl") {
		t.Errorf("The main packer.tpl wasn't properly read in, got: %s", template.PackerTpl)
	}
	if !strings.Contains(template.PackerTpl, "packer.vbox.tpl") {
		t.Errorf("The packer.vbox.tpl wasn't properly read in, got: %s", template.PackerTpl)
	}
	if !strings.Contains(template.PackerTpl, "packer.vmware.tpl") {
		t.Errorf("The packer.vmware.tpl wasn't properly read in, got: %s", template.PackerTpl)
	}
}

func createTemplateFile(tmpDir string, filename string) string {
	tplPath := path.Join(tmpDir, filename)
	ioutil.WriteFile(tplPath, []byte(filename), 0644)
	return tplPath
}
