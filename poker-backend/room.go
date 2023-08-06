package main

import (
	"log"
)

var RoomArray = []Room{}

func NewRoom(buyIn int) Room {
	var roomNumber int
GetRoomNumerLoop:
	for {
		roomNumber = RangeIn(1000, 9999)
		for _, room := range RoomArray {
			if room.RoomId == roomNumber {
				continue GetRoomNumerLoop
			}
		}
		break GetRoomNumerLoop
	}
	return Room{
		RoomId:                 roomNumber,
		RoomSize:               0,
		RoomBuyIn:              buyIn,
		RoomPot:                0,
		RoomBigBlindPosition:   2,
		RoomSmallBlindPosition: 1,
		RoomDealerPosition:     0,
		Players:                make([]Player, 0),
	}
}

type Room struct { // stores data about a room and its players
	RoomId                 int      `json:"room-id"`
	RoomSize               int      `json:"room-size"`
	RoomBuyIn              int      `json:"room-buy-in"`
	RoomPot                int      `json:"room-pot"`
	RoomBigBlindPosition   int      `json:"room-big-blind-position"`
	RoomSmallBlindPosition int      `json:"room-small-blind-position"`
	RoomDealerPosition     int      `json:"room-dealer-position"`
	Players                []Player `json:"players"`
}

func (room *Room) UpdateRoomForAllPlayers() {
	/*
		Goes through all the players and sends them the current state of the room
	*/
	log.Println("Updating all players")
	for _, p := range room.Players {
		log.Println(p)
		p.SendRoomState(room)
	}
}

type RoomSender struct {
	Type    string `json:"type"`
	Payload Room   `json:"payload"`
}

func FindRoomFromId(id int) (int, *Room, bool) {
	for i, room := range RoomArray {
		if room.RoomId == id {
			return i, &room, true
		}
	}
	return 0, &Room{}, false
}

type RoomNumberSender struct {
	Type    string     `json:"type"`
	Payload RoomNumber `json:"payload"`
}

type RoomNumber struct {
	RoomNo int `json:"room-number"`
}

type BuyInReciever struct {
	BuyIn int `json:"buy-in"`
}
