package mexchttpmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/consts"
	"github.com/shopspring/decimal"
	"net/http"
)

// GetAccountTradeList https://www.mexc.com/api-docs/spot-v3/spot-account-trade#account-trade-list
func (s *Service) GetAccountTradeList(ctx context.Context, req *GetAccountTradeListRequest) ([]*GetAccountTradeListResponse, error) {
	params := make(map[string]string)

	params["symbol"] = req.Symbol
	params["timestamp"] = s.getTimestamp()

	if req.OrderID != nil {
		params["orderId"] = *req.OrderID
	}
	if req.StartTime != nil {
		params["startTime"] = fmt.Sprintf("%d", *req.RecvWindow)
	}
	if req.EndTime != nil {
		params["endTime"] = fmt.Sprintf("%d", *req.EndTime)
	}
	if req.Limit != nil {
		params["limit"] = fmt.Sprintf("%d", *req.Limit)
	}
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", *req.RecvWindow)
	}

	res, err := s.client.SendRequest(ctx, http.MethodGet, consts.EndpointAccountTradeList, params)
	if err != nil {
		return nil, err
	}

	var response []*GetAccountTradeListResponse
	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type GetAccountTradeListRequest struct {
	Symbol     string  `json:"symbol"`
	OrderID    *string `json:"orderId,omitempty"`
	StartTime  *int64  `json:"startTime,omitempty"`
	EndTime    *int64  `json:"endTime,omitempty"`
	Limit      *int32  `json:"limit,omitempty"`
	RecvWindow *int64  `json:"recvWindow,omitempty"`
}

type GetAccountTradeListResponse struct {
	Symbol          string          `json:"symbol"`
	ID              string          `json:"id"`
	OrderID         string          `json:"orderId"`
	OrderListID     int64           `json:"orderListId"`
	Price           decimal.Decimal `json:"price"`
	Qty             decimal.Decimal `json:"qty"`
	QuoteQty        decimal.Decimal `json:"quoteQty"`
	Commission      decimal.Decimal `json:"commission"`
	CommissionAsset string          `json:"commissionAsset"`
	Time            int64           `json:"time"`
	IsBuyer         bool            `json:"isBuyer"`
	IsMaker         bool            `json:"isMaker"`
	IsBestMatch     bool            `json:"isBestMatch"`
	IsSelfTrade     bool            `json:"isSelfTrade"`
	ClientOrderID   *string         `json:"clientOrderId"`
}
