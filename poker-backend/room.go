package main

import "github.com/gorilla/websocket"

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

type Player struct {
	Client
	Position int             `json:"position"`
	Cards    []int           `json:"cards"`
	Conn     *websocket.Conn `json:"-"`
}

type Room struct {
	RoomId                 int      `json:"room-id"`
	RoomSize               int      `json:"room-size"`
	RoomBuyIn              int      `json:"room-buy-in"`
	RoomPot                int      `json:"room-pot"`
	RoomBigBlindPosition   int      `json:"room-big-blind-position"`
	RoomSmallBlindPosition int      `json:"room-small-blind-position"`
	RoomDealerPosition     int      `json:"room-dealer-position"`
	Players                []Player `json:"players"`
}

type RoomState struct {
	Type    string `json:"type"`
	Payload Room   `json:"payload"`
}
