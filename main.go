package main

import (
	"fmt"
	"log"
	"os"

	"github.com/quantumcycle/go-import-checks/validator"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func readYamlCfg(path string) (validator.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return validator.Config{}, err
	}

	cfg := validator.Config{}
	err = yaml.UnmarshalStrict(data, &cfg)
	if err != nil {
		return validator.Config{}, err
	}

	return cfg, nil
}

func validateImports(cfg validator.Config, path string, debug bool) {
	validationErrs, err := validator.Validate(path, cfg.Checks, debug)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if len(validationErrs) > 0 {
		for _, e := range validationErrs {
			if e.Reason == validator.ReasonNotAllow {
				fmt.Fprintf(os.Stderr, "File %s\n", e.Path)
				fmt.Fprintf(os.Stderr, "\tImport: %s\n", e.ImportPath)
				fmt.Fprintf(os.Stderr, "\tProblem: doesn't match anything in the allowed list %v\n", e.Rule.Allow)
				fmt.Fprintf(os.Stderr, "\tFolder definition: %s\n", e.Check.Folder)
			} else if e.Reason == validator.ReasonRejected {
				fmt.Fprintf(os.Stderr, "File %s\n", e.Path)
				fmt.Fprintf(os.Stderr, "\tImport: %s\n", e.ImportPath)
				fmt.Fprintf(os.Stderr, "\tProblem: is part the an explicit reject list %v\n", e.Rule.Reject)
				fmt.Fprintf(os.Stderr, "\tFolder definition: %s\n", e.Check.Folder)
			} else {
				panic("Unhandle validation reason")
			}
			fmt.Fprintf(os.Stderr, "\n")
		}
		os.Exit(1)
	}
}

func main() {
	app := &cli.App{
		Name:  "go-import-checks",
		Usage: "Validation import rules across a golang codebase",
		Action: func(c *cli.Context) error {
			var pathToCheck string = "."
			if c.String("root") != "" {
				pathToCheck = c.String("root")
			}

			cfg, err := readYamlCfg(c.String("config"))
			if err != nil {
				log.Fatalf("Cannot read yaml config file [%s]", c.String("config"))
			}

			validateImports(cfg, pathToCheck, c.Bool("debug"))
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Usage:    "Yaml config file for import rules",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "root",
				Usage:    "Root folder to recursively start checking golang source files (default to current folder)",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "debug",
				Usage:    "Debug output for checks",
				Required: false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
