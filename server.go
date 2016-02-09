package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/asadmshah/drawnearby_server/models"
	"github.com/asadmshah/drawnearby_server/names"
	"github.com/asadmshah/drawnearby_server/websocket"
)

var (
	addr           = flag.String("addr", "localhost:8080", "Defaults to localhost:8080")
	lobby          = models.NewLobby()
	namesGenerator = names.NewNameGenerator()
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, err := websocket.NewWebsocket(w, r)
	if err != nil {
		log.Printf("ERROR:Server.Index: %v\n", err)
		return
	}
	defer conn.Close()

	err = conn.Write(lobby.StatusMessage())
	if err != nil {
		return
	}

	lobby.InsertUpdateListener(conn)
	defer lobby.RemoveUpdateListener(conn)

	for {
		_, err := conn.Read()
		if err != nil {
			return
		}
	}
}

func RoomJoin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	room := lobby.GetRoom(name)
	if room == nil {
		room = models.NewRoom(name)
		lobby.InsertRoom(room)
		lobby.Write(lobby.StatusMessage())
	}
	defer ClearRoomIfEmpty(room.Name())

	conn, err := websocket.NewWebsocket(w, r)
	if err != nil {
		log.Printf("ERROR:RoomJoin:NewWebsocket %v\n", err)
		return
	}
	defer conn.Close()

	err = room.WriteHistory(conn)
	if err != nil {
		log.Printf("ERROR:RoomJoin:WriteHistory: %v\n", err)
		return
	}

	room.Join(conn)
	defer room.Exit(conn)

	for {
		data, err := conn.Read()
		if err != nil {
			return
		}
		room.Write(data, conn)
	}
}

func ClearRoomIfEmpty(name string) {
	room := lobby.GetRoom(name)
	if room != nil && room.Size() == 0 {
		lobby.RemoveRoom(room)
		lobby.Write(lobby.StatusMessage())
	}
}

func main() {
	flag.Parse()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/join/:name", RoomJoin)

	log.Printf("Serving at %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, router))
}
