package mexcwsuser

import (
	"context"
	"fmt"
	"github.com/kattana-io/mexc-golang-sdk/websocket/dto"
	mexcwstypes "github.com/kattana-io/mexc-golang-sdk/websocket/types"
)

const (
	SpotOrdersChannel = "spot@private.orders.v3.api.pb"
)

// OrdersSubscribe subscribes to user`s spot orders events, starts listen key keep-alive routine
// https://mexcdevelop.github.io/apidocs/spot_v3_en/#spot-account-orders
func (s *Service) OrdersSubscribe(ctx context.Context, callback func(*dto.PrivateOrdersV3Api), errCallback mexcwstypes.OnError) error {
	listenKey, err := s.httpStream.CreateListenKey(ctx)
	if err != nil {
		return err
	}

	go func(ctx context.Context, listenKey string) {
		kErr := s.httpStream.RunKeyKeepAlive(ctx, listenKey)
		if kErr != nil {
			errCallback(true, err)
		}
	}(ctx, listenKey)

	lstnr := func(message *dto.PushDataV3ApiWrapper) {
		switch msg := message.Body.(type) {
		case *dto.PushDataV3ApiWrapper_PrivateOrders:
			callback(msg.PrivateOrders)
		default:
			fmt.Println("Order callback unknown type:", message.Body)
		}
	}

	params := map[string]string{
		"listenKey": listenKey,
	}
	if err := s.wsClient.Subscribe(ctx, SpotOrdersChannel, params, lstnr); err != nil {
		return err
	}
	return nil
}

func (s *Service) OrdersUnsubscribe() error {
	return s.wsClient.Unsubscribe(SpotOrdersChannel)
}
