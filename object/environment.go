package object

type Environment struct {
	Store map[string]Object
	Outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{Store: make(map[string]Object), Outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.Store[name]
	return val, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.Store[name] = val
	return val
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outer
	return env
}
