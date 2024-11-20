package analysis

type FunctionType int
type ClassType int

const (
	NONE_FUNCTION FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

const (
	NONE_CLASS ClassType = iota
	CLASS_RESOLVER
	SUBCLASS
)

type Resolver struct {
	analyser        *Analyser
	scopes          []map[string]bool
	currentFunction FunctionType
	currrntClass    ClassType
	locals          map[Expr]int
}

func NewResolver(analyser *Analyser) *Resolver {
	return &Resolver{
		analyser:        analyser,
		scopes:          []map[string]bool{},
		currentFunction: NONE_FUNCTION,
		currrntClass:    NONE_CLASS,
		locals:          map[Expr]int{},
	}
}

func (resolver *Resolver) Resolve(statemnts []Stmt) {
	for _, stmt := range statemnts {
		resolver.resolveStmt(stmt)
	}
}

func (resolver *Resolver) VisitSuperExpr(expr Super) any {
	if resolver.currrntClass == NONE_CLASS {
		resolver.analyser.Error(expr.Keyword, "can not use 'super' outside class")
		return nil
	}

	if resolver.currrntClass != SUBCLASS {
		resolver.analyser.Error(expr.Keyword, "can not use 'super' for non subclass")
		return nil
	}

	resolver.resolveLocal(expr, expr.Keyword)
	return nil
}

func (resolver *Resolver) VisitThisExpr(expr This) any {
	if resolver.currrntClass == NONE_CLASS {
		resolver.analyser.Error(expr.Keyword, "can not use 'this' outside class")
		return nil
	}

	resolver.resolveLocal(expr, expr.Keyword)
	return nil
}

func (resolver *Resolver) resolveFunctionStmt(function Function, kind FunctionType) {
	resolver.declare(function.Name)
	resolver.define(function.Name)

	resolver.resolveFunction(function, FUNCTION)
}

func (resolver *Resolver) resolveFunction(stmt Function, kind FunctionType) {
	enclosing := resolver.currentFunction
	resolver.currentFunction = kind

	resolver.beginScope()

	for _, param := range stmt.Params {
		resolver.declare(param)
		resolver.define(param)
	}

	resolver.Resolve(stmt.Body)
	resolver.endScope()

	resolver.currentFunction = enclosing
}

func (resolver *Resolver) VisitReturnStmt(stmt Return) any {
	if resolver.currentFunction == NONE_FUNCTION {
		resolver.analyser.Error(stmt.Keyword, "can not use 'return' outside function")
		return nil
	}

	if stmt.Value == nil {
		return nil
	}

	if resolver.currentFunction == INITIALIZER {
		resolver.analyser.Error(stmt.Keyword, "can not use 'return' in initilzier function")
		return nil
	}

	resolver.resolveExpr(stmt.Value)
	return nil
}

func (resolver *Resolver) VisitCallExpr(expr Call) any {
	resolver.resolveExpr(expr.Callee)

	for _, arg := range expr.Arguments {
		resolver.resolveExpr(arg)
	}
	return nil
}

func (resolver *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(resolver)
}

func (resolver *Resolver) resolveExpr(expr Expr) {
	expr.Accept(resolver)
}

func (resolver *Resolver) resolveLocal(expr Expr, token Token) {
	for i := len(resolver.scopes) - 1; i >= 0; i-- {
		if _, ok := resolver.scopes[i][token.Lexeme]; ok {
			resolver.locals[expr] = len(resolver.scopes) - 1 - i
			return
		}
	}
}

func (resolver *Resolver) define(token Token) {
	lenScops := len(resolver.scopes)
	if lenScops == 0 {
		return
	}

	resolver.scopes[lenScops-1][token.Lexeme] = true
}

func (resolver *Resolver) declare(token Token) {
	lenScops := len(resolver.scopes)
	if lenScops == 0 {
		return
	}

	scope := resolver.scopes[lenScops-1]
	_, ok := scope[token.Lexeme]
	if ok {
		resolver.analyser.Error(token, "a variable with this name has already been declared")
	}

	scope[token.Lexeme] = false
}

func (resolver *Resolver) beginScope() {
	resolver.scopes = append(resolver.scopes, map[string]bool{})
}

func (resolver *Resolver) endScope() {
	resolver.scopes = resolver.scopes[:len(resolver.scopes)-1]
}

func (resolver *Resolver) VisitAssignExpr(expr Assign) any {
	resolver.resolveExpr(expr.Value)
	resolver.resolveLocal(expr, expr.Name)
	return nil
}
func (resolver *Resolver) VisitBinaryExpr(expr Binary) any {
	resolver.resolveExpr(expr.Left)
	resolver.resolveExpr(expr.Right)
	return nil
}
func (resolver *Resolver) VisitGetExpr(expr Get) any {
	resolver.resolveExpr(expr.Object)
	return nil
}
func (resolver *Resolver) VisitGroupingExpr(expr Grouping) any {
	resolver.resolveExpr(expr.Expression)
	return nil
}
func (resolver *Resolver) VisitLiteralExpr(expr Literal) any {
	return nil
}
func (resolver *Resolver) VisitLogicalExpr(expr Logical) any {
	resolver.resolveExpr(expr.Left)
	resolver.resolveExpr(expr.Right)
	return nil
}
func (resolver *Resolver) VisitSetExpr(expr Set) any {
	resolver.resolveExpr(expr.Value)
	resolver.resolveExpr(expr.Object)
	return nil
}

func (resolver *Resolver) VisitUnaryExpr(expr Unary) any {
	resolver.resolveExpr(expr.Right)
	return nil
}
func (resolver *Resolver) VisitVariableExpr(expr Variable) any {
	lenScopes := len(resolver.scopes)
	if lenScopes > 0 {
		val, ok := resolver.scopes[lenScopes-1][expr.Name.Lexeme]
		if ok && !val {
			resolver.analyser.Error(expr.Name, "Can't read local variable in its own initializer.")
		}
	}

	resolver.resolveLocal(expr, expr.Name)
	return nil
}

func (resolver *Resolver) VisitBlockStmt(stmt Block) any {
	resolver.Resolve(stmt.Statements)
	return nil
}
func (resolver *Resolver) VisitClassStmt(stmt Class) any {
	enclosingClass := resolver.currrntClass
	resolver.currrntClass = CLASS_RESOLVER

	resolver.declare(stmt.Name)
	resolver.define(stmt.Name)

	var zeroSuperClass Variable
	if stmt.Superclass != zeroSuperClass && stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		resolver.analyser.Error(stmt.Superclass.Name,
			"A class can't inherit from itself.")
	}

	if stmt.Superclass != zeroSuperClass {
		resolver.currrntClass = SUBCLASS
		resolver.resolveExpr(stmt.Superclass)
	}

	if stmt.Superclass != zeroSuperClass {
		resolver.beginScope()
		resolver.scopes[len(resolver.scopes)-1]["super"] = true
	}

	resolver.beginScope()
	resolver.scopes[len(resolver.scopes)-1]["this"] = true

	for _, method := range stmt.Methods {
		declaration := METHOD
		if method.Name.Lexeme == "init" {
			declaration = INITIALIZER
		}

		resolver.resolveFunction(method, declaration)
	}

	resolver.endScope()
	if stmt.Superclass != zeroSuperClass {
		resolver.endScope()
	}

	resolver.currrntClass = enclosingClass
	return nil
}
func (resolver *Resolver) VisitExpressionStmt(stmt Expression) any {
	resolver.resolveExpr(stmt.Expression)
	return nil
}
func (resolver *Resolver) VisitFunctionStmt(stmt Function) any {

	return nil
}
func (resolver *Resolver) VisitIfStmt(stmt If) any {
	resolver.resolveExpr(stmt.Condition)
	resolver.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		resolver.resolveStmt(stmt.ElseBranch)
	}
	return nil
}
func (resolver *Resolver) VisitPrintStmt(stmt Print) any {
	resolver.resolveExpr(stmt.Expression)
	return nil
}
func (resolver *Resolver) VisitVarStmt(stmt Var) any {
	resolver.declare(stmt.Name)
	if stmt.Initializer != nil {
		resolver.resolveExpr(stmt.Initializer)
	}
	resolver.define(stmt.Name)
	return nil
}
func (resolver *Resolver) VisitWhileStmt(stmt While) any {
	resolver.resolveExpr(stmt.Condition)
	resolver.resolveStmt(stmt.Body)
	return nil
}
