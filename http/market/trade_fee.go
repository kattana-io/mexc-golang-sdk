package mexchttpmarket

import (
	"context"
	"encoding/json"
	"github.com/shopspring/decimal"
	"net/http"
)

// TradeFee https://mexcdevelop.github.io/apidocs/spot_v3_en/#query-symbol-commission
func (s *Service) TradeFee(ctx context.Context, symbol string) (*TradeFeeResponse, error) {
	endpoint := "api/v3/tradeFee"

	params := map[string]string{
		"symbol": symbol,
	}

	res, err := s.client.SendRequest(ctx, http.MethodGet, endpoint, params)
	if err != nil {
		return nil, err
	}

	var info TradeFeeResponse
	err = json.Unmarshal(res, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

type TradeFeeResponse struct {
	Data      TradeFeeData `json:"data"`
	Code      int32        `json:"code"`
	Message   string       `json:"message"`
	Timestamp int64        `json:"timestamp"`
}

type TradeFeeData struct {
	TakerCommission decimal.Decimal `json:"taker_commission"`
	MakerCommission decimal.Decimal `json:"maker_commission"`
}
