package cli

import (
	"encoding/json"
)

// FlagParser represents a parser for parsing flag
type FlagParser interface {
	Parse(s string) error
}

// FlagParserCreator represents factory function of FlagParser
type FlagParserCreator func(ptr interface{}) FlagParser

var parserCreators = map[string]FlagParserCreator{}

// RegisterFlagParser registers FlagParserCreator by name
func RegisterFlagParser(name string, creator FlagParserCreator) {
	if _, ok := parserCreators[name]; ok {
		panic("RegisterFlagParser has registered: " + name)
	}
	parserCreators[name] = creator
}

func init() {
	RegisterFlagParser("json", newJSONParser)
	RegisterFlagParser("jsonfile", newJSONFileParser)
}

// JSON parser
type JSONParser struct {
	ptr interface{}
}

func newJSONParser(ptr interface{}) FlagParser {
	return &JSONParser{ptr}
}

func (p JSONParser) Parse(s string) error {
	return json.Unmarshal([]byte(s), p.ptr)
}

// JSON file parser
type JSONFileParser struct {
	ptr interface{}
}

func newJSONFileParser(ptr interface{}) FlagParser {
	return &JSONFileParser{ptr}
}

func (p JSONFileParser) Parse(s string) error {
	return ReadJSONFromFile(s, p.ptr)
}
