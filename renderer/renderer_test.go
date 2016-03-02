package renderer_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/joefitzgerald/inductor/renderer"
	"github.com/joefitzgerald/inductor/tpl"
	"github.com/joefitzgerald/inductor/tpl/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Renderer", func() {
	var (
		err           error
		engine        renderer.Renderer
		templates     *fakes.FakeTemplateContainer
		renderOptions *renderer.RenderOptions
		outDir        string
	)

	Describe("Vagrantfile template", func() {
		BeforeEach(func() {
			outDir, err = ioutil.TempDir("", "inductor")
			Expect(err).NotTo(HaveOccurred())
			renderOptions = renderer.NewDefaultRenderOptions()
			vagrantTemplate := new(fakes.FakeTemplater)
			vagrantTemplate.ContentStub = writeVagrantfile
			vagrantTemplate.BaseFilenameReturns("Vagrantfile")
			templates = new(fakes.FakeTemplateContainer)
			templates.ListTemplatesReturns([]tpl.Templater{vagrantTemplate})
			engine = renderer.New(renderOptions, outDir)
			err = engine.Render(templates)
		})
		AfterEach(func() {
			os.RemoveAll(outDir)
		})
		It("should not have errored", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		// TODO: need to test actually rendering a file
		It("should render Vagrantfile in out dir", func() {
			vagrantfilePath := filepath.Join(outDir, "Vagrantfile")
			Expect(vagrantfilePath).To(BeARegularFile())
		})
	})
})

func writeVagrantfile(buffer io.Writer) error {
	content := `
# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.require_version ">= 1.6.2"

Vagrant.configure("2") do |config|
config.vm.box = "{{.OSName}}"
config.vm.communicator = "winrm"

# Admin user name and password
config.winrm.username = "{{.Username}}"
config.winrm.password = "{{.Password}}"

config.vm.network :forwarded_port, guest: 3389, host: 3389, id: "rdp", auto_correct: true
config.vm.network :forwarded_port, guest: 22, host: 2222, id: "ssh", auto_correct: true

config.vm.provider :virtualbox do |v, override|
v.customize ["modifyvm", :id, "--memory", {{.RAM}}]
v.customize ["modifyvm", :id, "--cpus", {{.CPU}}]
v.customize ["setextradata", "global", "GUI/SuppressMessages", "all" ]
end

config.vm.provider :vmware_fusion do |v, override|
v.vmx["memsize"] = "{{.RAM}}"
v.vmx["numvcpus"] = "{{.CPU}}"
v.vmx["ethernet0.virtualDev"] = "vmxnet3"
v.vmx["RemoteDisplay.vnc.enabled"] = "false"
v.vmx["RemoteDisplay.vnc.port"] = "5900"
v.vmx["scsi0.virtualDev"] = "lsisas1068"
end

config.vm.provider :vmware_workstation do |v, override|
v.vmx["memsize"] = "{{.RAM}}"
v.vmx["numvcpus"] = "{{.CPU}}"
v.vmx["ethernet0.virtualDev"] = "vmxnet3"
v.vmx["RemoteDisplay.vnc.enabled"] = "false"
v.vmx["RemoteDisplay.vnc.port"] = "5900"
v.vmx["scsi0.virtualDev"] = "lsisas1068"
end
end
`
	buffer.Write([]byte(content))
	return nil
}
