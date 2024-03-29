// Code generated by go-enum DO NOT EDIT.
// Version: example
// Revision: example
// Build Date: example
// Built By: example

//go:build example
// +build example

package example

import (
	"errors"
	"fmt"
)

const (
	// ForceUpperTypeDataSwap is a ForceUpperType of type DataSwap.
	ForceUpperTypeDataSwap ForceUpperType = iota
	// ForceUpperTypeBootNode is a ForceUpperType of type BootNode.
	ForceUpperTypeBootNode
)

var ErrInvalidForceUpperType = errors.New("not a valid ForceUpperType")

const _ForceUpperTypeName = "DATASWAPBOOTNODE"

var _ForceUpperTypeMap = map[ForceUpperType]string{
	ForceUpperTypeDataSwap: _ForceUpperTypeName[0:8],
	ForceUpperTypeBootNode: _ForceUpperTypeName[8:16],
}

// String implements the Stringer interface.
func (x ForceUpperType) String() string {
	if str, ok := _ForceUpperTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ForceUpperType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ForceUpperType) IsValid() bool {
	_, ok := _ForceUpperTypeMap[x]
	return ok
}

var _ForceUpperTypeValue = map[string]ForceUpperType{
	_ForceUpperTypeName[0:8]:  ForceUpperTypeDataSwap,
	_ForceUpperTypeName[8:16]: ForceUpperTypeBootNode,
}

// ParseForceUpperType attempts to convert a string to a ForceUpperType.
func ParseForceUpperType(name string) (ForceUpperType, error) {
	if x, ok := _ForceUpperTypeValue[name]; ok {
		return x, nil
	}
	return ForceUpperType(0), fmt.Errorf("%s is %w", name, ErrInvalidForceUpperType)
}
