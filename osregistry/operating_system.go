package osregistry

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// OperatingSystem has all the OS specific details required for Packer
type OperatingSystem struct {
	Name                  string
	IsoChecksum           string `json:"iso_checksum"`
	IsoChecksumType       string `json:"iso_checksum_type"`
	IsoURL                string `json:"iso_url"`
	VirtualboxGuestOsType string `json:"virtualbox_guest_os_type"`
	VmwareGuestOsType     string `json:"vmware_guest_os_type"`
	WindowsImageName      string `json:"windows_image_name"`
}

// OperatingSystems contains all configuration data loaded from disk
type OperatingSystems struct {
	All map[string]OperatingSystem
}

// New creates a fully initialized instance of an OperatingSystem using the
// stored data on disk.
func New(osName string, osRegistry io.Reader) (*OperatingSystem, error) {
	osKey := strings.ToLower(osName)

	// decode outermost map of values keyed by OS name
	var objmap map[string]*json.RawMessage
	dec := json.NewDecoder(osRegistry)
	if err := dec.Decode(&objmap); err != nil {
		return nil, err
	}

	// find the specified OS
	val, ok := objmap[osKey]
	if !ok {
		return nil, fmt.Errorf("Could not find %s", osName)
	}

	// decode the specified OS into an OperatingSystem object
	var os OperatingSystem
	err := json.Unmarshal(*val, &os)
	if err != nil {
		return nil, err
	}

	// manually populate the name since its the map key
	os.Name = osKey

	return &os, nil
}
