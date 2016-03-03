package configuration_test

import (
	"strings"

	"github.com/joefitzgerald/inductor/configuration"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {
	var (
		err    error
		config *configuration.InductorConfiguration
	)
	BeforeEach(func() {
		config, err = configuration.New(strings.NewReader(testData))
		Expect(err).NotTo(HaveOccurred())
	})
	It("should show the gui", func() {
		Expect(config.Headless).To(BeFalse())
	})
	It("should skip windows updates", func() {
		Expect(config.WindowsUpdates).To(BeFalse())
	})
	It("should use the SSH communicator", func() {
		Expect(config.Communicator).To(Equal("ssh"))
	})
	It("should have an output dir of /tmp/foo/bar", func() {
		Expect(config.OutDir).To(Equal("/tmp/foo/bar"))
	})
	It("should have a username of admin", func() {
		Expect(config.Username).To(Equal("admin"))
	})
	It("should a password of secret", func() {
		Expect(config.Password).To(Equal("secret"))
	})
	It("should have RAM set to 1024", func() {
		Expect(config.RAM).To(Equal(uint32(1024)))
	})
	It("should have 2 CPUs", func() {
		Expect(config.CPU).To(Equal(uint8(1)))
	})
	It("should have a disk size of 10000", func() {
		Expect(config.DiskSize).To(Equal(uint32(10000)))
	})
	Describe("List available operating systems", func() {
		var oses []string
		BeforeEach(func() {
			oses = config.List()
		})
		It("should return 2 operating systems", func() {
			Expect(oses).To(HaveLen(2))
		})
		It("should return windows10 as the first entry", func() {
			Expect(oses[0]).To(Equal("windows10"))
		})
		It("should return windows2008r2 as the second entry", func() {
			Expect(oses[1]).To(Equal("windows2008r2"))
		})
	})
	Describe("Get operating system configuration", func() {
		var (
			os      *configuration.OperatingSystem
			found   bool
			edition configuration.Edition
		)
		Context("windows 10", func() {
			BeforeEach(func() {
				os, found = config.Get("windows10")
			})
			It("should have found windows10", func() {
				Expect(found).To(BeTrue())
			})
			It("should have correct name", func() {
				Expect(os.Name).To(Equal("windows10"))
			})
			It("should have correct iso url", func() {
				Expect(os.IsoURL).To(Equal("http://care.dlservice.microsoft.com/dl/download/C/3/9/C399EEA8-135D-4207-92C9-6AAB3259F6EF/10240.16384.150709-1700.TH1_CLIENTENTERPRISEEVAL_OEMRET_X64FRE_EN-US.ISO"))
			})
			It("should have correct iso checksum", func() {
				Expect(os.IsoChecksum).To(Equal("56ab095075be28a90bc0b510835280975c6bb2ce"))
			})
			It("should have correct iso checksum type", func() {
				Expect(os.IsoChecksumType).To(Equal("sha1"))
			})
			It("should have correct vbox guest type", func() {
				Expect(os.VirtualboxGuestOsType).To(Equal("Windows81_64"))
			})
			It("should have correct vmware guest type", func() {
				Expect(os.VmwareGuestOsType).To(Equal("windows8srv-64"))
			})
			It("should have 1 edition", func() {
				Expect(os.Editions).To(HaveLen(1))
			})
			Context("enterprise edition", func() {
				BeforeEach(func() {
					edition = os.Editions["enterprise"]
				})
				It("should have correct windows image name", func() {
					Expect(edition.WindowsImageName).To(Equal("Windows 10 Enterprise Evaluation"))
				})
			})
		})
		Context("windows 2008r2", func() {
			BeforeEach(func() {
				os, found = config.Get("windows2008r2")
			})
			It("should have found windows2008r2", func() {
				Expect(found).To(BeTrue())
			})
			It("should have correct name", func() {
				Expect(os.Name).To(Equal("windows2008r2"))
			})
			It("should have correct iso url", func() {
				Expect(os.IsoURL).To(Equal("http://download.microsoft.com/download/7/5/E/75EC4E54-5B02-42D6-8879-D8D3A25FBEF7/7601.17514.101119-1850_x64fre_server_eval_en-us-GRMSXEVAL_EN_DVD.iso"))
			})
			It("should have correct iso checksum", func() {
				Expect(os.IsoChecksum).To(Equal("4263be2cf3c59177c45085c0a7bc6ca5"))
			})
			It("should have correct iso checksum type", func() {
				Expect(os.IsoChecksumType).To(Equal("md5"))
			})
			It("should have correct vbox guest type", func() {
				Expect(os.VirtualboxGuestOsType).To(Equal("Windows2008_64"))
			})
			It("should have correct vmware guest type", func() {
				Expect(os.VmwareGuestOsType).To(Equal("windows7srv-64"))
			})
			It("should have 2 editions", func() {
				Expect(os.Editions).To(HaveLen(2))
			})
			Context("enterprise edition", func() {
				BeforeEach(func() {
					edition = os.Editions["enterprise"]
				})
				It("should have correct windows image name", func() {
					Expect(edition.WindowsImageName).To(Equal("Windows Server 2008 R2 SERVERENTERPRISE"))
				})
			})
			Context("standard edition", func() {
				BeforeEach(func() {
					edition = os.Editions["standard"]
				})
				It("should have correct windows image name", func() {
					Expect(edition.WindowsImageName).To(Equal("Windows Server 2008 R2 SERVERSTANDARD"))
				})
			})
		})
		Context("linux", func() {
			BeforeEach(func() {
				os, found = config.Get("linux")
			})
			It("should not have found linux", func() {
				Expect(found).To(BeFalse())
			})
		})
	})
})

var testData = `
{
  "config":{
    "headless":false,
    "windows_updates":false,
    "communicator":"ssh",
    "out_dir":"/tmp/foo/bar",
    "username":"admin",
    "password":"secret",
    "ram":1024,
    "cpu":1,
    "disk_size":10000
  },
  "operating_systems":{
    "windows10":{
      "iso_url":"http://care.dlservice.microsoft.com/dl/download/C/3/9/C399EEA8-135D-4207-92C9-6AAB3259F6EF/10240.16384.150709-1700.TH1_CLIENTENTERPRISEEVAL_OEMRET_X64FRE_EN-US.ISO",
      "iso_checksum_type":"sha1",
      "iso_checksum":"56ab095075be28a90bc0b510835280975c6bb2ce",
      "virtualbox_guest_os_type":"Windows81_64",
      "vmware_guest_os_type":"windows8srv-64",
      "editions":{
        "enterprise":{
          "windows_image_name":"Windows 10 Enterprise Evaluation"
        }
      }
    },
    "windows2008r2":{
      "iso_url":"http://download.microsoft.com/download/7/5/E/75EC4E54-5B02-42D6-8879-D8D3A25FBEF7/7601.17514.101119-1850_x64fre_server_eval_en-us-GRMSXEVAL_EN_DVD.iso",
      "iso_checksum_type":"md5",
      "iso_checksum":"4263be2cf3c59177c45085c0a7bc6ca5",
      "virtualbox_guest_os_type":"Windows2008_64",
      "vmware_guest_os_type":"windows7srv-64",
      "editions":{
        "standard":{
          "windows_image_name":"Windows Server 2008 R2 SERVERSTANDARD"
        },
        "enterprise":{
          "windows_image_name":"Windows Server 2008 R2 SERVERENTERPRISE"
        }
      }
    }
  }
}
`
