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
	// AllNegativeUnknown is a AllNegative of type Unknown.
	AllNegativeUnknown AllNegative = iota + -5
	// AllNegativeGood is a AllNegative of type Good.
	AllNegativeGood
	// AllNegativeBad is a AllNegative of type Bad.
	AllNegativeBad
	// AllNegativeUgly is a AllNegative of type Ugly.
	AllNegativeUgly
)

const _AllNegativeName = "UnknownGoodBadUgly"

var _AllNegativeMap = map[AllNegative]string{
	AllNegativeUnknown: _AllNegativeName[0:7],
	AllNegativeGood:    _AllNegativeName[7:11],
	AllNegativeBad:     _AllNegativeName[11:14],
	AllNegativeUgly:    _AllNegativeName[14:18],
}

// String implements the Stringer interface.
func (x AllNegative) String() string {
	if str, ok := _AllNegativeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("AllNegative(%d)", x)
}

var _AllNegativeValue = map[string]AllNegative{
	_AllNegativeName[0:7]:                    AllNegativeUnknown,
	strings.ToLower(_AllNegativeName[0:7]):   AllNegativeUnknown,
	_AllNegativeName[7:11]:                   AllNegativeGood,
	strings.ToLower(_AllNegativeName[7:11]):  AllNegativeGood,
	_AllNegativeName[11:14]:                  AllNegativeBad,
	strings.ToLower(_AllNegativeName[11:14]): AllNegativeBad,
	_AllNegativeName[14:18]:                  AllNegativeUgly,
	strings.ToLower(_AllNegativeName[14:18]): AllNegativeUgly,
}

// ParseAllNegative attempts to convert a string to a AllNegative.
func ParseAllNegative(name string) (AllNegative, error) {
	if x, ok := _AllNegativeValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _AllNegativeValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return AllNegative(0), fmt.Errorf("%s is not a valid AllNegative", name)
}

const (
	// StatusUnknown is a Status of type Unknown.
	StatusUnknown Status = iota + -1
	// StatusGood is a Status of type Good.
	StatusGood
	// StatusBad is a Status of type Bad.
	StatusBad
)

const _StatusName = "UnknownGoodBad"

var _StatusMap = map[Status]string{
	StatusUnknown: _StatusName[0:7],
	StatusGood:    _StatusName[7:11],
	StatusBad:     _StatusName[11:14],
}

// String implements the Stringer interface.
func (x Status) String() string {
	if str, ok := _StatusMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Status(%d)", x)
}

var _StatusValue = map[string]Status{
	_StatusName[0:7]:                    StatusUnknown,
	strings.ToLower(_StatusName[0:7]):   StatusUnknown,
	_StatusName[7:11]:                   StatusGood,
	strings.ToLower(_StatusName[7:11]):  StatusGood,
	_StatusName[11:14]:                  StatusBad,
	strings.ToLower(_StatusName[11:14]): StatusBad,
}

// ParseStatus attempts to convert a string to a Status.
func ParseStatus(name string) (Status, error) {
	if x, ok := _StatusValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _StatusValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Status(0), fmt.Errorf("%s is not a valid Status", name)
}
