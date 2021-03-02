package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/abice/go-enum/generator"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli/v2"
)

type rootT struct {
	FileNames      cli.StringSlice
	NoPrefix       bool
	Lowercase      bool
	NoCase         bool
	Marshal        bool
	SQL            bool
	Flag           bool
	Prefix         string
	Names          bool
	LeaveSnakeCase bool
	SQLNullStr     bool
	SQLNullInt     bool
	Ptr            bool
}

func main() {
	var argv rootT

	clr := color.New()
	out := func(format string, args ...interface{}) {
		_, _ = fmt.Fprintf(clr.Output(), format, args...)
	}

	app := &cli.App{
		Name:            "go-enum",
		Usage:           "An enum generator for go",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "The file(s) to generate enums.  Use more than one flag for more files.",
				Required:    true,
				Destination: &argv.FileNames,
			},
			&cli.BoolFlag{
				Name:        "noprefix",
				Usage:       "Prevents the constants generated from having the Enum as a prefix.",
				Destination: &argv.NoPrefix,
			},
			&cli.BoolFlag{
				Name:        "lower",
				Usage:       "Adds lowercase variants of the enum strings for lookup.",
				Destination: &argv.Lowercase,
			},
			&cli.BoolFlag{
				Name:        "nocase",
				Usage:       "Adds case insensitive parsing to the enumeration (forces lower flag).",
				Destination: &argv.NoCase,
			},
			&cli.BoolFlag{
				Name:        "marshal",
				Usage:       "Adds text (and inherently json) marshalling functions.",
				Destination: &argv.Marshal,
			},
			&cli.BoolFlag{
				Name:        "sql",
				Usage:       "Adds SQL database scan and value functions.",
				Destination: &argv.SQL,
			},
			&cli.BoolFlag{
				Name:        "flag",
				Usage:       "Adds golang flag functions.",
				Destination: &argv.Flag,
			},
			&cli.StringFlag{
				Name:        "prefix",
				Usage:       "Replaces the prefix with a user one.",
				Destination: &argv.Prefix,
			},
			&cli.BoolFlag{
				Name:        "names",
				Usage:       "Generates a 'Names() []string' function, and adds the possible enum values in the error response during parsing",
				Destination: &argv.Names,
			},
			&cli.BoolFlag{
				Name:        "nocamel",
				Usage:       "Removes the snake_case to CamelCase name changing",
				Destination: &argv.LeaveSnakeCase,
			},
			&cli.BoolFlag{
				Name:        "ptr",
				Usage:       "Adds a pointer method to get a pointer from const values",
				Destination: &argv.Ptr,
			},
			&cli.BoolFlag{
				Name:        "sqlnullint",
				Usage:       "Adds a Null{{ENUM}} type for marshalling a nullable int value to sql",
				Destination: &argv.SQLNullInt,
			},
			&cli.BoolFlag{
				Name:        "sqlnullstr",
				Usage:       "Adds a Null{{ENUM}} type for marshalling a nullable string value to sql.  If sqlnullint is specified too, it will be Null{{ENUM}}Str",
				Destination: &argv.SQLNullStr,
			},
		},
		Action: func(ctx *cli.Context) error {
			for _, fileName := range argv.FileNames.Value() {

				g := generator.NewGenerator()

				if argv.NoPrefix {
					g.WithNoPrefix()
				}
				if argv.Lowercase {
					g.WithLowercaseVariant()
				}
				if argv.NoCase {
					g.WithCaseInsensitiveParse()
				}
				if argv.Marshal {
					g.WithMarshal()
				}
				if argv.SQL {
					g.WithSQLDriver()
				}
				if argv.Flag {
					g.WithFlag()
				}
				if argv.Names {
					g.WithNames()
				}
				if argv.LeaveSnakeCase {
					g.WithoutSnakeToCamel()
				}
				if argv.Prefix != "" {
					g.WithPrefix(argv.Prefix)
				}
				if argv.Ptr {
					g.WithPtr()
				}
				if argv.SQLNullInt {
					g.WithSQLNullInt()
				}
				if argv.SQLNullStr {
					g.WithSQLNullStr()
				}

				originalName := fileName

				out("go-enum started. file: %s\n", color.Cyan(originalName))
				fileName, _ = filepath.Abs(fileName)
				outFilePath := fmt.Sprintf("%s_enum.go", strings.TrimSuffix(fileName, filepath.Ext(fileName)))

				// Parse the file given in arguments
				raw, err := g.GenerateFromFile(fileName)
				if err != nil {
					return fmt.Errorf("failed generating enums\nInputFile=%s\nError=%s", color.Cyan(fileName), color.RedBg(err))
				}

				mode := int(0644)
				err = ioutil.WriteFile(outFilePath, raw, os.FileMode(mode))
				if err != nil {
					return fmt.Errorf("failed writing to file %s: %s", color.Cyan(outFilePath), color.Red(err))
				}
				out("go-enum finished. file: %s\n", color.Cyan(originalName))
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
