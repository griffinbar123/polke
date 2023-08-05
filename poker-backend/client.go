package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type Client struct {
	Name     string `json:"name"`
	PlayerId string `json:"player-id"`
	Balance  int    `json:"balance"`
}
type ClientWrapper struct {
	Type    string `json:"type"`
	Payload Client `json:"payload"`
}

type ClientInRoom struct {
	Client
	Conn   *websocket.Conn `json:"-"`
	RoomId int             `json:"room-id"`
}

var ActivePlayerArray = []ClientInRoom{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func EstablishWS(w http.ResponseWriter, r *http.Request) {
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
			log.Println(client)
			if player, ok := CheckForPlayerMetadata(ws, client); ok {
				JoinRoom(ws, player)
			}
		case "join-room":
			var client = ClientInRoom{}
			if json.Unmarshal([]byte(temp[11:len(temp)-1]), &client); err != nil {
				log.Printf("error: %v", err)
			}
			client.Conn = ws
			log.Println(client)
			JoinRoom(ws, client)
		case "create-room":
			var buyIn = CreateRoomPayload{}
			if json.Unmarshal([]byte(temp[11:len(temp)-1]), &buyIn); err != nil {
				log.Printf("error: %v", err)
			}
			log.Println(buyIn)
			room := NewRoom(buyIn.BuyIn)
			RoomArray = append(RoomArray, room)
			ws.WriteJSON(SendRoomNumber{
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
			log.Println(client)
			if player, ok := FindPlayerInList(client); ok {
				if room, ok := FindRoomFromId(player.RoomId); ok {
					for i, p := range room.Players {
						if p.PlayerId == player.PlayerId {
							room.Players = removePlayer(room.Players, i)
							room.RoomSize -= 1
							UpdateRoomForAllPlayer(*room)
							break
						}
					}
				}
			}
		}

	}
}

type SendRoomNumber struct {
	Type    string     `json:"type"`
	Payload RoomNumber `json:"payload"`
}

type RoomNumber struct {
	RoomNo int `json:"room-number"`
}

type CreateRoomPayload struct {
	BuyIn int `json:"buy-in"`
}

func JoinRoom(conn *websocket.Conn, player ClientInRoom) {
	for i, room := range RoomArray {
		if room.RoomId == player.RoomId {
			if len(room.Players) > 8 {
				log.Printf("room is already at max capacity")
				return
			}
			for j, p := range room.Players {
				if p.PlayerId == player.PlayerId {
					room.Players[j].Conn = conn
					HandleUpdate(conn, player, room)
					return
				}
			}
			AddPlayerToRoom(conn, player, &room)
			RoomArray[i] = room
			HandleUpdate(conn, player, room)
			// log.Println(player)
			// log.Println(room)
			return
		}
	}
}

func HandleUpdate(conn *websocket.Conn, player ClientInRoom, room Room) {
	StoreConnection(conn, player)
	UpdateRoomForAllPlayer(room)
}

func UpdateRoomForAllPlayer(room Room) {
	log.Println("Updating all players")
	for _, p := range room.Players {
		log.Println(p)
		SendRoomState(p.Conn, RoomState{
			Type:    "room-state",
			Payload: room,
		})
	}
}

func SendRoomState(conn *websocket.Conn, room RoomState) {
	conn.WriteJSON(room)
}

func StoreConnection(conn *websocket.Conn, player ClientInRoom) {
	RemoveClientFromClientArray(player)
	AddClientToClientArray(player)
	conn.WriteJSON(ClientWrapper{
		Type:    "player-metadata",
		Payload: player.Client,
	})
}
