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
	// Enum32bitUnkno is a Enum32bit of type Unkno.
	Enum32bitUnkno Enum32bit = iota
	// Enum32bitE2P15 is a Enum32bit of type E2P15.
	Enum32bitE2P15 Enum32bit = iota + 32767
	// Enum32bitE2P16 is a Enum32bit of type E2P16.
	Enum32bitE2P16 Enum32bit = iota + 65534
	// Enum32bitE2P17 is a Enum32bit of type E2P17.
	Enum32bitE2P17 Enum32bit = iota + 131069
	// Enum32bitE2P18 is a Enum32bit of type E2P18.
	Enum32bitE2P18 Enum32bit = iota + 262140
	// Enum32bitE2P19 is a Enum32bit of type E2P19.
	Enum32bitE2P19 Enum32bit = iota + 524283
	// Enum32bitE2P20 is a Enum32bit of type E2P20.
	Enum32bitE2P20 Enum32bit = iota + 1048570
	// Enum32bitE2P21 is a Enum32bit of type E2P21.
	Enum32bitE2P21 Enum32bit = iota + 2097145
	// Enum32bitE2P22 is a Enum32bit of type E2P22.
	Enum32bitE2P22 Enum32bit = iota + 33554424
	// Enum32bitE2P23 is a Enum32bit of type E2P23.
	Enum32bitE2P23 Enum32bit = iota + 67108855
	// Enum32bitE2P28 is a Enum32bit of type E2P28.
	Enum32bitE2P28 Enum32bit = iota + 536870902
	// Enum32bitE2P30 is a Enum32bit of type E2P30.
	Enum32bitE2P30 Enum32bit = iota + 1073741813
)

const _Enum32bitName = "UnknoE2P15E2P16E2P17E2P18E2P19E2P20E2P21E2P22E2P23E2P28E2P30"

var _Enum32bitNames = []string{
	_Enum32bitName[0:5],
	_Enum32bitName[5:10],
	_Enum32bitName[10:15],
	_Enum32bitName[15:20],
	_Enum32bitName[20:25],
	_Enum32bitName[25:30],
	_Enum32bitName[30:35],
	_Enum32bitName[35:40],
	_Enum32bitName[40:45],
	_Enum32bitName[45:50],
	_Enum32bitName[50:55],
	_Enum32bitName[55:60],
}

// Enum32bitNames returns a list of possible string values of Enum32bit.
func Enum32bitNames() []string {
	tmp := make([]string, len(_Enum32bitNames))
	copy(tmp, _Enum32bitNames)
	return tmp
}

var _Enum32bitMap = map[Enum32bit]string{
	Enum32bitUnkno: _Enum32bitName[0:5],
	Enum32bitE2P15: _Enum32bitName[5:10],
	Enum32bitE2P16: _Enum32bitName[10:15],
	Enum32bitE2P17: _Enum32bitName[15:20],
	Enum32bitE2P18: _Enum32bitName[20:25],
	Enum32bitE2P19: _Enum32bitName[25:30],
	Enum32bitE2P20: _Enum32bitName[30:35],
	Enum32bitE2P21: _Enum32bitName[35:40],
	Enum32bitE2P22: _Enum32bitName[40:45],
	Enum32bitE2P23: _Enum32bitName[45:50],
	Enum32bitE2P28: _Enum32bitName[50:55],
	Enum32bitE2P30: _Enum32bitName[55:60],
}

// String implements the Stringer interface.
func (x Enum32bit) String() string {
	if str, ok := _Enum32bitMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Enum32bit(%d)", x)
}

var _Enum32bitValue = map[string]Enum32bit{
	_Enum32bitName[0:5]:   Enum32bitUnkno,
	_Enum32bitName[5:10]:  Enum32bitE2P15,
	_Enum32bitName[10:15]: Enum32bitE2P16,
	_Enum32bitName[15:20]: Enum32bitE2P17,
	_Enum32bitName[20:25]: Enum32bitE2P18,
	_Enum32bitName[25:30]: Enum32bitE2P19,
	_Enum32bitName[30:35]: Enum32bitE2P20,
	_Enum32bitName[35:40]: Enum32bitE2P21,
	_Enum32bitName[40:45]: Enum32bitE2P22,
	_Enum32bitName[45:50]: Enum32bitE2P23,
	_Enum32bitName[50:55]: Enum32bitE2P28,
	_Enum32bitName[55:60]: Enum32bitE2P30,
}

// ParseEnum32bit attempts to convert a string to a Enum32bit.
func ParseEnum32bit(name string) (Enum32bit, error) {
	if x, ok := _Enum32bitValue[name]; ok {
		return x, nil
	}
	return Enum32bit(0), fmt.Errorf("%s is not a valid Enum32bit, try [%s]", name, strings.Join(_Enum32bitNames, ", "))
}
