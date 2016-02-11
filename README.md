[![Build Status](https://travis-ci.org/joefitzgerald/inductor.svg?branch=master)](https://travis-ci.org/joefitzgerald/inductor)

# Inductor

Inductor is a command line tool used in conjunction with packer-windows to
create Windows Vagrant boxes.

## Introduction

Inductor uses Go [text/template](http://golang.org/pkg/text/template/) templates
to generate the files necessary to create Windows Vagrant boxes. Inductor is
meant to be used alongside the scripts and templates which comes with
[packer-windows](https://github.com/joefitzgerald/packer-windows).

## Basic Usage

Install inductor:

`go get github.com/joefitzgerald/inductor`

Inductor only supports one action by default, which is to generate the files
necessary to build out a Windows Vagrant box via Packer:

- packer.json
- Autounattend.xml
- Vagrantfile

The OSs available are driven by the osregistry.json in the packer-windows
repository. You see which OSs are configured by running inductor with no
arguments:

```
$ inductor

Available OSs:

  windows2012
  windows81
  windows2012r2core
  windows10
  windows2008r2
  windows2012r2hyperv
  windows2012r2
  windows2008r2core
  windows7

You must specify an operating system argument
```

To execute inductor in preparation for a Packer build just pass in the OS you'd
like to use, for example:

```
inductor windows10
```

This will generate an Autounattend.xml, packer.json, and Vagrantfile in the
current directory. To chain the inductor output into Packer do this:

```
packer build $(inductor windows10)
```

This will execute inductor creating all the required artifacts for Packer and
then execute Packer using the generated templates.

## Inductor Options

Inductor uses a lot of sane defaults to make the happy path very easy,
however when you want to iterate on a box and/or need to build a production box
you'll need some flexibility. Inductor supports the following command line
options:

- `--osregistry <file.json>` This specifies the file path to a json file which
contains all the metadata for various Windows OSs. See OS Registry below.
- `--autounattend <autounattend.xml>` The file path where inductor will write
out the generated Autounattend.xml.
- `--packer <packer.json>` The file path where inductor will write out the
generated packer.json.
- `--vagrantfile <Vagantfile>` The file path where inductor will write out the
generated vagrantfile.
- `--productkey <key>` The Windows product key to be inserted into the
Autounattend.xml
- `--skipwindowsupdates` When specified the Windows Update step will be skipped.
- `--gui` When specified Packer will run the VM in GUI mode (headless=false).
- `--ssh` When specified Packer will use the SSH communicator with OpenSSH
instead of WinRM. WinRM will still be configured on the box for Vagrant.

## Templates

All input templates are standard Golang text/templates. By default inductor will
attempt to use the following template files in the current working directory:

- Autounattend.tpl
- packer.tpl
- Vagrantfile.tpl

Inductor supports templates spread across multiple files as well as OS specific
templates. There is a one/many to one relationship from an input template to an
output inductor generated file. Inductor also allows you to specialize or
override generic input templates by OS. This works for all 3 input templates.

### Template Loading Convention

(Autounattend|packer|Vagrantfile).tpl
Example: Autounattend.tpl

(Autounattend|packer|Vagrantfile)-OS.tpl
Example: Autounattend-windows10.tpl

(Autounattend|packer|Vagrantfile).subsection.tpl
Example: Autounattend.oobe.tpl

(Autounattend|packer|Vagrantfile)-OS.subsection.tpl
Example: Autounattend-windows10.oobe.tpl

Anything with an operating system in the template name that matches the
current system you're building will take precedence over the same named template
without an OS in the file name. Any template with an OS in the name that
doesn't match the current system you're building is ignored.

Given the following files in the current directory, the bold files will be
automatically loaded and merged together to be rendered to produce the final
Autounattend.xml file for Windows2012r2:

- Autounattend.tpl
- packer.tpl
- __Autounattend-windows2012r2.tpl__
- Autounattend-windows2008.tpl
- __Autounattend-windows2012r2.windowsPE.tpl__
- Autounattend-windows2008.windowsPE.tpl
- Autounattend.windowsPE.tpl
- __Autounattend.offlineServicing.tpl__

### Template Variables
- OSName
-	ProductKey
- WindowsImageName
-	VirtualboxGuestOsType
-	VmwareGuestOsType
-	IsoURL
-	IsoChecksumType
-	IsoChecksum
-	Communicator
-	Username
-	Password
-	DiskSize
-	RAM
-	CPU
-	Headless
-	WindowsUpdates

### Template Functions
- Contains
- Replace
- ToUpper
- ToLower
- SafeComputerName

## OS Registry

The OS registry contains predefined attributes for each OS that inductor can
generate Packer templates for. Packer-windows has its own registry which by
default contains various Windows OS trial editions. This is perfect if you
quickly need a Vagrant box for testing Windows.

You may find that you want to build your own box using your own ISO or OS flavor
not supported by packer-windows. This is where creating your own custom OS
registry file comes into play. The JSON file format is pretty self explanatory
so here it is:

```json
{
  "windows10": {
    "iso_url": "./iso/CLIENTENTERPRISEEVAL_OEMRET_X64FRE_EN-US.ISO",
    "iso_checksum_type": "sha1",
    "iso_checksum": "56ab095075be28a90bc0b510835280975c6bb2ce",
    "windows_image_name": "Windows 10 Enterprise",
    "virtualbox_guest_os_type": "Windows81_64",
    "vmware_guest_os_type": "windows8srv-64",
    "product_key": "FEED-ME2D"
  }
}
```

Except for product_key all other fields are required.

By default inductor looks in the current working directory for a file named
osregistry.json. If you name it something else or is in another directory
you can specify the location using the --osregistry flag.

## Contributing

Pull requests welcomed. Please ensure you create your edits in a branch off of
the `master` branch.
