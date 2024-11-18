package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	Id     int    `json:"id"`
	Method string `json:"method"`

	//params ...
}

type Response struct {
	RPC string `json:"jsonrpc"`
	Id  *int   `json:"id"`

	//reults
	//error
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
