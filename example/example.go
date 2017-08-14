//go:generate go-enum -f=example.go --marshal --lower

package example

// X is doc'ed
type X struct {
}

// Animal x ENUM(
// Cat,
// Dog,
// Fish
// )
type Animal int32

// Model x ENUM(Toyota,_,Chevy,_,Ford,_,Tesla,_,Hyundai,_,Nissan,_,Jaguar,_,Audi,_,BMW,_,Mercedes,_,Volkswagon)
type Model int32
