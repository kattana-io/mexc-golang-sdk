package mexcwsuser

import (
	mexchttpstream "github.com/kattana-io/mexc-golang-sdk/http/stream"
	"github.com/kattana-io/mexc-golang-sdk/websocket"
)

type Service struct {
	wsClient   *mexcws.MEXCWebSocket
	httpStream *mexchttpstream.Service
}

func New(wsClient *mexcws.MEXCWebSocket,
	httpStream *mexchttpstream.Service) *Service {
	return &Service{
		wsClient:   wsClient,
		httpStream: httpStream,
	}
}
