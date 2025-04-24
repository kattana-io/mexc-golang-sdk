package mexchttpmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/consts"
	"net/http"
)

type WithdrawRequest struct {
	Coin            string  // required
	WithdrawOrderId *string // optional
	Network         *string // optional
	ContractAddress *string // optional
	Address         string  // required
	Memo            *string // optional
	Amount          string  // required
	Remark          *string // optional
	RecvWindow      *int64  // optional
}

type WithdrawResponse struct {
	Id string `json:"id"`
}

func (s *Service) Withdraw(ctx context.Context, req WithdrawRequest) (*WithdrawResponse, error) {
	params := map[string]string{
		"coin":      req.Coin,
		"address":   req.Address,
		"amount":    req.Amount,
		"timestamp": s.getTimestamp(),
	}

	if req.Network != nil {
		params["netWork"] = *req.Network
	}
	if req.Memo != nil {
		params["memo"] = *req.Memo
	}
	if req.ContractAddress != nil {
		params["contractAddress"] = *req.ContractAddress
	}
	if req.WithdrawOrderId != nil {
		params["withdrawOrderId"] = *req.WithdrawOrderId
	}
	if req.Remark != nil {
		params["remark"] = *req.Remark
	}
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", *req.RecvWindow)
	}

	body, err := s.client.SendRequest(ctx, http.MethodPost, consts.EndpointWithdraw, params)
	if err != nil {
		return nil, err
	}

	var resp WithdrawResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse withdraw response: %w", err)
	}

	return &resp, nil
}

type WithdrawHistoryRequest struct {
	Coin       *string // optional
	Status     *string // optional
	Limit      *string // optional
	StartTime  *string // optional
	EndTime    *string // optional
	RecvWindow *int64  // optional
}

type WithdrawRecord struct {
	ID              string         `json:"id"`
	Coin            string         `json:"coin"`
	Amount          string         `json:"amount"`
	Address         string         `json:"address"`
	Network         string         `json:"network"`
	Status          WithdrawStatus `json:"status"`
	TransactionFee  string         `json:"transactionFee"`
	TransHash       string         `json:"transHash"`
	TxID            *string        `json:"txId"`
	Memo            string         `json:"memo"`
	Remark          string         `json:"remark"`
	ApplyTime       int64          `json:"applyTime"`
	TransferType    TransferType   `json:"transferType"`
	WithdrawOrderID *string        `json:"withdrawOrderId,omitempty"`
	ConfirmNo       *string        `json:"confirmNo,omitempty"`
	CoinID          string         `json:"coinId"`
	VcoinID         string         `json:"vcoinId"`
}

func (s *Service) GetWithdrawsHistory(ctx context.Context, req WithdrawHistoryRequest) ([]WithdrawRecord, error) {
	// https://mexcdevelop.github.io/apidocs/spot_v3_en/#withdraw-history-supporting-network
	params := map[string]string{
		"timestamp": s.getTimestamp(),
	}

	if req.Coin != nil {
		params["coin"] = *req.Coin
	}
	if req.Status != nil {
		params["status"] = *req.Status
	}
	if req.StartTime != nil {
		params["startTime"] = *req.StartTime
	}
	if req.EndTime != nil {
		params["endTime"] = *req.EndTime
	}
	if req.Limit != nil {
		params["limit"] = *req.Limit
	}
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", *req.RecvWindow)
	}

	body, err := s.client.SendRequest(ctx, http.MethodGet, consts.EndpointWithdrawHistory, params)
	if err != nil {
		return nil, err
	}

	var withdraws []WithdrawRecord
	if err := json.Unmarshal(body, &withdraws); err != nil {
		return nil, fmt.Errorf("failed to parse withdraw history: %w", err)
	}

	return withdraws, nil
}

func (s *Service) GetWithdrawHistoryByID(ctx context.Context, withdrawId string) (*WithdrawRecord, error) {
	req := WithdrawHistoryRequest{}
	withdraws, err := s.GetWithdrawsHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, w := range withdraws {
		if w.ID == withdrawId {
			return &w, nil
		}
	}

	return nil, fmt.Errorf("withdraw record with id %s not found", withdrawId)
}
