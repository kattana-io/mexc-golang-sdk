package connection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kattana-io/mexc-golang-sdk/websocket/dto"
	"github.com/kattana-io/mexc-golang-sdk/websocket/types"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)

const (
	MaxMEXCWebSocketSubscriptions = 30
	keepAliveInterval             = 30 * time.Second
)

type MEXCWebSocketConnection struct {
	Subs          *Subscribes
	Conn          *websocket.Conn
	ErrorListener mexcwstypes.OnError
	sendMutex     *sync.Mutex
	subMtx        *sync.Mutex
	url           string
	readCancel    context.CancelFunc
	ctx           context.Context
	id            string
	logger        *log.Logger
}

func NewMEXCWebSocketConnection(url string, errorListener mexcwstypes.OnError) *MEXCWebSocketConnection {
	return &MEXCWebSocketConnection{
		sendMutex:     &sync.Mutex{},
		subMtx:        &sync.Mutex{},
		url:           url,
		ErrorListener: errorListener,
		Subs:          NewSubs(),
		id:            uuid.NewString(),
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

	m.ctx = ctx
	m.run(ctx)
	return nil
}

func (m *MEXCWebSocketConnection) Send(message *mexcwstypes.WsReq) error {
	if m.Conn == nil {
		return fmt.Errorf("no available connection id: %s", m.id)
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
	go m.reconnectLoop(readCtx)
}

func (m *MEXCWebSocketConnection) reconnectLoop(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(23 * time.Hour):
		m.logger.Printf("run scheduled reconnect. id: %s", m.id)
		if err := m.reconnect(); err != nil {
			m.ErrorListener(true, fmt.Errorf("schedulled reconnect error: %v", err))
		}
	}
}

// keepAlive sends a ping message to the server every 30 seconds to keep the connection alive
func (m *MEXCWebSocketConnection) keepAlive(ctx context.Context) {
	pingTicker := time.NewTicker(keepAliveInterval)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			err := m.Send(&mexcwstypes.WsReq{Method: "PING"})
			if err != nil {
				m.ErrorListener(false, fmt.Errorf("ping error: %v", err))
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

	msgType, buf, err := m.Conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			m.ErrorListener(true, fmt.Errorf("connection closed: %v", err))
			return
		}

		if rErr := m.reconnect(); rErr != nil {
			m.ErrorListener(true, fmt.Errorf("reconnect error: %v", err))
		}
		log.Printf("readLoop error for id %s: %v", m.id, err)
		return
	}

	switch msgType {
	case websocket.TextMessage:
		if string(buf) == "PONG" {
			return
		}
		data := make(map[string]any)
		err := json.Unmarshal(buf, &data)
		if err != nil {
			m.ErrorListener(false, fmt.Errorf("unmarshal error: %v", err))
			return
		}
		m.ErrorListener(false, fmt.Errorf("received unprocessed text message: %v", data))
	case websocket.BinaryMessage:
		update := &dto.PushDataV3ApiWrapper{}
		err = proto.Unmarshal(buf, update)
		if err != nil {
			m.ErrorListener(false, fmt.Errorf("unmarshal error: %v", err))
			return
		}

		listener := m.getListener(update.Channel)
		if listener != nil {
			listener(update)
			return
		}
	case websocket.PingMessage, websocket.PongMessage:
		return
	case websocket.CloseMessage:
		m.ErrorListener(true, fmt.Errorf("received websocket close message: %v", string(buf)))
	default:
		m.ErrorListener(false, fmt.Errorf("unhandled id %s: %v", m.id, string(buf)))
	}
}

func (m *MEXCWebSocketConnection) reconnect() error {
	// stop reading from old connection
	m.readCancel()

	m.subMtx.Lock()
	defer m.subMtx.Unlock()

	newConn, _, err := websocket.DefaultDialer.DialContext(m.ctx, m.url, nil)
	if err != nil {
		return fmt.Errorf("connect error: %v", err)
	}

	if err = m.Disconnect(); err != nil {
		log.Printf("closing old websocket connection error for id %s: %v", m.id, err)
	}

	m.Conn = newConn
	// run new connection read loop
	m.run(m.ctx)

	for _, ch := range m.Subs.GetAllChannels() {
		req := &mexcwstypes.WsReq{
			Method: "SUBSCRIPTION",
			Params: []string{ch},
		}
		if err = m.Send(req); err != nil {
			return fmt.Errorf("resubscription error for channel [%s]: %v", ch, err)
		}
	}

	log.Printf("reconnect successful %s", m.id)
	return nil
}

func (m *MEXCWebSocketConnection) getListener(channel string) mexcwstypes.OnReceive {
	v, _ := m.Subs.Load(channel)
	return v
}

func (m *MEXCWebSocketConnection) Disconnect() error {
	if err := m.Conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(60*time.Second)); err != nil && !errors.Is(err, websocket.ErrCloseSent) {
		return err
	}
	return m.Conn.Close()
}
