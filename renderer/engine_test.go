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
	opts := NewDefaultRenderOptions()
	tpl := NewPackerTemplate()
	err := tpl.Render(opts, &packer, &autounattend, &vagrantfile)
	if err != nil {
		t.Error("Failed to render the default Windows10 template:", err)
	} else {
		if !strings.Contains(packer.String(), "\"winrm_username\": \"vagrant\",") {
			t.Error("Expected packer.json to contain: \"winrm_username\": \"vagrant\",")
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
	opts := NewDefaultRenderOptions()
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
	opts := NewDefaultRenderOptions()
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
			fmt.Sprintf("\"winrm_username\": \"%s\",", opts.Username),
			"\"communicator\": \"winrm\"",
			fmt.Sprintf("\"winrm_password\": \"%s\",", opts.Password),
			fmt.Sprintf("\"guest_os_type\": \"%s\",", opts.VmwareGuestOsType),
			fmt.Sprintf("\"guest_os_type\": \"%s\",", opts.VirtualboxGuestOsType),
			fmt.Sprintf("\"output\": \"%s_{{.Provider}}.box\",", opts.OSName),
			"\"type\": \"windows-shell\"",
		}
		for _, e := range expected {
			if !strings.Contains(packer.String(), e) {
				t.Log(packer.String())
				t.Errorf("Expected packer.json to contain: %s", e)
			}
		}
	}
}

func TestPackerJSONValuesArePopulatedWhenUsingSSH(t *testing.T) {
	var packer, autounattend, vagrantfile bytes.Buffer
	opts := NewDefaultRenderOptions()
	opts.Communicator = "ssh"
	tpl := NewPackerTemplate()
	err := tpl.Render(opts, &packer, &autounattend, &vagrantfile)
	if err != nil {
		t.Error("Failed to render the default Windows10 template:", err)
	} else {
		var expected = [...]string{
			fmt.Sprintf("\"ssh_username\": \"%s\",", opts.Username),
			fmt.Sprintf("\"ssh_password\": \"%s\",", opts.Password),
			"\"type\": \"shell\"",
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
	opts := NewDefaultRenderOptions()
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
		if !strings.Contains(autounattend.String(), "a:\\winrm.ps1") {
			t.Log(autounattend.String())
			t.Error("WinRM was not started when Windows updates were skipped")
		}
		if strings.Contains(autounattend.String(), "a:\\openssh.ps1") {
			t.Log(autounattend.String())
			t.Error("OpenSSH should not be started with the WinRM communicator")
		}
	}
}

func TestCanSkipWindowsUpdatesWhenUsingSSH(t *testing.T) {
	var packer, autounattend, vagrantfile bytes.Buffer
	opts := NewDefaultRenderOptions()
	opts.WindowsUpdates = false
	opts.Communicator = "ssh"
	tpl := NewPackerTemplate()
	err := tpl.Render(opts, &packer, &autounattend, &vagrantfile)
	if err != nil {
		t.Error("Failed to render the default Windows10 template:", err)
	} else {
		if strings.Contains(autounattend.String(), "a:\\win-updates.ps1") {
			t.Log(autounattend.String())
			t.Error("Windows updates were not skipped!")
		}
		if !strings.Contains(autounattend.String(), "a:\\winrm.ps1") {
			t.Log(autounattend.String())
			t.Error("WinRM was not started when Windows updates were skipped")
		}
		if !strings.Contains(autounattend.String(), "a:\\openssh.ps1") {
			t.Log(autounattend.String())
			t.Error("OpenSSH was not started when Windows updates were skipped")
		}
	}
}

func TestShouldInstallWindowsUpdates(t *testing.T) {
	var packer, autounattend, vagrantfile bytes.Buffer
	opts := NewDefaultRenderOptions()
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
		if strings.Contains(autounattend.String(), "a:\\winrm.ps1") {
			t.Log(autounattend.String())
			t.Error("WinRM was started before Windows updates were applied")
		}
		if strings.Contains(autounattend.String(), "a:\\openssh.ps1") {
			t.Log(autounattend.String())
			t.Error("OpenSSH should not be started with the WinRM communicator")
		}
	}
}

func TestShouldInstallWindowsUpdatesWhenUsingSSH(t *testing.T) {
	var packer, autounattend, vagrantfile bytes.Buffer
	opts := NewDefaultRenderOptions()
	opts.WindowsUpdates = true
	opts.Communicator = "ssh"
	tpl := NewPackerTemplate()
	err := tpl.Render(opts, &packer, &autounattend, &vagrantfile)
	if err != nil {
		t.Error("Failed to render the default Windows10 template:", err)
	} else {
		if !strings.Contains(autounattend.String(), "a:\\win-updates.ps1") {
			t.Log(autounattend.String())
			t.Error("Windows updates were not applied!")
		}
		if strings.Contains(autounattend.String(), "a:\\winrm.ps1") {
			t.Log(autounattend.String())
			t.Error("WinRM was started before Windows updates were applied")
		}
		if strings.Contains(autounattend.String(), "a:\\openssh.ps1") {
			t.Log(autounattend.String())
			t.Error("OpenSSH was started before Windows updates were applied")
		}
	}
}

// NewPackerTemplate creates a new fully functioning PackerTemplate with
// hardcoded default templates.
func NewPackerTemplate() *PackerTemplate {
	pt := &PackerTemplate{}
	pt.PackerTpl = `{
  "builders": [
    {
      "type": "vmware-iso",
      "iso_url": "{{.IsoURL}}",
      "iso_checksum_type": "{{.IsoChecksumType}}",
      "iso_checksum": "{{.IsoChecksum}}",
      "headless": {{.Headless}},
      "boot_wait": "2m",
      {{ if eq .Communicator "ssh" }}
			"ssh_username": "{{.Username}}",
      "ssh_password": "{{.Password}}",
			"ssh_wait_timeout": "8h",
			{{ else }}
			"communicator": "winrm",
			"winrm_username": "{{.Username}}",
      "winrm_password": "{{.Password}}",
			"winrm_timeout": "8h",
			{{ end }}
      "shutdown_command": "shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"",
      "guest_os_type": "{{.VmwareGuestOsType}}",
      "tools_upload_flavor": "windows",
      "disk_size": {{.DiskSize}},
      "vnc_port_min": 5900,
      "vnc_port_max": 5980,
      "floppy_files": [
        "./Autounattend.xml",
        "./scripts/hotfix-KB3102810.bat",
        "./scripts/fixnetwork.ps1",
        "./scripts/microsoft-updates.bat",
        "./scripts/win-updates.ps1",
        "./scripts/openssh.ps1",
        "./scripts/winrm.ps1"
      ],
      "vmx_data": {
        "RemoteDisplay.vnc.enabled": "false",
        "RemoteDisplay.vnc.port": "5900",
        "memsize": "{{.RAM}}",
        "numvcpus": "{{.CPU}}",
        "scsi0.virtualDev": "lsisas1068"
      }
    },
    {
      "type": "virtualbox-iso",
      "iso_url": "{{.IsoURL}}",
      "iso_checksum_type": "{{.IsoChecksumType}}",
      "iso_checksum": "{{.IsoChecksum}}",
      "headless": {{.Headless}},
      "boot_wait": "2m",
      {{ if eq .Communicator "ssh" }}
			"ssh_username": "{{.Username}}",
      "ssh_password": "{{.Password}}",
			"ssh_wait_timeout": "8h",
			{{ else }}
			"communicator": "winrm",
			"winrm_username": "{{.Username}}",
      "winrm_password": "{{.Password}}",
			"winrm_timeout": "8h",
			{{ end }}
      "shutdown_command": "shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"",
      "guest_os_type": "{{.VirtualboxGuestOsType}}",
      "disk_size": {{.DiskSize}},
      "floppy_files": [
        "./Autounattend.xml",
        "./scripts/hotfix-KB3102810.bat",
        "./scripts/fixnetwork.ps1",
        "./scripts/microsoft-updates.bat",
        "./scripts/win-updates.ps1",
        "./scripts/openssh.ps1",
        "./scripts/winrm.ps1",
        "./scripts/oracle-cert.cer"
      ],
      "vboxmanage": [
        [
          "modifyvm",
          "{{"{{"}}.Name{{"}}"}}",
          "--memory",
          "{{.RAM}}"
        ],
        [
          "modifyvm",
          "{{"{{"}}.Name{{"}}"}}",
          "--cpus",
          "{{.CPU}}"
        ]
      ]
    }
  ],
  "provisioners": [
    {{ if eq .Communicator "ssh" }}
    {
      "type": "shell",
      "remote_path": "/tmp/script.bat",
      "execute_command": "{{"{{"}}.Vars{{"}}"}} cmd /c C:/Windows/Temp/script.bat",
      "scripts": [
        "./scripts/vm-guest-tools.bat",
        "./scripts/vagrant-ssh.bat",
        "./scripts/disable-auto-logon.bat",
        "./scripts/enable-rdp.bat",
        "./scripts/compile-dotnet-assemblies.bat",
        "./scripts/compact.bat"
      ]
    }
    {{ else }}
    {
      "type": "windows-shell",
      "scripts": [
        "./scripts/vm-guest-tools.bat",
        "./scripts/disable-auto-logon.bat",
        "./scripts/enable-rdp.bat",
        "./scripts/compile-dotnet-assemblies.bat",
        "./scripts/compact.bat"
      ]
    }
    {{ end }}
  ],
  "post-processors": [
    {
      "type": "vagrant",
      "keep_input_artifact": false,
      "output": "{{.OSName}}_{{"{{"}}.Provider{{"}}"}}.box",
      "vagrantfile_template": "Vagrantfile"
    }
  ]
}
`
	pt.AutounattendTpl = `<?xml version="1.0" encoding="utf-8"?>
<unattend xmlns="urn:schemas-microsoft-com:unattend">
    <servicing/>
    <settings pass="windowsPE">
        <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <DiskConfiguration>
                <Disk wcm:action="add">
                    <CreatePartitions>
                        <CreatePartition wcm:action="add">
                            <Order>1</Order>
                            <Type>Primary</Type>
                            <Extend>true</Extend>
                        </CreatePartition>
                    </CreatePartitions>
                    <ModifyPartitions>
                        <ModifyPartition wcm:action="add">
                            <Extend>false</Extend>
                            <Format>NTFS</Format>
                            <Letter>C</Letter>
                            <Order>1</Order>
                            <PartitionID>1</PartitionID>
                            <Label>{{.OSName}}</Label>
                        </ModifyPartition>
                    </ModifyPartitions>
                    <DiskID>0</DiskID>
                    <WillWipeDisk>true</WillWipeDisk>
                </Disk>
                <WillShowUI>OnError</WillShowUI>
            </DiskConfiguration>
            <UserData>
                <AcceptEula>true</AcceptEula>
                <FullName>{{.Username}}</FullName>
                <Organization></Organization>
                <ProductKey>
                    {{ if .ProductKey }}
                    <Key>{{.ProductKey}}</Key>
                    {{ end }}
                    <WillShowUI>Never</WillShowUI>
                </ProductKey>
            </UserData>
            <ImageInstall>
                <OSImage>
                    <InstallTo>
                        <DiskID>0</DiskID>
                        <PartitionID>1</PartitionID>
                    </InstallTo>
                    <WillShowUI>OnError</WillShowUI>
                    <InstallToAvailablePartition>false</InstallToAvailablePartition>
                    <InstallFrom>
                        <MetaData wcm:action="add">
                            <Key>/IMAGE/NAME</Key>
                            <Value>{{.WindowsImageName}}</Value>
                        </MetaData>
                    </InstallFrom>
                </OSImage>
            </ImageInstall>
        </component>
        <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-International-Core-WinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <SetupUILanguage>
                <UILanguage>en-US</UILanguage>
            </SetupUILanguage>
            <InputLocale>en-US</InputLocale>
            <SystemLocale>en-US</SystemLocale>
            <UILanguage>en-US</UILanguage>
            <UILanguageFallback>en-US</UILanguageFallback>
            <UserLocale>en-US</UserLocale>
        </component>
    </settings>
    <settings pass="offlineServicing">
        <component name="Microsoft-Windows-LUA-Settings" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <EnableLUA>false</EnableLUA>
        </component>
    </settings>
    <settings pass="oobeSystem">
        <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <UserAccounts>
                <AdministratorPassword>
                    <Value>{{.Password}}</Value>
                    <PlainText>true</PlainText>
                </AdministratorPassword>
                <LocalAccounts>
                    <LocalAccount wcm:action="add">
                        <Password>
                            <Value>{{.Password}}</Value>
                            <PlainText>true</PlainText>
                        </Password>
                        <Description>Vagrant User</Description>
                        <DisplayName>vagrant</DisplayName>
                        <Group>administrators</Group>
                        <Name>{{.Username}}</Name>
                    </LocalAccount>
                </LocalAccounts>
            </UserAccounts>
            <OOBE>
                <HideEULAPage>true</HideEULAPage>
                <HideWirelessSetupInOOBE>true</HideWirelessSetupInOOBE>
                <NetworkLocation>Home</NetworkLocation>
                <ProtectYourPC>1</ProtectYourPC>
            </OOBE>
            <AutoLogon>
                <Password>
                    <Value>{{.Password}}</Value>
                    <PlainText>true</PlainText>
                </Password>
                <Username>{{.Username}}</Username>
                <Enabled>true</Enabled>
            </AutoLogon>
            <FirstLogonCommands>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c powershell -Command "Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Force"</CommandLine>
                    <Description>Set Execution Policy 64 Bit</Description>
                    <Order>1</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>C:\Windows\SysWOW64\cmd.exe /c powershell -Command "Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Force"</CommandLine>
                    <Description>Set Execution Policy 32 Bit</Description>
                    <Order>2</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c reg add "HKLM\System\CurrentControlSet\Control\Network\NewNetworkWindowOff"</CommandLine>
                    <Description>Network prompt</Description>
                    <Order>3</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe -File a:\fixnetwork.ps1</CommandLine>
                    <Description>Fix public network</Description>
                    <Order>4</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced\ /v HideFileExt /t REG_DWORD /d 0 /f</CommandLine>
                    <Order>18</Order>
                    <Description>Show file extensions in Explorer</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKCU\Console /v QuickEdit /t REG_DWORD /d 1 /f</CommandLine>
                    <Order>19</Order>
                    <Description>Enable QuickEdit mode</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced\ /v Start_ShowRun /t REG_DWORD /d 1 /f</CommandLine>
                    <Order>20</Order>
                    <Description>Show Run command in Start Menu</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced\ /v StartMenuAdminTools /t REG_DWORD /d 1 /f</CommandLine>
                    <Order>21</Order>
                    <Description>Show Administrative Tools in Start Menu</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKLM\SYSTEM\CurrentControlSet\Control\Power\ /v HibernateFileSizePercent /t REG_DWORD /d 0 /f</CommandLine>
                    <Order>22</Order>
                    <Description>Zero Hibernation File</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKLM\SYSTEM\CurrentControlSet\Control\Power\ /v HibernateEnabled /t REG_DWORD /d 0 /f</CommandLine>
                    <Order>23</Order>
                    <Description>Disable Hibernation Mode</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c wmic useraccount where "name='vagrant'" set PasswordExpires=FALSE</CommandLine>
                    <Order>24</Order>
                    <Description>Disable password expiration for vagrant user</Description>
                </SynchronousCommand>
                {{ if .WindowsUpdates }}
                <!-- Fix high CPU utilization on Windows7 when installing updates -->
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c a:\hotfix-KB3102810.bat</CommandLine>
                    <Order>98</Order>
                    <Description>KB3102810</Description>
                </SynchronousCommand>
                <!-- Include non-Windows MS updates -->
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c a:\microsoft-updates.bat</CommandLine>
                    <Order>99</Order>
                    <Description>Enable Microsoft Updates</Description>
                </SynchronousCommand>
                <!-- Install Windows Updates, win-updates.ps1 will start SSH/WinRM when done -->
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe -File a:\win-updates.ps1 -Communicator {{.Communicator}}</CommandLine>
                    <Description>Install Windows Updates</Description>
                    <Order>100</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                {{ else }}
                <!-- Skipping Windows Updates, directly start SSH/WinRM -->
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe -File a:\winrm.ps1</CommandLine>
                    <Description>Configure and start WinRM</Description>
                    <Order>99</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                {{ if eq .Communicator "ssh" }}
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe -File a:\openssh.ps1 -AutoStart</CommandLine>
                    <Description>Install OpenSSH</Description>
                    <Order>100</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                {{ end }}
                {{ end }}
            </FirstLogonCommands>
            <ShowWindowsLive>false</ShowWindowsLive>
        </component>
    </settings>
    <settings pass="specialize">
        <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <OEMInformation>
                <HelpCustomized>false</HelpCustomized>
            </OEMInformation>
            <ComputerName>{{ SafeComputerName ( printf "vagrant-%s" ( Replace .OSName "windows" "win" -1 )) }}</ComputerName>
            <TimeZone>Pacific Standard Time</TimeZone>
            <RegisteredOwner/>
        </component>
        <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Security-SPP-UX" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <SkipAutoActivation>true</SkipAutoActivation>
        </component>
    </settings>
    <cpi:offlineImage xmlns:cpi="urn:schemas-microsoft-com:cpi" cpi:source="catalog:d:/sources/install_windows 7 ENTERPRISE.clg"/>
</unattend>
`
	pt.VagrantfileTpl = `# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.require_version ">= 1.6.2"

Vagrant.configure("2") do |config|
	config.vm.box = "{{.OSName}}"
	config.vm.communicator = "winrm"

	# Admin user name and password
	config.winrm.username = "{{.Username}}"
	config.winrm.password = "{{.Password}}"

	config.vm.network :forwarded_port, guest: 3389, host: 3389, id: "rdp", auto_correct: true
	config.vm.network :forwarded_port, guest: 22, host: 2222, id: "ssh", auto_correct: true

	config.vm.provider :virtualbox do |v, override|
		v.customize ["modifyvm", :id, "--memory", {{.RAM}}]
		v.customize ["modifyvm", :id, "--cpus", {{.CPU}}]
		v.customize ["setextradata", "global", "GUI/SuppressMessages", "all" ]
	end

  config.vm.provider :vmware_fusion do |v, override|
		v.vmx["memsize"] = "{{.RAM}}"
		v.vmx["numvcpus"] = "{{.CPU}}"
		v.vmx["ethernet0.virtualDev"] = "vmxnet3"
		v.vmx["RemoteDisplay.vnc.enabled"] = "false"
		v.vmx["RemoteDisplay.vnc.port"] = "5900"
		v.vmx["scsi0.virtualDev"] = "lsisas1068"
	end

	config.vm.provider :vmware_workstation do |v, override|
		v.vmx["memsize"] = "{{.RAM}}"
		v.vmx["numvcpus"] = "{{.CPU}}"
		v.vmx["ethernet0.virtualDev"] = "vmxnet3"
		v.vmx["RemoteDisplay.vnc.enabled"] = "false"
		v.vmx["RemoteDisplay.vnc.port"] = "5900"
		v.vmx["scsi0.virtualDev"] = "lsisas1068"
	end
end
`
	return pt
}
