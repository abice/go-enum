//go:generate ../bin/go-enum -f=$GOFILE -t user_template.tmpl -t *user_glob*.tmpl

package example

// OceanColor is an enumeration of ocean colors that are allowed.
/*
ENUM(
Cerulean
Blue
Green
)
*/
type OceanColor int
