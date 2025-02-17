package mexchttpmarket

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	mexchttp "github.com/kattana-io/mexc-golang-sdk/http"

	"github.com/stretchr/testify/assert"
)

func TestService_Ping(t *testing.T) {
	client := mexchttp.NewClient("", "", &http.Client{})

	service := &Service{
		client: client,
	}

	ctx := context.Background()

	response, err := service.Ping(ctx)

	assert.NoError(t, err)

	expectedResponse := "{}"
	assert.Equal(t, expectedResponse, response)
}

func TestService_ExchangeInfo(t *testing.T) {
	client := mexchttp.NewClient("", "", &http.Client{})

	service := &Service{
		client: client,
	}

	ctx := context.Background()

	response, err := service.ExchangeInfo(ctx, []string{"BTCUSDT", "ETHUSDT"})

	assert.NoError(t, err)

	fmt.Println(response)
}

func TestService_Time(t *testing.T) {
	client := mexchttp.NewClient("", "", &http.Client{})

	service := &Service{
		client: client,
	}

	ctx := context.Background()

	response, err := service.Time(ctx)

	assert.NoError(t, err)

	fmt.Println(response)
}
