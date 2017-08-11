package generator

import (
	"strconv"
)

// Stringify returns a string that is all of the enum value names concatenated without a separator
func Stringify(e Enum) (ret string, err error) {
	for _, val := range e.Values {
		if val.Name != skipHolder {
			ret = ret + val.Name
		}
	}
	return
}

// Indexify returns a string that is all of the indexes for a string value lookup
func Indexify(e Enum) (ret string, err error) {
	ret = `[...]uint8{`
	index := 0
	for _, val := range e.Values {
		if val.Name != skipHolder {
			ret = ret + strconv.Itoa(index) + `,`
			index = index + len(val.Name)
		}
	}

	ret = ret + strconv.Itoa(index) + `}`
	return
}
