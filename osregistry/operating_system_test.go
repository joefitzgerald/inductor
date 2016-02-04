package osregistry

import (
	"strings"
	"testing"
)

var testData = `
{
  "windows10": {
    "iso_url": "http://care.dlservice.microsoft.com/dl/download/C/3/9/C399EEA8-135D-4207-92C9-6AAB3259F6EF/10240.16384.150709-1700.TH1_CLIENTENTERPRISEEVAL_OEMRET_X64FRE_EN-US.ISO",
    "iso_checksum_type": "sha1",
    "iso_checksum": "56ab095075be28a90bc0b510835280975c6bb2ce",
    "windows_image_name": "Windows 10 Enterprise Evaluation",
    "virtualbox_guest_os_type": "Windows81_64",
    "vmware_guest_os_type": "windows8srv-64"
  },
  "windows2008r2": {
    "iso_url": "http://download.microsoft.com/download/7/5/E/75EC4E54-5B02-42D6-8879-D8D3A25FBEF7/7601.17514.101119-1850_x64fre_server_eval_en-us-GRMSXEVAL_EN_DVD.iso",
    "iso_checksum_type": "md5",
    "iso_checksum": "4263be2cf3c59177c45085c0a7bc6ca5",
    "windows_image_name": "Windows Server 2008 R2 SERVERSTANDARD",
    "virtualbox_guest_os_type": "Windows2008_64",
    "vmware_guest_os_type": "windows7srv-64"
  }
}
`

func TestCanLoadWindows10Config(t *testing.T) {
	r, err := New("Windows10", strings.NewReader(testData))
	if err != nil {
		t.Error("Failed to load the Windows OS registry for Windows 10:", err)
	} else {
		var tests = []struct {
			actual   string
			expected string
		}{
			{r.Name, "windows10"},
			{r.IsoURL, "http://care.dlservice.microsoft.com/dl/download/C/3/9/C399EEA8-135D-4207-92C9-6AAB3259F6EF/10240.16384.150709-1700.TH1_CLIENTENTERPRISEEVAL_OEMRET_X64FRE_EN-US.ISO"},
			{r.IsoChecksum, "56ab095075be28a90bc0b510835280975c6bb2ce"},
			{r.IsoChecksumType, "sha1"},
			{r.VirtualboxGuestOsType, "Windows81_64"},
			{r.VmwareGuestOsType, "windows8srv-64"},
			{r.WindowsImageName, "Windows 10 Enterprise Evaluation"},
		}
		for _, ts := range tests {
			if ts.actual != ts.expected {
				t.Logf("%#v", r)
				t.Errorf("Expected \"%s\" but got \"%s\"", ts.expected, ts.actual)
			}
		}
	}
}

func TestCanLoadWindows2008R2Config(t *testing.T) {
	r, err := New("Windows2008R2", strings.NewReader(testData))
	if err != nil {
		t.Error("Failed to load the Windows OS registry for Windows 2008 R2:", err)
	} else {
		var tests = []struct {
			actual   string
			expected string
		}{
			{r.Name, "windows2008r2"},
			{r.IsoURL, "http://download.microsoft.com/download/7/5/E/75EC4E54-5B02-42D6-8879-D8D3A25FBEF7/7601.17514.101119-1850_x64fre_server_eval_en-us-GRMSXEVAL_EN_DVD.iso"},
			{r.IsoChecksum, "4263be2cf3c59177c45085c0a7bc6ca5"},
			{r.IsoChecksumType, "md5"},
			{r.VirtualboxGuestOsType, "Windows2008_64"},
			{r.VmwareGuestOsType, "windows7srv-64"},
			{r.WindowsImageName, "Windows Server 2008 R2 SERVERSTANDARD"},
		}
		for _, ts := range tests {
			if ts.actual != ts.expected {
				t.Logf("%#v", r)
				t.Errorf("Expected \"%s\" but got \"%s\"", ts.expected, ts.actual)
			}
		}
	}
}

func TestMissingOS(t *testing.T) {
	_, err := New("Linux", strings.NewReader(testData))
	if err == nil || err.Error() != "Could not find Linux" {
		t.Error("Loading Linux should have returned an error")
	}
}
