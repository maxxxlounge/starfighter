package main

import (
	"encoding/json"
	"flag"
	"image"
	_ "image/png"
	"net/url"
	"os"
	"os/signal"

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
var P Player

type Player struct {
	*game.Player
	Sprite *pixel.Sprite
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

	sprite := pixel.NewSprite(pic, pic.Bounds())
	spriteX := 0.0
	spriteY := 0.0

	P.Sprite = sprite

	//last := time.Now()
	for !win.Closed() {
		_, message, err := conn.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
		err = json.Unmarshal(message, &g)
		if err != nil {
			err = errors.Wrap(err, "error unmarshalling game object")
			log.Error(err)
			win.SetClosed(true)
			return
		}
		win.Clear(colornames.Darkblue)
		sprite.Draw(win, pixel.IM.Moved(pixel.V(spriteX, spriteY)))
		if win.Pressed(pixelgl.KeyLeft) {
			SendInput(conn, pixelgl.KeyLeft.String())
		}
		if win.Pressed(pixelgl.KeyRight) {
			SendInput(conn, pixelgl.KeyRight.String())
		}
		if win.Pressed(pixelgl.KeyDown) {
			SendInput(conn, pixelgl.KeyDown.String())
		}
		if win.Pressed(pixelgl.KeyUp) {
			SendInput(conn, pixelgl.KeyUp.String())
		}
		if win.Pressed(pixelgl.KeySpace) {
			SendInput(conn, pixelgl.KeySpace.String())
		}
		UpdateGame(win, &g)
		win.Update()
	}
}

func main() {
	var err error
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8888", Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	conn = &CustomConn{}
	conn.Conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Conn.Close()

	pixelgl.Run(run)
}

func SendInput(c *CustomConn, input string) {
	log.Printf("client send command %s", input)
	err := c.Conn.WriteMessage(websocket.TextMessage, []byte(input))
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func UpdateGame(win *pixelgl.Window, g *game.Game) {
	for _, p := range g.Players {
		P.Sprite.Draw(win, pixel.IM.Moved(pixel.V(p.X, p.Y)))
	}
}
