package command

import (
  "errors"
  "fmt"
  "strings"

  "gopkg.in/urfave/cli.v1"
  "github.com/sebastianm/dockerfile2nomad/generator"
	"os"
  "io/ioutil"
)

var Generate = cli.Command{
  Name:    "generate",
  Aliases: []string{"g"},
  Usage:   "Generates a Nomad HCL file out of a Dockerfile",
  Action:  generateAction,
  Flags: []cli.Flag {
    cli.StringFlag{
      Name: "j, job",
      Usage: "Nomand job name",
    },
    cli.StringFlag{
      Name: "o, output",
      Usage: "Output file location",
    },
  },
}

func generateAction(c *cli.Context) error {
  if c.NArg() < 1 || strings.Trim(c.Args().First(), "") == "" {
    fmt.Println("Please provide a Dockerfile as the first argument")
    return errors.New("Please provide a Dockerfile as the first argument")
  }

  dockerfile := c.Args().First()
  file, err := os.Open(dockerfile)
  if err != nil {
    fmt.Printf("Error opening Dockerfile %s: %s\n", dockerfile, err.Error())
    return err
  }

  jobName := c.String("o")
  if jobName == "" {
    jobName = "default"
  }
  output, err := generator.GenerateFromReader(file, jobName)
  if err != nil {
    fmt.Printf("Error generating nomad file: %s\n", err.Error())
    return err
  }

  // write content to file, if option provided
  if c.String("o") != "" {
    if err := ioutil.WriteFile(c.String("o"), output, 0666); err != nil {
      fmt.Printf("Error writing to output file %s: %s\n", c.String("o"), err.Error())
      return err
    }
    return nil
  }

  // when no output file (-o) is specified, we write to stdout
  fmt.Print(string(output))
  return nil
}
