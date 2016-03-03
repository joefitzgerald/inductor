package configuration

import (
	"sort"
	"strings"
	"testing"
)

var testData = `
{
	"config": {
		"headless": false,
		"windows_updates": false,
		"communicator": "ssh",
		"out_dir": "/tmp/foo/bar"
	},
	"operating_systems": {
	  "windows10": {
	    "iso_url": "http://care.dlservice.microsoft.com/dl/download/C/3/9/C399EEA8-135D-4207-92C9-6AAB3259F6EF/10240.16384.150709-1700.TH1_CLIENTENTERPRISEEVAL_OEMRET_X64FRE_EN-US.ISO",
	    "iso_checksum_type": "sha1",
	    "iso_checksum": "56ab095075be28a90bc0b510835280975c6bb2ce",
	    "virtualbox_guest_os_type": "Windows81_64",
	    "vmware_guest_os_type": "windows8srv-64",
			"editions": {
        "enterprise": {
          "windows_image_name": "Windows 10 Enterprise Evaluation"
        }
      }
	  },
	  "windows2008r2": {
	    "iso_url": "http://download.microsoft.com/download/7/5/E/75EC4E54-5B02-42D6-8879-D8D3A25FBEF7/7601.17514.101119-1850_x64fre_server_eval_en-us-GRMSXEVAL_EN_DVD.iso",
	    "iso_checksum_type": "md5",
	    "iso_checksum": "4263be2cf3c59177c45085c0a7bc6ca5",
	    "virtualbox_guest_os_type": "Windows2008_64",
	    "vmware_guest_os_type": "windows7srv-64",
			"editions": {
				"standard": {
					"windows_image_name": "Windows Server 2008 R2 SERVERSTANDARD"
				},
        "enterprise": {
          "windows_image_name": "Windows Server 2008 R2 SERVERENTERPRISE"
        }
      }
	  }
	}
}
`

func TestCanLoadGlobalConfig(t *testing.T) {
	configuration := createConfiguration(t)
	if configuration.Headless {
		t.Error("Expected Gui to be false")
	}
	if configuration.WindowsUpdates {
		t.Error("Expected windows updates to be false")
	}
	if configuration.Communicator != "ssh" {
		t.Error("Expected communicator to be ssh")
	}
	if configuration.OutDir != "/tmp/foo/bar" {
		t.Error("Expected out dir to be /tmp/foo/bar")
	}
}

func TestCanListAllOSs(t *testing.T) {
	configuration := createConfiguration(t)
	os := configuration.List()
	sort.Strings(os)
	if len(os) != 2 {
		t.Errorf("Expected 2 OS entries, but got %d instead", len(os))
	} else {
		if os[0] != "windows10" {
			t.Errorf("Expected the first OS entry to be windows10, but was %s", os[0])
		}
		if os[1] != "windows2008r2" {
			t.Errorf("Expected the second OS entry to be windows2008r2, but was %s", os[0])
		}
	}
}

func TestCanLoadWindows10Config(t *testing.T) {
	configuration := createConfiguration(t)
	c, ok := configuration.Get("windows10")
	if !ok {
		t.Error("Failed to load Windows10 from the inductor configuration")
	} else {
		var tests = []struct {
			actual   string
			expected string
		}{
			{c.Name, "windows10"},
			{c.IsoURL, "http://care.dlservice.microsoft.com/dl/download/C/3/9/C399EEA8-135D-4207-92C9-6AAB3259F6EF/10240.16384.150709-1700.TH1_CLIENTENTERPRISEEVAL_OEMRET_X64FRE_EN-US.ISO"},
			{c.IsoChecksum, "56ab095075be28a90bc0b510835280975c6bb2ce"},
			{c.IsoChecksumType, "sha1"},
			{c.VirtualboxGuestOsType, "Windows81_64"},
			{c.VmwareGuestOsType, "windows8srv-64"},
			{c.Editions["enterprise"].WindowsImageName, "Windows 10 Enterprise Evaluation"},
		}
		for _, ts := range tests {
			if ts.actual != ts.expected {
				t.Logf("%#v", c)
				t.Errorf("Expected \"%s\" but got \"%s\"", ts.expected, ts.actual)
			}
		}
	}
}

func TestCanLoadWindows2008R2Config(t *testing.T) {
	configuration := createConfiguration(t)
	c, ok := configuration.Get("windows2008r2")
	if !ok {
		t.Error("Failed to load Windows 2008 R2 from the inductor metadata")
	} else {
		var tests = []struct {
			actual   string
			expected string
		}{
			{c.Name, "windows2008r2"},
			{c.IsoURL, "http://download.microsoft.com/download/7/5/E/75EC4E54-5B02-42D6-8879-D8D3A25FBEF7/7601.17514.101119-1850_x64fre_server_eval_en-us-GRMSXEVAL_EN_DVD.iso"},
			{c.IsoChecksum, "4263be2cf3c59177c45085c0a7bc6ca5"},
			{c.IsoChecksumType, "md5"},
			{c.VirtualboxGuestOsType, "Windows2008_64"},
			{c.VmwareGuestOsType, "windows7srv-64"},
			{c.Editions["standard"].WindowsImageName, "Windows Server 2008 R2 SERVERSTANDARD"},
			{c.Editions["enterprise"].WindowsImageName, "Windows Server 2008 R2 SERVERENTERPRISE"},
		}
		for _, ts := range tests {
			if ts.actual != ts.expected {
				t.Logf("%#v", c)
				t.Errorf("Expected \"%s\" but got \"%s\"", ts.expected, ts.actual)
			}
		}
	}
}

func TestMissingOS(t *testing.T) {
	configuration := createConfiguration(t)
	_, ok := configuration.Get("Linux")
	if ok {
		t.Error("Should have return !ok for Linux")
	}
}

func createConfiguration(t *testing.T) *InductorConfiguration {
	configuration, err := New(strings.NewReader(testData))
	if err != nil {
		t.Error("Failed to load the inductor configuration:", err)
	}
	return configuration
}
