package expr

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

// eval the expression
func eval(e *Expr, getter VarGetter, node ast.Expr) (Value, error) {
	switch n := node.(type) {
	case *ast.Ident:
		if getter == nil {
			return e.pool.onVarMissing(n.Name)
		}
		val, ok := getter.GetVar(n.Name)
		if !ok {
			return e.pool.onVarMissing(n.Name)
		}
		return val, nil

	case *ast.BasicLit:
		switch n.Kind {
		case token.INT:
			i, err := strconv.ParseInt(n.Value, 10, 64)
			if err != nil {
				return Zero(), err
			}
			return Int(i), nil
		case token.FLOAT:
			f, err := strconv.ParseFloat(n.Value, 64)
			if err != nil {
				return Zero(), err
			}
			return Float(f), nil
		case token.CHAR, token.STRING:
			s, err := strconv.Unquote(n.Value)
			if err != nil {
				return Zero(), err
			}
			return String(s), nil
		default:
			return Zero(), fmt.Errorf("unsupported token: %s(%v)", n.Value, n.Kind)
		}

	case *ast.ParenExpr:
		return eval(e, getter, n.X)

	case *ast.CallExpr:
		args := make([]Value, 0, len(n.Args))
		for _, arg := range n.Args {
			if val, err := eval(e, getter, arg); err != nil {
				return Zero(), err
			} else {
				args = append(args, val)
			}
		}
		if fnIdent, ok := n.Fun.(*ast.Ident); ok {
			if fn, ok := e.pool.fn(fnIdent.Name); !ok {
				return Zero(), fmt.Errorf("undefined function `%v`", fnIdent.Name)
			} else {
				return fn(args...)
			}
		}
		return Zero(), fmt.Errorf("unexpected func type: %T", n.Fun)

	case *ast.UnaryExpr:
		switch n.Op {
		case token.ADD:
			return eval(e, getter, n.X)
		case token.SUB:
			x, err := eval(e, getter, n.X)
			if err == nil {
				x, err = Zero().Sub(x)
			}
			return x, err
		case token.NOT:
			x, err := eval(e, getter, n.X)
			if err == nil {
				x = x.Not()
			}
			return x, err
		default:
			return Zero(), fmt.Errorf("unsupported unary op: %v", n.Op)
		}

	case *ast.BinaryExpr:
		x, err := eval(e, getter, n.X)
		if err != nil {
			return Zero(), err
		}
		y, err := eval(e, getter, n.Y)
		if err != nil {
			return Zero(), err
		}
		switch n.Op {
		case token.ADD:
			return x.Add(y)
		case token.SUB:
			return x.Sub(y)
		case token.MUL:
			return x.Mul(y)
		case token.QUO:
			return x.Quo(y)
		case token.REM:
			return x.Rem(y)
		case token.XOR:
			return x.Pow(y)
		case token.LAND:
			return x.And(y), nil
		case token.LOR:
			return x.Or(y), nil
		case token.EQL:
			return x.Eq(y)
		case token.NEQ:
			return x.Ne(y)
		case token.GTR:
			return x.Gt(y)
		case token.GEQ:
			return x.Ge(y)
		case token.LSS:
			return x.Lt(y)
		case token.LEQ:
			return x.Le(y)
		default:
			return Zero(), fmt.Errorf("unexpected binary operator: %v", n.Op)
		}

	default:
		return Zero(), fmt.Errorf("unexpected node type %T", n)
	}
}
