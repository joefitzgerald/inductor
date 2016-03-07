package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/joefitzgerald/inductor/configuration"
	"github.com/joefitzgerald/inductor/cpy"
	"github.com/joefitzgerald/inductor/renderer"
	"github.com/joefitzgerald/inductor/tpl"
)

// runError is called when an error occurs during Run
type runError func(error)

type app struct {
	ctx        *cli.Context
	errHandler runError
}

// newInductorApp creates a new inductor app instance
func newInductorApp(errHandler runError) *app {
	return &app{
		errHandler: errHandler,
	}
}

func (a *app) Run(ctx *cli.Context) {
	a.ctx = ctx
	config, err := a.loadConfiguration()
	if err != nil {
		a.errHandler(fmt.Errorf("Couldn't load the inductor.json configuration file. %s", err))
	}
	opts, err := a.createRenderOpts(config)
	if err != nil {
		a.errHandler(err)
	}
	outDir, err := a.outDir(config)
	if err != nil {
		a.errHandler(err)
	}

	// find all templates
	cwd, err := os.Getwd()
	if err != nil {
		a.errHandler(err)
	}
	templates := tpl.New(cwd, opts.OSName)

	// render all the templates to the output directory
	renderer := renderer.New(opts, outDir)
	err = renderer.Render(templates)
	if err != nil {
		a.errHandler(err)
	}

	// copy over any non-templates to the output directory
	copier := cpy.New()
	err = copier.Copy(cwd, outDir)
	if err != nil {
		a.errHandler(err)
	}
}

func (a *app) loadConfiguration() (config *configuration.InductorConfiguration, err error) {
	configFile, err := os.Open(a.ctx.String("config"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := configFile.Close(); cerr != nil {
			err = cerr
		}
	}()
	return configuration.New(configFile)
}

func (a *app) createRenderOpts(config *configuration.InductorConfiguration) (*renderer.RenderOptions, error) {
	var osname string
	if len(a.ctx.Args()) > 0 {
		osname = a.ctx.Args()[0]
	} else {
		fmt.Println("Available Operating Systems:")
		fmt.Println()
		for _, s := range config.List() {
			fmt.Println(fmt.Sprintf("  %s", s))
		}
		fmt.Println()
		return nil, errors.New("You must specify an operating system argument")
	}

	// create the default options set based on the inductor config
	var edition string
	if len(a.ctx.String("edition")) > 0 {
		edition = a.ctx.String("edition")
	}
	opts, err := renderer.NewRenderOptions(osname, edition, config)
	if err != nil {
		return nil, err
	}

	// apply any command line overrides to the options set
	opts.WindowsUpdates = !a.ctx.Bool("skipwindowsupdates")
	opts.Headless = !a.ctx.Bool("gui")
	if len(a.ctx.String("productkey")) > 0 {
		opts.ProductKey = a.ctx.String("productkey")
	}
	if a.ctx.Bool("ssh") {
		opts.Communicator = "ssh"
	}

	return opts, err
}

func (a *app) outDir(config *configuration.InductorConfiguration) (string, error) {
	outDir := config.OutDir
	if len(a.ctx.String("outdir")) > 0 {
		outDir = a.ctx.String("outdir")
	}
	return filepath.Abs(outDir)
}
