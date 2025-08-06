package mexcwsmarket

import (
	"context"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/websocket/dto"
)

type BookDepth int

const (
	MinBookDepth BookDepth = 5
	MidBookDepth BookDepth = 10
	MaxBookDepth BookDepth = 20

	PartialBooksDepthRequestPattern = "spot@public.limit.depth.v3.api.pb@%s@%d"
)

func (s *Service) OrderBookSubscribe(ctx context.Context, symbols []string, level BookDepth, callback func(api *dto.PublicLimitDepthsV3Api)) error {
	lstnr := func(message *dto.PushDataV3ApiWrapper) {
		switch msg := message.Body.(type) {
		case *dto.PushDataV3ApiWrapper_PublicLimitDepths:
			callback(msg.PublicLimitDepths)
		default:
			fmt.Println("OrderBook callback unknown type:", message.Body)
		}
	}

	for _, symbol := range symbols {
		channel := fmt.Sprintf(PartialBooksDepthRequestPattern, symbol, level)
		if err := s.client.Subscribe(ctx, channel, nil, lstnr); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) OrderBookUnsubscribe(symbols []string, level BookDepth) error {
	for _, symbol := range symbols {
		channel := fmt.Sprintf(PartialBooksDepthRequestPattern, symbol, level)
		if err := s.client.Unsubscribe(channel); err != nil {
			return err
		}
	}

	return nil
}
