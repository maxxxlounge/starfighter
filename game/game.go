package game

type Player struct {
	ID string
	X  float64
	Y  float64
}

type Game struct {
	Players []Player
}
