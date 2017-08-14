package generator

import (
	"fmt"
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

// Mapify returns a map that is all of the indexes for a string value lookup
func Mapify(e Enum) (ret string, err error) {
	strName := fmt.Sprintf(`_%sName`, e.Name)
	ret = fmt.Sprintf("map[%s]string{\n", e.Name)
	index := 0
	for _, val := range e.Values {
		if val.Name != skipHolder {
			nextIndex := index + len(val.Name)
			ret = fmt.Sprintf("%s%d: %s[%d:%d],\n", ret, val.Value, strName, index, nextIndex)
			index = nextIndex
		}
	}
	ret = ret + `}`
	return
}

// Unmapify returns a map that is all of the indexes for a string value lookup
func Unmapify(e Enum, lowercase bool) (ret string, err error) {
	strName := fmt.Sprintf(`_%sName`, e.Name)
	ret = fmt.Sprintf("map[string]%s{\n", e.Name)
	index := 0
	for _, val := range e.Values {
		if val.Name != skipHolder {
			nextIndex := index + len(val.Name)
			ret = fmt.Sprintf("%s%s[%d:%d]: %d,\n", ret, strName, index, nextIndex, val.Value)
			if lowercase {
				ret = fmt.Sprintf("%sstrings.ToLower(%s[%d:%d]): %d,\n", ret, strName, index, nextIndex, val.Value)
			}
			index = nextIndex
		}
	}
	ret = ret + `}`
	return
}
