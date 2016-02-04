package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/joefitzgerald/inductor/renderer"
)

// Version of the CLI
var Version = "0.0.0"

func main() {
	app := newApp()
	app.Run(os.Args)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Inductor"
	app.Usage = "Generate Packer Templates"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "osregistry, os",
			Value: "osregistry.json",
			Usage: "packer-windows Windows data",
		},
		cli.StringFlag{
			Name:  "autounattend, a",
			Value: "Autounattend.xml",
			Usage: "The output Autounattend.xml file path",
		},
		cli.StringFlag{
			Name:  "packer, p",
			Value: "packer.json",
			Usage: "The output packer.json file path",
		},
		cli.StringFlag{
			Name:  "vagrantfile, vf",
			Value: "Vagrantfile",
			Usage: "The output Vagrantfile file path",
		},
		cli.StringFlag{
			Name:  "autounattendtpl, atpl",
			Value: "Autounattend.xml.tpl",
			Usage: "The input Autounattend.xml template file path",
		},
		cli.StringFlag{
			Name:  "packertpl, ptpl",
			Value: "packer.json.tpl",
			Usage: "The input packer.json template file path",
		},
		cli.StringFlag{
			Name:  "vagrantfiletpl, vtpl",
			Value: "Vagrantfile.tpl",
			Usage: "The input Vagrantfile template file path",
		},
		cli.StringFlag{
			Name:  "productkey, pk",
			Usage: "The MS Windows product key if you have one",
		},
		cli.BoolFlag{
			Name:  "skipwindowsupdates, swu",
			Usage: "Skips running Windows updates on first boot",
		},
		cli.BoolFlag{
			Name:  "gui, g",
			Usage: "Run the VM with a GUI",
		},
	}
	app.Action = func(c *cli.Context) {
		// get our required windows version/edition string, e.g. 'windows10'
		var windowsEdition string
		if len(c.Args()) > 0 {
			windowsEdition = c.Args()[0]
		} else {
			die("You must specify a Windows edition/version argument")
		}

		// create the default options set based on the OS registry info
		opts, err := renderer.NewRenderOptionsWithOverrides(windowsEdition, c.String("osregistry"))
		if err != nil {
			die(err)
		}

		// apply any command line overrides to the options set
		opts.WindowsUpdates = !c.Bool("skipwindowsupdates")
		opts.Headless = !c.Bool("gui")
		opts.ProductKey = c.String("productkey")

		// read in the packer.json.tpl
		packerJSON, err := os.Create(c.String("packer"))
		if err != nil {
			die(err)
		}
		defer packerJSON.Close()
		packerJSONWriter := bufio.NewWriter(packerJSON)
		defer packerJSONWriter.Flush()

		// read in the Autounattend.xml.tpl
		autounattendXML, err := os.Create(c.String("autounattend"))
		if err != nil {
			die(err)
		}
		defer autounattendXML.Close()
		autounattendXMLWriter := bufio.NewWriter(autounattendXML)
		defer autounattendXMLWriter.Flush()

		// read in the Vagrantfile.tpl
		vagrantfile, err := os.Create(c.String("vagrantfile"))
		if err != nil {
			die(err)
		}
		defer vagrantfile.Close()
		vagrantfileWriter := bufio.NewWriter(vagrantfile)
		defer vagrantfileWriter.Flush()

		// finally render the packer.json and Autounattend.xml
		packerTplPath := c.String("packertpl")
		autounattendTplPath := c.String("autounattendtpl")
		vagrantfileTplPath := c.String("vagrantfiletpl")
		tpl := renderer.NewPackerTemplateWithOverrides(packerTplPath, autounattendTplPath, vagrantfileTplPath)
		err = tpl.Render(opts, packerJSONWriter, autounattendXMLWriter, vagrantfileWriter)
		if err != nil {
			die(err)
		}
	}
	return app
}

func die(vals ...interface{}) {
	if len(vals) > 1 || vals[0] != nil {
		os.Stderr.WriteString(fmt.Sprintln(vals...))
		os.Exit(1)
	}
}
