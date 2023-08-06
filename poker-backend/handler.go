package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func EstablishWS(w http.ResponseWriter, r *http.Request) {
	/*
	 upgrades the http connection to a websocket connection, then listens for messages and handles them
	*/
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		go HandleMessage(ws, message)
	}
}

func HandleMessage(ws *websocket.Conn, message []byte) {
	/*
		 handles the different type od messages the client sends. messages are formatted:
		 {
			"type": "do-something",
			"payload": {"data": "data"}
		 }
	*/
	messageAsString := string(message)
	j := objx.MustFromJSON(messageAsString)
	payload := j.Exclude([]string{"type"})
	temp, err := payload.JSON()
	if err != nil {
		log.Printf("error: %v", err)
	}
	t := j.Get("type").Str()
	log.Println("message: " + t)

	switch t {
	case "player-metadata":
		var client = Client{}
		if json.Unmarshal([]byte(temp[11:len(temp)-1]), &client); err != nil {
			log.Printf("error: %v", err)
		}
		if player, ok := client.CheckForClientMetadata(ws); ok {
			player.JoinRoom()
		}
	case "join-room":
		var client = ClientInRoom{}
		if json.Unmarshal([]byte(temp[11:len(temp)-1]), &client); err != nil {
			log.Printf("error: %v", err)
		}
		client.Conn = ws
		client.JoinRoom()
	case "create-room":
		var buyIn = BuyInReciever{}
		if json.Unmarshal([]byte(temp[11:len(temp)-1]), &buyIn); err != nil {
			log.Printf("error: %v", err)
		}
		room := NewRoom(buyIn.BuyIn)
		RoomArray = append(RoomArray, room)
		ws.WriteJSON(RoomNumberSender{
			Type: "send-room-number",
			Payload: RoomNumber{
				RoomNo: room.RoomId,
			},
		})
	case "disconnect":
		var client = Client{}
		if json.Unmarshal([]byte(temp[11:len(temp)-1]), &client); err != nil {
			log.Printf("error: %v", err)
		}
		client.DisconnectClientFromTheirRoom()
	}
}
