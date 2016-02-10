package renderer

import (
	"io/ioutil"
	"os"
	"path"
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
	packerTplPath := path.Join(tmpDir, "packer.json.tpl")
	autounattendTplPath := path.Join(tmpDir, "Autounattend.xml.tpl")
	vagrantTplPath := path.Join(tmpDir, "Vagrantfile.tpl")

	ioutil.WriteFile(packerTplPath, []byte("packer.json contents"), 0644)
	ioutil.WriteFile(autounattendTplPath, []byte("Autounattend.xml contents"), 0644)
	ioutil.WriteFile(vagrantTplPath, []byte("Vagrantfile contents"), 0644)

	template := NewPackerTemplateWithOverrides(packerTplPath, autounattendTplPath, vagrantTplPath)
	if template.PackerTpl != "packer.json contents" {
		t.Errorf("The packer.json.tpl wasn't properly read in, got: %s", template.PackerTpl)
	}
	if template.AutounattendTpl != "Autounattend.xml contents" {
		t.Errorf("The Autounattend.xml.tpl wasn't properly read in, got: %s", template.AutounattendTpl)
	}
	if template.VagrantfileTpl != "Vagrantfile contents" {
		t.Errorf("The Vagrantfile.tpl wasn't properly read in, got: %s", template.VagrantfileTpl)
	}
}
