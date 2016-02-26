package renderer

import (
	"fmt"

	"github.com/joefitzgerald/inductor/configuration"
)

// RenderOptions for packer.json and Autounattend.xml
type RenderOptions struct {
	OSName                string
	Edition               string
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

// NewRenderOptions creates render options using the base OS
// registry with any provided overrides.
func NewRenderOptions(osname string, edition string, config *configuration.InductorConfiguration) (*RenderOptions, error) {
	os, ok := config.Get(osname)
	if !ok {
		return nil, fmt.Errorf("Couldn't find OS configuration for '%s'", osname)
	}

	// set global config
	opts := NewDefaultRenderOptions()
	opts.Communicator = config.Communicator
	opts.Headless = config.Headless
	opts.WindowsUpdates = config.WindowsUpdates
	// TODO: Username, Password, DiskSize, RAM, CPU

	// default all rendering options to values in the OS registry
	opts.OSName = os.Name
	opts.IsoChecksum = os.IsoChecksum
	opts.IsoChecksumType = os.IsoChecksumType
	opts.IsoURL = os.IsoURL
	opts.VirtualboxGuestOsType = os.VirtualboxGuestOsType
	opts.VmwareGuestOsType = os.VmwareGuestOsType
	// TODO: Username, Password, DiskSize, RAM, CPU per OS

	// edition specific attributes
	if len(edition) == 0 {
		for k := range os.Editions {
			edition = k
			break
		}
	}

	// TODO: validate that the edition is valid
	opts.Edition = edition
	opts.WindowsImageName = os.Editions[edition].WindowsImageName
	opts.ProductKey = os.Editions[edition].ProductKey

	return opts, nil
}

// NewDefaultRenderOptions creates a new ready to use RenderOptions instance which
// defaults to Windows10 trial values
func NewDefaultRenderOptions() *RenderOptions {
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
