package mexcwsuser

import (
	"context"
	"encoding/json"
	mexcwstypes "github.com/kattana-io/mexc-golang-sdk/websocket/types"
	"github.com/shopspring/decimal"
)

const (
	SpotOrdersChannel = "spot@private.orders.v3.api"
)

// OrdersSubscribe subscribes to user`s spot orders events, starts listen key keep-alive routine
// https://mexcdevelop.github.io/apidocs/spot_v3_en/#spot-account-orders
func (s *Service) OrdersSubscribe(ctx context.Context, callback func(*OrderEvent), errCallback mexcwstypes.OnError) error {
	listenKey, err := s.httpStream.CreateListenKey(ctx)
	if err != nil {
		return err
	}

	go func(ctx context.Context, listenKey string) {
		kErr := s.httpStream.RunKeyKeepAlive(ctx, listenKey)
		if kErr != nil {
			errCallback(err)
		}
	}(ctx, listenKey)

	lstnr := func(message string) {
		var book OrderEvent
		mErr := json.Unmarshal([]byte(message), &book)
		if mErr != nil {
			errCallback(mErr)
			return
		}

		callback(&book)
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

type Side int32

const (
	SideBuy Side = iota + 1
	SideSell
)

type Type int32

const (
	TypeLimitOrder Type = iota + 1
	TypePostOnly
	TypeImmediateOrCancel
	TypeFillOrKill
	TypeMarketOrder
	TypeStopLimit
)

type Status int32

const (
	StatusNew Status = iota + 1
	StatusFilled
	StatusPartiallyFilled
	StatusCancelled
	StatusPartiallyCancelled
)

type OrderEvent struct {
	Channel string `json:"c"`
	Data    struct {
		RemainAmount       decimal.Decimal `json:"A"`
		CreateTime         int64           `json:"O"`
		Side               Side            `json:"S"`
		RemainQuantity     decimal.Decimal `json:"V"`
		Amount             decimal.Decimal `json:"a"`
		ClientOrderID      string          `json:"c"`
		OrderID            string          `json:"i"`
		IsMaker            byte            `json:"m"`
		Type               Type            `json:"o"`
		Price              decimal.Decimal `json:"p"`
		Status             Status          `json:"s"`
		Quantity           decimal.Decimal `json:"v"`
		AveragePrice       decimal.Decimal `json:"ap"`
		CumulativeQuantity decimal.Decimal `json:"cv"`
		CumulativeAmount   decimal.Decimal `json:"ca"`
	} `json:"d"`
	Symbol    string `json:"s"`
	Timestamp int64  `json:"t"`
}
