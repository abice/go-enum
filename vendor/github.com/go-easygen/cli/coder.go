package cli

type Decoder interface {
	Decode(s string) error
}

type SliceDecoder interface {
	Decoder
	DecodeSlice()
}

type Encoder interface {
	Encode() string
}

type CounterDecoder interface {
	Decoder
	IsCounter()
}

type Counter struct {
	value int
}

func (c Counter) Value() int { return c.value }

func (c *Counter) Decode(s string) error {
	c.value++
	return nil
}

func (c Counter) IsCounter() {}
