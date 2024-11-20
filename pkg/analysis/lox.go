package analysis

import (
	"fmt"
	"github.com/neet-007/lox_lsp_first/internal/lsp"
	"log"
)

func Analyse(source []byte, logger *log.Logger) ([]Token, []lsp.Diagnostic) {
	scanner := NewScanner(source)

	scanner.Scan()
	tokens := scanner.Tokens
	_ = scanner.Erros

	/*
		for _, v := range errors {
			logger.Printf("message:%s, seviry:%d, startLine:%d, startChar:%d\n", v.Message, v.Severity, v.Range.Start.Line, v.Range.Start.Character)
		}
	*/

	fmt.Println("before parse")
	parser := NewParser(tokens)

	astPrinter := NewAstPrinter()

	statements := parser.Parse()

	fmt.Println("after parse")

	for _, stmt := range statements {
		if stmt == nil{
			logger.Println("nil")
			continue
		}
		logger.Printf("%s\n", astPrinter.print(stmt))
	}

	return []Token{}, []lsp.Diagnostic{}
}

func Error(token Token, message string) {

}
