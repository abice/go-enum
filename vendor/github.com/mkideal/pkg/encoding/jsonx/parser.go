package jsonx

import (
	"fmt"
	"text/scanner"

	"github.com/mkideal/pkg/encoding"
)

const (
	opLBrace = '{'
	opRBrace = '}'
	opLBrack = '['
	opRBrack = ']'
	opComma  = ','
	opColon  = ':'
	opAdd    = '+'
	opSub    = '-'
)

// parser parses json
type parser struct {
	encoding.Parser
	opt options
}

func (p *parser) init(s *scanner.Scanner, opt options) error {
	p.opt = opt
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
	node.value = string(pfxTok) + node.value
	return node, err
}

func (p *parser) parseKey() (key string, err error) {
	lit := "`" + p.Lit + "`"
	if p.Tok == scanner.EOF {
		lit = "EOF"
	}
	if p.opt.unquotedKey {
		if p.Tok != scanner.Ident {
			err = fmt.Errorf("expect a identifier or `}`, but got %s at %v", lit, p.Pos)
		}
	} else {
		if p.Tok != scanner.String {
			err = fmt.Errorf("expect a string or `}`, but got %s at %v", lit, p.Pos)
		}
	}
	if err == nil {
		key = p.Lit
		err = p.Next()
	}
	return
}

func (p *parser) parseObjectNode() (Node, error) {
	doc := p.LeadComment
	pos := p.Pos
	if err := p.expect(opLBrace); err != nil {
		return nil, err
	}
	obj := newObjectNode()
	obj.doc = doc
	obj.pos = pos
	for p.Tok != scanner.EOF && p.Tok != opRBrace {
		doc := p.LeadComment
		key, err := p.parseKey()
		if err != nil {
			return nil, err
		}
		if err := p.expect(opColon); err != nil {
			return nil, err
		}
		value, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		value.setDoc(doc)
		comment := p.LineComment
		if p.Tok != opRBrace {
			pos := p.Pos
			if err := p.expect(opComma); err != nil {
				return nil, err
			}
			comment = p.LineComment
			// extra comma not allowed at last node of object but found
			if !p.opt.extraComma && p.Tok == opRBrace {
				return nil, fmt.Errorf("extra comma found at %v", pos)
			}
		}
		value.setComment(comment)
		obj.addChild(key, value)
	}
	if err := p.expect(opRBrace); err != nil {
		return nil, err
	}
	obj.comment = p.LineComment
	return obj, nil
}

func (p *parser) parseArrayNode() (Node, error) {
	doc := p.LeadComment
	pos := p.Pos
	if err := p.expect(opLBrack); err != nil {
		return nil, err
	}
	arr := newArrayNode()
	arr.doc = doc
	arr.pos = pos
	for p.Tok != scanner.EOF && p.Tok != opRBrack {
		doc := p.LeadComment
		value, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		value.setDoc(doc)
		arr.addChild(value)
		if p.Tok != opRBrack {
			pos := p.Pos
			if err := p.expect(opComma); err != nil {
				return nil, err
			}
			// extra comma not allowed at last node of array but found
			if !p.opt.extraComma && p.Tok == opRBrack {
				return nil, fmt.Errorf("extra comma found at %v", pos)
			}
		}
	}
	if err := p.expect(opRBrack); err != nil {
		return nil, err
	}
	arr.comment = p.LineComment
	return arr, nil
}
