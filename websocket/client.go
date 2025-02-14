package mexcws

import (
	"container/heap"
	"context"
	"sync"
)

const (
	MaxMEXCWebSocketSubscriptions = 30
)

// MEXCWebSocket is a WebSocket client for the MEXC exchange
type MEXCWebSocket struct {
	url           string
	mtx           *sync.Mutex
	Connections   *MEXCWebSocketConnections
	ErrorListener OnError
	subscribeMap  map[string]*MEXCWebSocketConnection
}

// NewMEXCWebSocket returns a new MEXCWebSocket instance
func NewMEXCWebSocket(errorListener OnError) *MEXCWebSocket {
	return &MEXCWebSocket{
		url:           "wss://wbs.mexc.com/ws",
		Connections:   NewMEXCWebSocketConnections(),
		ErrorListener: errorListener,
		subscribeMap:  make(map[string]*MEXCWebSocketConnection),
		mtx:           new(sync.Mutex),
	}
}

// Send sends a message to the server
func (m *MEXCWebSocket) Send(ctx context.Context, message *WsReq) error {
	conn, err := m.getWsConnection(ctx, false)
	if err != nil {
		return err
	}

	return conn.Send(message)
}

// Connect establishes a WebSocket connection to the MEXC exchange
func (m *MEXCWebSocket) Connect(ctx context.Context) error {
	_, err := m.getWsConnection(ctx, false)
	return err
}

func (m *MEXCWebSocket) Subscribe(ctx context.Context, channel string, callback OnReceive) error {
	conn, err := m.getWsConnection(ctx, true)
	if err != nil {
		return err
	}

	if err = conn.Subscribe(channel, callback); err != nil {
		return err
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.subscribeMap[channel] = conn
	return nil
}

func (m *MEXCWebSocket) Unsubscribe(channel string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	conn, ok := m.subscribeMap[channel]
	if !ok {
		return nil
	}

	if err := conn.Unsubscribe(channel); err != nil {
		return err
	}
	delete(m.subscribeMap, channel)
	return nil
}

func (m *MEXCWebSocket) getWsConnection(ctx context.Context, isSubscribe bool) (*MEXCWebSocketConnection, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.Connections.Len() == 0 {
		newConn, err := m.connectWs(ctx)
		if err != nil {
			return nil, err
		}

		heap.Push(m.Connections, newConn)
		return newConn, nil
	}

	lastConn := heap.Pop(m.Connections).(*MEXCWebSocketConnection)
	defer heap.Push(m.Connections, lastConn)
	if isSubscribe && lastConn.Subs.Len() < MaxMEXCWebSocketSubscriptions {
		return lastConn, nil
	}

	newConn, err := m.connectWs(ctx)
	if err != nil {
		return nil, err
	}
	heap.Push(m.Connections, newConn)
	return newConn, nil
}

func (m *MEXCWebSocket) connectWs(ctx context.Context) (*MEXCWebSocketConnection, error) {
	newConn := NewMEXCWebSocketConnection(m.url, m.ErrorListener)
	if err := newConn.Connect(ctx); err != nil {
		return nil, err
	}

	return newConn, nil
}

// Disconnect closes the WebSocket connection
func (m *MEXCWebSocket) Disconnect() error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for c := m.Connections.Pop(); m.Connections.Len() > 0; {
		conn := c.(*MEXCWebSocketConnection)
		err := conn.Disconnect()
		if err != nil {
			return err
		}
	}
	return nil
}
