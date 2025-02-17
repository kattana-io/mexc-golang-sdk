package mexchttpmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// QueryOrder https://mexcdevelop.github.io/apidocs/spot_v3_en/#query-order
func (s *Service) QueryOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error) {
	endpoint := "/api/v3/order"

	params := make(map[string]string)

	params["symbol"] = req.Symbol
	params["timestamp"] = s.getTimestamp()

	if req.OrderID != nil {
		params["orderId"] = *req.OrderID
	}
	if req.OrigClientOrderId != nil {
		params["origClientOrderId"] = *req.OrigClientOrderId
	}
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", *req.RecvWindow)
	}

	res, err := s.client.SendRequest(ctx, http.MethodGet, endpoint, params)
	if err != nil {
		return nil, err
	}

	var orderResponse GetOrderResponse
	err = json.Unmarshal(res, &orderResponse)
	if err != nil {
		return nil, err
	}

	return &orderResponse, nil
}

type GetOrderRequest struct {
	Symbol            string  `json:"symbol"`
	OrigClientOrderId *string `json:"origClientOrderId,omitempty"`
	OrderID           *string `json:"orderId,omitempty"`
	RecvWindow        *int64  `json:"recvWindow,omitempty"`
}

type GetOrderResponse struct {
	Symbol              string `json:"symbol"`
	OrderId             string `json:"orderId"`
	OrigClientOrderId   string `json:"origClientOrderId,omitempty"`
	ClientOrderID       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              Status `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                Type   `json:"type"`
	Side                Side   `json:"side"`
	StopPrice           string `json:"stopPrice"`
	CreateTime          int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
}
