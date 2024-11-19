package analysis

type VisitStmt interface {
	VisitBlockStmt(stmt Block) any
	VisitClassStmt(stmt Class) any
	VisitExpressionStmt(stmt Expression) any
	VisitFunctionStmt(stmt Function) any
	VisitIfStmt(stmt If) any
	VisitPrintStmt(stmt Print) any
	VisitReturnStmt(stmt Return) any
	VisitVarStmt(stmt Var) any
	VisitWhileStmt(stmt While) any
}

type Stmt interface {
	Accept(visitor VisitStmt) any
}

// Block
type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt) Block {
	return Block{statements}
}

func (b Block) Accept(visitor VisitStmt) any {
	return visitor.VisitBlockStmt(b)
}

// Class
type Class struct {
	Name       Token
	Superclass Variable
	Methods    []Function
}

func NewClass(name Token, superclass Variable, methods []Function) Class {
	return Class{name, superclass, methods}
}

func (c Class) Accept(visitor VisitStmt) any {
	return visitor.VisitClassStmt(c)
}

// Expression
type Expression struct {
	Expression Expr
}

func NewExpression(expression Expr) Expression {
	return Expression{expression}
}

func (e Expression) Accept(visitor VisitStmt) any {
	return visitor.VisitExpressionStmt(e)
}

// Function
type Function struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func NewFunction(name Token, params []Token, body []Stmt) Function {
	return Function{name, params, body}
}

func (f Function) Accept(visitor VisitStmt) any {
	return visitor.VisitFunctionStmt(f)
}

// If
type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) If {
	return If{condition, thenBranch, elseBranch}
}

func (i If) Accept(visitor VisitStmt) any {
	return visitor.VisitIfStmt(i)
}

// Print
type Print struct {
	Expression Expr
}

func NewPrint(expression Expr) Print {
	return Print{expression}
}

func (p Print) Accept(visitor VisitStmt) any {
	return visitor.VisitPrintStmt(p)
}

// Return
type Return struct {
	Keyword Token
	Value   Expr
}

func NewReturn(keyword Token, value Expr) Return {
	return Return{keyword, value}
}

func (r Return) Accept(visitor VisitStmt) any {
	return visitor.VisitReturnStmt(r)
}

// Var
type Var struct {
	Name        Token
	Initializer Expr
}

func NewVar(name Token, initializer Expr) Var {
	return Var{name, initializer}
}

func (v Var) Accept(visitor VisitStmt) any {
	return visitor.VisitVarStmt(v)
}

// While
type While struct {
	Condition Expr
	Body      Stmt
}

func NewWhile(condition Expr, body Stmt) While {
	return While{condition, body}
}

func (w While) Accept(visitor VisitStmt) any {
	return visitor.VisitWhileStmt(w)
}
