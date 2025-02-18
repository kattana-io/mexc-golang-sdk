package connection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kattana-io/mexc-golang-sdk/websocket/types"
	"log"
	"sync"
	"time"
)

const (
	MaxMEXCWebSocketSubscriptions = 30
)

type MEXCWebSocketConnection struct {
	Subs          *Subscribes
	Conn          *websocket.Conn
	ErrorListener mexcwstypes.OnError
	sendMutex     *sync.Mutex
	subMtx        *sync.Mutex
	url           string
}

func NewMEXCWebSocketConnection(url string, errorListener mexcwstypes.OnError) *MEXCWebSocketConnection {
	return &MEXCWebSocketConnection{
		sendMutex:     &sync.Mutex{},
		subMtx:        &sync.Mutex{},
		url:           url,
		ErrorListener: errorListener,
		Subs:          NewSubs(),
	}
}

// Connect establishes a WebSocket connection to the MEXC exchange
func (m *MEXCWebSocketConnection) Connect(ctx context.Context) error {
	var err error

	m.Conn, _, err = websocket.DefaultDialer.DialContext(ctx, m.url, nil)
	if err != nil {
		return err
	}

	go m.keepAlive(ctx)
	go m.readLoop(ctx)

	return nil
}

func (m *MEXCWebSocketConnection) Send(message *mexcwstypes.WsReq) error {
	if m.Conn == nil {
		return fmt.Errorf("no available connection")
	}

	m.sendMutex.Lock()
	defer m.sendMutex.Unlock()

	return m.Conn.WriteJSON(message)
}

func (m *MEXCWebSocketConnection) Subscribe(channel string, callback mexcwstypes.OnReceive) error {
	m.subMtx.Lock()
	defer m.subMtx.Unlock()

	if m.Subs.Len() >= MaxMEXCWebSocketSubscriptions {
		return errors.New("max subscriptions exceeded")
	}

	err := m.Send(&mexcwstypes.WsReq{
		Method: "SUBSCRIPTION",
		Params: []string{channel},
	})
	if err != nil {
		return err
	}

	m.Subs.Add(channel, callback)
	return nil
}

func (m *MEXCWebSocketConnection) Unsubscribe(channel string) error {
	m.subMtx.Lock()
	defer m.subMtx.Unlock()

	m.Subs.Remove(channel)
	return m.Send(&mexcwstypes.WsReq{
		Method: "UNSUBSCRIPTION",
		Params: []string{channel},
	})
}

// keepAlive sends a ping message to the server every 30 seconds to keep the connection alive
func (m *MEXCWebSocketConnection) keepAlive(ctx context.Context) {
	pingTicker := time.NewTicker(30 * time.Second)
	reconnectTicker := time.NewTicker(23 * time.Hour) // mexc terminate connection after 24h
	defer pingTicker.Stop()
	defer reconnectTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			err := m.Send(&mexcwstypes.WsReq{Method: "PING"})
			if err != nil {
				m.ErrorListener(err)
			}
		case <-reconnectTicker.C:
			if err := m.reconnect(ctx); err != nil {
				m.ErrorListener(err)
			}
		}
	}
}

// readLoop read messages and resolve handlers
func (m *MEXCWebSocketConnection) readLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			m.handleLoop()
		}
	}
}

func (m *MEXCWebSocketConnection) handleLoop() {
	if m.Conn == nil {
		return
	}

	_, buf, err := m.Conn.ReadMessage()
	if err != nil {
		m.ErrorListener(err)

		return
	}

	message := string(buf)

	var update map[string]interface{}
	err = json.Unmarshal(buf, &update)
	if err != nil {
		m.ErrorListener(err)

		return
	}

	if update["msg"] == "PONG" {
		return
	}

	listener := m.getListener(update)
	if listener != nil {
		listener(message)

		return
	}

	log.Println(fmt.Sprintf("Unhandled: %v", update))
}

func (m *MEXCWebSocketConnection) reconnect(ctx context.Context) error {
	m.subMtx.Lock()
	defer m.subMtx.Unlock()

	newConn, _, err := websocket.DefaultDialer.DialContext(ctx, m.url, nil)
	if err != nil {
		return err
	}
	oldConn := m.Conn
	m.Conn = newConn

	req := &mexcwstypes.WsReq{
		Method: "SUBSCRIPTION",
		Params: m.Subs.GetAllChannels(),
	}
	if err = m.Send(req); err != nil {
		return err
	}

	return oldConn.Close()
}

func (m *MEXCWebSocketConnection) getListener(argJson interface{}) mexcwstypes.OnReceive {
	mapData := argJson.(map[string]interface{})

	v, _ := m.Subs.Load(fmt.Sprintf("%s", mapData["c"]))
	return v
}

func (m *MEXCWebSocketConnection) Disconnect() error {
	return m.Conn.Close()
}
