package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/maxxxlounge/websocket/game"

	guuid "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type CustomConn struct {
	Conn *websocket.Conn
	ID   guuid.UUID
}

var mainGame game.Game
var connections map[guuid.UUID]*CustomConn

func main() {
	l := &log.Logger{}
	connections = make(map[guuid.UUID]*CustomConn)
	//server
	r := mux.NewRouter()
	srv := &http.Server{
		Addr:         ":8888",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	r.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		Connect(w, r, l)
	})

	mainGame = game.Game{
		Players: make(map[guuid.UUID]*game.Player),
	}

	go func() {
		Execute()
	}()

	l.Infof("start listening on %s", srv.Addr)
	l.Fatal(srv.ListenAndServe())
}

func Connect(w http.ResponseWriter, r *http.Request, l *log.Logger) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		l.Fatal(err)
		return
	}
	g := guuid.New()
	cc := CustomConn{
		ID:   g,
		Conn: c,
	}
	connections[g] = &cc
	defer func(conn *websocket.Conn, g guuid.UUID, game *game.Game) {
		delete(connections, g)
		delete(game.Players, g)
		c.Close()
	}(c, g, &mainGame)
	p := mainGame.NewPlayer(g)

	for {
		mType, m, err := cc.Conn.ReadMessage()
		if err != nil {
			l.Error(err)
			return
		}
		if mType != websocket.TextMessage {
			continue
		}

		switch string(m) {
		case "Leftup":
			p.Left = false
			break
		case "Leftdown":
			p.Left = true
			break
		case "Rightup":
			p.Right = false
			break
		case "Rightdown":
			p.Right = true
			break
		case "Downup":
			p.Down = false
			break
		case "Downdown":
			p.Down = true
			break
		case "Upup":
			p.Up = false
			break
		case "Updown":
			p.Up = true
			break
		case "shoot":
			mainGame.AddBullet(p.X, p.Y, g, p.Rotation, p.Power)
		}

	}
}

func Execute() {
	fmt.Print("executing")
	last := time.Now()
	for {
		dt := time.Since(last).Seconds()
		last = time.Now()
		mainGame.MovePlayers(dt)
		mainGame.MoveBullets(dt)
		//smainGame.Collision()

		for _, c := range connections {
			mainGame.You = mainGame.Players[c.ID]
			msg, err := json.Marshal(mainGame)
			if err != nil {
				fmt.Println(err.Error())
			}
			//msg := "send message to conn " + u.String()
			err = c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println(err.Error())
			}
			//fmt.Printf("%v", string(msg))
		}
		time.Sleep(20 * time.Millisecond)
	}
}
