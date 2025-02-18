package mexchttpmarket

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// ExchangeInfo https://mexcdevelop.github.io/apidocs/spot_v3_en/#exchange-information
func (s *Service) ExchangeInfo(ctx context.Context, symbols []string) (*ExchangeInfo, error) {
	endpoint := "/api/v3/exchangeInfo"

	params := make(map[string]string)
	if len(symbols) > 0 {
		params["symbols"] = strings.Join(symbols, ",")
	}

	res, err := s.client.SendRequest(ctx, http.MethodGet, endpoint, params)
	if err != nil {
		return nil, err
	}

	var info ExchangeInfo
	err = json.Unmarshal(res, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

type ExchangeInfo struct {
	Timezone        string           `json:"timezone"`
	ServerTime      int64            `json:"serverTime"`
	RateLimits      []RateLimit      `json:"rateLimits"`
	ExchangeFilters []ExchangeFilter `json:"exchangeFilters"`
	Symbols         []Symbol         `json:"symbols"`
}

type Symbol struct {
	Symbol                     string   `json:"symbol"`
	Status                     string   `json:"status"`
	BaseAsset                  string   `json:"baseAsset"`
	BaseAssetPrecision         int      `json:"baseAssetPrecision"`
	QuoteAsset                 string   `json:"quoteAsset"`
	QuotePrecision             int      `json:"quotePrecision"`
	QuoteAssetPrecision        int      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
	OrderTypes                 []string `json:"orderTypes"`
	IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
	QuoteAmountPrecision       string   `json:"quoteAmountPrecision"`
	BaseSizePrecision          string   `json:"baseSizePrecision"`
	Permissions                []string `json:"permissions"`
	Filters                    []Filter `json:"filters"`
	MaxQuoteAmount             string   `json:"maxQuoteAmount"`
	MakerCommission            string   `json:"makerCommission"`
	TakerCommission            string   `json:"takerCommission"`
	QuoteAmountPrecisionMarket string   `json:"quoteAmountPrecisionMarket"`
	MaxQuoteAmountMarket       string   `json:"maxQuoteAmountMarket"`
	FullName                   string   `json:"fullName"`
	TradeSideType              int      `json:"tradeSideType"`
}

type Filter struct {
	// Add fields as needed
}

type RateLimit struct {
	// Add fields as needed
}

type ExchangeFilter struct {
	// Add fields as needed
}
