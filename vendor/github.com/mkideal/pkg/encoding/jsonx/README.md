# jsonx - json pkg supports comments,extraComma,unquotedKey

example 1: WithComment,WithExtraComma

```js
{
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
}
```

example 2: WithComment,WithExtraComma,WithUnquotedKey

```js
{
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
}
```
