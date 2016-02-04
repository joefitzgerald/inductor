package renderer

import (
	"os"

	"github.com/joefitzgerald/inductor/osregistry"
)

// RenderOptions for packer.json and Autounattend.xml
type RenderOptions struct {
	OSName                string
	ProductKey            string
	WindowsImageName      string
	VirtualboxGuestOsType string
	VmwareGuestOsType     string
	IsoURL                string
	IsoChecksumType       string
	IsoChecksum           string
	Communicator          string
	Username              string
	Password              string
	DiskSize              uint32
	RAM                   uint32
	CPU                   uint8
	Headless              bool
	WindowsUpdates        bool
}

// NewRenderOptionsWithOverrides creates render options using the base OS
// registry with any provided overrides.
func NewRenderOptionsWithOverrides(windowsEdition string, osregistryFilePath string) (*RenderOptions, error) {
	file, err := os.Open(osregistryFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r, err := osregistry.New(windowsEdition, file)
	if err != nil {
		return nil, err
	}

	// default all rendering options to values in the OS registry
	opts := NewRenderOptions()
	opts.OSName = r.Name
	opts.IsoChecksum = r.IsoChecksum
	opts.IsoChecksumType = r.IsoChecksumType
	opts.IsoURL = r.IsoURL
	opts.VirtualboxGuestOsType = r.VirtualboxGuestOsType
	opts.VmwareGuestOsType = r.VmwareGuestOsType
	opts.WindowsImageName = r.WindowsImageName

	return opts, nil
}

// NewRenderOptions creates a new ready to use RenderOptions instance which
// defaults to Windows10 trial values
func NewRenderOptions() *RenderOptions {
	ro := &RenderOptions{
		OSName:                "windows10",
		ProductKey:            "",
		WindowsImageName:      "Windows 10 Enterprise Evaluation",
		VirtualboxGuestOsType: "Windows81_64",
		VmwareGuestOsType:     "windows8srv-64",
		IsoURL:                "http://care.dlservice.microsoft.com/dl/download/C/3/9/C399EEA8-135D-4207-92C9-6AAB3259F6EF/10240.16384.150709-1700.TH1_CLIENTENTERPRISEEVAL_OEMRET_X64FRE_EN-US.ISO",
		IsoChecksumType:       "sha1",
		IsoChecksum:           "56ab095075be28a90bc0b510835280975c6bb2ce",
		Communicator:          "ssh",
		Username:              "vagrant",
		Password:              "vagrant",
		DiskSize:              61400,
		RAM:                   2048,
		CPU:                   2,
		Headless:              true,
		WindowsUpdates:        true,
	}
	return ro
}
