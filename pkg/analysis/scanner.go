package analysis

import (
	"fmt"
	"strconv"

	"github.com/neet-007/lox_lsp_first/internal/lsp"
)

type Scanner struct {
	Tokens    []Token
	Erros     []lsp.Diagnostic
	Source    []byte
	Start     int
	Current   int
	Line      int
	Length    int
	StartChar int
	EndChar   int
	keyWords  map[string]TokenType
}

func NewScanner(source []byte) Scanner {
	return Scanner{
		Source:    source,
		Erros:     []lsp.Diagnostic{},
		Tokens:    []Token{},
		Start:     0,
		Current:   0,
		StartChar: 0,
		EndChar:   0,
		Line:      1,
		Length:    len(source),
		keyWords: map[string]TokenType{
			"and":    AND,
			"class":  CLASS,
			"else":   ELSE,
			"false":  FALSE,
			"for":    FOR,
			"fun":    FUN,
			"if":     IF,
			"nil":    NIL,
			"or":     OR,
			"print":  PRINT,
			"return": RETURN,
			"super":  SUPER,
			"this":   THIS,
			"true":   TRUE,
			"var":    VAR,
			"while":  WHILE,
		},
	}
}

func (scanner *Scanner) Scan() {
	for !scanner.isAtEnd() {
		scanner.Start = scanner.Current
		scanner.scanToken()
	}

	scanner.addToken(EOF, nil)
}

func (scanner *Scanner) scanToken() {
	c := scanner.advance()
	switch c {
	case '(':
		{
			scanner.addToken(LEFT_PAREN, nil)
			break
		}
	case ')':
		{
			scanner.addToken(RIGHT_PAREN, nil)
			break
		}
	case '{':
		{
			scanner.addToken(LEFT_BRACE, nil)
			break
		}
	case '}':
		{
			scanner.addToken(RIGHT_BRACE, nil)
			break
		}
	case '+':
		{
			scanner.addToken(PLUS, nil)
			break
		}
	case '-':
		{
			scanner.addToken(MINUS, nil)
			break
		}
	case '*':
		{
			scanner.addToken(STAR, nil)
			break
		}
	case '/':
		{
			if scanner.match('/') {
				for !scanner.match('\n') {
					scanner.advance()
				}
				scanner.Line++
				scanner.StartChar = 0
				scanner.EndChar = 0
				break
			}
			scanner.addToken(SLASH, nil)
			break
		}
	case ',':
		{
			scanner.addToken(COMMA, nil)
			break
		}
	case '.':
		{
			scanner.addToken(DOT, nil)
			break
		}
	case ';':
		{
			scanner.addToken(SEMICOLON, nil)
			break
		}
	case ' ':
	case 'r':
	case '\t':
		{
			break
		}
	case '\n':
		{
			scanner.Line++
			scanner.StartChar = 0
			scanner.EndChar = 0
			break
		}

	case '!':
		{
			if scanner.match('=') {
				scanner.addToken(BANG_EQUAL, nil)
				break
			}
			scanner.addToken(BANG, nil)
			break
		}
	case '=':
		{
			if scanner.match('=') {
				scanner.addToken(EQUAL_EQUAL, nil)
				break
			}
			scanner.addToken(EQUAL, nil)
			break
		}
	case '>':
		{
			if scanner.match('=') {
				scanner.addToken(GREATER_EQUAL, nil)
				break
			}
			scanner.addToken(GREATER, nil)
			break
		}
	case '<':
		{
			if scanner.match('=') {
				scanner.addToken(LESS_EQUAL, nil)
				break
			}
			scanner.addToken(LESS, nil)
			break
		}

	case '"':
		{
			scanner.string()
			break
		}

	default:
		{
			if scanner.isDigit(c) {
				scanner.number()
				break
			}
			if scanner.isAlpha(c) {
				scanner.identifier()
				break
			}

			scanner.addDaiagnostic(fmt.Sprintf("Unexpected token %s", c))
		}
	}

	scanner.StartChar = scanner.EndChar
}

func (scanner *Scanner) advance() byte {
	ret := scanner.Source[scanner.Current]
	scanner.Current++
	scanner.EndChar++

	return ret
}

func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return 0
	}

	return scanner.Source[scanner.Current]
}

func (scanner *Scanner) peekAhead() byte {
	if scanner.Current+1 >= scanner.Length {
		return 0
	}

	return scanner.Source[scanner.Current+1]
}

func (scanner *Scanner) match(c byte) bool {
	if scanner.peek() != c {
		return false
	}

	scanner.Current++
	scanner.EndChar++
	return true
}

func (scanner *Scanner) identifier() {
	for !scanner.isAtEnd() && scanner.isAlphaNumircal() {
		scanner.advance()
	}

	text := string(scanner.Source[scanner.Start:scanner.Current])
	identifier, ok := scanner.keyWords[text]
	if !ok {
		identifier = IDENTIFIER
	}

	scanner.addToken(identifier, nil)
}

func (scanner *Scanner) string() {
	for !scanner.isAtEnd() && scanner.peek() != '"' {
		if scanner.peek() == '\n' {
			scanner.Line++
			scanner.StartChar = 0
			scanner.EndChar = 0
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		scanner.addDaiagnostic("Unterminated string")
		return
	}

	scanner.advance()

	text := string(scanner.Source[scanner.Start+1 : scanner.Current-1])
	scanner.addToken(STRING, text)

	return
}

func (scanner *Scanner) number() {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	if scanner.peek() == '.' && scanner.isDigit(scanner.peekAhead()) {
		scanner.advance()

		for scanner.isDigit(scanner.peek()) {
			scanner.advance()
		}
	}

	number, err := strconv.ParseFloat(string(scanner.Source[scanner.Start:scanner.Current]), 64)
	if err != nil {
		scanner.addDaiagnostic("not valid number represntaion")
		return
	}

	scanner.addToken(NUMBER, number)
	return
}

func (scanner *Scanner) isDigit(c byte) bool {
	return !scanner.isAtEnd() && ('0' <= c && c <= '9')
}

func (scanner *Scanner) isAlpha(c byte) bool {
	return !scanner.isAtEnd() && ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_'
}

func (scanner *Scanner) isAlphaNumircal() bool {
	return scanner.isAlpha(scanner.Source[scanner.Current]) || scanner.isDigit(scanner.Source[scanner.Current])
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.Current >= scanner.Length
}

func (scanner *Scanner) addToken(tokenType TokenType, literal any) {
	scanner.Tokens = append(scanner.Tokens, Token{
		Type:    tokenType,
		Lexeme:  string(scanner.Source[scanner.Start:scanner.Current]),
		Literal: literal,
		Line:    scanner.Line,
	})
}

func (scanner *Scanner) addDaiagnostic(msg string) {
	if len(scanner.Tokens) > 0 {
		diagnostic := lsp.NewDiagnostic(
			lsp.Range{
				Start: lsp.Position{
					Line:      scanner.Line,
					Character: scanner.StartChar,
				},
				End: lsp.Position{
					Line:      scanner.Line,
					Character: scanner.EndChar,
				},
			},
			1,
			scanner.Tokens[0].Lexeme,
			msg,
		)
		scanner.Erros = append(scanner.Erros, diagnostic)
	} else {
		diagnostic := lsp.NewDiagnostic(
			lsp.Range{
				Start: lsp.Position{
					Line:      scanner.Line,
					Character: scanner.StartChar,
				},
				End: lsp.Position{
					Line:      scanner.Line,
					Character: scanner.EndChar,
				},
			},
			1,
			"",
			msg,
		)
		scanner.Erros = append(scanner.Erros, diagnostic)

	}
}
