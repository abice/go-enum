// +build tools

package main

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/kevinburke/go-bindata/go-bindata"
	_ "github.com/mattn/goveralls"
	_ "golang.org/x/tools/cmd/cover"
	_ "golang.org/x/tools/cmd/goimports"
)
