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
	"strconv"
)

const (
	// ProjectStatusPending is a ProjectStatus of type Pending.
	ProjectStatusPending ProjectStatus = iota
	// ProjectStatusInWork is a ProjectStatus of type InWork.
	ProjectStatusInWork
	// ProjectStatusCompleted is a ProjectStatus of type Completed.
	ProjectStatusCompleted
	// ProjectStatusRejected is a ProjectStatus of type Rejected.
	ProjectStatusRejected
)

const _ProjectStatusName = "pendinginWorkcompletedrejected"

var _ProjectStatusMap = map[ProjectStatus]string{
	ProjectStatusPending:   _ProjectStatusName[0:7],
	ProjectStatusInWork:    _ProjectStatusName[7:13],
	ProjectStatusCompleted: _ProjectStatusName[13:22],
	ProjectStatusRejected:  _ProjectStatusName[22:30],
}

// String implements the Stringer interface.
func (x ProjectStatus) String() string {
	if str, ok := _ProjectStatusMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ProjectStatus(%d)", x)
}

var _ProjectStatusValue = map[string]ProjectStatus{
	_ProjectStatusName[0:7]:   ProjectStatusPending,
	_ProjectStatusName[7:13]:  ProjectStatusInWork,
	_ProjectStatusName[13:22]: ProjectStatusCompleted,
	_ProjectStatusName[22:30]: ProjectStatusRejected,
}

// ParseProjectStatus attempts to convert a string to a ProjectStatus.
func ParseProjectStatus(name string) (ProjectStatus, error) {
	if x, ok := _ProjectStatusValue[name]; ok {
		return x, nil
	}
	return ProjectStatus(0), fmt.Errorf("%s is not a valid ProjectStatus", name)
}

func (x ProjectStatus) Ptr() *ProjectStatus {
	return &x
}

// MarshalText implements the text marshaller method.
func (x ProjectStatus) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ProjectStatus) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseProjectStatus(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var _ProjectStatusErrNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *ProjectStatus) Scan(value interface{}) (err error) {
	if value == nil {
		*x = ProjectStatus(0)
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case int64:
		*x = ProjectStatus(v)
	case string:
		*x, err = ParseProjectStatus(v)
		if err != nil {
			// try parsing the integer value as a string
			if val, verr := strconv.Atoi(v); verr == nil {
				*x, err = ProjectStatus(val), nil
			}
		}
	case []byte:
		*x, err = ParseProjectStatus(string(v))
		if err != nil {
			// try parsing the integer value as a string
			if val, verr := strconv.Atoi(string(v)); verr == nil {
				*x, err = ProjectStatus(val), nil
			}
		}
	case ProjectStatus:
		*x = v
	case int:
		*x = ProjectStatus(v)
	case *ProjectStatus:
		if v == nil {
			return _ProjectStatusErrNilPtr
		}
		*x = *v
	case uint:
		*x = ProjectStatus(v)
	case uint64:
		*x = ProjectStatus(v)
	case *int:
		if v == nil {
			return _ProjectStatusErrNilPtr
		}
		*x = ProjectStatus(*v)
	case *int64:
		if v == nil {
			return _ProjectStatusErrNilPtr
		}
		*x = ProjectStatus(*v)
	case float64: // json marshals everything as a float64 if it's a number
		*x = ProjectStatus(v)
	case *float64: // json marshals everything as a float64 if it's a number
		if v == nil {
			return _ProjectStatusErrNilPtr
		}
		*x = ProjectStatus(*v)
	case *uint:
		if v == nil {
			return _ProjectStatusErrNilPtr
		}
		*x = ProjectStatus(*v)
	case *uint64:
		if v == nil {
			return _ProjectStatusErrNilPtr
		}
		*x = ProjectStatus(*v)
	case *string:
		if v == nil {
			return _ProjectStatusErrNilPtr
		}
		*x, err = ParseProjectStatus(*v)
		if err != nil {
			// try parsing the integer value as a string
			if val, verr := strconv.Atoi(*v); verr == nil {
				*x, err = ProjectStatus(val), nil
			}
		}
	}

	return
}

// Value implements the driver Valuer interface.
func (x ProjectStatus) Value() (driver.Value, error) {
	return x.String(), nil
}

type NullProjectStatus struct {
	ProjectStatus ProjectStatus
	Valid         bool
	Set           bool
}

func NewNullProjectStatus(val interface{}) (x NullProjectStatus) {
	x.Scan(val) // yes, we ignore this error, it will just be an invalid value.
	return
}

// Scan implements the Scanner interface.
func (x *NullProjectStatus) Scan(value interface{}) (err error) {
	x.Set = true
	if value == nil {
		x.ProjectStatus, x.Valid = ProjectStatus(0), false
		return
	}

	err = x.ProjectStatus.Scan(value)
	x.Valid = (err == nil)
	return
}

// Value implements the driver Valuer interface.
func (x NullProjectStatus) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	// driver.Value accepts int64 for int values.
	return int64(x.ProjectStatus), nil
}

// MarshalJSON correctly serializes a NullProjectStatus to JSON.
func (n NullProjectStatus) MarshalJSON() ([]byte, error) {
	const nullStr = "null"
	if n.Valid {
		return json.Marshal(n.ProjectStatus)
	}
	return []byte(nullStr), nil
}

// UnmarshalJSON correctly deserializes a NullProjectStatus from JSON.
func (n *NullProjectStatus) UnmarshalJSON(b []byte) error {
	n.Set = true
	var x interface{}
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	err = n.Scan(x)
	return err
}

type NullProjectStatusStr struct {
	NullProjectStatus
}

func NewNullProjectStatusStr(val interface{}) (x NullProjectStatusStr) {
	x.Scan(val) // yes, we ignore this error, it will just be an invalid value.
	return
}

// Value implements the driver Valuer interface.
func (x NullProjectStatusStr) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	return x.ProjectStatus.String(), nil
}

// MarshalJSON correctly serializes a NullProjectStatus to JSON.
func (n NullProjectStatusStr) MarshalJSON() ([]byte, error) {
	const nullStr = "null"
	if n.Valid {
		return json.Marshal(n.ProjectStatus)
	}
	return []byte(nullStr), nil
}

// UnmarshalJSON correctly deserializes a NullProjectStatus from JSON.
func (n *NullProjectStatusStr) UnmarshalJSON(b []byte) error {
	n.Set = true
	var x interface{}
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	err = n.Scan(x)
	return err
}
