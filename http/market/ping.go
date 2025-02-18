package mexchttpmarket

import (
	"context"
	"net/http"
)

// Ping https://mexcdevelop.github.io/apidocs/spot_v3_en/#test-connectivity
func (s *Service) Ping(ctx context.Context) (string, error) {
	endpoint := "/api/v3/ping"

	res, err := s.client.SendRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
