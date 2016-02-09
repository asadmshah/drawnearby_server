package websocket

import (
	"net/http"

	ws "github.com/gorilla/websocket"
)

var (
	upgrader = ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Websocket implements the Connection interface.
type Websocket struct {
	conn *ws.Conn
}

// NewWebsocket returns a new Connection using the http request/response.
func NewWebsocket(w http.ResponseWriter, r *http.Request) (Connection, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return Connection(&Websocket{conn}), nil
}

// Write writes the data to the connection.
func (c *Websocket) Write(data []byte) error {
	return c.conn.WriteMessage(ws.TextMessage, data)
}

// Read blocks until a message is received from the connection.
func (c *Websocket) Read() ([]byte, error) {
	_, p, err := c.conn.ReadMessage()
	return p, err
}

// Close closes the connection.
func (c *Websocket) Close() error {
	return c.conn.Close()
}
