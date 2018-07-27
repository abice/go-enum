package namemapper

import (
	"bytes"
	"strings"
)

type NameStyle int

type Mapper func(string) string

const (
	UNKNOWN     NameStyle = -1
	NATIVE                = 0 // native
	UPPER                 = 1 // uppercase
	LOWER                 = 2 // lowercase
	UPPER_CAMEL           = 3 // upper camel
	LOWER_CAMEL           = 4 // lower camel
	UNDER_SCORE           = 5 // under score
	TITLE                 = 6 // upper first charactor
)

func GetStyleByName(stylename string) NameStyle {
	switch stylename {
	case "title":
		return TITLE
	case "upper":
		return UPPER
	case "lower":
		return LOWER
	case "upper-camel":
		return UPPER_CAMEL
	case "lower-camel":
		return LOWER_CAMEL
	case "under-score":
		return UNDER_SCORE
	case "native":
		return NATIVE
	default:
		return UNKNOWN
	}
}

// Style detects name style
// return NATIVE if style is unknown
// return UNDER_SCORE if s contain _
// return LOWER_CAMEL if first charactor is lowercase
// return UPPER_CAMEL // retst charactor is uppercase
// otherwise return NATIVE
func Style(s string) NameStyle {
	if s == "" {
		return NATIVE
	}
	if strings.Contains(s, "_") {
		return UNDER_SCORE
	}
	fields, ok := split(s)
	if !ok {
		return NATIVE
	}
	if len(fields) == 0 {
		return NATIVE
	}
	firstChar := fields[0][0]
	if firstChar >= 'a' && firstChar <= 'z' {
		return LOWER_CAMEL
	}
	if firstChar >= 'A' && firstChar <= 'Z' {
		return UPPER_CAMEL
	}
	return NATIVE
}

func Convert(s string, style NameStyle) string {
	switch style {
	case UPPER:
		return Upper(s)
	case LOWER:
		return Lower(s)
	case UPPER_CAMEL:
		return UpperCamel(s)
	case LOWER_CAMEL:
		return LowerCamel(s)
	case UNDER_SCORE:
		return UnderScore(s)
	case TITLE:
		return Title(s)
	default:
		return Native(s)
	}
}

func Native(s string) string {
	return s
}

func Upper(s string) string {
	return strings.ToUpper(s)
}

func Lower(s string) string {
	return strings.ToLower(s)
}

func UpperCamel(s string) string {
	fields, ok := split(s)
	if !ok || len(fields) == 0 {
		return s
	}
	for i := range fields {
		fields[i] = strings.Title(fields[i])
	}
	return strings.Join(fields, "")
}

func LowerCamel(s string) string {
	fields, ok := split(s)
	if !ok || len(fields) == 0 {
		return s
	}
	for i := range fields {
		if i == 0 {
			fields[0] = strings.ToLower(fields[0])
		} else {
			fields[i] = strings.Title(fields[i])
		}
	}
	return strings.Join(fields, "")
}

func UnderScore(s string) string {
	fields, ok := split(s)
	if !ok || len(fields) == 0 {
		return s
	}
	for i := range fields {
		fields[i] = strings.ToLower(fields[i])
	}
	return strings.Join(fields, "_")
}

func Title(s string) string {
	return strings.Title(s)
}

// split into slice
func split(s string) ([]string, bool) {
	if s == "" {
		return nil, false
	}
	ret := make([]string, 0)
	ccase := 0
	const UPPER_CASE = 1
	const LOWER_CASE = 2
	word := bytes.NewBufferString("")
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b >= 'A' && b <= 'Z' {
			if ccase != UPPER_CASE {
				if word.Len() > 0 {
					ret = append(ret, word.String())
					word.Reset()
				}
			}
			word.WriteByte(b)
			ccase = UPPER_CASE
		} else if b >= 'a' && b <= 'z' {
			word.WriteByte(b)
			ccase = LOWER_CASE
		} else if b == '_' {
			if word.Len() > 0 {
				ret = append(ret, word.String())
			}
			ccase = 0
			word.Reset()
		} else if b >= '0' && b <= '9' {
			word.WriteByte(b)
		} else {
			return nil, false
		}
		if i+1 == len(s) && word.Len() > 0 {
			ret = append(ret, word.String())
		}
	}
	return ret, true
}
