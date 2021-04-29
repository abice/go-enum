//go:generate ../bin/go-enum -f=$GOFILE --template user_template.tmpl

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