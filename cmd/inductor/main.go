package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
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
	app.Action = newInductorApp(die).Run
	return app
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}
