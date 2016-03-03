package tpl_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/joefitzgerald/inductor/tpl"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tpl", func() {
	var (
		err       error
		osName    string
		tmpDir    string
		templates tpl.TemplateContainer
	)
	BeforeEach(func() {
		tmpDir, err = ioutil.TempDir("", "inductor")
		Expect(err).NotTo(HaveOccurred())
		createTemplateFile(tmpDir, "nano/Autounattend.xml.disks.ptpl")
		createTemplateFile(tmpDir, "nano/packer.json.provisioners.ptpl")
		createTemplateFile(tmpDir, "windowsxp/Autounattend.xml.tpl")
		createTemplateFile(tmpDir, "scripts/win-updates.ps1")
		createTemplateFile(tmpDir, "scripts/nano/SetupComplete.cmd")
		createTemplateFile(tmpDir, "inductor.json")
		createTemplateFile(tmpDir, "README.md")
		createTemplateFile(tmpDir, "Autounattend.xml.tpl")
		createTemplateFile(tmpDir, "Autounattend.xml.oobe.ptpl")
		createTemplateFile(tmpDir, "Autounattend.xml.disks.ptpl")
		createTemplateFile(tmpDir, "packer.json.tpl")
		createTemplateFile(tmpDir, "packer.json.builders.ptpl")
		createTemplateFile(tmpDir, "packer.json.provisioners.ptpl")
		createTemplateFile(tmpDir, "Vagrantfile.tpl")
	})
	JustBeforeEach(func() {
		templates = tpl.New(tmpDir, osName)
	})
	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	Context("Windows2012r2", func() {
		BeforeEach(func() {
			osName = "windows2012r2"
		})
		Describe("Autounattend.xml.tpl root template", func() {
			var rootTemplate tpl.Templater
			JustBeforeEach(func() {
				rootTemplate = templates.FindTemplate(filepath.Join(tmpDir, "Autounattend.xml.tpl"))
			})
			It("should be found", func() {
				Expect(rootTemplate).ToNot(BeNil())
			})
			It("should have 2 partial templates", func() {
				Expect(rootTemplate.ListTemplates()).To(HaveLen(2))
			})
			It("should include Autounattend.xml.oobe.ptpl", func() {
				path := filepath.Join(tmpDir, "Autounattend.xml.oobe.ptpl")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
			It("should include Autounattend.xml.disks.ptpl", func() {
				path := filepath.Join(tmpDir, "Autounattend.xml.disks.ptpl")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
			It("should have template content", func() {
				expected := `Autounattend.xml.tpl
{{define "disks"}}
Autounattend.xml.disks.ptpl
{{end}}
{{define "oobe"}}
Autounattend.xml.oobe.ptpl
{{end}}`
				var buffer bytes.Buffer
				Expect(rootTemplate.Content(&buffer)).NotTo(HaveOccurred())
				Expect(buffer.String()).To(Equal(expected))
			})
		})
	})

	Context("Nano", func() {
		BeforeEach(func() {
			osName = "nano"
		})
		Describe("Autounattend.xml.tpl root template", func() {
			var rootTemplate tpl.Templater
			JustBeforeEach(func() {
				rootTemplate = templates.FindTemplate(filepath.Join(tmpDir, "Autounattend.xml.tpl"))
			})
			It("should be found", func() {
				Expect(rootTemplate).ToNot(BeNil())
			})
			It("should have 2 partial templates", func() {
				Expect(rootTemplate.ListTemplates()).To(HaveLen(2))
			})
			It("should include Autounattend.xml.oobe.ptpl", func() {
				path := filepath.Join(tmpDir, "Autounattend.xml.oobe.ptpl")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
			It("should include nano/Autounattend.xml.disks.ptpl", func() {
				path := filepath.Join(tmpDir, "nano/Autounattend.xml.disks.ptpl")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
		})

		Describe("packer.json.tpl root template", func() {
			var rootTemplate tpl.Templater
			JustBeforeEach(func() {
				rootTemplate = templates.FindTemplate(filepath.Join(tmpDir, "packer.json.tpl"))
			})
			It("should be found", func() {
				Expect(rootTemplate).ToNot(BeNil())
			})
			It("should have 2 partial templates", func() {
				Expect(rootTemplate.ListTemplates()).To(HaveLen(2))
			})
			It("should include packer.json.builders.ptpl", func() {
				path := filepath.Join(tmpDir, "packer.json.builders.ptpl")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
			It("should include nano/packer.json.provisioners.ptpl", func() {
				path := filepath.Join(tmpDir, "nano/packer.json.provisioners.ptpl")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
		})

		Describe("Vagrantfile.tpl root template", func() {
			var rootTemplate tpl.Templater
			JustBeforeEach(func() {
				rootTemplate = templates.FindTemplate(filepath.Join(tmpDir, "Vagrantfile.tpl"))
			})
			It("should be found", func() {
				Expect(rootTemplate).ToNot(BeNil())
			})
			It("should have zero partial templates", func() {
				Expect(rootTemplate.ListTemplates()).To(BeEmpty())
			})
			It("should have basefilename", func() {
				Expect(rootTemplate.BaseFilename()).To(Equal("Vagrantfile"))
			})
		})
	})
})

func createTemplateFile(baseDir, path string) {
	fullPath := filepath.Join(baseDir, path)
	fileName := filepath.Base(fullPath)
	fullDir := strings.TrimSuffix(fullPath, fileName)
	os.MkdirAll(fullDir, 0744)
	ioutil.WriteFile(fullPath, []byte(fileName), 0644)
}
