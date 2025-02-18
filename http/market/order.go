package mexchttpmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
)

// CreateOrder https://mexcdevelop.github.io/apidocs/spot_v3_en/#new-order
func (s *Service) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	endpoint := "/api/v3/order"

	params := make(map[string]string)

	params["symbol"] = req.Symbol
	params["side"] = string(req.Side)
	params["type"] = string(req.Type)
	params["timestamp"] = s.getTimestamp()

	if req.Quantity != nil {
		params["quantity"] = *req.Quantity
	}
	if req.QuoteOrderQty != nil {
		params["quoteOrderQty"] = *req.QuoteOrderQty
	}
	if req.Price != nil {
		params["price"] = *req.Price
	}
	if req.NewClientOrderId != nil {
		params["newClientOrderId"] = *req.NewClientOrderId
	}
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", *req.RecvWindow)
	}

	res, err := s.client.SendRequest(ctx, http.MethodPost, endpoint, params)
	if err != nil {
		return nil, err
	}

	var orderResponse CreateOrderResponse
	err = json.Unmarshal(res, &orderResponse)
	if err != nil {
		return nil, err
	}

	return &orderResponse, nil
}

type CreateOrderRequest struct {
	Symbol           string  `json:"symbol"`
	Side             Side    `json:"side"`
	Type             Type    `json:"type"`
	Quantity         *string `json:"quantity,omitempty"`
	QuoteOrderQty    *string `json:"quoteOrderQty,omitempty"`
	Price            *string `json:"price,omitempty"`
	NewClientOrderId *string `json:"newClientOrderId,omitempty"`
	RecvWindow       *int64  `json:"recvWindow,omitempty"`
}

type CreateOrderResponse struct {
	Symbol       string          `json:"symbol"`
	OrderId      string          `json:"orderId"`
	OrderListId  int             `json:"orderListId"`
	Price        decimal.Decimal `json:"price"`
	OrigQty      decimal.Decimal `json:"origQty"`
	Type         Type            `json:"type"`
	Side         Side            `json:"side"`
	TransactTime int64           `json:"transactTime"`
}
