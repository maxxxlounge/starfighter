package game

import (
	"math"
	"math/rand"
	"sync"

	"github.com/faiface/pixel"

	guuid "github.com/google/uuid"
)

type Player struct {
	ID                          string
	UUID                        guuid.UUID
	Name                        string
	X                           float64
	Y                           float64
	Left, Right, Up, Down, Fire bool
	Acceleration                float64
	Velocity                    float64
	Rotation                    RotationDegree
	Life                        float64
	Power                       float64
	ReloadTime                  float64
	Status                      PlayerStatus
	Score                       int
	You                         bool
}

type PlayerStatus string

const WaitForPlay PlayerStatus = "WaitForPlay"
const Pause PlayerStatus = "Pause"
const Resume PlayerStatus = "Resume"
const Ready PlayerStatus = "Ready"
const Died PlayerStatus = "Died"
const Idle PlayerStatus = "Idle"
const Respawn PlayerStatus = "Respawn"

type Camera struct {
	Pos       pixel.Vec
	Speed     float64
	Zoom      float64
	ZoomSpeed float64
}

type GameStatus string

const WaitForPlayer GameStatus = "WaitForPlayer"
const Playing GameStatus = "Playing"
const Scoreboard GameStatus = "Scoreboard"

type Game struct {
	playerMap map[guuid.UUID]*Player
	Players   []*Player
	Bullets   []*Bullet
	Status    GameStatus
	You       *Player
	Bounds    Bounds
}
type Bounds struct {
	Width  float64
	Height float64
}

type Bullet struct {
	ID        guuid.UUID
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

const GameWidth float64 = 1024
const GameHeight float64 = 768

const RotationUp RotationDegree = math.Pi * 2
const RotationDown RotationDegree = math.Pi
const RotationLeft RotationDegree = math.Pi / 2
const RotationRight RotationDegree = 3 * math.Pi / 2

const RotationRightUp RotationDegree = (3 * math.Pi / 2) + (math.Pi / 4)
const RotationRightDown RotationDegree = (3 * math.Pi / 2) - (math.Pi / 4)
const RotationLeftUp RotationDegree = (math.Pi / 2) - (math.Pi / 4)
const RotationLeftDown RotationDegree = (math.Pi / 2) + (math.Pi / 4)

func New() *Game {
	g := Game{
		playerMap: make(map[guuid.UUID]*Player),
		Bounds: Bounds{
			Width:  GameWidth,
			Height: GameHeight,
		},
	}
	return &g
}

func (p *Player) MovePlayer(dt float64) {
	if p.Life <= 0 {
		return
	}
	if p.ReloadTime > 0 {
		p.ReloadTime--
	}

	if p.Right && p.X < GameWidth {
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

	if p.Up && p.Y < GameHeight {
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
		if v.Fire && v.ReloadTime <= 0 {
			g.AddBullet(v.X, v.Y, v.UUID, v.Rotation, v.Power)
			v.ReloadTime = 25
		}
		m.Unlock()
	}
}

func (g *Game) Collision() {
	tolerance := 5.0
	for k, p := range g.playerMap {
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
			g.playerMap[b.Owner].Score++
		}
	}
	for i := len(g.Bullets) - 1; i >= 0; i-- {
		if g.Bullets[i].Exhausted {
			g.Bullets = append(g.Bullets[:i], g.Bullets[i+1:]...)
			continue
		}
		if g.Bullets[i].Y > GameHeight || g.Bullets[i].X > GameWidth || g.Bullets[i].X < 0 || g.Bullets[i].Y < 0 {
			g.Bullets[i].Exhausted = true
		}
	}
}

func (g *Game) SetYou(id guuid.UUID) {
	g.You = g.playerMap[id]
}

func (g *Game) DeletePlayer(id guuid.UUID) {
	for i, p := range g.Players {
		if p == nil {
			continue
		}
		if p.UUID != id {
			continue
		}
		g.Players[i] = g.Players[len(g.Players)-1]
		g.Players[len(g.Players)-1] = nil
		g.Players = g.Players[:len(g.Players)-1]
	}
	delete(g.playerMap, id)
}

func (g *Game) GetPlayer(connID guuid.UUID) *Player {
	return g.playerMap[connID]
}

func (g *Game) NewPlayer(id guuid.UUID) *Player {
	p := &Player{
		UUID:         id,
		X:            rand.Float64() * GameWidth,
		Y:            rand.Float64() * GameHeight,
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
		Status:       WaitForPlay,
		Score:        0,
	}
	g.playerMap[id] = p
	g.Players = append(g.Players, p)
	return p
}

func (g *Game) AddBullet(x, y float64, owner guuid.UUID, rotation RotationDegree, damage float64) {
	bulletID := guuid.New()
	g.Bullets = append(g.Bullets, &Bullet{
		ID:        bulletID,
		X:         x,
		Y:         y,
		Owner:     owner,
		Damage:    damage,
		Rotation:  rotation,
		Speed:     2,
		Exhausted: false,
	})
}

func (g *Game) MoveBullets() {
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
