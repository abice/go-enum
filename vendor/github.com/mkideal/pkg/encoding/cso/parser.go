package cso

import (
	"fmt"
	"text/scanner"

	"github.com/mkideal/pkg/encoding"
)

const (
	opLBrace = rune('{')
	opRBrace = rune('}')
	opLBrack = rune('[')
	opRBrack = rune(']')
	opComma  = rune(',')
	opAdd    = rune('+')
	opSub    = rune('-')
)

// parser parses json
type parser struct {
	encoding.Parser
}

func (p *parser) init(s *scanner.Scanner) error {
	p.Init(s)
	return p.Next()
}

func (p *parser) expect(tok rune) error {
	if p.Tok == tok {
		return p.Next() // make progress
	}
	lit := "`" + p.Lit + "`"
	if p.Tok == scanner.EOF {
		lit = "EOF"
	}
	err := fmt.Errorf("expect `%s`, but got %s at %s", string(tok), lit, p.Pos)
	return err
}

func (p *parser) parseNode() (Node, error) {
	switch p.Tok {
	case opLBrace:
		return p.parseObjectNode()
	case opLBrack:
		return p.parseArrayNode()
	case opAdd:
		return p.parseSignNode(opAdd)
	case opSub:
		return p.parseSignNode(opSub)
	default:
		n, err := newLiteralNode(p.Pos, p.Tok, p.Lit)
		if err != nil {
			return nil, err
		}
		err = p.Next()
		return n, err
	}
}

func (p *parser) parseSignNode(pfxTok rune) (Node, error) {
	if err := p.Next(); err != nil {
		return nil, err
	}
	lit := "`" + p.Lit + "`"
	if p.Tok == scanner.EOF {
		lit = "EOF"
	}
	if p.Tok != scanner.Float && p.Tok != scanner.Int {
		return nil, fmt.Errorf("expect float or integer, but got %v at %v", lit, p.Pos)
	}
	node, err := newLiteralNode(p.Pos, p.Tok, p.Lit)
	if err != nil {
		return nil, err
	}
	err = p.Next()
	node.LiteralNode.Value = string(pfxTok) + node.LiteralNode.Value
	return node, err
}

func (p *parser) parseObjectNode() (Node, error) {
	return p.parseListNode(encoding.ObjectNode, opLBrace, opRBrace, true)
}

func (p *parser) parseArrayNode() (Node, error) {
	return p.parseListNode(encoding.ArrayNode, opLBrack, opRBrack, true)
}

func (p *parser) parseListNode(kind encoding.NodeKind, openTok, closeTok rune, bound bool) (Node, error) {
	pos := p.Pos
	if bound {
		if err := p.expect(openTok); err != nil {
			return nil, err
		}
	}
	list := newListNode(pos, kind)
	for p.Tok != scanner.EOF && p.Tok != closeTok {
		// add empty node
		// e.g. {,,,}
		for p.Tok == opComma {
			child, err := newLiteralNode(p.Pos, scanner.Ident, "")
			if err != nil {
				return nil, err
			}
			list.addChild(child)
			p.Next()
		}
		if p.Tok == closeTok {
			child, err := newLiteralNode(p.Pos, scanner.Ident, "")
			if err != nil {
				return nil, err
			}
			list.addChild(child)
			break
		}
		child, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if p.Tok != closeTok {
			if err := p.expect(opComma); err != nil {
				if bound || p.Tok != scanner.EOF {
					return nil, err
				}
			}
		}
		list.addChild(child)
	}
	if bound {
		if err := p.expect(closeTok); err != nil {
			return nil, err
		}
	}
	return list, nil
}
