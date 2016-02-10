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
	packerTplPath := createTemplateFile(tmpDir, "packer.json.tpl")
	autounattendTplPath := createTemplateFile(tmpDir, "Autounattend.xml.tpl")
	vagrantTplPath := createTemplateFile(tmpDir, "Vagrantfile.tpl")

	template := NewPackerTemplateWithOverrides(packerTplPath, autounattendTplPath, vagrantTplPath)
	if template.PackerTpl != "packer.json.tpl" {
		t.Errorf("The packer.json.tpl wasn't properly read in, got: %s", template.PackerTpl)
	}
	if template.AutounattendTpl != "Autounattend.xml.tpl" {
		t.Errorf("The Autounattend.xml.tpl wasn't properly read in, got: %s", template.AutounattendTpl)
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
	packerTplPath := createTemplateFile(tmpDir, "packer.json.tpl")
	vagrantTplPath := createTemplateFile(tmpDir, "Vagrantfile.tpl")
	autounattendTplPath := createTemplateFile(tmpDir, "Autounattend.xml.tpl")
	createTemplateFile(tmpDir, "Autounattend.xml.tpl.disks")
	createTemplateFile(tmpDir, "Autounattend.xml.tpl.oobe")

	template := NewPackerTemplateWithOverrides(packerTplPath, autounattendTplPath, vagrantTplPath)
	if !strings.Contains(template.AutounattendTpl, "Autounattend.xml.tpl") {
		t.Errorf("The main Autounattend.xml.tpl wasn't properly read in, got: %s", template.AutounattendTpl)
	}
	if !strings.Contains(template.AutounattendTpl, "Autounattend.xml.tpl.disks") {
		t.Errorf("The Autounattend.xml.tpl.disks wasn't properly read in, got: %s", template.AutounattendTpl)
	}
	if !strings.Contains(template.AutounattendTpl, "Autounattend.xml.tpl.oobe") {
		t.Errorf("The Autounattend.xml.tpl.oobe wasn't properly read in, got: %s", template.AutounattendTpl)
	}
}

func createTemplateFile(tmpDir string, filename string) string {
	tplPath := path.Join(tmpDir, filename)
	ioutil.WriteFile(tplPath, []byte(filename), 0644)
	return tplPath
}
