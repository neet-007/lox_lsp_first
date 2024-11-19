package analysis

type VisitExpr interface {
	VisitAssignExpr(expr Assign) any
	VisitBinaryExpr(expr Binary) any
	VisitCallExpr(expr Call) any
	VisitGetExpr(expr Get) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitLogicalExpr(expr Logical) any
	VisitSetExpr(expr Set) any
	VisitSuperExpr(expr Super) any
	VisitThisExpr(expr This) any
	VisitUnaryExpr(expr Unary) any
	VisitVariableExpr(expr Variable) any
}

type Expr interface {
	Accept(visitor VisitExpr) any
}

// Assign
type Assign struct {
	Name  Token
	Value Expr
}

func NewAssign(name Token, value Expr) Assign {
	return Assign{name, value}
}

func (a Assign) Accept(visitor VisitExpr) any {
	return visitor.VisitAssignExpr(a)
}

// Binary
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func NewBinary(left Expr, operator Token, right Expr) Binary {
	return Binary{left, operator, right}
}

func (b Binary) Accept(visitor VisitExpr) any {
	return visitor.VisitBinaryExpr(b)
}

// Call
type Call struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func NewCall(callee Expr, paren Token, arguments []Expr) Call {
	return Call{callee, paren, arguments}
}

func (c Call) Accept(visitor VisitExpr) any {
	return visitor.VisitCallExpr(c)
}

// Get
type Get struct {
	Object Expr
	Name   Token
}

func NewGet(object Expr, name Token) Get {
	return Get{object, name}
}

func (g Get) Accept(visitor VisitExpr) any {
	return visitor.VisitGetExpr(g)
}

// Grouping
type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) Grouping {
	return Grouping{expression}
}

func (g Grouping) Accept(visitor VisitExpr) any {
	return visitor.VisitGroupingExpr(g)
}

// Literal
type Literal struct {
	Value any
}

func NewLiteral(value any) Literal {
	return Literal{value}
}

func (l Literal) Accept(visitor VisitExpr) any {
	return visitor.VisitLiteralExpr(l)
}

// Logical
type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func NewLogical(left Expr, operator Token, right Expr) Logical {
	return Logical{left, operator, right}
}

func (l Logical) Accept(visitor VisitExpr) any {
	return visitor.VisitLogicalExpr(l)
}

// Set
type Set struct {
	Object Expr
	Name   Token
	Value  Expr
}

func NewSet(object Expr, name Token, value Expr) Set {
	return Set{object, name, value}
}

func (s Set) Accept(visitor VisitExpr) any {
	return visitor.VisitSetExpr(s)
}

// Super
type Super struct {
	Keyword Token
	Method  Token
}

func NewSuper(keyword Token, method Token) Super {
	return Super{keyword, method}
}

func (s Super) Accept(visitor VisitExpr) any {
	return visitor.VisitSuperExpr(s)
}

// This
type This struct {
	Keyword Token
}

func NewThis(keyword Token) This {
	return This{keyword}
}

func (t This) Accept(visitor VisitExpr) any {
	return visitor.VisitThisExpr(t)
}

// Unary
type Unary struct {
	Operator Token
	Right    Expr
}

func NewUnary(operator Token, right Expr) Unary {
	return Unary{operator, right}
}

func (u Unary) Accept(visitor VisitExpr) any {
	return visitor.VisitUnaryExpr(u)
}

// Variable
type Variable struct {
	Name Token
}

func NewVariable(name Token) Variable {
	return Variable{name}
}

func (v Variable) Accept(visitor VisitExpr) any {
	return visitor.VisitVariableExpr(v)
}
