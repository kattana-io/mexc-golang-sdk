package mexcwsmarket

import (
	"context"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/websocket/dto"
)

const (
	MinInterval                  = "10ms"
	MaxInterval                  = "100ms"
	DiffBooksDepthRequestPattern = "spot@public.aggre.depth.v3.api.pb@%s@%s"
)

func (s *Service) OrderBookDiffSubscribe(ctx context.Context, symbols []string, interval string, callback func(api *dto.PublicLimitDepthsV3Api)) error {
	lstnr := func(message *dto.PushDataV3ApiWrapper) {
		switch msg := message.Body.(type) {
		case *dto.PushDataV3ApiWrapper_PublicLimitDepths:
			callback(msg.PublicLimitDepths)
		default:
			fmt.Println("OrderBook callback unknown type:", message.Body)
		}
	}

	for _, symbol := range symbols {
		channel := fmt.Sprintf(DiffBooksDepthRequestPattern, interval, symbol)
		if err := s.client.Subscribe(ctx, channel, nil, lstnr); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) OrderBookDiffUnsubscribe(symbols []string, interval string) error {
	for _, symbol := range symbols {
		channel := fmt.Sprintf(DiffBooksDepthRequestPattern, interval, symbol)
		if err := s.client.Unsubscribe(channel); err != nil {
			return err
		}
	}

	return nil
}
