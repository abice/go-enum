//go:build tools
// +build tools

package main

import (
	_ "github.com/mattn/goveralls"
	_ "go.uber.org/mock/mockgen"
	_ "go.uber.org/mock/mockgen/model"
	_ "golang.org/x/tools/cmd/cover"
	_ "golang.org/x/tools/cmd/goimports"
)
