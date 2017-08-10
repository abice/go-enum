# go-enum
An enum generator for go

## How it works

The goal of go-enum is to create an easy to use enum generator that will take a decorated type declaration like `type EnumName int` and create the associated constant values and funcs that will make life a little easier for adding new values.
It's not perfect, but I think it's useful.

I took the output of the [Stringer](golang.org/x/tools/cmd/stringer) command as the `String()` method, and added a way to parse a string value.


### Syntax
The parser looks for comments on your type defs and parse the enum declarations from it.  
The parser will look for `ENUM(` and continue to look for comma separated values until it finds a `)`.  You can put values on the same line, or on multiple lines.
Here are a few examples of decorated values that will be parsed into enums:

```
// SingleLine enumeration...
// ENUM(One, two, three)
type SingleLine int
```

```
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

```
const (
	// ColorBlack is a Color of type Black
	ColorBlack Color = iota
	// ColorWhite is a Color of type White
	ColorWhite Color = iota
	// ColorRed is a Color of type Red
	ColorRed Color = iota
	// ColorGreen is a Color of type Green
	ColorGreen Color = iota
	// ColorBlue is a Color of type Blue
	ColorBlue Color = iota
	// ColorGrey is a Color of type Grey
	ColorGrey Color = iota
	// ColorYellow is a Color of type Yellow
	ColorYellow Color = iota
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
}

func ParseColor(name string) Color {
	val := Color(0)
	if x, ok := _ColorValue[name]; ok {
		val = Color(x)
	}
	return val
}
```


## Adding it to your project

1. Add a go:generate line to your file like so... `// go:generate go-enum -f=thisfile.go`

2. Run go generate like so `go generate ./...`

3. Enjoy your newly created Enumeration
