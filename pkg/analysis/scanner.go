package analysis

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	analyser  *Analyser
	tokens    []Token
	source    []byte
	start     int
	current   int
	line      int
	length    int
	startChar int
	endChar   int
	keyWords  map[string]TokenType
}

func NewScanner(source []byte, analyser *Analyser) Scanner {
	return Scanner{
		analyser:  analyser,
		source:    source,
		tokens:    []Token{},
		start:     0,
		current:   0,
		startChar: 0,
		endChar:   0,
		line:      1,
		length:    len(source),
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

func (scanner *Scanner) Scan() []Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.startChar = scanner.endChar
		scanner.scanToken()
	}

	scanner.addToken(EOF, nil)
	return scanner.tokens
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
				scanner.line++
				scanner.startChar = 0
				scanner.endChar = 0
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
	case '\r':
	case '\t':
		{
			break
		}
	case '\n':
		{
			scanner.line++
			scanner.startChar = 0
			scanner.endChar = 0
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

			scanner.analyser.Error(Token{
				StartLine: scanner.line,
				EndLine:   scanner.line,
				StartChar: scanner.startChar,
				EndChar:   scanner.endChar,
				Lexeme:    "@",
			}, fmt.Sprintf("Unexpected token %s", c))
		}
	}

}

func (scanner *Scanner) advance() byte {
	ret := scanner.source[scanner.current]
	scanner.current++
	scanner.endChar++

	return ret
}

func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return 0
	}

	return scanner.source[scanner.current]
}

func (scanner *Scanner) peekAhead() byte {
	if scanner.current+1 >= scanner.length {
		return 0
	}

	return scanner.source[scanner.current+1]
}

func (scanner *Scanner) match(c byte) bool {
	if scanner.peek() != c {
		return false
	}

	scanner.current++
	scanner.endChar++
	return true
}

func (scanner *Scanner) identifier() {
	for !scanner.isAtEnd() && scanner.isAlphaNumircal() {
		scanner.advance()
	}

	text := string(scanner.source[scanner.start:scanner.current])
	identifier, ok := scanner.keyWords[text]
	if !ok {
		identifier = IDENTIFIER
	}

	scanner.addToken(identifier, nil)
}

func (scanner *Scanner) string() {
	for !scanner.isAtEnd() && scanner.peek() != '"' {
		if scanner.peek() == '\n' {
			scanner.line++
			scanner.startChar = 0
			scanner.endChar = 0
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		scanner.analyser.Error(Token{
			StartLine: scanner.line,
			EndLine:   scanner.line,
			StartChar: scanner.startChar,
			EndChar:   scanner.endChar,
			Lexeme:    "@",
		}, "Unterminated string")
		return
	}

	scanner.advance()

	text := string(scanner.source[scanner.start+1 : scanner.current-1])
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

	number, err := strconv.ParseFloat(string(scanner.source[scanner.start:scanner.current]), 64)
	if err != nil {
		scanner.analyser.Error(Token{
			StartLine: scanner.line,
			EndLine:   scanner.line,
			StartChar: scanner.startChar,
			EndChar:   scanner.endChar,
			Lexeme:    "@",
		}, "not valid number represntaion")
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
	return scanner.isAlpha(scanner.peek()) || scanner.isDigit(scanner.peek())
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= scanner.length
}

func (scanner *Scanner) addToken(tokenType TokenType, literal any) {
	scanner.tokens = append(scanner.tokens, Token{
		Type:      tokenType,
		Lexeme:    string(scanner.source[scanner.start:scanner.current]),
		Literal:   literal,
		StartLine: scanner.line,
		StartChar: scanner.startChar,
		EndChar:   scanner.endChar,
	})
}
