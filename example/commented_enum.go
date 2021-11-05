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
	// CommentedValue1 is a Commented of type Value1.
	// Commented value 1
	CommentedValue1 Commented = iota
	// CommentedValue2 is a Commented of type Value2.
	CommentedValue2
	// CommentedValue3 is a Commented of type Value3.
	// Commented value 3
	CommentedValue3
)

const _CommentedName = "value1value2value3"

var _CommentedMap = map[Commented]string{
	CommentedValue1: _CommentedName[0:6],
	CommentedValue2: _CommentedName[6:12],
	CommentedValue3: _CommentedName[12:18],
}

// String implements the Stringer interface.
func (x Commented) String() string {
	if str, ok := _CommentedMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Commented(%d)", x)
}

var _CommentedValue = map[string]Commented{
	_CommentedName[0:6]:                    CommentedValue1,
	strings.ToLower(_CommentedName[0:6]):   CommentedValue1,
	_CommentedName[6:12]:                   CommentedValue2,
	strings.ToLower(_CommentedName[6:12]):  CommentedValue2,
	_CommentedName[12:18]:                  CommentedValue3,
	strings.ToLower(_CommentedName[12:18]): CommentedValue3,
}

// ParseCommented attempts to convert a string to a Commented
func ParseCommented(name string) (Commented, error) {
	if x, ok := _CommentedValue[name]; ok {
		return x, nil
	}
	return Commented(0), fmt.Errorf("%s is not a valid Commented", name)
}

// MarshalText implements the text marshaller method
func (x Commented) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method
func (x *Commented) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseCommented(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

const (
	// Skipped value.
	// Placeholder with a ','  in it. (for harder testing)
	_ ComplexCommented = iota
	// ComplexCommentedValue1 is a ComplexCommented of type Value1.
	// Commented value 1
	ComplexCommentedValue1
	// ComplexCommentedValue2 is a ComplexCommented of type Value2.
	ComplexCommentedValue2
	// ComplexCommentedValue3 is a ComplexCommented of type Value3.
	// Commented value 3
	ComplexCommentedValue3
)

const _ComplexCommentedName = "value1value2value3"

var _ComplexCommentedMap = map[ComplexCommented]string{
	ComplexCommentedValue1: _ComplexCommentedName[0:6],
	ComplexCommentedValue2: _ComplexCommentedName[6:12],
	ComplexCommentedValue3: _ComplexCommentedName[12:18],
}

// String implements the Stringer interface.
func (x ComplexCommented) String() string {
	if str, ok := _ComplexCommentedMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ComplexCommented(%d)", x)
}

var _ComplexCommentedValue = map[string]ComplexCommented{
	_ComplexCommentedName[0:6]:                    ComplexCommentedValue1,
	strings.ToLower(_ComplexCommentedName[0:6]):   ComplexCommentedValue1,
	_ComplexCommentedName[6:12]:                   ComplexCommentedValue2,
	strings.ToLower(_ComplexCommentedName[6:12]):  ComplexCommentedValue2,
	_ComplexCommentedName[12:18]:                  ComplexCommentedValue3,
	strings.ToLower(_ComplexCommentedName[12:18]): ComplexCommentedValue3,
}

// ParseComplexCommented attempts to convert a string to a ComplexCommented
func ParseComplexCommented(name string) (ComplexCommented, error) {
	if x, ok := _ComplexCommentedValue[name]; ok {
		return x, nil
	}
	return ComplexCommented(0), fmt.Errorf("%s is not a valid ComplexCommented", name)
}

// MarshalText implements the text marshaller method
func (x ComplexCommented) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method
func (x *ComplexCommented) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseComplexCommented(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
