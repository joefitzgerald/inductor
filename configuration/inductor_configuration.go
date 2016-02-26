package configuration

import (
	"encoding/json"
	"io"
)

// InductorConfiguration contains all OS details
type InductorConfiguration struct {
	Gui              bool   `json:"gui"`
	WindowsUpdates   bool   `json:"windows_updates"`
	Communicator     string `json:"communicator"`
	OutDir           string `json:"out_dir"`
	OperatingSystems map[string]OperatingSystem
}

// OperatingSystem has all the OS specific details required for Packer
type OperatingSystem struct {
	Name                  string
	IsoChecksum           string             `json:"iso_checksum"`
	IsoChecksumType       string             `json:"iso_checksum_type"`
	IsoURL                string             `json:"iso_url"`
	VirtualboxGuestOsType string             `json:"virtualbox_guest_os_type"`
	VmwareGuestOsType     string             `json:"vmware_guest_os_type"`
	Editions              map[string]Edition `json:"editions"`
}

// Edition is the Windows edition, e.g. Enterprise, Home
type Edition struct {
	WindowsImageName string `json:"windows_image_name"`
	ProductKey       string `json:"product_key"`
}

// List all available OS names
func (reg *InductorConfiguration) List() []string {
	keys := make([]string, len(reg.OperatingSystems))
	i := 0
	for k := range reg.OperatingSystems {
		keys[i] = k
		i++
	}
	return keys
}

// Get returns the OS details for the named OS if found
func (reg *InductorConfiguration) Get(osName string) (*OperatingSystem, bool) {
	os, ok := reg.OperatingSystems[osName]
	if ok {
		return &os, ok
	}
	return nil, ok
}

// New creates an initialized InductorConfiguration
func New(configSrc io.Reader) (*InductorConfiguration, error) {
	configuration := InductorConfiguration{
		Gui:              false,
		WindowsUpdates:   true,
		Communicator:     "winrm",
		OutDir:           "out",
		OperatingSystems: make(map[string]OperatingSystem),
	}

	// decode outermost defaults and operations_systems into a map
	var topObjMap map[string]*json.RawMessage
	dec := json.NewDecoder(configSrc)
	if err := dec.Decode(&topObjMap); err != nil {
		return nil, err
	}

	// decode global config into main config
	err := json.Unmarshal(*topObjMap["config"], &configuration)
	if err != nil {
		return nil, err
	}

	// decode operating_systems into another map keyed by OS name
	var osObjMap map[string]*json.RawMessage
	err = json.Unmarshal(*topObjMap["operating_systems"], &osObjMap)
	if err != nil {
		return nil, err
	}

	// decode each OS each entry
	for k, v := range osObjMap {
		var os OperatingSystem
		err := json.Unmarshal(*v, &os)
		if err != nil {
			return nil, err
		}
		os.Name = k
		configuration.OperatingSystems[k] = os
	}

	return &configuration, nil
}
