package osregistry

import (
	"encoding/json"
	"io"
)

// OperatingSystemRegistry contains all OS details
type OperatingSystemRegistry struct {
	OperatingSystems map[string]OperatingSystem
}

// OperatingSystem has all the OS specific details required for Packer
type OperatingSystem struct {
	Name                  string
	IsoChecksum           string `json:"iso_checksum"`
	IsoChecksumType       string `json:"iso_checksum_type"`
	IsoURL                string `json:"iso_url"`
	VirtualboxGuestOsType string `json:"virtualbox_guest_os_type"`
	VmwareGuestOsType     string `json:"vmware_guest_os_type"`
	WindowsImageName      string `json:"windows_image_name"`
	ProductKey            string `json:"product_key"`
}

// List all available OS names
func (reg *OperatingSystemRegistry) List() []string {
	keys := make([]string, len(reg.OperatingSystems))
	i := 0
	for k := range reg.OperatingSystems {
		keys[i] = k
		i++
	}
	return keys
}

// Get returns the OS details for the named OS if found
func (reg *OperatingSystemRegistry) Get(osName string) (*OperatingSystem, bool) {
	os, ok := reg.OperatingSystems[osName]
	if ok {
		return &os, ok
	}
	return nil, ok
}

// New creates an initialized OS Registry instance
func New(osRegistry io.Reader) (*OperatingSystemRegistry, error) {
	registry := OperatingSystemRegistry{}
	registry.OperatingSystems = make(map[string]OperatingSystem)

	// decode outermost map of values keyed by OS name
	var objmap map[string]*json.RawMessage
	dec := json.NewDecoder(osRegistry)
	if err := dec.Decode(&objmap); err != nil {
		return nil, err
	}

	// decode each entry
	for k, v := range objmap {
		var os OperatingSystem
		err := json.Unmarshal(*v, &os)
		if err != nil {
			return nil, err
		}
		os.Name = k
		registry.OperatingSystems[k] = os
	}

	return &registry, nil
}
