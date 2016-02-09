package models

import (
	"sync"

	"github.com/asadmshah/drawnearby_server/websocket"
)

type Room struct {
	conns map[websocket.Connection]struct{}
	name  string
	mu    sync.Mutex
	hist  [][]byte
}

func NewRoom(name string) *Room {
	return &Room{
		conns: make(map[websocket.Connection]struct{}),
		name:  name,
		hist:  make([][]byte, 0),
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Size() int {
	return len(r.conns)
}

func (r *Room) Join(conn websocket.Connection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.conns[conn] = struct{}{}
}

func (r *Room) Exit(conn websocket.Connection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.conns[conn]; ok {
		delete(r.conns, conn)
	}
}

func (r *Room) Write(data []byte, source websocket.Connection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.hist = append(r.hist, data)

	for conn, _ := range r.conns {
		if conn != source {
			err := conn.Write(data)
			if err != nil {
				// TODO: Remove user?
			}
		}
	}
}

func (r *Room) WriteHistory(conn websocket.Connection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, data := range r.hist {
		err := conn.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
