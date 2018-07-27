package cli

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag(t *testing.T) {
	type argT struct {
		Cli    string `cli:"cli"`
		Pw     string `pw:"pw"`
		Edit   string `edit:"edit"`
		Usage  string `cli:"usage" usage:"hello,usage"`
		Dft    string `cli:"dft" dft:"hello,dft"`
		Name   string `cli:"name" name:"hello-name"`
		Prompt string `cli:"prompt" prompt:"hello,prompt"`
		Parser string `cli:"parser" parser:"json"`
		Sep    string `cli:"sep" sep:":"`

		Required     string `cli:"*r"`
		Force        string `cli:"!f"`
		EditFile     string `edit:"Filename:file"`
		ShortAndLong string `cli:"x,y,z,xy,yz,xyz"`

		Empty string `cli:"-"`

		// errors
		Err string `cli:"e" pw:"e" edit:"e" usage:"cli,pw,edit both are cli-like flag"`
	}
	argv := new(argT)
	typ := reflect.TypeOf(argv).Elem()
	val := reflect.ValueOf(argv).Elem()
	for i, size := 0, val.NumField(); i < size; i++ {
		typField, _ := typ.Field(i), val.Field(i)
		tag, isEmpty, err := parseTag(typField.Name, typField.Tag)
		if err != nil {
			assert.Equal(t, typField.Name, "Err")
			assert.Equal(t, err, errCliTagTooMany)
			continue
		}
		if isEmpty {
			assert.Equal(t, typField.Name, "Empty")
			continue
		}

		switch typField.Name {
		case "Cli":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
			assert.Equal(t, tag.longNames, []string{"--cli"})
		case "Pw":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.True(t, tag.isPassword)
			assert.False(t, tag.isEdit)
			assert.Equal(t, tag.longNames, []string{"--pw"})
		case "Edit":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.True(t, tag.isEdit)
			assert.Equal(t, tag.longNames, []string{"--edit"})
		case "Usage":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
			assert.Equal(t, tag.longNames, []string{"--usage"})
			assert.Equal(t, tag.usage, "hello,usage")
		case "Dft":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
			assert.Equal(t, tag.longNames, []string{"--dft"})
			assert.Equal(t, tag.dft, "hello,dft")
		case "Name":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
			assert.Equal(t, tag.longNames, []string{"--name"})
			assert.Equal(t, tag.name, "hello-name")
		case "Prompt":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
			assert.Equal(t, tag.longNames, []string{"--prompt"})
			assert.Equal(t, tag.prompt, "hello,prompt")
		case "Parser":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
			assert.Equal(t, tag.longNames, []string{"--parser"})
			assert.Equal(t, typField.Tag.Get(tagParser), "json")
		case "Sep":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
			assert.Equal(t, tag.longNames, []string{"--sep"})
			assert.Equal(t, tag.sep, ":")
		case "Required":
			assert.True(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
		case "Force":
			assert.False(t, tag.isRequired)
			assert.True(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
		case "EditFile":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.True(t, tag.isEdit)
			assert.Equal(t, tag.editFile, "Filename")
		case "ShortAndLong":
			assert.False(t, tag.isRequired)
			assert.False(t, tag.isForce)
			assert.False(t, tag.isPassword)
			assert.False(t, tag.isEdit)
			assert.Equal(t, tag.shortNames, []string{"-x", "-y", "-z"})
			assert.Equal(t, tag.longNames, []string{"--xy", "--yz", "--xyz"})
		}
	}
}
