package analysis

type LoxClass struct {
	Name       string
	superClass *LoxClass
	methods    map[string]*LoxFunction
}

func NewLoxClass(name string, superClass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		Name:       name,
		superClass: superClass,
		methods:    methods,
	}
}

func (class *LoxClass) findMethod(name string) *LoxFunction {
	if val, ok := class.superClass.methods[name]; ok {
		return val
	}

	if class.superClass != nil {
		return class.superClass.findMethod(name)
	}

	return nil
}

func (class *LoxClass) Arity() int {
	initializer := class.findMethod("init")
	if initializer != nil {
		return initializer.Arity()
	}

	return 0
}

func (class *LoxClass) Call(args ...any) any {
	instance := NewInstance(class)

	initializer := class.findMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(args...)
	}

	return instance
}
