package example

//go:generate ../bin/go-enum --forcelower -b example

/*
ENUM(

	B3 = 03
	B4 = 04
	B5 = 5
	B6 = 0b110
	B7 = 0b111
	B8 = 0x08
	B9 = 0x09

)
*/
type DiffBase int
