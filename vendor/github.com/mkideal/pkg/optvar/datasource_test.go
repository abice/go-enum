package optvar_test

import (
	"testing"

	"github.com/mkideal/pkg/optvar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataSource(t *testing.T) {
	var (
		a int
		b bool
		c string
		d int
		e bool
		f string
	)
	tree := optvar.Values("",
		optvar.Int("a", &a, 2),
		optvar.Bool("b", &b, true),
		optvar.String("c", &c, "abc"),
		optvar.RequiredInt("d", &d),
		optvar.RequiredBool("e", &e),
		optvar.RequiredString("f", &f),
	)
	err := optvar.Map(map[string]string{
		"a": "1",
		"b": "false",
		"d": "4",
		"e": "true",
		"f": "world",
	}).Apply(tree)
	require.Nil(t, err)
	assert.Equal(t, 1, a)
	assert.Equal(t, false, b)
	assert.Equal(t, "abc", c)
	assert.Equal(t, 4, d)
	assert.Equal(t, true, e)
	assert.Equal(t, "world", f)
}
