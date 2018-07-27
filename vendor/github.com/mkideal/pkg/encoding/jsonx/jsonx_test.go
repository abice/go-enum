package jsonx

import (
	"os"
	"strings"
	"testing"
)

func ExampleDecode() {
	r := strings.NewReader(`{"a":1,"b":true,"c":[{"x":1.2},{"y":2.3}],"d":{},"e":-1,"f":+1}`)
	node, err := Read(r)
	if err != nil {
		return
	}
	Write(os.Stdout, node, WithIndent("  "))
	// Output:
	// {
	//   "a": 1,
	//   "b": true,
	//   "c": [
	//     {
	//       "x": 1.2
	//     },
	//     {
	//       "y": 2.3
	//     }
	//   ],
	//   "d": {},
	//   "e": -1,
	//   "f": +1
	// }
}

func ExampleDecodeExtraComma() {
	r := strings.NewReader(`{"a":1,"b":true,"c":[{"x":1.2},{"y":2.3},],"d":{},}`)
	node, err := Read(r, WithExtraComma())
	if err != nil {
		return
	}
	Write(os.Stdout, node, WithIndent("  "))
	Write(os.Stdout, node, WithIndent("  "), WithExtraComma())
	// Output:
	// {
	//   "a": 1,
	//   "b": true,
	//   "c": [
	//     {
	//       "x": 1.2
	//     },
	//     {
	//       "y": 2.3
	//     }
	//   ],
	//   "d": {}
	// }{
	//   "a": 1,
	//   "b": true,
	//   "c": [
	//     {
	//       "x": 1.2,
	//     },
	//     {
	//       "y": 2.3,
	//     },
	//   ],
	//   "d": {},
	// }
}

func ExampleDecodeWithUnquotedKey() {
	r := strings.NewReader(`{a:1,b:true,c:[{x:1.2},{y:2.3}],d:{}}`)
	node, err := Read(r, WithUnquotedKey())
	if err != nil {
		return
	}
	Write(os.Stdout, node, WithIndent("  "))
	Write(os.Stdout, node, WithIndent("  "), WithUnquotedKey())
	// Output:
	// {
	//   "a": 1,
	//   "b": true,
	//   "c": [
	//     {
	//       "x": 1.2
	//     },
	//     {
	//       "y": 2.3
	//     }
	//   ],
	//   "d": {}
	// }{
	//   a: 1,
	//   b: true,
	//   c: [
	//     {
	//       x: 1.2
	//     },
	//     {
	//       y: 2.3
	//     }
	//   ],
	//   d: {}
	// }
}

func ExampleDecodeWithComment() {
	r := strings.NewReader(`{
	// doc a
	"a":1, // line a
	// doc b
	"b":true, // line b

	// doc c1
	// doc c2
	"c":[
		{"x":1.2},
		// doc y
		{"y":2.3}
	], // line c
	// doc d
	"d":{
		// doc e
		"e": {
			// doc f
			"f": 1
		}
	}
}`)
	node, err := Read(r, WithComment())
	if err != nil {
		return
	}
	Write(os.Stdout, node, WithIndent("  "), WithComment())
	// Output:
	// {
	//   // doc a
	//   "a": 1,// line a
	//   // doc b
	//   "b": true,// line b
	//   // doc c1
	//   // doc c2
	//   "c": [
	//     {
	//       "x": 1.2
	//     },
	//     // doc y
	//     {
	//       "y": 2.3
	//     }
	//   ],// line c
	//   // doc d
	//   "d": {
	//     // doc e
	//     "e": {
	//       // doc f
	//       "f": 1
	//     }
	//   }
	// }
}

func TestParser(t *testing.T) {
	type argt struct {
		src  string
		err  string
		kind NodeKind
		opt  options
	}
	for i, ts := range []argt{
		{``, "unexpected begin of json node  at <input>:1:1", InvalidNode, options{}},
		{`%`, "unexpected begin of json node % at <input>:1:2", InvalidNode, options{}},
		{`(`, "unexpected begin of json node ( at <input>:1:2", InvalidNode, options{}},
		{`{]`, "expect a string or `}`, but got `]` at <input>:1:3", InvalidNode, options{}},
		{`//comment`, "unexpected begin of json node / at <input>:1:2", InvalidNode, options{}},
		{`/*comment*/`, "unexpected begin of json node / at <input>:1:2", InvalidNode, options{}},
		{`1`, "", IntNode, options{}},
		{`1.2`, "", FloatNode, options{}},
		{`/*comment*/1.2`, "", FloatNode, options{supportComment: true}},
		{`abc`, "", IdentNode, options{}},
		{`abc//comment`, "", IdentNode, options{supportComment: true}},
		{`'a'`, "", CharNode, options{}},
		{`''`, "illegal char literal at <input>:1:2", InvalidNode, options{}},
		{`'xxx'`, "illegal char literal at <input>:1:5", InvalidNode, options{}},
		{`""`, "", StringNode, options{}},
		{`"abcd"`, "", StringNode, options{}},
		{`'abcd"`, "illegal char literal at <input>:1:7", InvalidNode, options{}},
		{`// doc
		"abcd"`, "", StringNode, options{supportComment: true}},
		{`{"x":1}`, "", ObjectNode, options{}},
		{`{"x":1,}`, "extra comma found at <input>:1:8", InvalidNode, options{}},
		{`{"x":1,}`, "", ObjectNode, options{extraComma: true}},
		{`{"x":1,"y":{}}`, "", ObjectNode, options{}},
		{`{"x":1,"y":{]}`, "expect a string or `}`, but got `]` at <input>:1:14", InvalidNode, options{}},
		{`{"x":1,"y":{/**/}}`, "", ObjectNode, options{supportComment: true}},
		{`{"x":1,"y":{//}}`, "expect `}`, but got EOF at <input>:1:17", InvalidNode, options{supportComment: true}},
		{`[]`, "", ArrayNode, options{}},
		{`[x]`, "", ArrayNode, options{}},
		{`[x, y, z]`, "", ArrayNode, options{}},
		{`["x", "y", z]`, "", ArrayNode, options{}},
		{`[{}]`, "", ArrayNode, options{}},
		{`[1,{}]`, "", ArrayNode, options{}},
		{`[-1,{}]`, "", ArrayNode, options{}},
		{`{x:1}`, "expect a string or `}`, but got `x` at <input>:1:3", InvalidNode, options{}},
		{`{x:1}`, "", ObjectNode, options{unquotedKey: true}},
	} {
		r := strings.NewReader(ts.src)
		node, err := Read(r, ts.opt.clone)
		if err != nil && ts.err == "" {
			t.Errorf("%dth: want nil, got error %v", i, err)
			continue
		}
		if err == nil && ts.err != "" {
			t.Errorf("%dth: want error, but got nil", i)
			continue
		}
		if err != nil && ts.err != err.Error() {
			t.Errorf("%dth: want error %v, but got %v", i, ts.err, err)
			continue
		}
		if err == nil {
			if node.Kind() != ts.kind {
				t.Errorf("%dth: want node kind %v, but got %v", i, ts.kind, node.Kind())
				continue
			}
		}
	}
}

func TestUnmarshal(t *testing.T) {
	type A struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	type B struct {
		Key   string `json:"key"`
		Value int    `json:"value"`
	}
	type C struct {
		A A       `json:"a"`
		B B       `json:"b"`
		X int     `json:"x"`
		Y float32 `json:"y"`
		Z bool    `json:"z"`
		S []B     `json:"s"`
	}
	value := `{"a":{"id":1,"name":"xxx"},"b":{"key":"key","value":222},"x":2,"y":1.5,"z":true,"s":[{"key":"k1","value":1},{"key":"k2","value":2},{"key":"","value":0}]}`
	data1 := []byte(`{
		/* c-style comment supported if contains option WithComment
		 * extra comma allowed if contains option WithExtraComma
		 */
		// a is an object
		"a": {
			// id is an integer
			"id": 1,
			// name is a string
			"name": "xxx", // name of A
		},
		"b": {
			"key": "key",
			"value": 222, // NOTE: here is an extra comma, but no problem if contains option WithExtraComma
		},
		"x": 2,
		"y": 1.5,
		"z": true,
		"s": [
			// array element 0
			{"key":"k1", "value": 1},
			// array element 1
			{
				"key":"k2",
				"value":2
			},{} // array element 2
		]
	}`)
	v1 := new(C)
	err := Unmarshal(data1, v1, WithExtraComma(), WithComment())
	if err != nil {
		t.Errorf("Unmarshal data1 error: %v", err)
		return
	}
	marshaled, err := Marshal(v1)
	if err != nil {
		t.Errorf("Unmarshal v1 error: %v", err)
		return
	}
	if string(marshaled) != value {
		t.Errorf("marshaled v1 not equal to value: `%s` vs `%s`", string(marshaled), value)
		return
	}

	data2 := []byte(`{
		/* c-style comment supported if contains option WithComment
		 * extra comma allowed if contains option WithExtraComma
		 * key not quoted with " allowed if contains option WithUnquotedKey
		 */
		// a is an object
		a: {
			// id is an integer
			id: 1,
			// name is a string
			name: "xxx", // name of A
		},
		b: {
			key: "key",
			value: 222, // NOTE: here is an extra comma, but no problem if contains option WithExtraComma
		},
		x: 2,
		y: 1.5,
		z: true,
		s: [
			// array element 0
			{key:"k1", value: 1},
			// array element 1
			{
				key:"k2",
				value:2
			},{} // array element 2
		]
	}`)
	v2 := new(C)
	err = Unmarshal(data2, v2, WithExtraComma(), WithComment(), WithUnquotedKey())
	if err != nil {
		t.Errorf("Unmarshal data2 error: %v", err)
		return
	}
	marshaled, err = Marshal(v2)
	if err != nil {
		t.Errorf("Unmarshal v2 error: %v", err)
		return
	}
	if string(marshaled) != value {
		t.Errorf("marshaled v2 not equal to value: `%s` vs `%s`", string(marshaled), value)
		return
	}
}
