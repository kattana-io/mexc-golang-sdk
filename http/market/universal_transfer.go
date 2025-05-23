package mexchttpmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/consts"
	"net/http"
)

type UniversalTransferRequest struct {
	FromAccount     *string // optional, the sender’s identifier (email, UID or mobile)
	ToAccount       *string // optional, the recipient’s identifier (email, UID or mobile)
	FromAccountType string  // required, e.g. "SPOT" / "FUTURES"
	ToAccountType   string  // required, e.g. "SPOT" / "FUTURES"
	Asset           string  // required, e.g. “BNB”, “USDT”, “BTC” — must be a supported token
	Amount          string  // required, decimal string, e.g. “0.002” — ≤ available balance, respects on-chain precision
	RecvWindow      *int64  // optional, request validity window in ms (default 5000)
}

type TransferResponse struct {
	TranId string `json:"tranId"`
}

func (s *Service) NewUniversalTransfer(ctx context.Context, req UniversalTransferRequest) (*TransferResponse, error) {
	// https://mexcdevelop.github.io/apidocs/spot_v3_en/#universal-transfer-for-master-account
	if req.FromAccount == nil && req.ToAccount == nil {
		return nil, fmt.Errorf("fromAccount and toAccount both are empty, at least one of them must be specified")
	}

	params := map[string]string{
		"asset":           req.Asset,
		"amount":          req.Amount,
		"fromAccountType": req.FromAccountType,
		"toAccountType":   req.ToAccountType,
		"timestamp":       s.getTimestamp(),
	}

	if req.FromAccount != nil {
		params["fromAccount"] = *req.FromAccount
	}
	if req.ToAccount != nil {
		params["toAccount"] = *req.ToAccount
	}
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}

	body, err := s.client.SendRequest(ctx, http.MethodPost, consts.EndpointUniversalTransfer, params)
	if err != nil {
		return nil, err
	}

	var resp TransferResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

type TransferHistoryRequest struct {
	FromAccount     *string // optional
	ToAccount       *string // optional
	FromAccountType string  // required
	ToAccountType   string  // required
	StartTime       *string // optional
	EndTime         *string // optional
	Page            *string // optional
	Limit           *string // optional
	RecvWindow      *int64  // optional
}

type TransferRecord struct {
	TranId          string `json:"tranId"`
	FromAccount     string `json:"fromAccount"`
	ToAccount       string `json:"toAccount"`
	ClientTranId    string `json:"clientTranId"`
	Asset           string `json:"asset"`
	Amount          string `json:"amount"`
	FromAccountType string `json:"fromAccountType"`
	ToAccountType   string `json:"toAccountType"`
	FromSymbol      string `json:"fromSymbol"`
	ToSymbol        string `json:"toSymbol"`
	Status          string `json:"status"`
	Timestamp       int64  `json:"timestamp"`
}

type TransferHistoryResponse struct {
	Transfers []TransferRecord `json:"transfers"`
}

func (s *Service) GetUniversalTransferHistory(ctx context.Context, req TransferHistoryRequest) (*TransferHistoryResponse, error) {
	// https://mexcdevelop.github.io/apidocs/spot_v3_en/#query-universal-transfer-history-for-master-account
	params := map[string]string{
		"fromAccountType": req.FromAccountType,
		"toAccountType":   req.ToAccountType,
		"timestamp":       s.getTimestamp(),
	}

	if req.FromAccount != nil {
		params["fromAccount"] = fmt.Sprintf("%s", *req.FromAccount)
	}
	if req.ToAccount != nil {
		params["toAccount"] = fmt.Sprintf("%s", *req.ToAccount)
	}
	if req.StartTime != nil {
		params["startTime"] = fmt.Sprintf("%s", *req.StartTime)
	}
	if req.EndTime != nil {
		params["endTime"] = fmt.Sprintf("%s", *req.EndTime)
	}
	if req.Page != nil {
		params["page"] = fmt.Sprintf("%s", *req.Page)
	}
	if req.Limit != nil {
		params["limit"] = fmt.Sprintf("%s", *req.Limit)
	}
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", *req.RecvWindow)
	}

	body, err := s.client.SendRequest(ctx, http.MethodGet, consts.EndpointUniversalTransfer, params)
	if err != nil {
		return nil, err
	}

	var resp TransferHistoryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse transfer history: %w", err)
	}

	return &resp, nil
}
