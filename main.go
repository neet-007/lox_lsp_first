package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/neet-007/lox_lsp_first/internal/lsp"
	"github.com/neet-007/lox_lsp_first/pkg/analysis"
	"github.com/neet-007/lox_lsp_first/pkg/rpc"
)

func main() {
	logger := getLogger("/home/moayed/personal/lox_lsp_first/logs.txt")
	logger.Println("Starting...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error:%v", err)
		}

		handleMessage(logger, writer, "", method, content)

	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state string, method string, content []byte) {
	logger.Printf("Message with method:%s\n", method)
	switch method {
	case "initialize":
		{
			var request lsp.InitializeRequest

			if err := json.Unmarshal(content, &request); err != nil {
				logger.Printf("Error: %v", err)
				return
			}

			logger.Printf("Connected to: %s %s",
				request.Params.ClientInfo.Version, request.Params.ClientInfo.Name)

			response := lsp.NewInitializeResponse(request.Id)
			writeResponse(writer, response)

			logger.Println("reply sent")
		}
	case "textDocument/didOpen":
		{
			var didOpenTextDocumentNotification lsp.DidOpenTextDocumentNotification
			if err := json.Unmarshal(content, &didOpenTextDocumentNotification); err != nil {
				logger.Printf("Error did open:%s\n", err)
				return
			}
			logger.Printf("text document with uri:%s\n", didOpenTextDocumentNotification.Params.TextDocument.URI)
			_, diagnostics := analysis.Analyse([]byte(didOpenTextDocumentNotification.Params.TextDocument.Text))
			writeResponse(writer, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         didOpenTextDocumentNotification.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}
	case "textDocument/didChange":
		{
			var didChangeTextDocumentNotification lsp.TextDocumentDidChangeNotification
			if err := json.Unmarshal(content, &didChangeTextDocumentNotification); err != nil {
				logger.Printf("textDocument/didChange: %s", err)
				return
			}

			logger.Printf("Changed: %s", didChangeTextDocumentNotification.Params.TextDocument.URI)
			for _, change := range didChangeTextDocumentNotification.Params.ContentChanges {
				_, diagnostics := analysis.Analyse([]byte(change.Text))
				writeResponse(writer, lsp.PublishDiagnosticsNotification{
					Notification: lsp.Notification{
						RPC:    "2.0",
						Method: "textDocument/publishDiagnostics",
					},
					Params: lsp.PublishDiagnosticsParams{
						URI:         didChangeTextDocumentNotification.Params.TextDocument.URI,
						Diagnostics: diagnostics,
					},
				})
			}
		}
	}

}

func writeResponse(writer io.Writer, msg any) error {
	reply := rpc.EncodeMessage(msg)

	if _, err := writer.Write([]byte(reply)); err != nil {
		return err
	}

	return nil
}

func getLogger(filePath string) *log.Logger {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	return log.New(file, "[LOX_LSP] ", log.Ldate|log.Lshortfile)
}
