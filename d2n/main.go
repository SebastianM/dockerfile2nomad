package main

import (
  "os"

  "gopkg.in/urfave/cli.v1"
  "github.com/sebastianm/dockerfile2nomad/d2n/command"
)

const version = "0.1.0"

func main() {
  app := cli.NewApp()
  app.Name = "d2n"
  app.Usage = "dockerfile2nomad"
  app.Version = version

  app.Commands = []cli.Command{
    command.Generate,
  }
  app.Run(os.Args)
}