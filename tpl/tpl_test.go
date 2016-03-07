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
		createTemplateFile(tmpDir, "nano/Autounattend.xml.disks.partial")
		createTemplateFile(tmpDir, "nano/packer.json.provisioners.partial")
		createTemplateFile(tmpDir, "windowsxp/Autounattend.xml.template")
		createTemplateFile(tmpDir, "scripts/win-updates.ps1")
		createTemplateFile(tmpDir, "scripts/nano/SetupComplete.cmd")
		createTemplateFile(tmpDir, "inductor.json")
		createTemplateFile(tmpDir, "README.md")
		createTemplateFile(tmpDir, "Autounattend.xml.template")
		createTemplateFile(tmpDir, "Autounattend.xml.oobe.partial")
		createTemplateFile(tmpDir, "Autounattend.xml.disks.partial")
		createTemplateFile(tmpDir, "packer.json.template")
		createTemplateFile(tmpDir, "packer.json.builders.partial")
		createTemplateFile(tmpDir, "packer.json.provisioners.partial")
		createTemplateFile(tmpDir, "Vagrantfile.template")
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
		Describe("Autounattend.xml.template root template", func() {
			var rootTemplate tpl.Templater
			JustBeforeEach(func() {
				rootTemplate = templates.FindTemplate(filepath.Join(tmpDir, "Autounattend.xml.template"))
			})
			It("should be found", func() {
				Expect(rootTemplate).ToNot(BeNil())
			})
			It("should have 2 partial templates", func() {
				Expect(rootTemplate.ListTemplates()).To(HaveLen(2))
			})
			It("should include Autounattend.xml.oobe.partial", func() {
				path := filepath.Join(tmpDir, "Autounattend.xml.oobe.partial")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
			It("should include Autounattend.xml.disks.partial", func() {
				path := filepath.Join(tmpDir, "Autounattend.xml.disks.partial")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
			It("should have template content", func() {
				expected := `Autounattend.xml.template
{{define "disks"}}
Autounattend.xml.disks.partial
{{end}}
{{define "oobe"}}
Autounattend.xml.oobe.partial
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
		Describe("Autounattend.xml.template root template", func() {
			var rootTemplate tpl.Templater
			JustBeforeEach(func() {
				rootTemplate = templates.FindTemplate(filepath.Join(tmpDir, "Autounattend.xml.template"))
			})
			It("should be found", func() {
				Expect(rootTemplate).ToNot(BeNil())
			})
			It("should have 2 partial templates", func() {
				Expect(rootTemplate.ListTemplates()).To(HaveLen(2))
			})
			It("should include Autounattend.xml.oobe.partial", func() {
				path := filepath.Join(tmpDir, "Autounattend.xml.oobe.partial")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
			It("should include nano/Autounattend.xml.disks.partial", func() {
				path := filepath.Join(tmpDir, "nano/Autounattend.xml.disks.partial")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
		})

		Describe("packer.json.template root template", func() {
			var rootTemplate tpl.Templater
			JustBeforeEach(func() {
				rootTemplate = templates.FindTemplate(filepath.Join(tmpDir, "packer.json.template"))
			})
			It("should be found", func() {
				Expect(rootTemplate).ToNot(BeNil())
			})
			It("should have 2 partial templates", func() {
				Expect(rootTemplate.ListTemplates()).To(HaveLen(2))
			})
			It("should include packer.json.builders.partial", func() {
				path := filepath.Join(tmpDir, "packer.json.builders.partial")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
			It("should include nano/packer.json.provisioners.partial", func() {
				path := filepath.Join(tmpDir, "nano/packer.json.provisioners.partial")
				Expect(rootTemplate.FindTemplate(path)).ToNot(BeNil())
			})
		})

		Describe("Vagrantfile.template root template", func() {
			var rootTemplate tpl.Templater
			JustBeforeEach(func() {
				rootTemplate = templates.FindTemplate(filepath.Join(tmpDir, "Vagrantfile.template"))
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
