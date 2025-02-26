package mexchttpstream

import (
	"context"
	"encoding/json"
	"fmt"
	mexchttp "github.com/kattana-io/mexc-golang-sdk/http"
	"net/http"
	"strconv"
	"time"
)

const (
	ListenKeyInterval = time.Minute * 30 // key alive for 60 minutes. Documentation recommends to send a ping about every 30 minutes
)

// Service implements user data stream apis https://mexcdevelop.github.io/apidocs/spot_v3_en/#websocket-user-data-streams
type Service struct {
	client *mexchttp.Client
}

func New(client *mexchttp.Client) *Service {
	return &Service{client: client}
}

// CreateListenKey create new key for user data stream
// https://mexcdevelop.github.io/apidocs/spot_v3_en/#listen-key
func (s *Service) CreateListenKey(ctx context.Context) (string, error) {
	endpoint := "/api/v3/userDataStream"

	res, err := s.client.SendRequest(ctx, http.MethodPost, endpoint, map[string]string{
		"timestamp": strconv.FormatInt(time.Now().UnixMilli(), 10),
	})
	if err != nil {
		return "", err
	}

	var response ListenKeyResponse
	err = json.Unmarshal(res, &response)
	if err != nil {
		return "", err
	}

	return response.Key, nil
}

// KeepAliveKey extends its validity for 60 minutes
// https://mexcdevelop.github.io/apidocs/spot_v3_en/#listen-key
func (s *Service) KeepAliveKey(ctx context.Context, key string) error {
	endpoint := "/api/v3/userDataStream"

	params := map[string]string{
		"listenKey": key,
		"timestamp": strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
	_, err := s.client.SendRequest(ctx, http.MethodPut, endpoint, params)
	return err
}

// RunKeyKeepAlive send PUT request every 30 min to avoid key invalidating
func (s *Service) RunKeyKeepAlive(ctx context.Context, key string) error {
	for {
		select {
		case <-time.After(ListenKeyInterval):
			if err := s.KeepAliveKey(ctx, key); err != nil {
				return fmt.Errorf("failed to keep alive key: %s", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
