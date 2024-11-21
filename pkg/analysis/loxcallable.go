package analysis

type LoxCallable interface {
	Arity() int
	Call(args ...any) any
}
