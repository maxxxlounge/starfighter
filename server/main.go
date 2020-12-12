package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

var mainGame *game.Game
var connections map[guuid.UUID]*CustomConn

func main() {
	enablelog := false
	if len(os.Args) > 1 && os.Args[1] == "log" {
		enablelog = true
	}
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

	mainGame = game.New()

	go func() {
		Execute(enablelog)
	}()

	fmt.Printf("start listening on %s\n", srv.Addr)
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
	fmt.Printf("incoming connection %s from %s\n", g.String(), cc.Conn.LocalAddr().String())
	defer func(conn *websocket.Conn, g guuid.UUID, game *game.Game) {
		delete(connections, g)
		game.DeletePlayer(g)
		c.Close()
	}(c, g, mainGame)
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
		msg := string(m)
		//fmt.Println(cc.ID.String() + " send message " + msg)
		if strings.Contains(msg, "setup|") {
			p.Name = strings.Replace(msg, "setup|", "", -1)
		}
		switch string(m) {
		case "pause":
			p.Status = game.Pause
			break
		case "resume":
			p.Status = game.Resume
			break
		case "Leftrelease":
			p.Left = false
			break
		case "Leftpressed":
			p.Left = true
			break
		case "LeftUprelease":
			p.Left = false
			p.Up = false
			break
		case "LeftUppressed":
			p.Left = true
			p.Up = true
			break
		case "LeftDownrelease":
			p.Left = false
			p.Down = false
			break
		case "LeftDownpressed":
			p.Left = true
			p.Down = true
			break
		case "Rightrelease":
			p.Right = false
			break
		case "Rightpressed":
			p.Right = true
			break
		case "RightUprelease":
			p.Right = false
			p.Up = false
			break
		case "RightUppressed":
			p.Right = true
			p.Up = true
			break
		case "RightDownrelease":
			p.Right = false
			p.Down = false
			break
		case "RightDownpressed":
			p.Right = true
			p.Down = true
			break
		case "Downrelease":
			p.Down = false
			break
		case "Downpressed":
			p.Down = true
			break
		case "Uprelease":
			p.Up = false
			break
		case "Uppressed":
			p.Up = true
			break
		case "Firepressed":
			p.Fire = true
			break
		case "Firerelease":
			p.Fire = false
			break
		}

	}
}

func Execute(enablelog bool) {
	fmt.Println("executing")
	last := time.Now()
	for {
		dt := time.Since(last).Seconds()
		last = time.Now()
		mainGame.MovePlayers(dt)
		mainGame.MoveBullets()
		mainGame.Collision()

		for _, c := range connections {
			p := mainGame.GetPlayer(c.ID)
			if p == nil {
				return
			}
			if p.Status == game.Pause {
				continue
			}
			mainGame.SetYou(c.ID)
			msg, err := json.Marshal(mainGame)
			if err != nil {
				fmt.Println(err.Error())
			}
			//msg := "send message to conn " + u.String()
			err = c.Conn.WriteMessage(websocket.TextMessage, msg) //[]byte(msg))
			if err != nil {
				fmt.Println(err.Error())
			}
			if enablelog {
				fmt.Printf("%v\n", string(msg))
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
}
