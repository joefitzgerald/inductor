package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/joefitzgerald/inductor/osregistry"
	"github.com/joefitzgerald/inductor/renderer"
)

// Version of the CLI
var Version = "0.0.0"

func main() {
	app := newApp()
	err := app.Run(os.Args)
	if err != nil {
		die(err)
	}
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
			Name:  "productkey, pk",
			Usage: "The MS Windows product key if you have one",
		},
		cli.BoolFlag{
			Name:  "skipwindowsupdates, swu",
			Usage: "Skips running Windows updates on first boot",
		},
		cli.BoolFlag{
			Name:  "ssh",
			Usage: "Uses the Packer SSH communicator instead of the default WinRM",
		},
		cli.BoolFlag{
			Name:  "gui, g",
			Usage: "Run the VM with a GUI",
		},
	}
	app.Action = func(c *cli.Context) {
		// load the OS registry data
		file, err := os.Open(c.String("osregistry"))
		if err != nil {
			die("Couldn't find the requried Windows OS registry file", err)
		}
		defer file.Close()

		osRegistry, err := osregistry.New(file)
		if err != nil {
			die("Couldn't load the Windows OS registry", err)
		}

		// get our required windows version/edition string, e.g. 'windows10'
		var windowsEdition string
		if len(c.Args()) > 0 {
			windowsEdition = c.Args()[0]
		} else {
			fmt.Println("Available OSs:")
			fmt.Println()
			for _, s := range osRegistry.List() {
				fmt.Println(fmt.Sprintf("  %s", s))
			}
			fmt.Println()
			die("You must specify an operating system argument")
		}

		// create the default options set based on the OS registry info
		windows, ok := osRegistry.Get(windowsEdition)
		if !ok {
			die("Couldn't find OS registration for '%s'", windowsEdition)
		}
		opts, err := renderer.NewRenderOptionsWithOverrides(windows)
		if err != nil {
			die(err)
		}

		// apply any command line overrides to the options set
		opts.WindowsUpdates = !c.Bool("skipwindowsupdates")
		opts.Headless = !c.Bool("gui")
		if len(c.String("productkey")) > 0 {
			opts.ProductKey = c.String("productkey")
		}
		if c.Bool("ssh") {
			opts.Communicator = "ssh"
		}

		// create output packer.json file writer
		packerJSONOutPath := c.String("packer")
		packerJSON, err := os.Create(packerJSONOutPath)
		if err != nil {
			die(err)
		}
		defer packerJSON.Close()
		packerJSONWriter := bufio.NewWriter(packerJSON)
		defer packerJSONWriter.Flush()

		// create output Autounattend.xml file writer
		autounattendXML, err := os.Create(c.String("autounattend"))
		if err != nil {
			die(err)
		}
		defer autounattendXML.Close()
		autounattendXMLWriter := bufio.NewWriter(autounattendXML)
		defer autounattendXMLWriter.Flush()

		// create output Vagrantfile file writer
		vagrantfile, err := os.Create(c.String("vagrantfile"))
		if err != nil {
			die(err)
		}
		defer vagrantfile.Close()
		vagrantfileWriter := bufio.NewWriter(vagrantfile)
		defer vagrantfileWriter.Flush()

		// finally render the packer.json and Autounattend.xml
		cwd, err := os.Getwd()
		if err != nil {
			die(err)
		}
		tpl, err := renderer.NewPackerTemplateWithOverrides(cwd, opts.OSName)
		if err != nil {
			die(err)
		}
		err = tpl.Render(opts, packerJSONWriter, autounattendXMLWriter, vagrantfileWriter)
		if err != nil {
			die(err)
		}

		// this allows us to do command substitution with Packer
		fmt.Print(packerJSONOutPath)
	}
	return app
}

func die(vals ...interface{}) {
	if len(vals) > 1 || vals[0] != nil {
		os.Stderr.WriteString(fmt.Sprintln(vals...))
		os.Exit(1)
	}
}
