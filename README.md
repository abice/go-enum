# go-enum

[![Actions Status](https://github.com/abice/go-enum/actions/workflows/build_and_test.yml/badge.svg)](https://github.com/abice/go-enum/actions/workflows/build_and_test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/abice/go-enum)](https://goreportcard.com/report/github.com/abice/go-enum)
[![Coverage Status](https://coveralls.io/repos/github/abice/go-enum/badge.svg)](https://coveralls.io/github/abice/go-enum)
[![GoDoc](https://godoc.org/github.com/abice/go-enum?status.svg)](https://godoc.org/github.com/abice/go-enum)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

An enum generator for Go that creates type-safe enumerations with useful methods from simple type declarations.

## Key Features

- 🎯 **Type-safe enums** - Generate const declarations with proper typing
- 🔄 **String conversion** - Automatic `String()` and `Parse()` methods
- 📦 **JSON marshaling** - Built-in JSON marshal/unmarshal support
- 🗄️ **SQL integration** - Database scan and value functions
- 🧪 **Custom templates** - Extend with your own code generation
- ✅ **Multiple types** - Support for `int`, `string`, and other base types
- 🏷️ **Flexible syntax** - Simple comment-based enum declarations

## Table of Contents

- [Key Features](#key-features)
- [Requirements](#requirements)
- [Quick Start](#quick-start)
- [How it works](#how-it-works)
- [String typed enums](#now-with-string-typed-enums)
- [Recent Features](#recent-features)
- [Goal](#goal)
- [Docker image](#docker-image)
- [Installation](#installation)
- [Adding it to your project](#adding-it-to-your-project)
- [Command options](#command-options)
- [Syntax](#syntax)
- [Examples](#example)

## Requirements

- Go 1.23.0 or later

## Quick Start

1. Install go-enum:

   ```shell
   go install github.com/abice/go-enum@latest
   ```

2. Create a Go file with an enum declaration:

   ```go
   package main
   
   // ENUM(red, green, blue)
   type Color int
   ```

3. Generate the enum:

   ```shell
   go-enum -f your_file.go
   ```

4. Use your generated enum:

   ```go
   color := ColorRed
   fmt.Println(color.String()) // prints "red"
   
   parsed, err := ParseColor("green")
   if err == nil {
       fmt.Println(parsed) // prints "green"
   }
   ```

## How it works

go-enum will take a commented type declaration like this:

```go
// ENUM(jpeg, jpg, png, tiff, gif)
type ImageType int
```

and generate a file with the iota definition along various optional niceties that you may need:

```go
const (
 // ImageTypeJpeg is a ImageType of type Jpeg.
 ImageTypeJpeg ImageType = iota
 // ImageTypeJpg is a ImageType of type Jpg.
 ImageTypeJpg
 // ImageTypePng is a ImageType of type Png.
 ImageTypePng
 // ImageTypeTiff is a ImageType of type Tiff.
 ImageTypeTiff
 // ImageTypeGif is a ImageType of type Gif.
 ImageTypeGif
)

// String implements the Stringer interface.
func (x ImageType) String() string

// ParseImageType attempts to convert a string to a ImageType.
func ParseImageType(name string) (ImageType, error)

// MarshalText implements the text marshaller method.
func (x ImageType) MarshalText() ([]byte, error)

// UnmarshalText implements the text unmarshaller method.
func (x *ImageType) UnmarshalText(text []byte) error
```

**Fear not the fact that the `MarshalText` and `UnmarshalText` are generated rather than JSON methods... they will still be utilized by the default JSON encoding methods.**

If you find that the options given are not adequate for your use case, there is an option to add a custom template (`-t` flag) to the processing engine so that your custom code can be created!

## Now with string typed enums

```go
// ENUM(pending, running, completed, failed)
type StrState string
```

```go
const (
 // StrStatePending is a StrState of type pending.
 StrStatePending StrState = "pending"
 // StrStateRunning is a StrState of type running.
 StrStateRunning StrState = "running"
 // StrStateCompleted is a StrState of type completed.
 StrStateCompleted StrState = "completed"
 // StrStateFailed is a StrState of type failed.
 StrStateFailed StrState = "failed"
)
```

If you would like to get integer values in sql, but strings elsewhere, you can assign an int value in the declaration
like always, and specify the `--sqlint` flag.  Those values will be then used to convey the int value to sql, while allowing you to use only strings elsewhere.
This might be helpful for things like swagger docs where you want the same type being used on the api layer, as you do in the
sql layer, and not have swagger assume that your enumerations are integers, but are in fact strings!

```go
// swagger:enum StrState
// ENUM(pending, running, completed, failed)
type StrState string
```

## Recent Features

### Custom JSON Package Support (v0.7.0+)

You can now specify a custom JSON package for imports instead of the standard `encoding/json`:

```shell
go-enum --marshal --jsonpkg="github.com/goccy/go-json" -f your_file.go
```

### Disable iota (v0.9.0+)

For cases where you don't want to use `iota` in your generated enums:

```shell
go-enum --no-iota -f your_file.go
```

### Custom Output Suffix

Change the default `_enum.go` suffix to something else:

```shell
go-enum --output-suffix="_generated" -f your_file.go  # Creates your_file_generated.go
```

## Goal

The goal of go-enum is to create an easy to use enum generator that will take a decorated type declaration like `type EnumName int` and create the associated constant values and funcs that will make life a little easier for adding new values.
It's not perfect, but I think it's useful.

I took the output of the [Stringer](https://godoc.org/golang.org/x/tools/cmd/stringer) command as the `String()` method, and added a way to parse a string value.

## Docker image

You can use the Docker image directly for running the command if you do not wish to install anything locally:

```shell
docker run -w /app -v $(pwd):/app abice/go-enum:latest
```

## Installation

### Using go install (recommended)

Install the latest version directly from source:

```shell
go install github.com/abice/go-enum@latest
```

### Using GitHub releases

You can download a pre-built binary from GitHub releases. (Thanks to [GoReleaser](https://github.com/goreleaser/goreleaser-action))

```shell
# Download the latest release for your platform
curl -fsSL "https://github.com/abice/go-enum/releases/download/v0.9.0/go-enum_$(uname -s)_$(uname -m)" -o go-enum
chmod +x go-enum

# Or use a specific version by replacing v0.9.0 with your desired version
# Example: curl -fsSL "https://github.com/abice/go-enum/releases/download/v0.8.0/go-enum_$(uname -s)_$(uname -m)" -o go-enum
```

### Using Docker

You can use the Docker image directly without installing anything locally:

```shell
# Use the latest version
docker run -w /app -v $(pwd):/app abice/go-enum:latest

# Or use a specific version
docker run -w /app -v $(pwd):/app abice/go-enum:v0.9.0
```

## Adding it to your project

### Using go generate

1. Add a go:generate line to your file like so... `//go:generate go-enum --marshal`
1. Run go generate like so `go generate ./...`
1. Enjoy your newly created Enumeration!

### Using Makefile

If you prefer makefile stuff, you can always do something like this:

```Makefile
STANDARD_ENUMS = ./example/animal_enum.go \
 ./example/color_enum.go

NULLABLE_ENUMS = ./example/sql_enum.go

$(STANDARD_ENUMS): GO_ENUM_FLAGS=--nocase --marshal --names --ptr
$(NULLABLE_ENUMS): GO_ENUM_FLAGS=--nocase --marshal --names --sqlnullint --ptr

enums: $(STANDARD_ENUMS) $(NULLABLE_ENUMS)

# The generator statement for go enum files.  Files that invalidate the
# enum file: source file, the binary itself, and this file (in case you want to generate with different flags)
%_enum.go: %.go $(GOENUM) Makefile
 $(GOENUM) -f $*.go $(GO_ENUM_FLAGS)
```

## Command options

``` shell
go-enum --help

NAME:
   go-enum - An enum generator for go

USAGE:
   go-enum [global options]

VERSION:
   example

GLOBAL OPTIONS:
   --file value, -f value [ --file value, -f value ]          The file(s) to generate enums.  Use more than one flag for more files. [$GOFILE]
   --noprefix                                                 Prevents the constants generated from having the Enum as a prefix. (default: false)
   --lower                                                    Adds lowercase variants of the enum strings for lookup. (default: false)
   --nocase                                                   Adds case insensitive parsing to the enumeration (forces lower flag). (default: false)
   --marshal                                                  Adds text (and inherently json) marshalling functions. (default: false)
   --sql                                                      Adds SQL database scan and value functions. (default: false)
   --sqlint                                                   Tells the generator that a string typed enum should be stored in sql as an integer value. (default: false)
   --flag                                                     Adds golang flag functions. (default: false)
   --jsonpkg value                                            Custom json package for imports instead encoding/json.
   --prefix value                                             Adds a prefix with a user one. If you would like to replace the prefix, then combine this option with --noprefix.
   --names                                                    Generates a 'Names() []string' function, and adds the possible enum values in the error response during parsing (default: false)
   --values                                                   Generates a 'Values() []{{ENUM}}' function. (default: false)
   --nocamel                                                  Removes the snake_case to CamelCase name changing (default: false)
   --ptr                                                      Adds a pointer method to get a pointer from const values (default: false)
   --sqlnullint                                               Adds a Null{{ENUM}} type for marshalling a nullable int value to sql (default: false)
   --sqlnullstr                                               Adds a Null{{ENUM}} type for marshalling a nullable string value to sql.  If sqlnullint is specified too, it will be Null{{ENUM}}Str (default: false)
   --template value, -t value [ --template value, -t value ]  Additional template file(s) to generate enums.  Use more than one flag for more files. Templates will be executed in alphabetical order.
   --alias value, -a value [ --alias value, -a value ]        Adds or replaces aliases for a non alphanumeric value that needs to be accounted for. [Format should be "key:value,key2:value2", or specify multiple entries, or both!]
   --mustparse                                                Adds a Must version of the Parse that will panic on failure. (default: false)
   --forcelower                                               Forces a camel cased comment to generate lowercased names. (default: false)
   --forceupper                                               Forces a camel cased comment to generate uppercased names. (default: false)
   --nocomments                                               Removes auto generated comments.  If you add your own comments, these will still be created. (default: false)
   --buildtag value, -b value [ --buildtag value, -b value ]  Adds build tags to a generated enum file.
   --output-suffix .go                                        Changes the default filename suffix of _enum to something else.  .go will be appended to the end of the string no matter what, so that `_test.go` cases can be accommodated
   --no-iota                                                  Disables the use of iota in generated enums. (default: false)
   --help, -h                                                 show help
   --version, -v                                              print the version
```

### Syntax

The parser looks for comments on your type defs and parse the enum declarations from it.
The parser will look for `ENUM(` and continue to look for comma separated values until it finds a `)`.  You can put values on the same line, or on multiple lines.\
If you need to have a specific value jump in the enum, you can now specify that by adding `=numericValue` to the enum declaration.  Keep in mind, this resets the data for all following values.  So if you specify `50` in the middle of an enum, each value after that will be `51, 52, 53...`

[Examples can be found in the example folder](./example/)

#### Comments

You can use comments inside enum that start with `//`\
The comment must be at the end of the same line as the comment value, only then it will be added as a comment to the generated constant.

```go
// Commented is an enumeration of commented values
/*
ENUM(
value1 // Commented value 1
value2
value3 // Commented value 3
)
*/
type Commented int
```

The generated comments in code will look something like:

```go
...
const (
    // CommentedValue1 is a Commented of type Value1
    // Commented value 1
    CommentedValue1 Commented = iota
    // CommentedValue2 is a Commented of type Value2
    CommentedValue2
    // CommentedValue3 is a Commented of type Value3
    // Commented value 3
    CommentedValue3
)
...
```

#### Example

There are a few examples in the `example` [directory](./example/).
I've included one here for easy access, but can't guarantee it's up to date.

``` go
// Color is an enumeration of colors that are allowed.
/* ENUM(
Black, White, Red
Green = 33 // Green starts with 33
*/
// Blue
// grey=
// yellow
// blue-green
// red-orange
// yellow_green
// red-orange-blue
// )
type Color int32
```

The generated code will look something like:

``` go
// Code generated by go-enum DO NOT EDIT.
// Version: example
// Revision: example
// Build Date: example
// Built By: example

package example

import (
 "fmt"
 "strings"
)

const (
 // ColorBlack is a Color of type Black.
 ColorBlack Color = iota
 // ColorWhite is a Color of type White.
 ColorWhite
 // ColorRed is a Color of type Red.
 ColorRed
 // ColorGreen is a Color of type Green.
 // Green starts with 33
 ColorGreen Color = iota + 30
 // ColorBlue is a Color of type Blue.
 ColorBlue
 // ColorGrey is a Color of type Grey.
 ColorGrey
 // ColorYellow is a Color of type Yellow.
 ColorYellow
 // ColorBlueGreen is a Color of type Blue-Green.
 ColorBlueGreen
 // ColorRedOrange is a Color of type Red-Orange.
 ColorRedOrange
 // ColorYellowGreen is a Color of type Yellow_green.
 ColorYellowGreen
 // ColorRedOrangeBlue is a Color of type Red-Orange-Blue.
 ColorRedOrangeBlue
)

const _ColorName = "BlackWhiteRedGreenBluegreyyellowblue-greenred-orangeyellow_greenred-orange-blue"

var _ColorMap = map[Color]string{
 ColorBlack:         _ColorName[0:5],
 ColorWhite:         _ColorName[5:10],
 ColorRed:           _ColorName[10:13],
 ColorGreen:         _ColorName[13:18],
 ColorBlue:          _ColorName[18:22],
 ColorGrey:          _ColorName[22:26],
 ColorYellow:        _ColorName[26:32],
 ColorBlueGreen:     _ColorName[32:42],
 ColorRedOrange:     _ColorName[42:52],
 ColorYellowGreen:   _ColorName[52:64],
 ColorRedOrangeBlue: _ColorName[64:79],
}

// String implements the Stringer interface.
func (x Color) String() string {
 if str, ok := _ColorMap[x]; ok {
  return str
 }
 return fmt.Sprintf("Color(%d)", x)
}

var _ColorValue = map[string]Color{
 _ColorName[0:5]:                    ColorBlack,
 strings.ToLower(_ColorName[0:5]):   ColorBlack,
 _ColorName[5:10]:                   ColorWhite,
 strings.ToLower(_ColorName[5:10]):  ColorWhite,
 _ColorName[10:13]:                  ColorRed,
 strings.ToLower(_ColorName[10:13]): ColorRed,
 _ColorName[13:18]:                  ColorGreen,
 strings.ToLower(_ColorName[13:18]): ColorGreen,
 _ColorName[18:22]:                  ColorBlue,
 strings.ToLower(_ColorName[18:22]): ColorBlue,
 _ColorName[22:26]:                  ColorGrey,
 strings.ToLower(_ColorName[22:26]): ColorGrey,
 _ColorName[26:32]:                  ColorYellow,
 strings.ToLower(_ColorName[26:32]): ColorYellow,
 _ColorName[32:42]:                  ColorBlueGreen,
 strings.ToLower(_ColorName[32:42]): ColorBlueGreen,
 _ColorName[42:52]:                  ColorRedOrange,
 strings.ToLower(_ColorName[42:52]): ColorRedOrange,
 _ColorName[52:64]:                  ColorYellowGreen,
 strings.ToLower(_ColorName[52:64]): ColorYellowGreen,
 _ColorName[64:79]:                  ColorRedOrangeBlue,
 strings.ToLower(_ColorName[64:79]): ColorRedOrangeBlue,
}

// ParseColor attempts to convert a string to a Color
func ParseColor(name string) (Color, error) {
 if x, ok := _ColorValue[name]; ok {
  return x, nil
 }
 return Color(0), fmt.Errorf("%s is not a valid Color", name)
}

func (x Color) Ptr() *Color {
 return &x
}

// MarshalText implements the text marshaller method
func (x Color) MarshalText() ([]byte, error) {
 return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method
func (x *Color) UnmarshalText(text []byte) error {
 name := string(text)
 tmp, err := ParseColor(name)
 if err != nil {
  return err
 }
 *x = tmp
 return nil
}
```
