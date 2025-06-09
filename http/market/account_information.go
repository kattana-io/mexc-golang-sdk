package mexchttpmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/consts"
	"net/http"
)

type AccountInformationRequest struct {
	RecvWindow *int64 // optional
}

type AccountInformationResponse struct {
	CanTrade    bool      `json:"canTrade"`
	CanWithdraw bool      `json:"canWithdraw"`
	CanDeposit  bool      `json:"canDeposit"`
	UpdateTime  int64     `json:"updateTime"`
	AccountType string    `json:"accountType"`
	Balances    []Balance `json:"balances"`
	Permissions []string  `json:"permissions"`
}

type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

func (s *Service) GetAccountInformation(ctx context.Context, req AccountInformationRequest) (*AccountInformationResponse, error) {
	// https://mexcdevelop.github.io/apidocs/spot_v3_en/#account-information

	// build required params
	params := map[string]string{
		"timestamp": s.getTimestamp(),
	}

	// optional:recvWindow
	if req.RecvWindow != nil {
		params["recvWindow"] = fmt.Sprintf("%d", *req.RecvWindow)
	}

	body, err := s.client.SendRequest(ctx, http.MethodGet, consts.EndpointAccountInformation, params)
	if err != nil {
		return nil, fmt.Errorf("account information failed: %w", err)
	}

	var resp AccountInformationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse account information response: %w", err)
	}

	return &resp, nil
}
