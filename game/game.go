package game

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"

	guuid "github.com/google/uuid"
)

type Player struct {
	X                     float64
	Y                     float64
	Left, Right, Up, Down bool
	Acceleration          float64
	Velocity              float64
	Rotation              float64
	Life                  int
}

type Camera struct {
	Pos       pixel.Vec
	Speed     float64
	Zoom      float64
	ZoomSpeed float64
}

type Game struct {
	Players map[guuid.UUID]*Player
	Bullets []Bullet
	You     *Player
}

type Bullet struct {
	X      float64
	Y      float64
	active bool
	Owner  guuid.UUID
}

func (p *Player) MovePlayer(dt float64) {

	if p.Right {
		p.X += p.Acceleration * p.Velocity
		p.Rotation = 3 * math.Pi / 2
		if p.Up {
			p.Rotation += math.Pi / 4
		}
		if p.Down {
			p.Rotation -= math.Pi / 4
		}
	}
	if p.Left {
		p.X -= p.Acceleration * p.Velocity
		p.Rotation = math.Pi / 2
		if p.Down {
			p.Rotation += math.Pi / 4
		}
		if p.Up {
			p.Rotation -= math.Pi / 4
		}
	}
	if p.Up {
		p.Y += p.Acceleration * p.Velocity
		if !p.Left && !p.Right {
			p.Rotation = math.Pi * 2
		}
	}
	if p.Down {
		p.Y -= p.Acceleration * p.Velocity
		if !p.Left && !p.Right {
			p.Rotation = math.Pi
		}
	}
}

func (g *Game) MovePlayers(dt float64) {
	for _, v := range g.Players {
		v.MovePlayer(dt)
	}
}

func (g *Game) Collision() {
	tolerance := 0.5
	for k, p := range g.Players {
		for _, b := range g.Bullets {
			if b.Owner == k {
				continue
			}
			if math.Abs(b.X-p.X) > tolerance {
				continue
			}
			if math.Abs(b.Y-p.Y) > tolerance {
				continue
			}
			p.Life -= 1
		}
	}
}

func (g *Game) NewPlayer(id guuid.UUID) *Player {
	p := &Player{
		X:            rand.Float64() * 500,
		Y:            rand.Float64() * 500,
		Left:         false,
		Right:        false,
		Up:           false,
		Down:         false,
		Acceleration: 1.5,
		Velocity:     1,
		Rotation:     math.Pi / 2,
		Life:         10,
	}
	g.Players[id] = p
	return p
}

func (g *Game) AddBullet(x, y float64, owner guuid.UUID) {
	g.Bullets = append(g.Bullets, Bullet{
		X:     x,
		Y:     y,
		Owner: owner,
	})
}

func (g *Game) MoveBullets(dt float64) {
	return
}
