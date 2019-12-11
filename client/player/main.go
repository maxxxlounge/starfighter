package main

import (
	"encoding/json"
	"flag"
	"image"
	_ "image/png"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	guuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/maxxxlounge/websocket/game"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
)

type CustomConn struct {
	Conn *websocket.Conn
	ID   guuid.UUID
}

var conn *CustomConn
var bulletSprite *pixel.Sprite
var sprite *pixel.Sprite

type Player struct {
	*game.Player
	Sprite *pixel.Sprite
	ID     guuid.UUID
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func ReceiveMessage(g *game.Game) {
	_, message, err := conn.Conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	//log.Printf("recv: %s", message)
	err = json.Unmarshal(message, &g)
	if err != nil {
		err = errors.Wrap(err, "error unmarshalling game object")
		log.Error(err)
		return
	}
}

func run() {
	var g game.Game
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	cfg := pixelgl.WindowConfig{
		Title:  "Starfighter",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	pic, err := loadPicture("./pig.png")
	if err != nil {
		log.Fatal(err)
	}
	bg, err := loadPicture("./bg.png")
	if err != nil {
		log.Fatal(err)
	}
	bullet, err := loadPicture("./bullet.png")
	if err != nil {
		log.Fatal(err)
	}

	sprite = pixel.NewSprite(pic, pic.Bounds())
	bgsprite := pixel.NewSprite(bg, bg.Bounds())
	bulletSprite = pixel.NewSprite(bullet, bullet.Bounds())

	for !win.Closed() {
		win.Clear(colornames.Black)
		bgsprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		ReceiveMessage(&g)
		if g.You != nil {
			if g.You.Life < 0 {
				log.Println("you died!")
				return
			}
			if win.Pressed(pixelgl.KeyLeft) {
				SendInput(conn, pixelgl.KeyLeft.String()+"down")
			}
			if win.JustReleased(pixelgl.KeyLeft) {
				SendInput(conn, pixelgl.KeyLeft.String()+"up")
			}
			if win.Pressed(pixelgl.KeyRight) {
				SendInput(conn, pixelgl.KeyRight.String()+"down")
			}
			if win.JustReleased(pixelgl.KeyRight) {
				SendInput(conn, pixelgl.KeyRight.String()+"up")
			}
			if win.Pressed(pixelgl.KeyDown) {
				SendInput(conn, pixelgl.KeyDown.String()+"down")
			}
			if win.JustReleased(pixelgl.KeyDown) {
				SendInput(conn, pixelgl.KeyDown.String()+"up")
			}
			if win.Pressed(pixelgl.KeyUp) {
				SendInput(conn, pixelgl.KeyUp.String()+"down")
			}
			if win.JustReleased(pixelgl.KeyUp) {
				SendInput(conn, pixelgl.KeyUp.String()+"up")
			}
			if win.Pressed(pixelgl.KeySpace) {
				SendInput(conn, "shoot")
			}
		}

		UpdateGame(win, &g)
		win.Update()
	}
}

func main() {
	var err error
	flag.Parse()

	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8888", Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	conn = &CustomConn{}
	conn.Conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Conn.Close()
	//doing
	pixelgl.Run(run)
}

func SendInput(c *CustomConn, input string) {
	//log.Printf("client send command %s", input)
	err := c.Conn.WriteMessage(websocket.TextMessage, []byte(input))
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func UpdateGame(win *pixelgl.Window, g *game.Game) {
	camPos := pixel.ZV
	if g.You != nil {
		camPos = pixel.V(g.You.X, g.You.Y)
		cam := pixel.IM.Scaled(camPos, 4).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		mat := pixel.IM.Moved(pixel.V(g.You.X, g.You.Y))
		mat = mat.Rotated(pixel.V(g.You.X, g.You.Y), float64(g.You.Rotation))
		sprite.Draw(win, mat)
	}

	for _, p := range g.Players {
		if p != g.You {
			mat := pixel.IM.Moved(pixel.V(p.X, p.Y))
			mat = mat.Rotated(pixel.V(p.X, p.Y), float64(p.Rotation))
			sprite.Draw(win, mat)
		}
	}

	for _, b := range g.Bullets {
		mat := pixel.IM.Moved(pixel.V(b.X, b.Y))
		mat = mat.Rotated(pixel.V(b.X, b.Y), float64(b.Rotation))
		bulletSprite.Draw(win, mat)
	}
}
