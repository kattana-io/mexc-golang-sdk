package mexchttpmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	DefaultOrderBookDepth = 100
	MaxOrderBookDepth     = 5000
)

// OrderBook https://mexcdevelop.github.io/apidocs/spot_v3_en/#order-book
func (s *Service) OrderBook(ctx context.Context, symbol string, limit int32) (*OrderBookResponse, error) {
	endpoint := "/api/v3/depth"

	if limit <= 0 || limit > MaxOrderBookDepth {
		limit = DefaultOrderBookDepth
	}

	params := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprintf("%d", limit),
	}

	res, err := s.client.SendRequest(ctx, http.MethodGet, endpoint, params)
	if err != nil {
		return nil, err
	}

	var info OrderBookResponse
	err = json.Unmarshal(res, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

type OrderBookResponse struct {
	LastUpdateID int64       `json:"last_update_id"`
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
}
