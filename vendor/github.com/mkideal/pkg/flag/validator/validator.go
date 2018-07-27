package validator

import (
	"errors"
)

type Validator interface {
	Validate() error
}

// ValidatorFunc represents a function which implements Validator interface
type ValidatorFunc func() error

func (v ValidatorFunc) Validate() error { return v() }

// ValidatorList holds n Validators, and implements Validator interface too
type ValidatorList []Validator

func NewSet() *ValidatorList {
	set := ValidatorList(make([]Validator, 0))
	return &set
}

func (v ValidatorList) Validate() error {
	for _, x := range v {
		if err := x.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// AndRequire appends a boolValidator
func (v *ValidatorList) AndRequire(val bool, msg string) *ValidatorList {
	(*v) = append((*v), BoolValidator{val: val, msg: msg})
	return v
}

// And appends a Validator
func (v *ValidatorList) And(validator Validator) *ValidatorList {
	(*v) = append((*v), validator)
	return v
}

// BoolValidator implements Validator which used to valite a boolean value should be true
type BoolValidator struct {
	val bool
	msg string
}

func NewBoolValidator(val bool, msg string) Validator {
	return BoolValidator{val: val, msg: msg}
}

func (v BoolValidator) Validate() error {
	if !v.val {
		return errors.New(v.msg)
	}
	return nil
}

// Require returns a ValidatorList which holds a boolValidator created by (val,msg)
func Require(val bool, msg string) *ValidatorList {
	v := new(ValidatorList)
	return v.AndRequire(val, msg)
}
