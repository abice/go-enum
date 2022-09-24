// Copyright (c) 2017 Alex Bice

// go-enum is a utility for generating a more functional version of
// enumerations in go.
//
// The generator looks for the exact string `ENUM(` and will continue
// to store comma separated values until it finds a `)`.  Those values
// can be on one line, or separate lines; they can include `_` in order
// to skip a value in the enum that won't be an allowed value.
//
// Installation
//
//	go get github.com/abice/go-enum
//
// Usage:
// Sample File
//
//	//go:generate go-enum -f=myenum.go --marshal --
//
//	package mypackage
//
//	// MyEnum docs here
//	// ENUM(Value1, Value2
//	// Value3,_,
//	// Value4)
//	type MyEnum int
//
// Command to generate your enum
//
//	go generate ./
package main
