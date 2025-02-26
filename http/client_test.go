package mexchttp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Mock HTTP server for testing
func setupMockServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/v3/account", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"makerCommission":15,"takerCommission":15,"buyerCommission":0,"sellerCommission":0}`))
	})

	return httptest.NewServer(handler)
}

func TestMEXCClient_GenerateSignature(t *testing.T) {
	client := NewClient("test_api_key", "test_secret_key", &http.Client{})

	params := "timestamp=1617785000000&recvWindow=5000"
	expectedSignature := "3a4ebaba0e92657d7824e6e8d1ffcad39bd3f387c5d73da2768d8bf6d9baa062"
	signature := client.generateSignature(params)

	if signature != expectedSignature {
		t.Errorf("Expected signature %s, but got %s", expectedSignature, signature)
	}
}

func TestMEXCClient_NewRequest(t *testing.T) {
	client := NewClient("test_api_key", "test_secret_key", &http.Client{})
	ctx := context.Background()

	params := map[string]string{
		"recvWindow": "5000",
		"timestamp":  "1617785000000",
	}

	req, err := client.newRequest(ctx, "GET", "/api/v3/account", params)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	//nolint:lll
	expectedURL := "https://api.mexc.com/api/v3/account?recvWindow=5000&signature=e1009356af9025283af1fa79c5000c2a9350335e9a1818460aafb7bc3980abac&timestamp=1617785000000"
	if req.URL.String() != expectedURL {
		t.Errorf("Expected URL %s, but got %s", expectedURL, req.URL.String())
	}

	if req.Header.Get("X-MEXC-APIKEY") != "test_api_key" {
		t.Errorf("Expected X-MEXC-APIKEY to be 'test_api_key', but got %s", req.Header.Get("X-MEXC-APIKEY"))
	}
}

func TestMEXCClient_SendRequest(t *testing.T) {
	server := setupMockServer()
	defer server.Close()

	client := NewClient("test_api_key", "test_secret_key", &http.Client{})
	client.baseURL = server.URL // перенаправляем клиент на мок сервер

	ctx := context.Background()

	params := map[string]string{
		"recvWindow": "5000",
		"timestamp":  fmt.Sprintf("%d", time.Now().UnixMilli()),
	}

	response, err := client.SendRequest(ctx, "GET", "/api/v3/account", params)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	expectedResponse := `{"makerCommission":15,"takerCommission":15,"buyerCommission":0,"sellerCommission":0}`
	if strings.TrimSpace(string(response)) != expectedResponse {
		t.Errorf("Expected response %s, but got %s", expectedResponse, response)
	}
}
