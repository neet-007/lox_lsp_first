package analysis

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(environment *Environment) *Environment {
	return &Environment{
		enclosing: environment,
		values:    map[string]any{},
	}
}

func (env *Environment) Get(token Token) (any, error) {
	if val, ok := env.values[token.Lexeme]; ok {
		return val, nil
	}

	if env.enclosing != nil {
		return env.enclosing.Get(token)
	}

	return nil, &RunTimeError{Code: 1, Message: "cannot get undefined value"}
}

func (env *Environment) Assige(token Token, val any) error {
	if _, ok := env.values[token.Lexeme]; ok {
		env.values[token.Lexeme] = val
	}

	if env.enclosing != nil {
		return env.Assige(token, val)
	}

	return &RunTimeError{Code: 1, Message: "cannot assigne to undefined value"}
}

func (env *Environment) Define(name string, val any) {
	env.values[name] = val
}

func (env *Environment) ancestor(dist int) *Environment {
	newEnv := env
	for _ = range dist {
		newEnv = newEnv.enclosing
	}

	return newEnv
}

func (env *Environment) GetAT(token Token, dist int) (any, error) {
	return env.ancestor(dist).Get(token)
}

func (env *Environment) AssignAT(token Token, dist int, val any) error {
	return env.ancestor(dist).Assige(token, val)
}
