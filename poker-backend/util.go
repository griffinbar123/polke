package main

import (
	"math/rand"
)

func RangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func remove(s []ClientInRoom, i int) []ClientInRoom {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removePlayer(s []Player, i int) []Player {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
