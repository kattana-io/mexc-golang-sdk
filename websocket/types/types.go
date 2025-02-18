package mexcwstypes

type OnReceive func(message string)
type OnError func(err error)

type WsReq struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
}
