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
	// AcmeInt_SOME_PLACE_AWESOME is a IntShop of type SOME_PLACE_AWESOME.
	AcmeInt_SOME_PLACE_AWESOME IntShop = iota
	// AcmeInt_SomewhereElse is a IntShop of type SomewhereElse.
	AcmeInt_SomewhereElse
	// AcmeInt_LocationUnknown is a IntShop of type LocationUnknown.
	AcmeInt_LocationUnknown
)

var ErrInvalidIntShop = fmt.Errorf("not a valid IntShop, try [%s]", strings.Join(_IntShopNames, ", "))

const _IntShopName = "SOME_PLACE_AWESOMESomewhereElseLocationUnknown"

var _IntShopNames = []string{
	_IntShopName[0:18],
	_IntShopName[18:31],
	_IntShopName[31:46],
}

// IntShopNames returns a list of possible string values of IntShop.
func IntShopNames() []string {
	tmp := make([]string, len(_IntShopNames))
	copy(tmp, _IntShopNames)
	return tmp
}

var _IntShopMap = map[IntShop]string{
	AcmeInt_SOME_PLACE_AWESOME: _IntShopName[0:18],
	AcmeInt_SomewhereElse:      _IntShopName[18:31],
	AcmeInt_LocationUnknown:    _IntShopName[31:46],
}

// String implements the Stringer interface.
func (x IntShop) String() string {
	if str, ok := _IntShopMap[x]; ok {
		return str
	}
	return fmt.Sprintf("IntShop(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x IntShop) IsValid() bool {
	_, ok := _IntShopMap[x]
	return ok
}

var _IntShopValue = map[string]IntShop{
	_IntShopName[0:18]:  AcmeInt_SOME_PLACE_AWESOME,
	_IntShopName[18:31]: AcmeInt_SomewhereElse,
	_IntShopName[31:46]: AcmeInt_LocationUnknown,
}

// ParseIntShop attempts to convert a string to a IntShop.
func ParseIntShop(name string) (IntShop, error) {
	if x, ok := _IntShopValue[name]; ok {
		return x, nil
	}
	return IntShop(0), fmt.Errorf("%s is %w", name, ErrInvalidIntShop)
}

// MarshalText implements the text marshaller method.
func (x IntShop) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *IntShop) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseIntShop(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
