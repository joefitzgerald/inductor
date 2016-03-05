package cpy_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/joefitzgerald/inductor/cpy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cpy", func() {
	var (
		err    error
		copier cpy.Copier
		outDir string
		srcDir string
	)

	BeforeEach(func() {
		srcDir, err = ioutil.TempDir("", "inductor-src")
		Expect(err).NotTo(HaveOccurred())
		createFile(srcDir, "README.md")
		createFile(srcDir, "Vagrantfile")
		createFile(srcDir, "Autounattend.xml.template")
		createFile(srcDir, "Autounattend.xml.oobe.partial")
		createFile(srcDir, "scripts/winrm.ps1")
		createFile(srcDir, "scripts/windows-updates.ps1")
		createFile(srcDir, "scripts/nano/SetupComplete.cmd")
		outDir, err = ioutil.TempDir("", "inductor-out")
		Expect(err).NotTo(HaveOccurred())
		copier = cpy.New()
		err = copier.Copy(srcDir, outDir)
	})
	AfterEach(func() {
		os.RemoveAll(srcDir)
		os.RemoveAll(outDir)
	})

	Describe("Copy recursive", func() {
		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("should ignore .template files", func() {
			path := filepath.Join(outDir, "Autounattend.xml.template")
			Expect(path).ToNot(BeARegularFile())
		})
		It("should ignore .partial files", func() {
			path := filepath.Join(outDir, "Autounattend.xml.partial")
			Expect(path).ToNot(BeARegularFile())
		})
		It("should copy regular files in src dir", func() {
			Expect(filepath.Join(outDir, "README.md")).To(BeARegularFile())
			Expect(filepath.Join(outDir, "Vagrantfile")).To(BeARegularFile())
		})
		It("should recursively copy files in sub dirs", func() {
			Expect(filepath.Join(outDir, "scripts/winrm.ps1")).To(BeARegularFile())
			Expect(filepath.Join(outDir, "scripts/windows-updates.ps1")).To(BeARegularFile())
			Expect(filepath.Join(outDir, "scripts/nano/SetupComplete.cmd")).To(BeARegularFile())
		})
	})
})

func createFile(baseDir, path string) {
	fullPath := filepath.Join(baseDir, path)
	fileName := filepath.Base(fullPath)
	fullDir := strings.TrimSuffix(fullPath, fileName)
	os.MkdirAll(fullDir, 0744)
	err := ioutil.WriteFile(fullPath, []byte(fileName), 0644)
	if err != nil {
		panic(err)
	}
	Expect(fullPath).To(BeARegularFile())
}
