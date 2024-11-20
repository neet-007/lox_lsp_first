package analysis

import (
	"fmt"

	"github.com/neet-007/lox_lsp_first/internal/lsp"
)

type Parser struct {
	analyser    *Analyser
	tokens      []Token
	current     int
	diagnostics []lsp.Diagnostic
}

type ParseError struct {
	Code    int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Code %d: %s", e.Code, e.Message)
}

func NewParser(tokens []Token, analyser *Analyser) Parser {
	return Parser{
		analyser:    analyser,
		tokens:      tokens,
		current:     0,
		diagnostics: []lsp.Diagnostic{},
	}
}

func (parser *Parser) Parse() []Stmt {
	statments := []Stmt{}

	for !parser.isAtEnd() {
		statments = append(statments, parser.declaration())
	}

	return statments
}

func (parser *Parser) declaration() Stmt {
	if parser.match(VAR) {
		stmt, err := parser.varDeclaration()
		if err != nil {
			if _, ok := err.(*ParseError); ok {
				parser.synchronize()
				return nil
			}
			panic(err)
		}
		return stmt
	}

	if parser.match(CLASS) {
		stmt, err := parser.classDeclaration()
		if err != nil {
			if _, ok := err.(*ParseError); ok {
				parser.synchronize()
				return nil
			}
			panic(err)
		}
		return stmt
	}
	if parser.match(FUN) {
		stmt, err := parser.function("function")
		if err != nil {
			if _, ok := err.(*ParseError); ok {
				parser.synchronize()
				return nil
			}
			panic(err)
		}
		return stmt
	}

	stmt, err := parser.statement()
	if err != nil {
		if _, ok := err.(*ParseError); ok {
			parser.synchronize()
			return nil
		}
		panic(err)
	}
	return stmt
}

func (parser *Parser) function(kind string) (Function, error) {
	name, err := parser.consume(IDENTIFIER, fmt.Sprintf("Expect name for %s ", kind))
	if err != nil {
		return Function{}, err
	}

	_, err = parser.consume(LEFT_PAREN, fmt.Sprintf("Expect ( after name for %s ", kind))
	if err != nil {
		return Function{}, err
	}

	params := []Token{}
	if !parser.check(RIGHT_PAREN) {
		param, err := parser.consume(IDENTIFIER, fmt.Sprintf("Expect param for %s ", kind))
		if err != nil {
			return Function{}, err
		}

		params = append(params, *param)

		for parser.match(COMMA) && !parser.isAtEnd() {
			if len(params) >= 256 {
				return Function{}, &ParseError{
					Code:    1,
					Message: "cant have params for than 256",
				}
			}
			param, err = parser.consume(IDENTIFIER, fmt.Sprintf("Expect param for %s ", kind))
			if err != nil {
				return Function{}, err
			}

			params = append(params, *param)
		}
	}

	_, err = parser.consume(RIGHT_PAREN, fmt.Sprintf("Expect ) after params for %s ", kind))
	if err != nil {
		return Function{}, err
	}

	_, err = parser.consume(LEFT_BRACE, fmt.Sprintf("Expect { before body for %s ", kind))
	if err != nil {
		return Function{}, err
	}

	body, err := parser.block()
	if err != nil {
		return Function{}, err
	}

	return NewFunction(*name, params, body), nil
}

func (parser *Parser) classDeclaration() (Stmt, error) {
	name, err := parser.consume(IDENTIFIER, "Expect identifier after class")
	if err != nil {
		return nil, err
	}

	var superclass Variable
	if parser.match(LESS) {
		_, err = parser.consume(IDENTIFIER, "Expect superclass name after < ")
		if err != nil {
			return nil, err
		}
		superclass = NewVariable(*parser.previous())
	}

	_, err = parser.consume(LEFT_BRACE, "Expect { before class body")
	if err != nil {
		return nil, err
	}

	methods := []Function{}

	for !parser.check(RIGHT_BRACE) && !parser.isAtEnd() {
		method, err := parser.function("method")
		if err != nil {
			return nil, err
		}

		methods = append(methods, method)
	}

	_, err = parser.consume(RIGHT_BRACE, "Expect } after class body")
	if err != nil {
		return nil, err
	}

	return NewClass(*name, superclass, methods), nil
}

func (parser *Parser) varDeclaration() (Stmt, error) {
	name, err := parser.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr = nil
	if parser.match(EQUAL) {
		initializer, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = parser.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return NewVar(*name, initializer), nil
}

func (parser *Parser) statement() (Stmt, error) {
	if parser.match(FOR) {
		return parser.forStatement()
	}
	if parser.match(IF) {
		return parser.ifStatement()
	}
	if parser.match(PRINT) {
		return parser.printStatement()
	}
	if parser.match(RETURN) {
		return parser.returnStatement()
	}
	if parser.match(WHILE) {
		return parser.whileStatement()
	}
	if parser.match(LEFT_BRACE) {
		stmts, err := parser.block()
		if err != nil {
			return nil, err
		}

		return NewBlock(stmts), nil
	}

	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}
	_, err = parser.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return NewExpression(expr), nil
}

func (parser *Parser) block() ([]Stmt, error) {
	stmts := []Stmt{}

	for !parser.check(RIGHT_BRACE) && !parser.isAtEnd() {
		stmt, err := parser.statement()
		if err != nil {
			return []Stmt{}, err
		}

		stmts = append(stmts, stmt)
	}

	_, err := parser.consume(RIGHT_BRACE, "Expect '}' after block")
	if err != nil {
		return []Stmt{}, err
	}

	return stmts, nil
}

func (parser *Parser) whileStatement() (Stmt, error) {
	_, err := parser.consume(LEFT_PAREN, "Expect '(' before condition")
	if err != nil {
		return nil, err
	}

	condition, err := parser.expression()
	if err != nil {
		return nil, err
	}

	_, err = parser.consume(RIGHT_PAREN, "Expect ')' after condition")
	if err != nil {
		return nil, err
	}

	body, err := parser.statement()
	if err != nil {
		return nil, err
	}

	return NewWhile(condition, body), nil
}

func (parser *Parser) returnStatement() (Stmt, error) {
	token := parser.previous()
	var value Expr = nil
	var err error
	if !parser.check(SEMICOLON) {
		value, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = parser.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return NewReturn(*token, value), nil
}

func (parser *Parser) printStatement() (Stmt, error) {
	value, err := parser.expression()
	if err != nil {
		return nil, err
	}

	_, err = parser.consume(SEMICOLON, "Expect ';' after statement")
	if err != nil {
		return nil, err
	}

	return NewPrint(value), nil
}

func (parser *Parser) ifStatement() (Stmt, error) {
	_, err := parser.consume(LEFT_PAREN, "Expect '(' before statement")
	if err != nil {
		return nil, err
	}

	condition, err := parser.expression()
	if err != nil {
		return nil, err
	}

	_, err = parser.consume(RIGHT_PAREN, "Expect ')' before statement")
	if err != nil {
		return nil, err
	}

	thenBranch, err := parser.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt = nil
	if parser.match(ELSE) {
		elseBranch, err = parser.statement()
		if err != nil {
			return nil, err
		}
	}

	stmt := NewIf(condition, thenBranch, elseBranch)

	return stmt, nil
}

func (parser *Parser) forStatement() (Stmt, error) {
	_, err := parser.consume(LEFT_PAREN, "Expect '(' before initializer")
	if err != nil {
		return nil, err
	}

	var initializer Stmt
	if parser.match(SEMICOLON) {
		initializer = nil
	} else if parser.match(VAR) {
		initializer, err = parser.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = parser.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if !parser.check(SEMICOLON) {
		condition, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = parser.consume(SEMICOLON, "Expect ';' between for")
	if err != nil {
		return nil, err
	}

	var incerment Expr
	if !parser.check(RIGHT_PAREN) {
		incerment, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = parser.consume(RIGHT_PAREN, "Expect ')' between for")
	if err != nil {
		return nil, err
	}

	body, err := parser.statement()
	if err != nil {
		return nil, err
	}

	if incerment != nil {
		body = NewBlock(
			[]Stmt{body, NewExpression(incerment)},
		)
	}

	if condition == nil {
		condition = NewLiteral(true)
	}

	body = NewWhile(condition, body)

	if initializer != nil {
		body = NewBlock([]Stmt{initializer, body})
	}

	return body, nil
}

func (parser *Parser) expressionStatement() (Stmt, error) {
	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	parser.consume(SEMICOLON, "Expect ';' after expression.")

	return NewExpression(expr), nil
}

func (parser *Parser) expression() (Expr, error) {
	return parser.assignment()
}

func (parser *Parser) assignment() (Expr, error) {
	expr, err := parser.or()
	if err != nil {
		return nil, err
	}

	if parser.match(EQUAL) {
		equals := parser.previous()
		value, err := parser.assignment()
		if err != nil {
			return nil, err
		}

		varExpr, varOk := expr.(Variable)
		if varOk {
			name := varExpr.Name
			return NewAssign(name, value), nil
		}

		parser.error(*equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (parser *Parser) or() (Expr, error) {
	expr, err := parser.and()
	if err != nil {
		return nil, err
	}

	if parser.match(OR) {
		operator := parser.previous()
		right, err := parser.and()
		if err != nil {
			return nil, err
		}

		return NewLogical(expr, *operator, right), nil
	}

	return expr, nil
}

func (parser *Parser) and() (Expr, error) {
	expr, err := parser.equality()
	if err != nil {
		return nil, err
	}

	if parser.match(AND) {
		operator := parser.previous()
		right, err := parser.equality()
		if err != nil {
			return nil, err
		}

		return NewLogical(expr, *operator, right), nil
	}

	return expr, nil
}

func (parser *Parser) equality() (Expr, error) {
	expr, err := parser.comparission()
	if err != nil {
		return nil, err
	}

	if parser.match(EQUAL_EQUAL, BANG_EQUAL) {
		operator := parser.previous()
		rigth, err := parser.comparission()
		if err != nil {
			return nil, err
		}

		return NewBinary(expr, *operator, rigth), nil
	}
	return expr, nil
}

func (parser *Parser) comparission() (Expr, error) {
	expr, err := parser.term()
	if err != nil {
		return nil, err
	}

	if parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.previous()
		rigth, err := parser.term()
		if err != nil {
			return nil, err
		}

		return NewBinary(expr, *operator, rigth), nil
	}
	return expr, nil
}

func (parser *Parser) term() (Expr, error) {
	expr, err := parser.factor()
	if err != nil {
		return nil, err
	}

	if parser.match(MINUS, PLUS) {
		operator := parser.previous()
		rigth, err := parser.factor()
		if err != nil {
			return nil, err
		}

		return NewBinary(expr, *operator, rigth), nil
	}
	return expr, nil
}

func (parser *Parser) factor() (Expr, error) {
	expr, err := parser.unary()
	if err != nil {
		return nil, err
	}

	if parser.match(SLASH, STAR) {
		operator := parser.previous()
		rigth, err := parser.unary()
		if err != nil {
			return nil, err
		}

		return NewBinary(expr, *operator, rigth), nil
	}
	return expr, nil
}

func (parser *Parser) unary() (Expr, error) {
	if parser.match(BANG, MINUS) {
		operator := parser.previous()
		right, err := parser.unary()
		if err != nil {
			return nil, err
		}

		return NewUnary(*operator, right), nil
	}

	return parser.call()
}

func (parser *Parser) call() (Expr, error) {
	expr, err := parser.primary()
	if err != nil {
		return nil, err
	}

	for {
		if parser.match(LEFT_PAREN) {
			expr, err = parser.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if parser.match(DOT) {
			name, err := parser.consume(IDENTIFIER, "Expect proprety name after '.'")
			if err != nil {
				return nil, err
			}
			expr = NewGet(expr, *name)
		} else {
			break
		}
	}

	return expr, nil
}

func (parser *Parser) finishCall(callee Expr) (Expr, error) {
	args := []Expr{}

	if !parser.check(RIGHT_PAREN) {
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		for !parser.isAtEnd() && parser.match(COMMA) {
			if len(args) >= 256 {
				return nil, &ParseError{
					Code:    1,
					Message: "cant have more than 255 arguemnts ",
				}
			}
			expr, err = parser.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
	}

	rightParen, err := parser.consume(RIGHT_PAREN, "Expect ')' after call")
	if err != nil {
		return nil, err
	}

	return NewCall(callee, *rightParen, args), nil
}

func (parser *Parser) primary() (Expr, error) {
	if parser.match(FALSE) {
		return NewLiteral(false), nil
	}
	if parser.match(TRUE) {
		return NewLiteral(true), nil
	}
	if parser.match(NIL) {
		return NewLiteral(nil), nil
	}

	if parser.match(STRING, NUMBER) {
		prev := *parser.previous()
		return NewLiteral(prev.Literal), nil
	}

	if parser.match(IDENTIFIER) {
		return NewVariable(*parser.previous()), nil
	}

	if parser.match(LEFT_PAREN) {
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}

		_, err = parser.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return NewGrouping(expr), nil
	}

	return nil, parser.error(*parser.peek(), "Expect expression.")
}

func (parser *Parser) isAtEnd() bool {
	return parser.current >= len(parser.tokens)
}

func (parser *Parser) check(tokenType TokenType) bool {
	curr := parser.peek()
	if curr == nil {
		return false
	}

	return curr.Type == tokenType
}

func (parser *Parser) previous() *Token {
	return &parser.tokens[parser.current-1]
}

func (parser *Parser) advance() *Token {
	if !parser.isAtEnd() {
		parser.current++
	}

	return parser.previous()
}

func (parser *Parser) consume(tokenType TokenType, msg string) (*Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	}

	// !TODO make error
	err := parser.error(*parser.peek(), msg)
	return nil, err
}

func (parser *Parser) error(token Token, msg string) error {
	fmt.Println("parsing error")
	parser.analyser.Error(token, msg)
	return &ParseError{
		Code:    1,
		Message: msg,
	}
}

func (parser *Parser) match(tokenTypes ...TokenType) bool {
	for _, t := range tokenTypes {
		if parser.check(t) {
			parser.advance()
			return true
		}
	}

	return false
}

func (parser *Parser) peek() *Token {
	return &parser.tokens[parser.current]
}

func (parser *Parser) synchronize() {
	parser.advance()

	for !parser.isAtEnd() {
		if parser.previous().Type == SEMICOLON {
			return
		}

		switch parser.peek().Type {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		parser.advance()
	}
}
