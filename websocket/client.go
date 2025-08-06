package mexcws

import (
	"container/heap"
	"context"
	"github.com/kattana-io/mexc-golang-sdk/websocket/connection"
	"github.com/kattana-io/mexc-golang-sdk/websocket/types"
	"net/url"
	"sync"
)

// MEXCWebSocket is a WebSocket client for the MEXC exchange
type MEXCWebSocket struct {
	URL           string
	mtx           *sync.Mutex
	Connections   *connection.MEXCWebSocketConnections
	ErrorListener mexcwstypes.OnError
	subscribeMap  map[string]*connection.MEXCWebSocketConnection
}

// NewMEXCWebSocket returns a new MEXCWebSocket instance
func NewMEXCWebSocket(errorListener mexcwstypes.OnError) *MEXCWebSocket {
	return &MEXCWebSocket{
		URL:           "wss://wbs-api.mexc.com/ws",
		Connections:   connection.NewMEXCWebSocketConnections(),
		ErrorListener: errorListener,
		subscribeMap:  make(map[string]*connection.MEXCWebSocketConnection),
		mtx:           new(sync.Mutex),
	}
}

// Send sends a message to the server
func (m *MEXCWebSocket) Send(ctx context.Context, message *mexcwstypes.WsReq) error {
	conn, err := m.getWsConnection(ctx, nil, false)
	if err != nil {
		return err
	}

	return conn.Send(message)
}

// Connect establishes a WebSocket connection to the MEXC exchange
func (m *MEXCWebSocket) Connect(ctx context.Context, params map[string]string) (*connection.MEXCWebSocketConnection, error) {
	return m.getWsConnection(ctx, params, false)
}

func (m *MEXCWebSocket) Subscribe(ctx context.Context, channel string, params map[string]string,
	callback mexcwstypes.OnReceive) error {
	conn, err := m.getWsConnection(ctx, params, true)
	if err != nil {
		return err
	}

	if err := conn.Subscribe(channel, callback); err != nil {
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

func (m *MEXCWebSocket) getWsConnection(ctx context.Context, params map[string]string,
	isSubscribe bool) (*connection.MEXCWebSocketConnection, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if len(params) != 0 || m.Connections.Len() == 0 {
		newConn, err := m.connectWs(ctx, params)
		if err != nil {
			return nil, err
		}

		heap.Push(m.Connections, newConn)
		return newConn, nil
	}

	lastConn := heap.Pop(m.Connections).(*connection.MEXCWebSocketConnection)
	defer heap.Push(m.Connections, lastConn)
	if isSubscribe && lastConn.Subs.Len() < connection.MaxMEXCWebSocketSubscriptions {
		return lastConn, nil
	}

	newConn, err := m.connectWs(ctx, nil)
	if err != nil {
		return nil, err
	}
	heap.Push(m.Connections, newConn)
	return newConn, nil
}

func (m *MEXCWebSocket) connectWs(ctx context.Context, params map[string]string) (*connection.MEXCWebSocketConnection, error) {
	reqURL, err := url.Parse(m.URL)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}
	reqURL.RawQuery = query.Encode()

	newConn := connection.NewMEXCWebSocketConnection(reqURL.String(), m.ErrorListener)
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
		conn := c.(*connection.MEXCWebSocketConnection)
		err := conn.Disconnect()
		if err != nil {
			return err
		}
	}
	return nil
}
