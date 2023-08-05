package main

import (
	"math/rand"

	"github.com/gorilla/websocket"
)

func RangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func SendEvent(conn *websocket.Conn, t string, p string) error {
	return conn.WriteMessage(websocket.TextMessage, []byte("{\"type\": \""+t+"\", \"payload\": "+p+"}"))
}

func CheckForPlayerMetadata(conn *websocket.Conn, playerData Client) (ClientInRoom, bool) {
	if player, ok := FindPlayerInList(playerData); ok {
		player.Conn = conn
		return player, true
	}
	return ClientInRoom{}, false

}

func FindPlayerInList(playerData Client) (ClientInRoom, bool) {
	for _, player := range ActivePlayerArray {
		if player.PlayerId == playerData.PlayerId {
			return player, true
		}
	}
	return ClientInRoom{}, false
}

func FindRoomFromId(id int) (*Room, bool) {
	for _, room := range RoomArray {
		if room.RoomId == id {
			return &room, true
		}
	}
	return &Room{}, false
}

func AddPlayerToRoom(conn *websocket.Conn, player ClientInRoom, room *Room) {
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
	room.Players = append(room.Players, Player{
		Position: smalledUnfilledPosition,
		Cards:    []int{0, 0},
		Client:   player.Client,
		Conn:     player.Conn,
	})
}

func remove(s []ClientInRoom, i int) []ClientInRoom {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removePlayer(s []Player, i int) []Player {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveClientFromClientArray(player ClientInRoom) {
	for i, p := range ActivePlayerArray {
		if p.PlayerId == player.PlayerId {
			ActivePlayerArray = remove(ActivePlayerArray, i)
		}
	}
}

func AddClientToClientArray(player ClientInRoom) {
	ActivePlayerArray = append(ActivePlayerArray, player)
}
