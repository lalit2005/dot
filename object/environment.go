package object

type Environment struct {
	Store map[string]Object
	Outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{Store: s, Outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.Store[name]
	if !ok && e.Outer != nil {
		obj, ok = e.Outer.Get(name)
	}
	return obj, ok
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
