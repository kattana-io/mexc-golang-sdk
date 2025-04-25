package mexchttpmarket

import (
	"context"
	"encoding/json"
	"github.com/kattana-io/mexc-golang-sdk/consts"
	"github.com/shopspring/decimal"
	"net/http"
)

// CurrencyInformation https://mexcdevelop.github.io/apidocs/spot_v3_en/#query-the-currency-information
func (s *Service) CurrencyInformation(ctx context.Context) ([]*CoinInfoResponse, error) {
	params := make(map[string]string)
	params["timestamp"] = s.getTimestamp()

	res, err := s.client.SendRequest(ctx, http.MethodGet, consts.EndpointGetCurrencyInformation, params)
	if err != nil {
		return nil, err
	}

	info := make([]*CoinInfoResponse, 0)
	err = json.Unmarshal(res, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

type CoinInfoResponse struct {
	Coin        string             `json:"coin"`
	Name        string             `json:"Name"`
	NetworkList []CoinWithdrawInfo `json:"networkList"`
}

type CoinWithdrawInfo struct {
	Coin                    string          `json:"coin"`
	DepositDesc             string          `json:"depositDesc"`
	DepositEnable           bool            `json:"depositEnable"`
	MinConfirm              int             `json:"minConfirm"`
	Name                    string          `json:"Name"`
	Network                 string          `json:"network"`
	WithdrawEnable          bool            `json:"withdrawEnable"`
	WithdrawFee             decimal.Decimal `json:"withdrawFee"`
	WithdrawIntegerMultiple string          `json:"withdrawIntegerMultiple"`
	WithdrawMax             decimal.Decimal `json:"withdrawMax"`
	WithdrawMin             decimal.Decimal `json:"withdrawMin"`
	SameAddress             bool            `json:"sameAddress"`
	Contract                string          `json:"contract"`
	WithdrawTips            string          `json:"withdrawTips"`
	DepositTips             string          `json:"depositTips"`
	NetWork                 string          `json:"netWork"`
}
