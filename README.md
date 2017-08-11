# go-enum
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

var _ColorIndex = [...]uint8{0, 5, 10, 13, 18, 22, 26, 32}

func (i Color) String() string {
	if i < 0 || i >= Color(len(_ColorIndex)-1) {
		return fmt.Sprintf("Color(%d)", i)
	}
	return _ColorName[_ColorIndex[i]:_ColorIndex[i+1]]
}

var _ColorValue = map[string]int{

	"Black":  0,
	"White":  1,
	"Red":    2,
	"Green":  3,
	"Blue":   4,
	"Grey":   5,
	"Yellow": 6,
	"black":  0,
	"white":  1,
	"red":    2,
	"green":  3,
	"blue":   4,
	"grey":   5,
	"yellow": 6,
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

1. Add a go:generate line to your file like so... `//go:generate go-enum -f=thisfile.go`

2. Run go generate like so `go generate ./...`

3. Enjoy your newly created Enumeration


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