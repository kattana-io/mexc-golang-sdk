package mexchttpmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/consts"
	"net/http"
)

type InternalTransferRequest struct {
	ToAccount     string  // required, the recipient’s identifier (email, uid or mobile)
	ToAccountType string  // required, EMAIL / UID / MOBILE
	AreaCode      *string // optional, only if ToAccountType == "MOBILE"
	Asset         string  // required, e.g. "BNB", "USDT", "BTC"
	Amount        string  // required, e.g. "0.0097"
	RecvWindow    *int64  // optional
}

type InternalTransferResponse struct {
	TranId string `json:"tranId"`
}

func (s *Service) NewInternalTransfer(ctx context.Context, req InternalTransferRequest) (*InternalTransferResponse, error) {
	// https://mexcdevelop.github.io/apidocs/spot_v3_en/#internal-transfer

	// build required params
	params := map[string]string{
		"toAccountType": req.ToAccountType,
		"toAccount":     req.ToAccount,
		"asset":         req.Asset,
		"amount":        req.Amount,
		"timestamp":     s.getTimestamp(),
	}

	// optional: areaCode, recvWindow
	if req.AreaCode != nil {
		params["areaCode"] = *req.AreaCode
	}
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", *req.RecvWindow)
	}

	body, err := s.client.SendRequest(ctx, http.MethodPost, consts.EndpointInternalTransfer, params)
	if err != nil {
		return nil, fmt.Errorf("internal transfer failed: %w", err)
	}

	var resp InternalTransferResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse internal transfer response: %w", err)
	}

	return &resp, nil
}

type InternalTransferHistoryRequest struct {
	StartTime  *int64  // optional, milliseconds
	EndTime    *int64  // optional, milliseconds
	Page       *int    // optional, defaults to 1
	Limit      *int    // optional, defaults to 10
	TranId     *string // optional, specific tranId to query
	RecvWindow *int64  // optional
}

type InternalTransferRecord struct {
	TranId          string `json:"tranId"`
	Asset           string `json:"asset"`
	Amount          string `json:"amount"`
	FromAccountType string `json:"fromAccountType"`
	ToAccountType   string `json:"toAccountType"`
	FromAccount     string `json:"fromAccount"`
	ToAccount       string `json:"toAccount"`
	Status          string `json:"status"`
	Timestamp       int64  `json:"timestamp"`
}

type InternalTransferHistoryResponse struct {
	Page         int                      `json:"page"`
	TotalRecords int                      `json:"totalRecords"`
	TotalPageNum int                      `json:"totalPageNum"`
	Data         []InternalTransferRecord `json:"data"`
}

func (s *Service) GetInternalTransferHistory(ctx context.Context, req InternalTransferHistoryRequest) (*InternalTransferHistoryResponse, error) {
	// https://mexcdevelop.github.io/apidocs/spot_v3_en/#query-internal-transfer-history

	params := map[string]string{
		"timestamp": s.getTimestamp(),
	}

	if req.StartTime != nil {
		params["startTime"] = fmt.Sprintf("%d", *req.StartTime)
	}
	if req.EndTime != nil {
		params["endTime"] = fmt.Sprintf("%d", *req.EndTime)
	}
	if req.Page != nil {
		params["page"] = fmt.Sprintf("%d", *req.Page)
	}
	if req.Limit != nil {
		params["limit"] = fmt.Sprintf("%d", *req.Limit)
	}
	if req.TranId != nil {
		params["tranId"] = *req.TranId
	}
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", *req.RecvWindow)
	}

	body, err := s.client.SendRequest(ctx, http.MethodGet, consts.EndpointInternalTransfer, params)
	if err != nil {
		return nil, fmt.Errorf("failed to query internal transfer history: %w", err)
	}

	var resp InternalTransferHistoryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse internal transfer history: %w", err)
	}

	return &resp, nil
}
