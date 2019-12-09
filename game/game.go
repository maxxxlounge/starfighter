package game

import (
	guuid "github.com/google/uuid"
)

type Player struct {
	X float64
	Y float64
}

type Game struct {
	Players map[guuid.UUID]Player
}
