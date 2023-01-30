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
	// AnimalCat is a Animal of type Cat.
	AnimalCat Animal = iota
	// AnimalDog is a Animal of type Dog.
	AnimalDog
	// AnimalFish is a Animal of type Fish.
	AnimalFish
	// AnimalFishPlusPlus is a Animal of type Fish++.
	AnimalFishPlusPlus
	// AnimalFishSharp is a Animal of type Fish#.
	AnimalFishSharp
)

var ErrInvalidAnimal = errors.New("not a valid Animal")

const _animalName = "CatDogFishFish++Fish#"

var _animalMap = map[Animal]string{
	AnimalCat:          _animalName[0:3],
	AnimalDog:          _animalName[3:6],
	AnimalFish:         _animalName[6:10],
	AnimalFishPlusPlus: _animalName[10:16],
	AnimalFishSharp:    _animalName[16:21],
}

// String implements the Stringer interface.
func (x Animal) String() string {
	if str, ok := _animalMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Animal(%d)", x)
}

var _animalValue = map[string]Animal{
	_animalName[0:3]:   AnimalCat,
	_animalName[3:6]:   AnimalDog,
	_animalName[6:10]:  AnimalFish,
	_animalName[10:16]: AnimalFishPlusPlus,
	_animalName[16:21]: AnimalFishSharp,
}

// ParseAnimal attempts to convert a string to a Animal.
func ParseAnimal(name string) (Animal, error) {
	if x, ok := _animalValue[name]; ok {
		return x, nil
	}
	return Animal(0), fmt.Errorf("%s is %w", name, ErrInvalidAnimal)
}
