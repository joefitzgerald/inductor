package renderer

import "github.com/joefitzgerald/inductor/osregistry"

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
func NewRenderOptionsWithOverrides(os *osregistry.OperatingSystem) (*RenderOptions, error) {
	// default all rendering options to values in the OS registry
	opts := NewRenderOptions()
	opts.OSName = os.Name
	opts.IsoChecksum = os.IsoChecksum
	opts.IsoChecksumType = os.IsoChecksumType
	opts.IsoURL = os.IsoURL
	opts.VirtualboxGuestOsType = os.VirtualboxGuestOsType
	opts.VmwareGuestOsType = os.VmwareGuestOsType
	opts.WindowsImageName = os.WindowsImageName
	opts.ProductKey = os.ProductKey
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
		Communicator:          "winrm",
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
