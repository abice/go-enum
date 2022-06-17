//go:generate ../bin/go-enum -f=$GOFILE --names

package example

/* ENUM(

Unkno 				= 0
E2P15					= 32768
E2P16					= 65536
E2P17					= 131072
E2P18					= 262144
E2P19					= 524288
E2P20					= 1048576
E2P21					= 2097152
E2P22					= 33554432
E2P23					= 67108864
E2P28					= 536870912
E2P30					= 1073741824
E2P31					= 2147483648
E2P32					= 4294967296
E2P33					= 8454967296
E2P63					= 18446744073709551615
)

*/
type Enum64bit uint64
