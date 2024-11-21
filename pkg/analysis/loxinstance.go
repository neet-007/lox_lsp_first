package analysis

type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func NewInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: map[string]any{},
	}
}

func (instance *LoxInstance) get(name Token) any {
	if val, ok := instance.fields[name.Lexeme]; ok {
		return val
	}

	method := instance.class.findMethod(name.Lexeme)
	if method != nil {
		return method.Bind(instance)
	}

	//!TDOD error
	return nil
}

func (instance *LoxInstance) put(name Token, val any) {
	instance.fields[name.Lexeme] = val
}
