package connection

type MEXCWebSocketConnections struct {
	data []*MEXCWebSocketConnection
}

func NewMEXCWebSocketConnections() *MEXCWebSocketConnections {
	return &MEXCWebSocketConnections{
		data: make([]*MEXCWebSocketConnection, 0),
	}
}

func (c *MEXCWebSocketConnections) Len() int {
	return len(c.data)
}

func (c *MEXCWebSocketConnections) Less(i, j int) bool {
	return c.data[i].Subs.Len() < c.data[j].Subs.Len()
}

func (c *MEXCWebSocketConnections) Swap(i, j int) {
	c.data[i], c.data[j] = c.data[j], c.data[i]
}

func (c *MEXCWebSocketConnections) Push(x any) {
	item := x.(*MEXCWebSocketConnection)
	c.data = append(c.data, item)
}

func (c *MEXCWebSocketConnections) Pop() any {
	old := c.data
	n := len(old)
	item := old[n-1]
	c.data = old[0 : n-1]
	return item
}
