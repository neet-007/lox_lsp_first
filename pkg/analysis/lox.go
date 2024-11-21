package analysis

import (
	"log"
	"os"

	"github.com/neet-007/lox_lsp_first/internal/lsp"
	"github.com/neet-007/lox_lsp_first/pkg/rpc"
)

type Analyser struct {
	hadError bool
	uri      string
}

func NewAnaylser() *Analyser {
	return &Analyser{
		hadError: true,
		uri:      "",
	}
}

func (analyser *Analyser) Analyse(source []byte, uri string, logger *log.Logger) {
	analyser.uri = uri
	scanner := NewScanner(source, analyser)

	tokens := scanner.Scan()

	parser := NewParser(tokens, analyser)

	statements := parser.Parse()

	resolver := NewResolver(analyser)

	resolver.Resolve(statements)

	interpreter := NewInterpreter(resolver.locals, analyser)
	interpreter.Interpert(statements)
}

func (analyser *Analyser) Error(token Token, message string) {
	analyser.hadError = true
	lexeme := ""
	if token.Lexeme != "@" {
		lexeme = token.Lexeme
	}

	diagnostic := lsp.NewDiagnostic(
		lsp.Range{
			Start: lsp.Position{
				Line:      token.StartLine,
				Character: token.StartChar,
			},
			End: lsp.Position{
				Line:      token.StartLine,
				Character: token.EndChar,
			},
		},
		1,
		lexeme,
		message,
	)

	reply := rpc.EncodeMessage(lsp.PublishDiagnosticsNotification{
		Notification: lsp.Notification{
			RPC:    "2.0",
			Method: "textDocument/publishDiagnostics",
		},
		Params: lsp.PublishDiagnosticsParams{
			URI:         analyser.uri,
			Diagnostics: []lsp.Diagnostic{diagnostic},
		},
	})

	if _, err := os.Stdout.Write([]byte(reply)); err != nil {

	}
}
