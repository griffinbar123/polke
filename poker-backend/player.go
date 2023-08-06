package main

import "github.com/gorilla/websocket"

type Player struct { // a player is a client with extra game specific data
	Client
	Position int             `json:"position"`
	Cards    []int           `json:"cards"`
	Conn     *websocket.Conn `json:"-"`
}

func (player *Player) SendRoomState(room *Room) {
	/*
		sends the room state to a player
	*/
	player.Conn.WriteJSON(RoomSender{
		Type:    "room-state",
		Payload: *room,
	})
}
