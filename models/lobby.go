package models

import (
	"sort"
	"sync"

	"github.com/asadmshah/drawnearby_server/messages"
	"github.com/asadmshah/drawnearby_server/websocket"
)

// Lobby hosts all the rooms and any user waiting to play.
type Lobby struct {
	rooms   map[string]*Room
	roomsMu sync.Mutex
	conns   map[websocket.Connection]struct{}
	connsMu sync.Mutex
}

// NewLobby creates a new empty lobby.
func NewLobby() *Lobby {
	return &Lobby{
		rooms: make(map[string]*Room),
		conns: make(map[websocket.Connection]struct{}),
	}
}

// InsertRoom generates a new room
func (l *Lobby) InsertRoom(room *Room) {
	l.roomsMu.Lock()
	defer l.roomsMu.Unlock()

	l.rooms[room.Name()] = room
}

func (l *Lobby) GetRoom(roomName string) *Room {
	l.roomsMu.Lock()
	defer l.roomsMu.Unlock()

	return l.rooms[roomName]
}

func (l *Lobby) RemoveRoom(room *Room) {
	l.roomsMu.Lock()
	defer l.roomsMu.Unlock()

	if _, ok := l.rooms[room.Name()]; ok {
		delete(l.rooms, room.Name())
	}
}

func (l *Lobby) InsertUpdateListener(conn websocket.Connection) {
	l.connsMu.Lock()
	defer l.connsMu.Unlock()

	l.conns[conn] = struct{}{}
}

func (l *Lobby) RemoveUpdateListener(conn websocket.Connection) {
	l.connsMu.Lock()
	defer l.connsMu.Unlock()

	if _, ok := l.conns[conn]; ok {
		delete(l.conns, conn)
	}
}

func (l *Lobby) Write(data []byte) {
	l.connsMu.Lock()
	defer l.connsMu.Unlock()
	l.roomsMu.Lock()
	defer l.roomsMu.Unlock()

	for conn, _ := range l.conns {
		err := conn.Write(data)

		// Remove any troublesome connections.
		if err != nil {
			delete(l.conns, conn)
		}
	}
}

func (l *Lobby) StatusMessage() []byte {
	l.roomsMu.Lock()
	defer l.roomsMu.Unlock()

	rooms := make([]string, len(l.rooms), len(l.rooms))
	i := 0
	for room, _ := range l.rooms {
		rooms[i] = room
		i++
	}
	sort.Strings(rooms)

	data, _ := messages.CreateLobbyStatusMessage(rooms)
	return data
}
