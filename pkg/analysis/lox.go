package analysis

import "github.com/neet-007/lox_lsp_first/internal/lsp"

func Analyse(source []byte) ([]Token, []lsp.Diagnostic) {
	scanner := NewScanner(source)

	scanner.Scan()
	tokens := scanner.Tokens
	errors := scanner.Erros

	return tokens, errors
}

func Error(token Token, message string) {

}
