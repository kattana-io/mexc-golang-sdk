package mexcwsmarket

import (
	"context"
	"encoding/json"
	"fmt"
)

type BookDepth int

const (
	MinBookDepth BookDepth = 5
	MidBookDepth BookDepth = 10
	MaxBookDepth BookDepth = 20

	PartialBooksDepthRequestPattern = "spot@public.limit.depth.v3.api@%s@%d"
)

type OrderBook struct {
	Channel string `json:"c"`
	Data    struct {
		Bids []struct {
			Price  string `json:"p"`
			Volume string `json:"v"`
		} `json:"bids"`
		Asks []struct {
			Price  string `json:"p"`
			Volume string `json:"v"`
		} `json:"asks"`
		Event     string `json:"e"`
		RequestID string `json:"r"`
	} `json:"d"`
	Symbol    string `json:"s"`
	Timestamp int64  `json:"t"`
}

func (s *Service) OrderBookSubscribe(ctx context.Context, symbols []string, level BookDepth, callback func(*OrderBook)) error {
	lstnr := func(message string) {
		var book OrderBook

		err := json.Unmarshal([]byte(message), &book)
		if err != nil {
			fmt.Println("OrderBook callback unmarshal error:", err)
			return
		}

		callback(&book)
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
