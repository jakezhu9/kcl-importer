package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"kcl-importer/pkg/convert"
	"kcl-importer/pkg/jsonschema"
	"kcl-importer/pkg/kclschema"
	"log"
	"os"
	"path/filepath"
)

func main() {
	app := &cli.App{
		Name:     "kcl-importer",
		Usage:    "a command-line tool aimed at transforming a variety of languages into KCL.",
		Commands: []*cli.Command{NewImportCommand()},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func NewImportCommand() *cli.Command {
	return &cli.Command{
		Name:    "import",
		Aliases: []string{"i"},
		Usage:   "Import a file and convert it into KCL",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "the `FILE` to be imported",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "target",
				Aliases:     []string{"t"},
				Usage:       "the target `DIRECTORY` to store the generated files",
				DefaultText: "./",
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Printf("file: %s target: %s\n", ctx.String("file"), ctx.String("target"))

			// TODO: support other file
			schemaData, err := os.ReadFile(ctx.String("file"))
			if err != nil {
				return err
			}

			jsonSch := &jsonschema.Schema{}
			if err := json.Unmarshal(schemaData, jsonSch); err != nil {
				return err
			}

			kclSch := convert.JsonSchemaToKclSchema(*jsonSch)
			file, err := os.Create(filepath.Join(ctx.String("target"), "output.k"))
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			return kclschema.NewGenerator().Generate(kclSch, file)
		},
	}
}
