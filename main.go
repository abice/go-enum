package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/abice/go-enum/generator"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli/v2"
)

var (
	versionOnce sync.Once
	version     string
	commit      string
	date        string
	builtBy     string
)

type rootT struct {
	FileNames         cli.StringSlice
	NoPrefix          bool
	NoIota            bool
	Lowercase         bool
	NoCase            bool
	Marshal           bool
	SQL               bool
	SQLInt            bool
	Flag              bool
	JsonPkg           string
	Prefix            string
	Names             bool
	Values            bool
	LeaveSnakeCase    bool
	SQLNullStr        bool
	SQLNullInt        bool
	Ptr               bool
	TemplateFileNames cli.StringSlice
	Aliases           cli.StringSlice
	BuildTags         cli.StringSlice
	MustParse         bool
	ForceLower        bool
	ForceUpper        bool
	NoComments        bool
	OutputSuffix      string
}

func initializeVersion() {
	versionOnce.Do(func() {
		if version != "" {
			return
		}
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			return
		}
		builtBy = "go install"
		version = buildInfo.Main.Version
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value
			case "vcs.time":
				date = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					commit += "-modified"
				}
			}
		}
	})
}

func main() {
	var argv rootT

	initializeVersion()

	clr := color.New()
	out := func(format string, args ...any) {
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
				Name:        "jsonpkg",
				Usage:       "Custom json package for imports instead encoding/json.",
				Destination: &argv.JsonPkg,
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
				Name:        "values",
				Usage:       "Generates a 'Values() []{{ENUM}}' function.",
				Destination: &argv.Values,
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
			&cli.BoolFlag{
				Name:        "forceupper",
				Usage:       "Forces a camel cased comment to generate uppercased names.",
				Destination: &argv.ForceUpper,
			},
			&cli.BoolFlag{
				Name:        "nocomments",
				Usage:       "Removes auto generated comments.  If you add your own comments, these will still be created.",
				Destination: &argv.NoComments,
			},
			&cli.StringSliceFlag{
				Name:        "buildtag",
				Aliases:     []string{"b"},
				Usage:       "Adds build tags to a generated enum file.",
				Destination: &argv.BuildTags,
			},
			&cli.StringFlag{
				Name:        "output-suffix",
				Usage:       "Changes the default filename suffix of _enum to something else.  `.go` will be appended to the end of the string no matter what, so that `_test.go` cases can be accommodated ",
				Destination: &argv.OutputSuffix,
			},
			&cli.BoolFlag{
				Name:        "no-iota",
				Usage:       "Disables the use of iota in generated enums.",
				Destination: &argv.NoIota,
			},
		},
		Action: func(ctx *cli.Context) error {
			aliases, err := generator.ParseAliases(argv.Aliases.Value())
			if err != nil {
				return err
			}
			for _, fileOption := range argv.FileNames.Value() {

				// Build configuration structure
				jsonPkg := argv.JsonPkg
				if jsonPkg == "" {
					jsonPkg = "encoding/json"
				}

				var templateFileNames []string
				if templates := []string(argv.TemplateFileNames.Value()); len(templates) > 0 {
					for _, t := range templates {
						if fn, err := globFilenames(t); err != nil {
							return err
						} else {
							templateFileNames = append(templateFileNames, fn...)
						}
					}
				}

				config := generator.GeneratorConfig{
					NoPrefix:          argv.NoPrefix,
					NoIota:            argv.NoIota,
					LowercaseLookup:   argv.Lowercase || argv.NoCase,
					CaseInsensitive:   argv.NoCase,
					Marshal:           argv.Marshal,
					SQL:               argv.SQL,
					SQLInt:            argv.SQLInt,
					Flag:              argv.Flag,
					Names:             argv.Names,
					Values:            argv.Values,
					LeaveSnakeCase:    argv.LeaveSnakeCase,
					JSONPkg:           jsonPkg,
					Prefix:            argv.Prefix,
					SQLNullInt:        argv.SQLNullInt,
					SQLNullStr:        argv.SQLNullStr,
					Ptr:               argv.Ptr,
					MustParse:         argv.MustParse,
					ForceLower:        argv.ForceLower,
					ForceUpper:        argv.ForceUpper,
					NoComments:        argv.NoComments,
					BuildTags:         argv.BuildTags.Value(),
					ReplacementNames:  aliases,
					TemplateFileNames: templateFileNames,
				}

				// Create generator with configuration
				g := generator.NewGeneratorWithConfig(config)
				g.Version = version
				g.Revision = commit
				g.BuildDate = date
				g.BuiltBy = builtBy

				var filenames []string
				if fn, err := globFilenames(fileOption); err != nil {
					return err
				} else {
					filenames = fn
				}

				outputSuffix := `_enum`
				if argv.OutputSuffix != "" {
					outputSuffix = argv.OutputSuffix
				}

				for _, fileName := range filenames {
					originalName := fileName

					out("go-enum started. file: %s\n", color.Cyan(originalName))
					fileName, _ = filepath.Abs(fileName)

					outFilePath := fmt.Sprintf("%s%s.go", strings.TrimSuffix(fileName, filepath.Ext(fileName)), outputSuffix)
					if strings.HasSuffix(fileName, "_test.go") {
						outFilePath = strings.Replace(outFilePath, "_test"+outputSuffix+".go", outputSuffix+"_test.go", 1)
					}

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

					mode := int(0o644)
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
