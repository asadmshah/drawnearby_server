package messages

import (
	"encoding/json"
)

type lobbyStatusMessage struct {
	Rooms []string `json:"rooms"`
}

// CreateLobbyStatusMessage returns the list of rooms serialized as JSON.
func CreateLobbyStatusMessage(rooms []string) ([]byte, error) {
	return json.Marshal(&lobbyStatusMessage{rooms})
}
