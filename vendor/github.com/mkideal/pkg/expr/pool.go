package expr

import (
	"fmt"
	"math/rand"
	"regexp"
	"sync"
)

type VarMissingFunc func(string) (Value, error)

func DefaultOnVarMissing(varName string) (Value, error) {
	return Zero(), fmt.Errorf("var `%s' missing", varName)
}

type Pool struct {
	locker sync.RWMutex
	pool   map[string]*Expr

	factory      map[string]Func
	onVarMissing VarMissingFunc
}

func MustNewPool(factories ...map[string]Func) *Pool {
	pool, err := NewPool(factories...)
	if err != nil {
		panic(err)
	}
	return pool
}

func NewPool(factories ...map[string]Func) (*Pool, error) {
	p := &Pool{
		pool:         make(map[string]*Expr),
		factory:      newDefaultFactory(),
		onVarMissing: DefaultOnVarMissing,
	}
	for _, factory := range factories {
		if factory == nil {
			continue
		}
		for name, fn := range factory {
			if !validateFuncName(name) {
				return nil, fmt.Errorf("illegal function name `%s`", name)
			}
			p.factory[name] = fn
		}
	}
	return p, nil
}

func (p *Pool) SetOnVarMissing(fn VarMissingFunc) {
	p.onVarMissing = fn
}

func (p *Pool) get(s string) (*Expr, bool) {
	p.locker.RLock()
	defer p.locker.RUnlock()
	e, ok := p.pool[s]
	return e, ok && e != nil
}

func (p *Pool) set(s string, e *Expr) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.pool[s] = e
}

func (p *Pool) fn(name string) (Func, bool) {
	fn, ok := p.factory[name]
	return fn, ok
}

// validate function name
var funcNameRegexp = regexp.MustCompile("[a-zA-Z_][a-z-A-Z_0-9]{0,254}")

func validateFuncName(name string) bool {
	return funcNameRegexp.MatchString(name)
}

// default Pool
var defaultPool = func() *Pool {
	p, err := NewPool()
	if err != nil {
		panic(err)
	}
	return p
}()

// default factory
var newDefaultFactory = func() map[string]Func {
	return map[string]Func{
		"min":  builtin_min,
		"max":  builtin_max,
		"rand": builtin_rand,
	}
}

//------------------
// builtin function
//------------------

func builtin_min(args ...Value) (Value, error) {
	if len(args) == 0 {
		return nilValue, fmt.Errorf("missing arguments for function `min`")
	}
	x := args[0]
	for i, size := 1, len(args); i < size; i++ {
		lt, err := args[i].Lt(x)
		if err != nil {
			return Zero(), err
		}
		if lt.Bool() {
			x = args[i]
		}
	}
	return x, nil
}

func builtin_max(args ...Value) (Value, error) {
	if len(args) == 0 {
		return nilValue, fmt.Errorf("missing arguments for function `max`")
	}
	x := args[0]
	for i, size := 1, len(args); i < size; i++ {
		gt, err := args[i].Gt(x)
		if err != nil {
			return Zero(), err
		}
		if gt.Bool() {
			x = args[i]
		}
	}
	return x, nil
}

func builtin_rand(args ...Value) (Value, error) {
	if len(args) == 0 {
		return Int(int64(rand.Intn(10000))), nil
	}
	if len(args) == 1 {
		x := args[0].Int()
		if x <= 0 {
			return Zero(), fmt.Errorf("bad argument for function `rand`: argument %v <= 0", x)
		}
	}
	if len(args) == 2 {
		x, y := int(args[0].Int()), int(args[1].Int())
		if x > y {
			return Zero(), fmt.Errorf("bad arguments for function `rand`: first > second")
		}
		return Int(int64(rand.Intn(y-x+1) + x)), nil
	}
	return Zero(), fmt.Errorf("too many arguments for function `rand`: arguments size=%d", len(args))
}
