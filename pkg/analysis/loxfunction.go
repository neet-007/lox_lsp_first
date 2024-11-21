package analysis

type LoxFunction struct {
	declaration   Function
	clouser       *Environment
	isInitializer bool
}

func NewLoxFunction(declaration Function, clouser *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		declaration:   declaration,
		clouser:       clouser,
		isInitializer: isInitializer,
	}
}

func (lFun LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	env := NewEnvironment(lFun.clouser)
	env.Define("this", instance)

	return NewLoxFunction(lFun.declaration, env, lFun.isInitializer)
}

func (lFun LoxFunction) Arity() int {
	return len(lFun.declaration.Params)
}

func (lFun LoxFunction) Call(args ...any) any {
	env := NewEnvironment(lFun.clouser)
	for i, param := range lFun.declaration.Params {
		env.Define(param.Lexeme, args[i])
	}

	//!TODO do lsp things
	return nil
}
