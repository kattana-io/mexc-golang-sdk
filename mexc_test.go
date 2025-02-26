package mexc

import (
	"context"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/websocket"
	"github.com/kattana-io/mexc-golang-sdk/websocket/market"
	"net/http"
	"testing"
	"time"

	mexchttp "github.com/kattana-io/mexc-golang-sdk/http"
)

func TestHttp(_ *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	cl := mexchttp.NewClient("", "", &http.Client{})

	rClient, _ := NewRest(ctx, cl)
	res, _ := rClient.MarketService.Ping(ctx)

	fmt.Println(res)
	cancel()
}

func TestWs(_ *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	wc := mexcws.NewMEXCWebSocket(func(_ bool, err error) {
		fmt.Println("Error: ", err)
	})

	wc.Connect(ctx, nil)

	ws := NewWs(wc)

	ws.MarketService.OrderBookSubscribe(
		ctx,
		[]string{
			"BTCUSDT",
			"ETHUSDT",
		},
		mexcwsmarket.MinBookDepth,
		func(book *mexcwsmarket.OrderBook) {
			fmt.Println("Symbol: ", book.Symbol)
			fmt.Println("ASKS: ", book.Data.Asks)
			fmt.Println("BIDS: ", book.Data.Bids)
			fmt.Println("-----------")
		},
	)

	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("END")
}
