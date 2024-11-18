package lsp

type ErrorResponse struct {
	Response
	Error Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data.omitempty"`
}
