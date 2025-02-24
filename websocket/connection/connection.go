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
	readCancel    context.CancelFunc
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
	if m.Conn != nil {
		// already connected
		return nil
	}

	var err error

	m.Conn, _, err = websocket.DefaultDialer.DialContext(ctx, m.url, nil)
	if err != nil {
		return err
	}

	m.run(ctx)
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

	m.Subs.Add(channel, callback)
	err := m.Send(&mexcwstypes.WsReq{
		Method: "SUBSCRIPTION",
		Params: []string{channel},
	})
	if err != nil {
		m.Subs.Remove(channel)
		return err
	}

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

func (m *MEXCWebSocketConnection) run(ctx context.Context) {
	readCtx, cancel := context.WithCancel(ctx)
	m.readCancel = cancel

	go m.keepAlive(readCtx)
	go m.readLoop(readCtx)
	go m.reconnectLoop(ctx)
}

func (m *MEXCWebSocketConnection) reconnectLoop(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(23 * time.Hour):
		if err := m.reconnect(ctx); err != nil {
			m.ErrorListener(fmt.Errorf("reconnect error: %v", err))
		}
	}
}

// keepAlive sends a ping message to the server every 30 seconds to keep the connection alive
func (m *MEXCWebSocketConnection) keepAlive(ctx context.Context) {
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			err := m.Send(&mexcwstypes.WsReq{Method: "PING"})
			if err != nil {
				m.ErrorListener(fmt.Errorf("ping error: %v", err))
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
		m.ErrorListener(fmt.Errorf("read error: %v", err))
		return
	}

	message := string(buf)

	var update map[string]interface{}
	err = json.Unmarshal(buf, &update)
	if err != nil {
		m.ErrorListener(fmt.Errorf("unmarshal error: %v", err))

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
		return fmt.Errorf("connect error: %v", err)
	}

	// stop reading from old connection
	m.readCancel()
	oldConn := m.Conn
	m.Conn = newConn
	// run new connection read loop
	m.run(ctx)

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
