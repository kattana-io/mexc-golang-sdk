package mexc

import (
	"context"
	"fmt"
	mexchttp "github.com/kattana-io/mexc-golang-sdk/http"
	"net/http"
	"testing"
)

func TestHttp(_ *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	cl := mexchttp.NewClient("", "", &http.Client{})

	rClient, _ := NewRest(ctx, cl)
	res, _ := rClient.MarketService.Ping(ctx)

	fmt.Println(res)
	cancel()
}
