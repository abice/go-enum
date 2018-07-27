package textutil

import (
	"bytes"
	"text/template"
)

func Tpl(text string, data map[string]string) string {
	t := template.New(text)
	t, err := t.Parse(text)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBufferString("")
	if err := t.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}
