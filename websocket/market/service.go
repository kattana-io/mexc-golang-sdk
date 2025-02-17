package mexcwsmarket

import (
	mexcws "github.com/kattana-io/mexc-golang-sdk/websocket"
)

type Service struct {
	client *mexcws.MEXCWebSocket
}

func New(client *mexcws.MEXCWebSocket) *Service {
	return &Service{
		client: client,
	}
}
