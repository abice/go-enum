package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"strings"
)

type (
	Func func(...Value) (Value, error)

	// VarGetter defines interface for getting value of variable
	VarGetter interface {
		GetVar(string) (Value, bool)
	}

	// Expr is top-level object of expr package
	Expr struct {
		root ast.Expr
		pool *Pool
	}
)

// default VarGetter implementation
type Getter map[string]Value

// GetVar gets the value of variable
func (getter Getter) GetVar(name string) (Value, bool) {
	if getter == nil {
		return nilValue, false
	}
	v, ok := getter[name]
	return v, ok
}

// New creates an Expr and parses string s, pool can be nil
func New(s string, pool *Pool) (*Expr, error) {
	s = strings.TrimSpace(s)
	if pool == nil {
		pool = defaultPool
	}
	if e, ok := pool.get(s); ok {
		return e, nil
	}
	e := new(Expr)
	e.pool = pool
	if err := e.parse(s); err != nil {
		return nil, err
	}
	pool.set(s, e)
	return e, nil
}

// parse parses string s
func (e *Expr) parse(s string) error {
	if s == "" {
		return nil
	}
	node, err := parser.ParseExpr(s)
	if err != nil {
		return err
	}
	e.root = node

	v := &visitor{pool: e.pool}
	ast.Walk(v, e.root)
	return v.err
}

type visitor struct {
	pool *Pool
	err  error
}

// Visit implements ast.Visitor Visit method
func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if n, ok := node.(*ast.CallExpr); ok {
		if fnIdent, ok := n.Fun.(*ast.Ident); ok {
			if _, ok := v.pool.fn(fnIdent.Name); ok {
				return v
			} else {
				v.err = fmt.Errorf("undefined function `%v`", fnIdent.Name)
			}
		} else {
			v.err = fmt.Errorf("unsupported call expr")
		}
		return nil
	}
	return v
}

// Eval calculate the expression
// getter maybe nil
func (e *Expr) Eval(getter VarGetter) (Value, error) {
	if e.root == nil {
		return Zero(), nil
	}
	v, err := eval(e, getter, e.root)
	if err != nil {
		return Zero(), err
	}
	return v, nil
}

// Eval calculate expression with pool(maybe nil)
func Eval(s string, getter map[string]Value, pool *Pool) (Value, error) {
	e, err := New(s, pool)
	if err != nil {
		return Zero(), err
	}
	return e.Eval(Getter(getter))
}
