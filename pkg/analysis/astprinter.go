package analysis

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() AstPrinter {
	return AstPrinter{}
}

func (a *AstPrinter) VisitAssignExpr(expr Assign) any {
	return a.parenthesize(fmt.Sprintf("assign %v", expr.Name.Lexeme), expr.Value)
}

func (a *AstPrinter) VisitBinaryExpr(expr Binary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitCallExpr(expr Call) any {
	return a.parenthesize("call", append([]Expr{expr.Callee}, expr.Arguments...)...)
}

func (a *AstPrinter) VisitGetExpr(expr Get) any {
	return a.parenthesize(fmt.Sprintf("get %v", expr.Name.Lexeme), expr.Object)
}

func (a *AstPrinter) VisitGroupingExpr(expr Grouping) any {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitLogicalExpr(expr Logical) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitSetExpr(expr Set) any {
	return a.parenthesize(fmt.Sprintf("set %v", expr.Name.Lexeme), expr.Object, expr.Value)
}

func (a *AstPrinter) VisitSuperExpr(expr Super) any {
	return fmt.Sprintf("(super %v)", expr.Method.Lexeme)
}

func (a *AstPrinter) VisitThisExpr(expr This) any {
	return "this"
}

func (a *AstPrinter) VisitUnaryExpr(expr Unary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitVariableExpr(expr Variable) any {
	return expr.Name.Lexeme
}

func (a *AstPrinter) VisitBlockStmt(stmt Block) any {
	var stmts []string
	for _, statement := range stmt.Statements {
		stmts = append(stmts, a.print(statement))
	}
	return fmt.Sprintf("(block %s)", strings.Join(stmts, " "))
}

func (a *AstPrinter) VisitClassStmt(stmt Class) any {
	superclass := stmt.Superclass
	var methods []string
	for _, method := range stmt.Methods {
		methods = append(methods, a.print(method))
	}
	return fmt.Sprintf("(class %s superclass [%s] %s)", stmt.Name.Lexeme, a.VisitVariableExpr(superclass), strings.Join(methods, " "))
}

func (a *AstPrinter) VisitFunctionStmt(stmt Function) any {
	var params []string
	for _, param := range stmt.Params {
		params = append(params, param.Lexeme)
	}
	bodyStatms := ""
	for _, bodyStmt := range stmt.Body {
		bodyStatms += a.print(bodyStmt)
	}
	return fmt.Sprintf("(fun %s (%s) %s)", stmt.Name.Lexeme, strings.Join(params, " "), bodyStatms)
}

func (a *AstPrinter) VisitIfStmt(stmt If) any {
	if stmt.ElseBranch != nil {
		return fmt.Sprintf("(if %s %s %s)", a.parenthesize("condition", stmt.Condition), a.print(stmt.ThenBranch), a.print(stmt.ElseBranch))
	}
	return fmt.Sprintf("(if %s %s)", a.parenthesize("condition", stmt.Condition), a.print(stmt.ThenBranch))
}

func (a *AstPrinter) VisitPrintStmt(stmt Print) any {
	return fmt.Sprintf("(print %s)", a.parenthesize("value", stmt.Expression))
}

func (a *AstPrinter) VisitReturnStmt(stmt Return) any {
	if stmt.Value != nil {
		return fmt.Sprintf("(return %s)", a.parenthesize("value", stmt.Value))
	}
	return "(return)"
}

func (a *AstPrinter) VisitVarStmt(stmt Var) any {
	if stmt.Initializer != nil {
		return fmt.Sprintf("(var %s %s)", stmt.Name.Lexeme, a.parenthesize("initializer", stmt.Initializer))
	}
	return fmt.Sprintf("(var %s)", stmt.Name.Lexeme)
}

func (a *AstPrinter) VisitWhileStmt(stmt While) any {
	return fmt.Sprintf("(while %s %s)", a.parenthesize("condition", stmt.Condition), a.print(stmt.Body))
}

func (a *AstPrinter) VisitExpressionStmt(stmt Expression) any {
	return a.parenthesize("expression", stmt.Expression)
}

func (a *AstPrinter) print(stmt Stmt) string {
	if stmt == nil {
		return "nil"
	}

	val := stmt.Accept(a)
	if val == nil {
		return "nil"
	}

	valStrting, ok := val.(string)
	if !ok {
		panic("not ok")
	}

	return valStrting
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		val := expr.Accept(a)
		if val == nil {
			builder.WriteString("nil")
			continue
		}
		valStr, ok := val.(string)
		if !ok {
			builder.WriteString("COULD_NOT_GET_STRING")
			continue
		}
		builder.WriteString(valStr)
	}
	builder.WriteString(")")

	return builder.String()
}
