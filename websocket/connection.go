package websocket

// Connection is the interface that wraps the Websocket connection.
type Connection interface {
	Write(data []byte) error
	Read() ([]byte, error)
	Close() error
}
