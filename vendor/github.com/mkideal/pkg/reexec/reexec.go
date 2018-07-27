package reexec

import (
	"fmt"
	"os"
	"path/filepath"
)

type Executor func()

type Registry struct {
	executors map[string]Executor
}

func New() *Registry {
	return &Registry{
		executors: map[string]Executor{},
	}
}

func (r *Registry) Register(name string, executor Executor) error {
	if _, ok := r.executors[name]; ok {
		return fmt.Errorf("%s have registered", name)
	}
	r.executors[name] = executor
	return nil
}

func (r *Registry) Exec(name string) bool {
	executor, ok := r.executors[name]
	if ok {
		executor()
		return true
	}
	return false
}

func (r *Registry) Startup() bool {
	name := filepath.Base(os.Args[0])
	return r.Exec(name)
}

// global registry
var registry = New()

func Register(name string, executor Executor) {
	if err := registry.Register(name, executor); err != nil {
		panic(err)
	}
}

func Startup() bool {
	return registry.Startup()
}
