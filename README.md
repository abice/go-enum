# go-enum
[![CircleCI](https://circleci.com/gh/abice/go-enum.svg?style=svg&circle-token=b44c10ce16bcef76e86da801d67811a5ff71fc72)](https://circleci.com/gh/abice/go-enum)
[![Go Report Card](https://goreportcard.com/badge/github.com/abice/go-enum)](https://goreportcard.com/report/github.com/abice/go-enum)
[![Coverage Status](https://coveralls.io/repos/github/abice/go-enum/badge.svg)](https://coveralls.io/github/abice/go-enum)
[![GoDoc](https://godoc.org/github.com/abice/go-enum?status.svg)](https://godoc.org/github.com/abice/go-enum)

An enum generator for go


## How it works

The goal of go-enum is to create an easy to use enum generator that will take a decorated type declaration like `type EnumName int` and create the associated constant values and funcs that will make life a little easier for adding new values.
It's not perfect, but I think it's useful.

I took the output of the [Stringer](golang.org/x/tools/cmd/stringer) command as the `String()` method, and added a way to parse a string value.


### Syntax
The parser looks for comments on your type defs and parse the enum declarations from it.  
The parser will look for `ENUM(` and continue to look for comma separated values until it finds a `)`.  You can put values on the same line, or on multiple lines.

There are a few examples in the `example` [directory](repo/blob/master/example).
I've included one here for easy access, but can't guarantee it's up to date.

``` go
// Color is an enumeration of colors that are allowed.
// ENUM(
// Black, White, Red
// Green 
// Blue
// grey
// yellow
// )
type Color int32
```

The generated code will look something like:

``` go

const (
	// ColorBlack is a Color of type Black
	ColorBlack Color = iota
	// ColorWhite is a Color of type White
	ColorWhite
	// ColorRed is a Color of type Red
	ColorRed
	// ColorGreen is a Color of type Green
	ColorGreen
	// ColorBlue is a Color of type Blue
	ColorBlue
	// ColorGrey is a Color of type Grey
	ColorGrey
	// ColorYellow is a Color of type Yellow
	ColorYellow
)

const _ColorName = "BlackWhiteRedGreenBlueGreyYellow"

var _ColorMap = map[Color]string{
	0: _ColorName[0:5],
	1: _ColorName[5:10],
	2: _ColorName[10:13],
	3: _ColorName[13:18],
	4: _ColorName[18:22],
	5: _ColorName[22:26],
	6: _ColorName[26:32],
}

func (i Color) String() string {
	if str, ok := _ColorMap[i]; ok {
		return str
	}
	return fmt.Sprintf("Color(%d)", i)
}

var _ColorValue = map[string]Color{
	_ColorName[0:5]:                    0,
	strings.ToLower(_ColorName[0:5]):   0,
	_ColorName[5:10]:                   1,
	strings.ToLower(_ColorName[5:10]):  1,
	_ColorName[10:13]:                  2,
	strings.ToLower(_ColorName[10:13]): 2,
	_ColorName[13:18]:                  3,
	strings.ToLower(_ColorName[13:18]): 3,
	_ColorName[18:22]:                  4,
	strings.ToLower(_ColorName[18:22]): 4,
	_ColorName[22:26]:                  5,
	strings.ToLower(_ColorName[22:26]): 5,
	_ColorName[26:32]:                  6,
	strings.ToLower(_ColorName[26:32]): 6,
}

// ParseColor attempts to convert a string to a Color
func ParseColor(name string) (Color, error) {
	if x, ok := _ColorValue[name]; ok {
		return Color(x), nil
	}
	return Color(0), fmt.Errorf("%s is not a valid Color", name)
}

func (x *Color) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

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


## Adding it to your project
1. `go get github.com/abice/go-enum`
1. Add a go:generate line to your file like so... `//go:generate go-enum -f=thisfile.go`
1. Run go generate like so `go generate ./...`
1. Enjoy your newly created Enumeration


## Options

``` shell
go-enum --help
Options:

  -h, --help       display help information
  -f, --file      *The file(s) to generate enums.  Use more than one flag for more files.
      --noprefix   Prevents the constants generated from having the Enum as a prefix.
      --lower      Adds lowercase variants of the enum strings for lookup.
      --marshal    Adds text marshalling functions.
```