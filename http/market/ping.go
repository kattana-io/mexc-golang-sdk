package mexchttpmarket

import (
	"context"
	"github.com/kattana-io/mexc-golang-sdk/consts"
	"net/http"
)

// Ping https://mexcdevelop.github.io/apidocs/spot_v3_en/#test-connectivity
func (s *Service) Ping(ctx context.Context) (string, error) {
	res, err := s.client.SendRequest(ctx, http.MethodGet, consts.EndpointPing, nil)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
