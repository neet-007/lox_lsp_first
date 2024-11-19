package analysis

import (
	"log"

	"github.com/neet-007/lox_lsp_first/internal/lsp"
)

func Analyse(source []byte, logger *log.Logger) ([]Token, []lsp.Diagnostic) {
	scanner := NewScanner(source)

	scanner.Scan(logger)
	tokens := scanner.Tokens
	errors := scanner.Erros

	for _, v := range errors {
		logger.Printf("message:%s, seviry:%d, startLine:%d, startChar:%d\n", v.Message, v.Severity, v.Range.Start.Line, v.Range.Start.Character)
	}

	return tokens, errors
}

func Error(token Token, message string) {

}
