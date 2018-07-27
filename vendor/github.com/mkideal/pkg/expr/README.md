expr [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/pkg/master/LICENSE)
===================================================================================================================================================

License
-------

[The MIT License (MIT)](https://raw.githubusercontent.com/mkideal/pkg/master/LICENSE)

Install
-------

```shell
go get github.com/mkideal/pkg/expr
```

Expr
----

`Expr` is top-level object of `expr` package.

`New` function new an Expr from string `s` and `pool`(using default if nil) or get from pool.

```go
func New(s string, pool *Pool) (*Expr, error)
````

`Eval` method evaluate expression by VarGetter

```go
func (e *Expr) Eval(getter VarGetter) (float64, error)
```

example:

```go
e, _ := expr.New("x+1", nil)
getter := expr.VarGetter(map[string]float64{"x": 1})
result, _ := e.Eval(getter) // result: 2
```

VarGetter
---------

`VarGetter` define an interface for getting variable by name.

`Getter` implements VarGetter using `map[string]float64`.

```go
type VarGetter interface {
	GetVar(string) (float64, error)
}

// default VarGetter implement
type Getter map[string]float64
```

Func
----

`Func` define function type used in expression.

```go
type Func func(...float64) (float64, error)
```

builtin functions:

-	min(arg1[, arg2, ...])
-	max(arg1[, arg2, ...])
-	rand([arg1[, arg2]])
-	iff(ok, x, y) <==> ok ? x : y

Pool
----

`Pool` used for cache expression objects.

```go
func NewPool(factories ...map[string]Func) (*Pool, error)
```

External
--------

-	[exp](https://github.com/mkideal/tools/tree/master/exp) - `a command line app for evaluate expression`
