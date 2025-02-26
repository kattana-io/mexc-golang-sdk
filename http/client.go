package mexchttp

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client представляет клиента для работы с API MEXC.
type Client struct {
	apiKey     string
	secretKey  string
	baseURL    string
	httpClient *http.Client
}

// NewClient создает новый экземпляр клиента для работы с API MEXC.
func NewClient(apiKey, secretKey string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    "https://api.mexc.com", // базовый URL для API MEXC
		apiKey:     apiKey,
		secretKey:  secretKey,
		httpClient: httpClient,
	}
}

// generateSignature генерирует подпись HMAC SHA256.
func (c *Client) generateSignature(query string) string {
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(query))
	return hex.EncodeToString(mac.Sum(nil))
}

// newRequest создает новый HTTP-запрос с контекстом и подписью.
func (c *Client) newRequest(ctx context.Context, method, endpoint string, params map[string]string) (*http.Request, error) {
	// Создание URL с параметрами
	reqURL, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}
	reqURL.RawQuery = query.Encode()

	// Signature generation
	signature := c.generateSignature(reqURL.RawQuery)
	query.Add("signature", signature)

	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	// Установка заголовков авторизации
	req.Header.Set("X-MEXC-APIKEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// SendRequest отправляет запрос к API и возвращает ответ.
func (c *Client) SendRequest(ctx context.Context, method, endpoint string, params map[string]string) ([]byte, error) {
	req, err := c.newRequest(ctx, method, endpoint, params)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	return body, nil
}
