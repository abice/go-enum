package cso

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"

	"github.com/mkideal/pkg/encoding"
)

type Node interface {
	// embed encoding.Node
	encoding.Node
	// NumChild returns number of children
	NumChild() int
	// ByIndex gets ith child node,panic if i out of range [0,NumChild)
	ByIndex(i int) Node
	// Value returns value of node as an interface
	Value() interface{}

	// output writes Node to writer
	output(w io.Writer) error
}

// literalNode implements Node interface
type literalNode struct {
	encoding.LiteralNode
}

func newLiteralNode(pos scanner.Position, tok rune, value string) (*literalNode, error) {
	kind := encoding.InvalidNode
	switch tok {
	case scanner.Char:
		kind = encoding.CharNode
	case scanner.String:
		kind = encoding.StringNode
	case scanner.Float:
		kind = encoding.FloatNode
	case scanner.Int:
		kind = encoding.IntNode
	case scanner.Ident:
		kind = encoding.IdentNode
	default:
		return nil, fmt.Errorf("unexpected begin of json node %v at %v", value, pos)
	}
	n := &literalNode{LiteralNode: encoding.NewLiteralNode(pos, kind)}
	n.LiteralNode.Value = value
	return n, nil
}

func (n literalNode) NumChild() int      { return 0 }
func (n literalNode) ByIndex(i int) Node { panic("no child") }

func (n literalNode) Value() interface{} {
	switch n.Kind() {
	case encoding.CharNode:
		value, _, _, _ := strconv.UnquoteChar(n.LiteralNode.Value, '\'')
		return value
	case encoding.StringNode:
		value, _ := strconv.Unquote(n.LiteralNode.Value)
		return value
	case encoding.FloatNode:
		value, _ := strconv.ParseFloat(n.LiteralNode.Value, 64)
		return value
	case encoding.IntNode:
		value, _ := strconv.ParseInt(n.LiteralNode.Value, 0, 64)
		return value
	case encoding.IdentNode:
		return n.LiteralNode.Value
	default:
		return nil
	}
}

func (n *literalNode) output(w io.Writer) error {
	_, err := fmt.Fprint(w, n.LiteralNode.Value)
	return err
}

// listNode represents object or array
type listNode struct {
	encoding.Nodebase
	kind     encoding.NodeKind // ObjectNode or ArrayNode
	children []Node
}

func newListNode(pos scanner.Position, kind encoding.NodeKind) *listNode {
	return &listNode{
		Nodebase: encoding.NewNodebase(pos),
		kind:     kind,
	}
}

func (n *listNode) addChild(child Node) {
	n.children = append(n.children, child)
}

func (n listNode) Kind() encoding.NodeKind { return n.kind }
func (n listNode) NumChild() int           { return len(n.children) }
func (n listNode) ByIndex(i int) Node      { return n.children[i] }

func (n listNode) Value() interface{} {
	size := len(n.children)
	if size == 0 {
		return []interface{}{}
	}
	s := make([]interface{}, 0, size)
	for _, child := range n.children {
		s = append(s, child.Value())
	}
	return s
}

func (n *listNode) output(w io.Writer) error {
	openTok, closeTok := opLBrack, opRBrack
	if n.kind == encoding.ObjectNode {
		openTok, closeTok = opLBrace, opRBrace
	}
	if _, err := fmt.Fprint(w, string(openTok)); err != nil {
		return err
	}
	for i, child := range n.children {
		if i > 0 {
			if _, err := fmt.Fprint(w, ","); err != nil {
				return err
			}
		}
		if err := child.output(w); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w, string(closeTok))
	return err
}
