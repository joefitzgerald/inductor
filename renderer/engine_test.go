package renderer

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestCanBuildDefaultWindows10Config(t *testing.T) {
	var packer, autounattend, vagrantfile bytes.Buffer
	opts := NewRenderOptions()
	tpl := NewPackerTemplate()
	err := tpl.Render(opts, &packer, &autounattend, &vagrantfile)
	if err != nil {
		t.Error("Failed to render the default Windows10 template:", err)
	} else {
		if !strings.Contains(packer.String(), "\"ssh_username\": \"vagrant\",") {
			t.Error("Expected packer.json to contain: \"ssh_username\": \"vagrant\",")
		}
		if !strings.Contains(autounattend.String(), "<Label>windows10</Label>") {
			t.Error("Expected Autounattend.xml to contain: <Label>windows10</Label>")
		}
		if !strings.Contains(vagrantfile.String(), "config.vm.box = \"windows10\"") {
			t.Error("Expected Vagrantfile to contain: config.vm.box = \"windows10\"")
		}
	}
}

func TestAutounattendValuesArePopulated(t *testing.T) {
	var packer, autounattend, vagrantfile bytes.Buffer
	opts := NewRenderOptions()
	opts.ProductKey = "FOOBAR-7Y6KF-2VJC9-XBBR8-HVTHH"
	tpl := NewPackerTemplate()
	err := tpl.Render(opts, &packer, &autounattend, &vagrantfile)
	if err != nil {
		t.Error("Failed to render the default Windows10 template:", err)
	} else {
		var expected = [...]string{
			fmt.Sprintf("<Key>%s</Key>", opts.ProductKey),
			fmt.Sprintf("<FullName>%s</FullName>", opts.Username),
			fmt.Sprintf("<Value>%s</Value>", opts.WindowsImageName),
			fmt.Sprintf("<Username>%s</Username>", opts.Username),
			"<ComputerName>vagrant-win10</ComputerName>",
		}
		for _, e := range expected {
			if !strings.Contains(autounattend.String(), e) {
				t.Log(autounattend.String())
				t.Errorf("Expected Autounattend.xml to contain: %s", e)
			}
		}
	}
}

func TestPackerJSONValuesArePopulated(t *testing.T) {
	var packer, autounattend, vagrantfile bytes.Buffer
	opts := NewRenderOptions()
	tpl := NewPackerTemplate()
	err := tpl.Render(opts, &packer, &autounattend, &vagrantfile)
	if err != nil {
		t.Error("Failed to render the default Windows10 template:", err)
	} else {
		var expected = [...]string{
			fmt.Sprintf("\"iso_url\": \"%s\",", opts.IsoURL),
			fmt.Sprintf("\"iso_checksum_type\": \"%s\",", opts.IsoChecksumType),
			fmt.Sprintf("\"iso_checksum\": \"%s\",", opts.IsoChecksum),
			fmt.Sprintf("\"headless\": %s,", strconv.FormatBool(opts.Headless)),
			fmt.Sprintf("\"ssh_username\": \"%s\",", opts.Username),
			fmt.Sprintf("\"ssh_password\": \"%s\",", opts.Password),
			fmt.Sprintf("\"guest_os_type\": \"%s\",", opts.VmwareGuestOsType),
			fmt.Sprintf("\"guest_os_type\": \"%s\",", opts.VirtualboxGuestOsType),
			fmt.Sprintf("\"output\": \"%s_{{.Provider}}.box\",", opts.OSName),
		}
		for _, e := range expected {
			if !strings.Contains(packer.String(), e) {
				t.Log(packer.String())
				t.Errorf("Expected packer.json to contain: %s", e)
			}
		}
	}
}

func TestCanSkipWindowsUpdates(t *testing.T) {
	var packer, autounattend, vagrantfile bytes.Buffer
	opts := NewRenderOptions()
	opts.WindowsUpdates = false
	tpl := NewPackerTemplate()
	err := tpl.Render(opts, &packer, &autounattend, &vagrantfile)
	if err != nil {
		t.Error("Failed to render the default Windows10 template:", err)
	} else {
		if strings.Contains(autounattend.String(), "a:\\win-updates.ps1") {
			t.Log(autounattend.String())
			t.Error("Windows updates were not skipped!")
		}
		if !strings.Contains(autounattend.String(), "a:\\openssh.ps1 -AutoStart") {
			t.Log(autounattend.String())
			t.Error("SSH was not started when Windows updates were skipped")
		}
	}
}

func TestShouldStartSSHWhenWindowsUpdatesAreSkipped(t *testing.T) {
	var packer, autounattend, vagrantfile bytes.Buffer
	opts := NewRenderOptions()
	opts.WindowsUpdates = true
	tpl := NewPackerTemplate()
	err := tpl.Render(opts, &packer, &autounattend, &vagrantfile)
	if err != nil {
		t.Error("Failed to render the default Windows10 template:", err)
	} else {
		if !strings.Contains(autounattend.String(), "a:\\win-updates.ps1") {
			t.Log(autounattend.String())
			t.Error("Windows updates were not applied!")
		}
		if strings.Contains(autounattend.String(), "a:\\openssh.ps1 -AutoStart") {
			t.Log(autounattend.String())
			t.Error("SSH was auto started when Windows updates were applied")
		}
	}
}