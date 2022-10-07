package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/abice/go-enum/generator"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli/v2"
)

var (
	version string
	commit  string
	date    string
	builtBy string
)

type rootT struct {
	FileNames         cli.StringSlice
	NoPrefix          bool
	Lowercase         bool
	NoCase            bool
	Marshal           bool
	SQL               bool
	SQLInt            bool
	Flag              bool
	Prefix            string
	Names             bool
	LeaveSnakeCase    bool
	SQLNullStr        bool
	SQLNullInt        bool
	Ptr               bool
	TemplateFileNames cli.StringSlice
	Aliases           cli.StringSlice
	MustParse         bool
	ForceLower        bool
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
		Version:         version,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				EnvVars:     []string{"GOFILE"},
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
				Name:        "sqlint",
				Usage:       "Tells the generator that a string typed enum should be stored in sql as an integer value.",
				Destination: &argv.SQLInt,
			},
			&cli.BoolFlag{
				Name:        "flag",
				Usage:       "Adds golang flag functions.",
				Destination: &argv.Flag,
			},
			&cli.StringFlag{
				Name:        "prefix",
				Usage:       "Adds a prefix with a user one. If you would like to replace the prefix, then combine this option with --noprefix.",
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
			&cli.StringSliceFlag{
				Name:        "template",
				Aliases:     []string{"t"},
				Usage:       "Additional template file(s) to generate enums.  Use more than one flag for more files. Templates will be executed in alphabetical order.",
				Destination: &argv.TemplateFileNames,
			},
			&cli.StringSliceFlag{
				Name:        "alias",
				Aliases:     []string{"a"},
				Usage:       "Adds or replaces aliases for a non alphanumeric value that needs to be accounted for. [Format should be \"key:value,key2:value2\", or specify multiple entries, or both!]",
				Destination: &argv.Aliases,
			},
			&cli.BoolFlag{
				Name:        "mustparse",
				Usage:       "Adds a Must version of the Parse that will panic on failure.",
				Destination: &argv.MustParse,
			},
			&cli.BoolFlag{
				Name:        "forcelower",
				Usage:       "Forces a camel cased comment to generate lowercased names.",
				Destination: &argv.ForceLower,
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := generator.ParseAliases(argv.Aliases.Value()); err != nil {
				return err
			}
			for _, fileOption := range argv.FileNames.Value() {

				g := generator.NewGenerator()
				g.Version = version
				g.Revision = commit
				g.BuildDate = date
				g.BuiltBy = builtBy

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
				if argv.SQLInt {
					g.WithSQLInt()
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
				if argv.MustParse {
					g.WithMustParse()
				}
				if argv.ForceLower {
					g.WithForceLower()
				}
				if templates := []string(argv.TemplateFileNames.Value()); len(templates) > 0 {
					for _, t := range templates {
						if fn, err := globFilenames(t); err != nil {
							return err
						} else {
							g.WithTemplates(fn...)
						}
					}
				}

				var filenames []string
				if fn, err := globFilenames(fileOption); err != nil {
					return err
				} else {
					filenames = fn
				}

				for _, fileName := range filenames {
					originalName := fileName

					out("go-enum started. file: %s\n", color.Cyan(originalName))
					fileName, _ = filepath.Abs(fileName)
					outFilePath := fmt.Sprintf("%s_enum.go", strings.TrimSuffix(fileName, filepath.Ext(fileName)))

					// Parse the file given in arguments
					raw, err := g.GenerateFromFile(fileName)
					if err != nil {
						return fmt.Errorf("failed generating enums\nInputFile=%s\nError=%s", color.Cyan(fileName), color.RedBg(err))
					}

					// Nothing was generated, ignore the output and don't create a file.
					if len(raw) < 1 {
						out(color.Yellow("go-enum ignored. file: %s\n"), color.Cyan(originalName))
						continue
					}

					mode := int(0644)
					err = os.WriteFile(outFilePath, raw, os.FileMode(mode))
					if err != nil {
						return fmt.Errorf("failed writing to file %s: %s", color.Cyan(outFilePath), color.Red(err))
					}
					out("go-enum finished. file: %s\n", color.Cyan(originalName))
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// globFilenames gets a list of filenames matching the provided filename.
// In order to maintain existing capabilities, only glob when a * is in the path.
// Leave execution on par with old method in case there are bad patterns in use that somehow
// work without the Glob method.
func globFilenames(filename string) ([]string, error) {
	if strings.Contains(filename, "*") {
		matches, err := filepath.Glob(filename)
		if err != nil {
			return []string{}, fmt.Errorf("failed parsing glob filepath\nInputFile=%s\nError=%s", color.Cyan(filename), color.RedBg(err))
		}
		return matches, nil
	} else {
		return []string{filename}, nil
	}
}
