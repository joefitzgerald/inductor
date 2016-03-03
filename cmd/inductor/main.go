package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/joefitzgerald/inductor/configuration"
	"github.com/joefitzgerald/inductor/renderer"
	"github.com/joefitzgerald/inductor/tpl"
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
			Name:  "config, c",
			Value: "inductor.json",
			Usage: "Inductor configuration source file",
		},
		cli.StringFlag{
			Name:  "outdir, o",
			Value: "out",
			Usage: "The root output directory for all rendered templates",
		},
		cli.StringFlag{
			Name:  "edition, e",
			Usage: "The optional operating system edition",
		},
		cli.StringFlag{
			Name:  "productkey, k",
			Usage: "The MS Windows product key if you have one",
		},
		cli.BoolFlag{
			Name:  "skipwindowsupdates, u",
			Usage: "Skips running Windows updates on first boot",
		},
		cli.BoolFlag{
			Name:  "ssh, s",
			Usage: "Uses the Packer SSH communicator instead of the default WinRM",
		},
		cli.BoolFlag{
			Name:  "gui, g",
			Usage: "Run the VM with a GUI",
		},
	}
	app.Action = func(c *cli.Context) {
		// load the inductor configuration
		configFile, err := os.Open(c.String("config"))
		if err != nil {
			die("Couldn't find the requried inductor.json configuration file", err)
		}
		defer configFile.Close()

		config, err := configuration.New(configFile)
		if err != nil {
			die("Couldn't load the inductor.json configuration file", err)
		}

		// get our required windows version/edition string, e.g. 'windows10'
		var osname string
		if len(c.Args()) > 0 {
			osname = c.Args()[0]
		} else {
			fmt.Println("Available Operating Systems:")
			fmt.Println()
			for _, s := range config.List() {
				fmt.Println(fmt.Sprintf("  %s", s))
			}
			fmt.Println()
			die("You must specify an operating system argument")
		}

		// create the default options set based on the inductor config
		var edition string
		if len(c.String("edition")) > 0 {
			edition = c.String("edition")
		}
		opts, err := renderer.NewRenderOptions(osname, edition, config)
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

		outDir := config.OutDir
		if len(c.String("outdir")) > 0 {
			outDir = c.String("outdir")
		}
		outDir, err = filepath.Abs(outDir)
		if err != nil {
			die(err)
		}

		// find all templates
		cwd, err := os.Getwd()
		if err != nil {
			die(err)
		}
		templates := tpl.New(cwd, osname)

		// finally render all the templates to the output directory
		renderer := renderer.New(opts, outDir)
		err = renderer.Render(templates)
		if err != nil {
			die(err)
		}

		// this allows us to do command substitution with Packer
		// TODO: output path to packer.json
	}
	return app
}

func die(vals ...interface{}) {
	if len(vals) > 1 || vals[0] != nil {
		os.Stderr.WriteString(fmt.Sprintln(vals...))
		os.Exit(1)
	}
}
