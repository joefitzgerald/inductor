package main

import (
    "os"
    "github.com/codegangsta/cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "Inductor"
    app.Usage = "Generate Packer Templates"
    app.Action = func(c *cli.Context) {
      println("Hello world!")
    }

    app.Run(os.Args)
}