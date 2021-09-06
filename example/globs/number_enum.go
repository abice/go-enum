// Code generated by go-enum DO NOT EDIT.
// Version: example
// Revision: example
// Build Date: example
// Built By: example

package globs

import (
	"fmt"
)

const (
	// Number0 is a Number of type 0.
	Number0 Number = iota
	// Number1 is a Number of type 1.
	Number1
	// Number2 is a Number of type 2.
	Number2
	// Number3 is a Number of type 3.
	Number3
	// Number4 is a Number of type 4.
	Number4
	// Number5 is a Number of type 5.
	Number5
	// Number6 is a Number of type 6.
	Number6
	// Number7 is a Number of type 7.
	Number7
	// Number8 is a Number of type 8.
	Number8
	// Number9 is a Number of type 9.
	Number9
)

const _NumberName = "0123456789"

var _NumberMap = map[Number]string{
	0: _NumberName[0:1],
	1: _NumberName[1:2],
	2: _NumberName[2:3],
	3: _NumberName[3:4],
	4: _NumberName[4:5],
	5: _NumberName[5:6],
	6: _NumberName[6:7],
	7: _NumberName[7:8],
	8: _NumberName[8:9],
	9: _NumberName[9:10],
}

// String implements the Stringer interface.
func (x Number) String() string {
	if str, ok := _NumberMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Number(%d)", x)
}

var _NumberValue = map[string]Number{
	_NumberName[0:1]:  0,
	_NumberName[1:2]:  1,
	_NumberName[2:3]:  2,
	_NumberName[3:4]:  3,
	_NumberName[4:5]:  4,
	_NumberName[5:6]:  5,
	_NumberName[6:7]:  6,
	_NumberName[7:8]:  7,
	_NumberName[8:9]:  8,
	_NumberName[9:10]: 9,
}

// ParseNumber attempts to convert a string to a Number
func ParseNumber(name string) (Number, error) {
	if x, ok := _NumberValue[name]; ok {
		return x, nil
	}
	return Number(0), fmt.Errorf("%s is not a valid Number", name)
}
