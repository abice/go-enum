package jsonx

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"text/scanner"
)

// Option represents a function for setting options
type Option func(*options)

type options struct {
	// prefix, indent used to pretty output json
	prefix, indent string
	// c-style comment supported if supportComment is true
	supportComment bool
	// key should be unquoted if unquotedKey is true
	unquotedKey bool
	// extra comma could be insert to end of last node of object or array if extraComma is true
	extraComma bool
}

func (opt options) clone(dst *options) {
	dst.prefix = opt.prefix
	dst.indent = opt.indent
	dst.supportComment = opt.supportComment
	dst.unquotedKey = opt.unquotedKey
	dst.extraComma = opt.extraComma
}

// WithComment returns an option which sets supportComment true
func WithComment() Option {
	return func(opt *options) {
		opt.supportComment = true
	}
}

// WithPrefix returns an option which with prefix while outputing
func WithPrefix(prefix string) Option {
	return func(opt *options) {
		opt.prefix = prefix
	}
}

// WithIndent returns an option which with indent while outputing
func WithIndent(indent string) Option {
	return func(opt *options) {
		opt.indent = indent
	}
}

// WithUnquotedKey returns an option which sets unquotedKey true
func WithUnquotedKey() Option {
	return func(opt *options) {
		opt.unquotedKey = true
	}
}

// WithExtraComma returns an option which sets extraComma true
func WithExtraComma() Option {
	return func(opt *options) {
		opt.extraComma = true
	}
}

func applyOptions(opts []Option) options {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Read reads a json node from reader r
func Read(r io.Reader, opts ...Option) (Node, error) {
	opt := applyOptions(opts)
	s := new(scanner.Scanner)
	s = s.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings
	if opt.supportComment {
		s.Mode |= scanner.ScanComments
	}
	p := new(parser)
	if err := p.init(s, opt); err != nil {
		return nil, err
	}
	return p.parseNode()
}

// ReadBytes reads a json node from bytes
func ReadBytes(data []byte, opts ...Option) (Node, error) {
	return Read(bytes.NewBuffer(data), opts...)
}

// ReadFile reads a json node from file
func ReadFile(filename string, opts ...Option) (Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Read(file, opts...)
}

// Write writes a json node to writer w
func Write(w io.Writer, node Node, opts ...Option) error {
	return node.output("", w, applyOptions(opts), true, true)
}

// WriteFile writer a json node to file
func WriteFile(filename string, node Node, perm os.FileMode, opts ...Option) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	return Write(file, node, opts...)
}

// Unmarshal wraps json.Unmarshal with options
func Unmarshal(data []byte, v interface{}, opts ...Option) error {
	if len(opts) == 0 {
		return json.Unmarshal(data, v)
	}
	buf := bytes.NewBuffer(data)
	node, err := Read(buf, opts...)
	if err != nil {
		return err
	}
	buf.Reset()
	if err := Write(buf, node); err != nil {
		return err
	}
	return json.Unmarshal(buf.Bytes(), v)
}

// Marshal marshal value v to json with options
func Marshal(v interface{}, opts ...Option) ([]byte, error) {
	return json.Marshal(v)
	/*opt := applyOptions(opts)
	node, err := reflectValue(v, opt)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = Write(&buf, node, opt.clone)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil*/
}

// Decoder wraps json.Decoder with options
type Decoder struct {
	r   io.Reader
	opt options
}

// NewDecoder creates a decoder with reader and options
func NewDecoder(r io.Reader, opts ...Option) *Decoder {
	return &Decoder{
		r:   r,
		opt: applyOptions(opts),
	}
}

// Decode decodes data and stores it in the valyala pointed to by v
func (decoder *Decoder) Decode(v interface{}) error {
	node, err := Read(decoder.r, decoder.opt.clone)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := Write(&buf, node); err != nil {
		return err
	}
	return json.Unmarshal(buf.Bytes(), v)
}

// NewEncoder wraps json.NewEncoder
func NewEncoder(w io.Writer) *json.Encoder {
	return json.NewEncoder(w)
}
