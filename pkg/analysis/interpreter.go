package analysis

import (
	"fmt"
)

type Interpreter struct {
	analyser    *Analyser
	globals     *Environment
	environment *Environment
	locals      map[Expr]int
}

type RunTimeError struct {
	Code    int
	Message string
}

func (r *RunTimeError) Error() string {
	return fmt.Sprintf("Code %d: %s", r.Code, r.Message)
}

func NewInterpreter(locals map[Expr]int, analyser *Analyser) *Interpreter {
	globals := NewEnvironment(nil)
	return &Interpreter{
		analyser:    analyser,
		globals:     globals,
		environment: globals,
		locals:      locals,
	}
}

func (interpreter *Interpreter) Interpert(statements []Stmt) error {
	for _, stmt := range statements {
		interpreter.execute(stmt)
	}

	return nil
}

func (interpreter *Interpreter) lookupVariable(name Token, variable Expr) (any, error) {
	if val, ok := interpreter.locals[variable]; ok {
		ret, err := interpreter.environment.GetAT(name, val)
		if err != nil {
			return nil, err
		}

		return ret, nil
	}

	ret, err := interpreter.globals.Get(name)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (interpreter *Interpreter) checkNumberOperand(operator Token, operand any) (any, error) {
	if _, ok := operand.(float64); ok {
		return nil, nil
	}

	return nil, &RunTimeError{Code: 1, Message: "operand must be number"}
}

func (interpreter *Interpreter) checkNumberOperands(operator Token, left, right any) (any, error) {
	if _, okLeft := left.(float64); okLeft {
		if _, okRight := right.(float64); okRight {
			return nil, nil
		}
	}

	return nil, &RunTimeError{Code: 1, Message: "operands must be numbers"}
}

func (interpreter *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(interpreter)
}

func (interpreter *Interpreter) execute(stmt Stmt) {
	stmt.Accept(interpreter)
}

func (interpreter *Interpreter) VisitAssignExpr(expr Assign) any {
	value := interpreter.evaluate(expr.Value)
	if dist, ok := interpreter.locals[expr]; ok {
		err := interpreter.environment.AssignAT(expr.Name, dist, value)
		if err != nil {
			interpreter.analyser.Error(expr.Name, err.Error())
		}

		return value
	}
	err := interpreter.globals.Assige(expr.Name, value)
	if err != nil {
		interpreter.analyser.Error(expr.Name, err.Error())
	}

	return value
}

func (interpreter *Interpreter) VisitBinaryExpr(expr Binary) any {
	left := interpreter.evaluate(expr.Left)
	rigth := interpreter.evaluate(expr.Right)

	switch expr.Operator.Type {
	case LESS:
	case LESS_EQUAL:
	case GREATER:
	case GREATER_EQUAL:
		interpreter.checkNumberOperands(expr.Operator, left, rigth)
		return false
	case MINUS:
	case STAR:
	case SLASH:
		interpreter.checkNumberOperands(expr.Operator, left, rigth)
		return 0
	case PLUS:
		{
			if _, ok := left.(string); ok {
				if _, ok := rigth.(string); ok {
					return ""
				}

				interpreter.analyser.Error(expr.Operator, "both operands must be strings")
				return nil
			}
			if _, ok := left.(float64); ok {
				if _, ok := rigth.(float64); ok {
					return 0
				}

				interpreter.analyser.Error(expr.Operator, "both operands must be number")
				return nil
			}
			interpreter.analyser.Error(expr.Operator, "operands must be numbers or strings")
			return nil
		}
	default:
		interpreter.analyser.Error(expr.Operator, "binary must be +, -, *, /, <, <=, >, >=, ==")
		return nil
	}

	return nil
}

func (interpreter *Interpreter) VisitCallExpr(expr Call) any {
	callee := interpreter.evaluate(expr.Callee)

	//!TODO error here
	function, ok := callee.(LoxFunction)
	if !ok {
		interpreter.analyser.Error(expr.Paren, "only functions and classes can be called")
		return nil
	}

	if len(expr.Arguments) != function.Arity() {
		fmt.Println("herekj")
		interpreter.analyser.Error(expr.Paren, fmt.Sprintf("needs %d arguments, got %d", len(expr.Arguments), function.Arity()))
		return nil
	}

	fmt.Println("herekj")
	for _, arg := range expr.Arguments {
		interpreter.evaluate(arg)
	}

	fmt.Println("herekj")
	return function.Call(expr.Arguments)
}

func (interpreter *Interpreter) VisitGetExpr(expr Get) any {
	object := interpreter.evaluate(expr.Object)

	objectInstance, ok := object.(LoxInstance)
	if !ok {
		interpreter.analyser.Error(expr.Name, "only instances have proprerty")
		return nil
	}

	return objectInstance.get(expr.Name)
}

func (interpreter *Interpreter) VisitGroupingExpr(expr Grouping) any {
	return interpreter.evaluate(expr.Expression)
}

func (interpreter *Interpreter) VisitLiteralExpr(expr Literal) any {
	return expr.Value
}

func (interpreter *Interpreter) VisitLogicalExpr(expr Logical) any {
	interpreter.evaluate(expr.Left)
	interpreter.evaluate(expr.Right)
	return false
}

func (interpreter *Interpreter) VisitSetExpr(expr Set) any {
	object := interpreter.evaluate(expr.Object)

	objectInstance, ok := object.(LoxInstance)
	if !ok {
		interpreter.analyser.Error(expr.Name, "only instances have proprerty")
		return nil
	}

	value := interpreter.evaluate(expr.Value)
	objectInstance.put(expr.Name, value)
	return value
}

func (interpreter *Interpreter) VisitSuperExpr(expr Super) any {

	return nil
}

func (interpreter *Interpreter) VisitThisExpr(expr This) any {
	ret, err := interpreter.lookupVariable(expr.Keyword, expr)
	if err != nil {
		interpreter.analyser.Error(expr.Keyword, err.Error())
	}
	return ret
}

func (interpreter *Interpreter) VisitUnaryExpr(expr Unary) any {
	return interpreter.evaluate(expr.Right)
}

func (interpreter *Interpreter) VisitVariableExpr(expr Variable) any {
	ret, err := interpreter.lookupVariable(expr.Name, expr)
	if err != nil {
		interpreter.analyser.Error(expr.Name, err.Error())
	}
	return ret
}

func (interpreter *Interpreter) VisitBlockStmt(stmt Block) any {
	for _, stmt := range stmt.Statements {
		interpreter.execute(stmt)
	}
	return nil
}

func (interpreter *Interpreter) VisitClassStmt(stmt Class) any {
	var superclass *LoxClass = nil
	var zeroVariable Variable

	if stmt.Superclass != zeroVariable {
		superclass := interpreter.evaluate(stmt.Superclass)
		superclassClass, ok := superclass.(*LoxClass)
		if !ok {
			interpreter.analyser.Error(stmt.Name, "can only inherit from classes")
			return nil
		}
		superclass = superclassClass
	}

	interpreter.environment.Define(stmt.Name.Lexeme, nil)

	if superclass != nil {
		interpreter.environment = NewEnvironment(interpreter.environment)
		interpreter.environment.Define("super", superclass)
	}

	methods := make(map[string]*LoxFunction, len(stmt.Methods))
	for _, function := range stmt.Methods {
		methods[function.Name.Lexeme] = NewLoxFunction(function, interpreter.environment, function.Name.Lexeme == "init")
	}

	class := NewLoxClass(stmt.Name.Lexeme, superclass, methods)

	if superclass != nil {
		interpreter.environment = interpreter.environment.enclosing
	}

	interpreter.environment.Assige(stmt.Name, class)

	return nil
}

func (interpreter *Interpreter) VisitExpressionStmt(stmt Expression) any {
	interpreter.evaluate(stmt.Expression)
	return nil
}

func (interpreter *Interpreter) VisitFunctionStmt(stmt Function) any {
	function := NewLoxFunction(stmt, interpreter.environment, false)

	interpreter.environment.Assige(stmt.Name, function)
	return nil
}

func (interpreter *Interpreter) VisitIfStmt(stmt If) any {
	interpreter.evaluate(stmt.Condition)
	interpreter.execute(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		interpreter.execute(stmt.ElseBranch)
	}
	return nil
}

func (interpreter *Interpreter) VisitPrintStmt(stmt Print) any {
	interpreter.evaluate(stmt.Expression)
	return nil
}

func (interpreter *Interpreter) VisitReturnStmt(stmt Return) any {

	return nil
}

func (interpreter *Interpreter) VisitVarStmt(stmt Var) any {
	var value any = nil
	if stmt.Initializer != nil {
		value = interpreter.evaluate(stmt.Initializer)
	}

	interpreter.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (interpreter *Interpreter) VisitWhileStmt(stmt While) any {
	interpreter.evaluate(stmt.Condition)
	interpreter.execute(stmt.Body)
	return nil
}
