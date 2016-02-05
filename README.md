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

In the current directory you will now have an Autounattend.xml, packer.json, and
Vagrantfile ready for Packer `packer build packer.json`.

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
- `--autounattendtpl <Autounattend.xml.tpl>` The file path to the input
Autounattend text/template used to generate the Autounattend.xml.
- `--packertpl <packer.json.tpl>` The file path to the input packer.json
text/template used to generate packer.json.
- `--vagrantfiletpl <Vagrantfile.tpl>` The file path to the input Vagrantfile
text/tepmlate used to generate the box Vagrantfile.
- `--productkey <key>` The Windows product key to be inserted into the
Autounattend.xml
- `--skipwindowsupdates` When specified the Windows Update step will be skipped.
- `--gui` When specified Packer will run the VM in GUI mode (headless=false).

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
  "windowsXX": {
    "iso_url": "./iso/CLIENTENTERPRISEEVAL_OEMRET_X64FRE_EN-US.ISO",
    "iso_checksum_type": "sha1",
    "iso_checksum": "56ab095075be28a90bc0b510835280975c6bb2ce",
    "windows_image_name": "Windows 10 Enterprise",
    "virtualbox_guest_os_type": "Windows81_64",
    "vmware_guest_os_type": "windows8srv-64"
  }
}
```

By default inductor looks in the current working directory for a file named
osregistry.json. If you name it something else or is in another directory
you can specify the location using the --osregistry flag.

## Contributing

Pull requests welcomed. Please ensure you create your edits in a branch off of
the `master` branch.
