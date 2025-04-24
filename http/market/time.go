package mexchttpmarket

import (
	"context"
	"encoding/json"
	"github.com/kattana-io/mexc-golang-sdk/consts"
	"net/http"
)

// Time https://mexcdevelop.github.io/apidocs/spot_v3_en/#check-server-time
func (s *Service) Time(ctx context.Context) (*TimeResponse, error) {
	res, err := s.client.SendRequest(ctx, http.MethodGet, consts.EndpointTime, nil)
	if err != nil {
		return nil, err
	}

	var timeResponse TimeResponse
	err = json.Unmarshal(res, &timeResponse)
	if err != nil {
		return nil, err
	}

	return &timeResponse, nil
}

type TimeResponse struct {
	ServerTime int64 `json:"serverTime"`
}
