package mexcwsmarket

import (
	"context"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/websocket/dto"
)

const (
	DiffBooksDepthRequestPattern = "spot@public.arrg.depth.v3.api@%s@%d"
)

func (s *Service) OrderBookDiffSubscribe(ctx context.Context, symbols []string, level BookDepth, callback func(api *dto.PublicLimitDepthsV3Api)) error {
	lstnr := func(message *dto.PushDataV3ApiWrapper) {
		switch msg := message.Body.(type) {
		case *dto.PushDataV3ApiWrapper_PublicLimitDepths:
			callback(msg.PublicLimitDepths)
		default:
			fmt.Println("OrderBook callback unknown type:", message.Body)
		}
	}

	for _, symbol := range symbols {
		channel := fmt.Sprintf(DiffBooksDepthRequestPattern, symbol, level)
		if err := s.client.Subscribe(ctx, channel, nil, lstnr); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) OrderBookDiffUnsubscribe(symbols []string, level BookDepth) error {
	for _, symbol := range symbols {
		channel := fmt.Sprintf(DiffBooksDepthRequestPattern, symbol, level)
		if err := s.client.Unsubscribe(channel); err != nil {
			return err
		}
	}

	return nil
}
