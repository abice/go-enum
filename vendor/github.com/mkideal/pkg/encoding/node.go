package encoding

import (
	"strconv"
	"text/scanner"
)

// NodeKind represents kind of json
type NodeKind int

const (
	InvalidNode NodeKind = iota
	IdentNode            // abc,true,false
	IntNode              // 1
	FloatNode            // 1.2
	CharNode             // 'c'
	StringNode           // "xyz"
	ObjectNode           // {}
	ArrayNode            // []
)

func (kind NodeKind) String() string {
	if kind >= 0 && kind < NodeKind(len(nodeKinds)) {
		return nodeKinds[kind]
	}
	return "Unknown kind(" + strconv.Itoa(int(kind)) + ")"
}

var nodeKinds = [...]string{
	InvalidNode: "InvalidNode",
	IdentNode:   "IdentNode",
	IntNode:     "IntNode",
	FloatNode:   "FloatNode",
	CharNode:    "CharNode",
	StringNode:  "StringNode",
	ObjectNode:  "ObjectNode",
	ArrayNode:   "ArrayNode",
}

type Node interface {
	// Pos returns position of node
	Pos() scanner.Position
	// Kind returns kind of node
	Kind() NodeKind
}

type Nodebase struct {
	pos scanner.Position
}

func NewNodebase(pos scanner.Position) Nodebase { return Nodebase{pos: pos} }

func (n Nodebase) Pos() scanner.Position { return n.pos }

type LiteralNode struct {
	Nodebase
	kind  NodeKind
	Value string
}

func NewLiteralNode(pos scanner.Position, kind NodeKind) LiteralNode {
	return LiteralNode{
		Nodebase: NewNodebase(pos),
		kind:     kind,
	}
}

func (n LiteralNode) Kind() NodeKind { return n.kind }
