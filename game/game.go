package game

import (
	"math"
	"math/rand"
	"sync"

	"github.com/faiface/pixel"

	guuid "github.com/google/uuid"
)

type Player struct {
	X                     float64
	Y                     float64
	Left, Right, Up, Down bool
	Acceleration          float64
	Velocity              float64
	Rotation              RotationDegree
	Life                  float64
	Power                 float64
	ReloadTime            float64
}

type Camera struct {
	Pos       pixel.Vec
	Speed     float64
	Zoom      float64
	ZoomSpeed float64
}

type Game struct {
	Players map[guuid.UUID]*Player
	Bullets []*Bullet
	You     *Player
}

type Bullet struct {
	X         float64
	Y         float64
	active    bool
	Owner     guuid.UUID
	Rotation  RotationDegree
	Damage    float64
	Speed     float64
	Exhausted bool
}

type RotationDegree float64

const RotationUp RotationDegree = math.Pi * 2
const RotationDown RotationDegree = math.Pi
const RotationLeft RotationDegree = math.Pi / 2
const RotationRight RotationDegree = 3 * math.Pi / 2

const RotationRightUp RotationDegree = (3 * math.Pi / 2) + (math.Pi / 4)
const RotationRightDown RotationDegree = (3 * math.Pi / 2) - (math.Pi / 4)
const RotationLeftUp RotationDegree = (math.Pi / 2) - (math.Pi / 4)
const RotationLeftDown RotationDegree = (math.Pi / 2) + (math.Pi / 4)

func (p *Player) MovePlayer(dt float64) {

	if p.Right && p.X < 1024 {
		p.X += p.Acceleration * p.Velocity
		if !p.Up && !p.Down {
			p.Rotation = RotationRight
		}
	}
	if p.Right && p.Up {
		p.Rotation = RotationRightUp
	}
	if p.Right && p.Down {
		p.Rotation = RotationRightDown
	}

	if p.Left && p.X > 0 {
		p.X -= p.Acceleration * p.Velocity
		if !p.Up && !p.Down {
			p.Rotation = RotationLeft
		}
	}
	if p.Left && p.Up {
		p.Rotation = RotationLeftUp
	}
	if p.Left && p.Down {
		p.Rotation = RotationLeftDown
	}

	if p.Up && p.Y < 768 {
		p.Y += p.Acceleration * p.Velocity
		if !p.Left && !p.Right {
			p.Rotation = RotationUp
		}
	}
	if p.Down && p.Y > 0 {
		p.Y -= p.Acceleration * p.Velocity
		if !p.Left && !p.Right {
			p.Rotation = RotationDown
		}
	}
}

func (g *Game) MovePlayers(dt float64) {
	for _, v := range g.Players {
		m := sync.Mutex{}
		m.Lock()
		v.MovePlayer(dt)
		m.Unlock()
	}
}

func (g *Game) Collision() {
	tolerance := 5.0
	for k, p := range g.Players {
		if p.Life <= 0 {
			continue
		}
		for _, b := range g.Bullets {
			if b.Owner == k {
				continue
			}
			if math.Abs(b.X-p.X) > tolerance {
				continue
			}
			if math.Abs(b.Y-p.Y) > tolerance*2 {
				continue
			}
			p.Life -= 1
			b.Exhausted = true
		}
	}

	for i := len(g.Bullets) - 1; i >= 0; i-- {
		if g.Bullets[i].Exhausted || g.Bullets[i].Y > 768 || g.Bullets[i].X > 1024 || g.Bullets[i].X < 0 || g.Bullets[i].Y < 0 {
			g.Bullets = append(g.Bullets[:i], g.Bullets[i+1:]...)
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
		Power:        1,
		ReloadTime:   50.0,
	}
	g.Players[id] = p
	return p
}

func (g *Game) AddBullet(x, y float64, owner guuid.UUID, rotation RotationDegree, damage float64) {
	g.Bullets = append(g.Bullets, &Bullet{
		X:         x,
		Y:         y,
		Owner:     owner,
		Damage:    damage,
		Rotation:  rotation,
		Speed:     2,
		Exhausted: false,
	})
}

func (g *Game) MoveBullets(dt float64) {
	for _, b := range g.Bullets {
		switch b.Rotation {
		case RotationUp:
			b.Y += b.Speed
			break
		case RotationDown:
			b.Y -= b.Speed
			break
		case RotationLeft:
			b.X -= b.Speed
			break
		case RotationRight:
			b.X += b.Speed
			break
		case RotationLeftDown:
			b.X -= b.Speed
			b.Y -= b.Speed
			break
		case RotationLeftUp:
			b.X -= b.Speed
			b.Y += b.Speed
			break
		case RotationRightDown:
			b.X += b.Speed
			b.Y -= b.Speed
			break
		case RotationRightUp:
			b.X += b.Speed
			b.Y += b.Speed
			break
		}
	}
	return
}
