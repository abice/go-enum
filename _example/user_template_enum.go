// Code generated by go-enum DO NOT EDIT.
// Version: example
// Revision: example
// Build Date: example
// Built By: example

package example

import (
	"errors"
	"fmt"
)

const (
	// OceanColorCerulean is a OceanColor of type Cerulean.
	OceanColorCerulean OceanColor = iota
	// OceanColorBlue is a OceanColor of type Blue.
	OceanColorBlue
	// OceanColorGreen is a OceanColor of type Green.
	OceanColorGreen
)

var ErrInvalidOceanColor = errors.New("not a valid OceanColor")

const _oceanColorName = "CeruleanBlueGreen"

var _oceanColorMap = map[OceanColor]string{
	OceanColorCerulean: _oceanColorName[0:8],
	OceanColorBlue:     _oceanColorName[8:12],
	OceanColorGreen:    _oceanColorName[12:17],
}

// String implements the Stringer interface.
func (x OceanColor) String() string {
	if str, ok := _oceanColorMap[x]; ok {
		return str
	}
	return fmt.Sprintf("OceanColor(%d)", x)
}

var _oceanColorValue = map[string]OceanColor{
	_oceanColorName[0:8]:   OceanColorCerulean,
	_oceanColorName[8:12]:  OceanColorBlue,
	_oceanColorName[12:17]: OceanColorGreen,
}

// ParseOceanColor attempts to convert a string to a OceanColor.
func ParseOceanColor(name string) (OceanColor, error) {
	if x, ok := _oceanColorValue[name]; ok {
		return x, nil
	}
	return OceanColor(0), fmt.Errorf("%s is %w", name, ErrInvalidOceanColor)
}

func ParseOceanColorGlobbedExample() bool {
	return true
}
func ParseOceanColorGlobbedExample2() bool {
	return true
}

// Additional template
func ParseOceanColorExample() bool {
	return true
}
