//go:generate go-enum -f=example.go --marshal --lower --flag

package example

// X is doc'ed
type X struct {
}

// Make x ENUM(Toyota,_,Chevy,_,Ford,_,Tesla,_,Hyundai,_,Nissan,_,Jaguar,_,Audi,_,BMW,_,Mercedes-Benz,_,Volkswagon)
type Make int32
