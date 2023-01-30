// Code generated by go-enum DO NOT EDIT.
// Version: example
// Revision: example
// Build Date: example
// Built By: example

package example

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	StrStatePending   StrState = "pending"
	StrStateRunning   StrState = "running"
	StrStateCompleted StrState = "completed"
	StrStateFailed    StrState = "error"
)

var _strStateNames = []string{
	string(StrStatePending),
	string(StrStateRunning),
	string(StrStateCompleted),
	string(StrStateFailed),
}

// StrStateNames returns a list of possible string values of StrState.
func StrStateNames() []string {
	tmp := make([]string, len(_strStateNames))
	copy(tmp, _strStateNames)
	return tmp
}

var ErrInvalidStrState = fmt.Errorf("not a valid StrState, try [%s]", strings.Join(_strStateNames, ", "))

// StrStateValues returns a list of the values for StrState
func StrStateValues() []StrState {
	return []StrState{
		StrStatePending,
		StrStateRunning,
		StrStateCompleted,
		StrStateFailed,
	}
}

// String implements the Stringer interface.
func (x StrState) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x StrState) IsValid() bool {
	_, err := ParseStrState(string(x))
	return err == nil
}

var _strStateValue = map[string]StrState{
	"pending":   StrStatePending,
	"running":   StrStateRunning,
	"completed": StrStateCompleted,
	"error":     StrStateFailed,
}

// ParseStrState attempts to convert a string to a StrState.
func ParseStrState(name string) (StrState, error) {
	if x, ok := _strStateValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _strStateValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return StrState(""), fmt.Errorf("%s is %w", name, ErrInvalidStrState)
}

// MustParseStrState converts a string to a StrState, and panics if is not valid.
func MustParseStrState(name string) StrState {
	val, err := ParseStrState(name)
	if err != nil {
		panic(err)
	}
	return val
}

func (x StrState) Ptr() *StrState {
	return &x
}

// MarshalText implements the text marshaller method.
func (x StrState) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *StrState) UnmarshalText(text []byte) error {
	tmp, err := ParseStrState(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var ErrStrStateNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *StrState) Scan(value interface{}) (err error) {
	if value == nil {
		*x = StrState("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseStrState(v)
	case []byte:
		*x, err = ParseStrState(string(v))
	case StrState:
		*x = v
	case *StrState:
		if v == nil {
			return ErrStrStateNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return ErrStrStateNilPtr
		}
		*x, err = ParseStrState(*v)
	default:
		return errors.New("invalid type for StrState")
	}

	return
}

// Value implements the driver Valuer interface.
func (x StrState) Value() (driver.Value, error) {
	return x.String(), nil
}

// Set implements the Golang flag.Value interface func.
func (x *StrState) Set(val string) error {
	v, err := ParseStrState(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *StrState) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *StrState) Type() string {
	return "StrState"
}

type NullStrState struct {
	StrState StrState
	Valid    bool
	Set      bool
}

func NewNullStrState(val interface{}) (x NullStrState) {
	err := x.Scan(val) // yes, we ignore this error, it will just be an invalid value.
	_ = err            // make any errcheck linters happy
	return
}

// Scan implements the Scanner interface.
func (x *NullStrState) Scan(value interface{}) (err error) {
	if value == nil {
		x.StrState, x.Valid = StrState(""), false
		return
	}

	err = x.StrState.Scan(value)
	x.Valid = (err == nil)
	return
}

// Value implements the driver Valuer interface.
func (x NullStrState) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	return x.StrState.String(), nil
}

// MarshalJSON correctly serializes a NullStrState to JSON.
func (n NullStrState) MarshalJSON() ([]byte, error) {
	const nullStr = "null"
	if n.Valid {
		return json.Marshal(n.StrState)
	}
	return []byte(nullStr), nil
}

// UnmarshalJSON correctly deserializes a NullStrState from JSON.
func (n *NullStrState) UnmarshalJSON(b []byte) error {
	n.Set = true
	var x interface{}
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	err = n.Scan(x)
	return err
}
