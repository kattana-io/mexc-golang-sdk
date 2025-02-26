package mexcwstypes

type OnReceive func(message string)
type OnError func(connClosed bool, err error)

type WsReq struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
}
