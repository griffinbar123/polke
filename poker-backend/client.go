package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Name     string `json:"name"`
	PlayerId string `json:"player-id"`
	Balance  int    `json:"balance"`
}

func (client *Client) CheckForClientMetadata(conn *websocket.Conn) (ClientInRoom, bool) {
	/*
		looks at the localhost data send by the client and checks if the player data it contains
		exists as an active player connection
	*/
	if player, ok := client.FindClientInList(); ok {
		player.Conn = conn
		return player, true
	}
	return ClientInRoom{}, false
}

func (client *Client) FindClientInList() (ClientInRoom, bool) {
	/*
		checks if a client is in the list of active players
	*/
	for _, player := range ActivePlayerArray {
		if player.PlayerId == client.PlayerId {
			return player, true
		}
	}
	return ClientInRoom{}, false
}

func (client *Client) DisconnectClientFromTheirRoom() {
	/*
		takes a client and removes them from the room in which they are apart of
	*/
	if clientInRoom, ok := client.FindClientInList(); ok {
		if _, room, ok := FindRoomFromId(clientInRoom.RoomId); ok {
			for i, p := range room.Players {
				if p.PlayerId == clientInRoom.PlayerId {
					room.Players = removePlayer(room.Players, i)
					room.RoomSize -= 1
					room.UpdateRoomForAllPlayers()
					break
				}
			}
		}
	}
}

type ClientSender struct {
	Type    string `json:"type"`
	Payload Client `json:"payload"`
}

type ClientInRoom struct { //struct to store a client and their associeted connection and room id
	Client
	Conn   *websocket.Conn `json:"-"`
	RoomId int             `json:"room-id"`
}

func (client *ClientInRoom) RemoveClientFromClientArray() {
	for i, p := range ActivePlayerArray {
		if p.PlayerId == client.PlayerId {
			ActivePlayerArray = remove(ActivePlayerArray, i)
		}
	}
}

func (client *ClientInRoom) AddClientToClientArray() {
	ActivePlayerArray = append(ActivePlayerArray, *client)
}

func (client *ClientInRoom) JoinRoom() {
	/*
		handles adding a client to a room (if they are not already in the room)
	*/
	if i, room, ok := FindRoomFromId(client.RoomId); ok {
		if len(room.Players) > 8 {
			log.Printf("warning: cannot add player - room is already at max capacity")
			return
		}
		for j, p := range room.Players {
			if p.PlayerId == client.PlayerId {
				room.Players[j].Conn = client.Conn
				client.HandlePlayerUpdate(room)
				return
			}
		}
		client.AddPlayerToRoom(room)
		RoomArray[i] = *room
		client.StoreConnection()
		client.HandlePlayerUpdate(room)
		return
	}
}

func (client *ClientInRoom) AddPlayerToRoom(room *Room) {
	/*
		handles adding a client to a room and making them a player
	*/
	smalledUnfilledPosition := 0
	for {
		z := false
		for _, p := range room.Players {
			if p.Position == smalledUnfilledPosition {
				smalledUnfilledPosition += 1
				z = true
				break
			}
		}
		if z == false {
			break
		}
	}
	room.Players = append(room.Players, client.NewPlayer(smalledUnfilledPosition, []int{0, 0}))
}

func (client *ClientInRoom) NewPlayer(smalledUnfilledPosition int, cards []int) Player {
	return Player{
		Position: smalledUnfilledPosition,
		Cards:    cards,
		Client:   client.Client,
		Conn:     client.Conn,
	}
}

func (client *ClientInRoom) StoreConnection() {
	/*
		stores player info on client side in localstorage
	*/
	client.RemoveClientFromClientArray()
	client.AddClientToClientArray()
	client.Conn.WriteJSON(ClientSender{
		Type:    "player-metadata",
		Payload: client.Client,
	})
}

func (client *ClientInRoom) HandlePlayerUpdate(room *Room) {
	/*
		when a player updates, we send the current roomdata to all the players so everyone
		is looking at the same room
	*/
	room.UpdateRoomForAllPlayers()
}

var ActivePlayerArray = []ClientInRoom{}
