package mexcwsmarket

import (
	"context"
	mexcws "github.com/kattana-io/mexc-golang-sdk/websocket"
)

func (s *Service) Ping(ctx context.Context) error {
	req := &mexcws.WsReq{
		Method: "PING",
	}

	return s.client.Send(ctx, req)
}
