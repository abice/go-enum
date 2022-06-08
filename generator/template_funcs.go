package generator

import (
	"fmt"
	"strconv"
	"strings"
)

// Stringify returns a string that is all of the enum value names concatenated without a separator
func Stringify(e Enum, forceLower bool) (ret string, err error) {
	for _, val := range e.Values {
		if val.Name != skipHolder {
			next := val.RawName
			if forceLower {
				next = strings.ToLower(next)
			}
			ret = ret + next
		}
	}
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
			ret = fmt.Sprintf("%s%s: %s[%d:%d],\n", ret, val.PrefixedName, strName, index, nextIndex)
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
			ret = fmt.Sprintf("%s%s[%d:%d]: %s,\n", ret, strName, index, nextIndex, val.PrefixedName)
			if lowercase {
				ret = fmt.Sprintf("%sstrings.ToLower(%s[%d:%d]): %s,\n", ret, strName, index, nextIndex, val.PrefixedName)
			}
			index = nextIndex
		}
	}
	ret = ret + `}`
	return
}

// Namify returns a slice that is all of the possible names for an enum in a slice
func Namify(e Enum) (ret string, err error) {
	strName := fmt.Sprintf(`_%sName`, e.Name)
	ret = "[]string{\n"
	index := 0
	for _, val := range e.Values {
		if val.Name != skipHolder {
			nextIndex := index + len(val.Name)
			ret = fmt.Sprintf("%s%s[%d:%d],\n", ret, strName, index, nextIndex)
			index = nextIndex
		}
	}
	ret = ret + "}"
	return
}

func Offset(index int, enumType string, val EnumValue) (strResult string) {
	if strings.HasPrefix(enumType, "u") {
		// Unsigned
		return strconv.FormatUint(val.Value.(uint64)-uint64(index), 10)
	} else {
		// Signed
		return strconv.FormatInt(val.Value.(int64)-int64(index), 10)
	}
}
